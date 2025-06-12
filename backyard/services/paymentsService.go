package services

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "sort"
    "time"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

type checkInstance struct {
    CheckAt       time.Time
    WalletAddress string
    OrderId       string
}

type PaymentsService interface {
    Start()
    CreateWallet() (encryptedPk string, address *models.Address, error error)
}

type PaymentsServiceImpl struct {
    updateChannel     chan []checkInstance
    ordersRepo        repositories.OrdersRepository
    paymentsRepo      repositories.PaymentsRepository
    blockchainService BlockchainService
    logger            log.Logger
    encryptionKey     string
    disabled          bool
    emailSender       EmailNotificationService
}

func NewPaymentService(
    ordersRepo repositories.OrdersRepository,
    paymentsRepo repositories.PaymentsRepository,
    Blockchain BlockchainService,
    encryptionKey string,
    disabled bool,
    emailSender EmailNotificationService,
) PaymentsService {
    logger := log.GetOrCreateLogger("PaymentsService", "INFO")
    return &PaymentsServiceImpl{
        ordersRepo:        ordersRepo,
        paymentsRepo:      paymentsRepo,
        logger:            logger,
        blockchainService: Blockchain,
        updateChannel:     make(chan []checkInstance),
        encryptionKey:     encryptionKey,
        disabled:          disabled,
        emailSender:       emailSender,
    }
}

func (p *PaymentsServiceImpl) Start() {
    if p.disabled {
        p.logger.Info("Skipping Payments Service initialization (disabled = true)")
    } else {
        go p.listUpdate()
        go p.paymentChecker()
        p.logger.Info("Payments service initialized")
    }
}

func (p *PaymentsServiceImpl) CreateWallet() (encryptedPk string, address *models.Address, error error) {
    pk, err := models.NewPrivateKey()
    if err != nil {
        p.logger.Error("Error when generating privateKey", err.Error())
        return "", nil, UnknownError
    }

    pkEncrypted, err := encryptAES([]byte(p.encryptionKey), pk.String())
    if err != nil {
        p.logger.Error("Error when encrypting privateKey", err.Error())
        return "", nil, UnknownError
    }

    pub, err := pk.GetPublicKey()
    if err != nil {
        p.logger.Error("Error casting public key to ECDSA", err.Error())
        return "", nil, UnknownError
    }

    addressStruct, err := pub.GetAddress()
    if err != nil {
        p.logger.Error("Cannot get address from public key", err.Error())
        return "", nil, UnknownError
    }
    return pkEncrypted, addressStruct, nil

}

func (p *PaymentsServiceImpl) RetrievePrivateKey(order_id string) (string, error) {
    key := []byte(p.encryptionKey)

    // recover encrypted private key from database
    encryptedPK, err := p.paymentsRepo.GetPkByOrderId(order_id)
    if err != nil {
        p.logger.Error("Error when retrieving encrypted privateKey", err.Error())
        return "", err
    }

    // Decrypt the private key
    privateKey, err := decryptAES(key, encryptedPK)
    if err != nil {
        p.logger.Error("Error when decrypting privateKey", err.Error())
        return "", err
    }

    return privateKey, nil
}

func (p *PaymentsServiceImpl) updateOrderStatusIfExpired(order *models.Order) error {
    now := time.Now().UTC()
    if now.Sub(order.CreatedAt) >= 1*time.Hour {
        // If an hour has been reached since order creation, we update it in the database as discarded
        status := p.ordersRepo.GetById(order.Id).Status
        if status == "ephemeral" {
            p.logger.Infof("Discarding order %v created at %v", order.Id, order.CreatedAt)
            discarded := "discarded"
            err := updateOrderStatus(p.ordersRepo, p.logger, order.Id, discarded)
            if err != nil {
                return err
            }
            return nil
        }
    }
    return nil
}

func (p *PaymentsServiceImpl) listUpdate() {
    for {
        orders := p.ordersRepo.GetEphemeralOrders()
        if len(orders) > 0 {
            p.logger.Infof("Got %v ephemeral orders from the DB", len(orders))
        } else {
            p.logger.Debugf("Got %v ephemeral orders from the DB", len(orders))
        }

        var checkInstances []checkInstance
        for _, order := range orders {
            checkInstances = p.generateCheckInstances(order, checkInstances)
        }

        couldSend := true
        if len(checkInstances) > 0 {
            couldSend = p.sendCheckInstances(checkInstances)
        }
        if couldSend {
            time.Sleep(65 * time.Second)
        }
    }
}

func (p *PaymentsServiceImpl) generateCheckInstances(order *models.Order, checkInstances []checkInstance) []checkInstance {
    now := time.Now().UTC()
    nextCheck := now
    for {
        nextCheck = nextCheck.Add(getNextCheckDuration(nextCheck, order.CreatedAt))
        if nextCheck.After(now.Add(1 * time.Minute)) {
            break // We don't send timestamps further than one minute in the future.
        }
        checkInstances = append(checkInstances, checkInstance{
            CheckAt:       nextCheck,
            WalletAddress: order.WalletAddress,
            OrderId:       order.Id,
        })
    }
    return checkInstances
}

func (p *PaymentsServiceImpl) sendCheckInstances(checkInstances []checkInstance) bool {
    sort.Slice(checkInstances, func(i, j int) bool {
        return checkInstances[i].CheckAt.Before(checkInstances[j].CheckAt)
    })
    p.logger.Debugf("checkInstances list prepared: %v", checkInstances)
    select {
    case p.updateChannel <- checkInstances:
        p.logger.Debugf("list of checkInstances sent: %v", checkInstances)
        // The list was sent successfully
        return true
    case <-time.After(65 * time.Second):
        p.logger.Debugf("Could not sent list through the channel after waiting for 1 minute. Moving on. List: %v", checkInstances)
        // Could not send the list because there are no readers,
        // then we do nothing and continue with the next loop
        return false
    }
}

// getNextCheckDuration returns the time lapse until the next check.
func getNextCheckDuration(lastCheck, orderCreatedAt time.Time) time.Duration {
    elapsed := lastCheck.Sub(orderCreatedAt)
    if elapsed < 20*time.Minute {
        return 15 * time.Second
    } else if elapsed < 40*time.Minute {
        return 30 * time.Second
    } else {
        return 1 * time.Minute
    }
}

func (p *PaymentsServiceImpl) updatePaymentStatus(orderId string, receivedAmountCents uint, receivedCurrency string, status *string) error {
    paymentEditRequest := models.PaymentEditRequestByOrderId{
        Id:                  orderId,
        ReceivedAmountCents: &receivedAmountCents,
        ReceivedCurrency:    &receivedCurrency,
        Status:              status,
    }
    _, err := p.paymentsRepo.EditByOrderId(paymentEditRequest)
    if err != nil {
        p.logger.Error("Error updating payment status:", err.Error())
    }
    return err
}

func (p *PaymentsServiceImpl) processCheckInstance(checkInstance checkInstance, wasOrderPaid map[string]bool) error {
    if _, ok := wasOrderPaid[checkInstance.OrderId]; ok {
        p.logger.Infof("Order %v already paid, skipping check", checkInstance.OrderId)
        return nil
    }
    p.logger.Infof("Processing checkInstance %v for order id %v", checkInstance.CheckAt, checkInstance.OrderId)
    balance_f, err := p.blockchainService.GetBalanceWithRetries(checkInstance.WalletAddress, 2)
    if err != nil {
        p.logger.Errorf("Could not check balance of wallet %v", checkInstance.WalletAddress)
        return err
    }

    order := p.ordersRepo.GetById(checkInstance.OrderId)
    if order == nil {
        p.logger.Errorf("Could not check balance because order with id %v was not found", checkInstance.OrderId)
        return ErrOrderNotFound
    }
    balanceCents := uint(balance_f * 100)
    if balanceCents >= order.TotalBilledCents {
        p.logger.Infof(
            "Balance obtained ($ %v) for order %v satisfied the total billed amount ($ %v). Updating order and payment data",
            balanceCents/100,
            checkInstance.OrderId,
            order.TotalBilledCents/100,
        )
        if err := updateOrderStatus(p.ordersRepo, p.logger, order.Id, "pending"); err != nil {
            p.logger.Errorf("Could not update status of order with id %v: %v", order.Id, err.Error())
            return err
        }
        confirmed_status := "confirmed"
        if err := p.updatePaymentStatus(order.Id, balanceCents, "USDT", &confirmed_status); err != nil {
            p.logger.Errorf("Could not update to confirmed status payment with order id %v: %v", order.Id, err.Error())
            return err
        }
        p.emailSender.SendOrderPendingOwnerNotification(order.Id, order.PropertyId)
        p.emailSender.SendPaymentReceivedNotification(order.Id, order.UserId)
        wasOrderPaid[order.Id] = true
    } else {
        p.logger.Infof(
            "Balance obtained ($ %v) for order %v did not satisfy the total billed amount ($ %v)",
            balanceCents/100,
            checkInstance.OrderId,
            order.TotalBilledCents/100,
        )
        if err := p.updatePaymentStatus(order.Id, balanceCents, "USDT", nil); err != nil {
            p.logger.Errorf("Could not update balance of payment with order id %v: %v", order.Id, err.Error())
            return err
        }
        err := p.updateOrderStatusIfExpired(order)
        if err != nil {
            p.logger.Errorf("Error editing status of order %v to expired: %v", order.Id, err.Error())
        }
    }
    return nil
}

func (p *PaymentsServiceImpl) paymentChecker() {
    for {
        select {
        case checkInstances := <-p.updateChannel:
            wasOrderPaid := make(map[string]bool)
            p.logger.Debugf("List of checkInstances received: %v", checkInstances)
            if len(checkInstances) > 0 {
                p.logger.Infof("Prepared to process %v checkInstances", len(checkInstances))
            } else {
                p.logger.Debugf("Prepared to process %v checkInstances", len(checkInstances))
            }

            for _, checkInstance := range checkInstances {
                remainingTime := checkInstance.CheckAt.Sub(time.Now().UTC())
                if remainingTime > 0 {
                    time.Sleep(remainingTime)
                }
                err := p.processCheckInstance(checkInstance, wasOrderPaid)
                if err != nil {
                    p.logger.Errorf("Could not process checkInstance %v: %v", checkInstance, err.Error())
                }
            }
        }
    }
}

func encryptAES(encryptionKey []byte, plaintext string) (string, error) {
    cipherBlock, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return "", err
    }

    plaintextBase64 := base64.StdEncoding.EncodeToString([]byte(plaintext))
    buffer := make([]byte, aes.BlockSize+len(plaintextBase64))

    initializationVector := buffer[:aes.BlockSize]

    cfbEncrypter := cipher.NewCFBEncrypter(cipherBlock, initializationVector)
    cfbEncrypter.XORKeyStream(buffer[aes.BlockSize:], []byte(plaintextBase64))

    return base64.StdEncoding.EncodeToString(buffer), nil
}

func decryptAES(key []byte, cryptoText string) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    decodeCryptoText, err := base64.StdEncoding.DecodeString(cryptoText)
    if err != nil {
        return "", err
    }
    if len(decodeCryptoText) < aes.BlockSize {
        return "", err
    }
    iv := decodeCryptoText[:aes.BlockSize]
    cryptoTextBytes := decodeCryptoText[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(cryptoTextBytes, cryptoTextBytes)
    decryptedText, err := base64.StdEncoding.DecodeString(string(cryptoTextBytes))
    if err != nil {
        return "", err
    }
    return string(decryptedText), nil
}
