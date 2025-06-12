package services

import (
    "bytes"
    "context"
    "encoding/hex"
    "encoding/json"
    "math"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/hospedate/backyard/log"
    "github.com/mr-tron/base58"
)

type BlockchainService interface {
    GetBalance(address string) (float64, error)
    GetBalanceWithRetries(address string, maxRetries int) (float64, error)
}

type BlockchainServiceImp struct {
    logger          log.Logger
    USDTContract    string
    MethodBalanceOf string
    url             string
    api_key_tron    string
}

type PayloadTronGridBalance struct {
    OwnerAddress     string `json:"owner_address"`
    ContractAddress  string `json:"contract_address"`
    FunctionSelector string `json:"function_selector"`
    Parameter        string `json:"parameter"`
    Visible          bool   `json:"visible"`
}

type ResponseDataTronGridBalance struct {
    Result struct {
        Result bool `json:"result"`
    } `json:"result"`
    ConstantResult []string `json:"constant_result"`
}

func NewBlockchainService(USDTContract string, MethodBalanceOf string, url string, api_key_tron string) BlockchainService {
    logger := log.GetOrCreateLogger("BlockchainService", "INFO")
    return &BlockchainServiceImp{
        logger:          logger,
        USDTContract:    USDTContract,
        MethodBalanceOf: MethodBalanceOf,
        url:             url,
        api_key_tron:    api_key_tron,
    }
}

func (b *BlockchainServiceImp) GetBalance(address string) (float64, error) {
    ctx, cancel := createContextWithTimeout(5 * time.Second)
    defer cancel()

    payloadBytes, err := createPayload(address, b.USDTContract, b.MethodBalanceOf)
    if err != nil {
        b.logger.Error("Error marshaling payload:", err.Error())
        return 0, err
    }

    response, err := makeHTTPRequest(ctx, b.url, payloadBytes)
    if err != nil {
        b.logger.Error("Error:", err.Error())
        return 0, err
    }

    balance, err := parseResponseAndCalculateBalance(b.logger, response)
    if err != nil {
        return 0, err
    }

    return balance, nil
}

func createContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), timeout)
}

func createPayload(address, contractAddress, functionSelector string) ([]byte, error) {
    payload := PayloadTronGridBalance{
        OwnerAddress:     address,
        ContractAddress:  contractAddress,
        FunctionSelector: functionSelector,
        Parameter:        addressToParameter(address),
        Visible:          true,
    }

    return json.Marshal(payload)
}

func makeHTTPRequest(ctx context.Context, url string, payloadBytes []byte) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    return client.Do(req)
}

func parseResponseAndCalculateBalance(logger log.Logger, response *http.Response) (float64, error) {
    if response.StatusCode != http.StatusOK {
        logger.Error("Error:", ErrHttpRequestStatus.Error())
        return 0, ErrHttpRequestStatus
    }

    defer response.Body.Close()

    data := new(ResponseDataTronGridBalance)
    json.NewDecoder(response.Body).Decode(data)

    if data.Result.Result {
        val := data.ConstantResult[0]
        intVal, err := strconv.ParseInt(val, 16, 64)
        if err != nil {
            logger.Error("Error:", err.Error())
            return 0, err
        }

        final := float64(intVal) * math.Pow(10, -6)
        return math.Round(final*1e6) / 1e6, nil
    } else {
        logger.Error("Error:", string(ErrKeyNotExistInJson))
        return 0, ErrKeyNotExistInJson
    }
}

func (b *BlockchainServiceImp) GetBalanceWithRetries(address string, maxRetries int) (float64, error) {
    var balance float64
    var err error

    backoffTime := 2 * time.Second
    for i := 0; i < maxRetries; i++ {
        balance, err = b.GetBalance(address)
        if err == nil {
            b.logger.Debugf("Got balance from address %v: %v", address, balance)
            return balance, nil
        }
        b.logger.Errorf("Error retrieving balance (attempt %v/%v): %v\n", i, maxRetries, err.Error())

        // Sleep for backoffTime and increase it
        time.Sleep(backoffTime)
        backoffTime *= 2
    }

    // If we've reached maxRetries, return the error
    return 0, err
}

func addressToParameter(addr string) string {
    addrBytes, _ := base58.Decode(addr)
    return strings.Repeat("0", 24) + hex.EncodeToString((addrBytes[1:]))
}
