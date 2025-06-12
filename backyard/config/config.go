package config

import (
    "log"
    "os"
    "strconv"
    "strings"
    "time"
)

type Config struct {
    Server               *ServerConfig
    Db                   *DBConfig
    BlockchainService    *BlockchainServiceConfig
    LogLevel             string
    BackyardVersion      string
    ProductVersion       string
    AirbnbFetcherService *AirbnbFetcherServiceConfig
    PaymentsService      *PaymentsServiceConfig
    OrderUpdateService   *OrderUpdateServiceConfig
    EmailService         *EmailServiceConfig
    InPlatformOrderFees  *FeesConfig
    Invitations          *InvitationsConfig
}

type ServerConfig struct {
    Host string
    Port string
}

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Database string
    Init_db  InitDBOption
}

type AirbnbFetcherServiceConfig struct {
    Disabled   bool
    PathPrefix string
}

type BlockchainServiceConfig struct {
    USDTContract    string
    MethodBalanceOf string
    Url             string
    Api_key_tron    string
}

type OrderUpdateServiceConfig struct {
    SleepCheckUpdate time.Duration
    SleepCheckCancel time.Duration
    TimeToCancel     time.Duration
    TimeToInProgress time.Duration
    TimeToCompleted  time.Duration
}

type EmailServiceConfig struct {
    PathPrefix              string
    SenderAddress           string
    EMAILS_CHANNEL_CAPACITY int32
    TimeOutEmailSend        time.Duration
    Disabled                bool
}

type PaymentsServiceConfig struct {
    EncryptionKey string
    Disabled      bool
}

type InvitationsConfig struct {
    CreditForVerifiedProperty float32
    CreditForInPlatformOrder  float32
}

type FeesConfig struct {
    TravelerOrderFee float32
    OwnerOrderFee    float32
}

type InitDBOption string

const (
    InitDBOptionMigration InitDBOption = "MIGRATION"
    InitDBOptionHard      InitDBOption = "HARD"
    InitDBOptionNone      InitDBOption = "NONE"
)

func ConvertToBool(value string) bool {
    result, err := strconv.ParseBool(value)
    if err != nil {
        log.Fatal("Error converting string to bool:", err)
    }
    return result
}

func ConvertToFloat(s string) float32 {
    f, err := strconv.ParseFloat(s, 32)
    if err != nil {
        log.Fatal("Error converting string to float:", err)
    }
    return float32(f)
}

func ConvertToDuration(s string, unit time.Duration) time.Duration {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        log.Fatal("Error converting string to time.Duration:", err)
    }
    return time.Duration(i) * unit
}

func GetConfig() *Config {

    return &Config{
        Server: &ServerConfig{
            Host: getOrElse("BACKYARD_HOST", "localhost"),
            Port: getOrElse("BACKYARD_PORT", "8080"),
        },
        Db: &DBConfig{
            Host:     getOrElse("POSTGRES_HOST", "localhost"),
            Port:     getOrElse("POSTGRES_PORT", "5432"),
            User:     getOrElse("POSTGRES_USER", "backyard"),
            Password: getOrElse("POSTGRES_PASS", "backyard"),
            Database: getOrElse("POSTGRES_DB_NAME", "backyarddb"),
            Init_db:  loadInitDBOption(getOrElse("POSTGRES_INIT_DB", "NONE")),
        },
        AirbnbFetcherService: &AirbnbFetcherServiceConfig{
            Disabled:   ConvertToBool(getOrElse("AIRBNB_FETCHER_SERVICE_DISABLED", "false")),
            PathPrefix: getOrElse("AIRBNB_FETCHER_SERVICE_PATH_PREFIX", "/opt"),
        },
        BlockchainService: &BlockchainServiceConfig{
            USDTContract:    "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
            MethodBalanceOf: "balanceOf(address)",
            Url:             "https://api.trongrid.io/wallet/triggerconstantcontract",
            Api_key_tron:    "25f66928-0b70-48cd-9ac6-da6f8247c663",
        },
        LogLevel:        getOrElse("LOG_LEVEL", "INFO"),
        BackyardVersion: getOrElse("BACKYARD_VERSION", "unknown"),
        ProductVersion:  getOrElse("PRODUCT_VERSION", "unknown"),
        PaymentsService: &PaymentsServiceConfig{
            EncryptionKey: loadEncryptionKey("ENCRYPTION_KEY", "DEV_ENCRYPTION_KEY"),
            Disabled:      ConvertToBool(getOrElse("PAYMENTS_SERVICE_DISABLED", "false")),
        },
        OrderUpdateService: &OrderUpdateServiceConfig{
            SleepCheckUpdate: ConvertToDuration(getOrElse("SLEEP_CHECK_UPDATE", "86400"), time.Second),
            SleepCheckCancel: ConvertToDuration(getOrElse("SLEEP_CHECK_CANCEL", "1200"), time.Second),
            TimeToCancel:     ConvertToDuration(getOrElse("TIME_TO_CANCEL", "5"), time.Hour),
            TimeToInProgress: ConvertToDuration(getOrElse("TIME_TO_IN_PROGRESS", "36"), time.Hour),
            TimeToCompleted:  ConvertToDuration(getOrElse("TIME_TO_COMPLETED", "0"), time.Hour),
        },
        EmailService: &EmailServiceConfig{
            PathPrefix:              getOrElse("EMAIL_SERVICE_PATH_PREFIX", "/opt"),
            SenderAddress:           "Hospedate <info@hospedate.app>",
            EMAILS_CHANNEL_CAPACITY: 1000,
            TimeOutEmailSend:        10 * time.Second,
            Disabled:                ConvertToBool(getOrElse("EMAIL_SERVICE_DISABLED", "false")),
        },
        InPlatformOrderFees: &FeesConfig{
            TravelerOrderFee: ConvertToFloat(getOrElse("TRAVELER_ORDER_FEE", "0.07")),
            OwnerOrderFee:    ConvertToFloat(getOrElse("OWNER_ORDER_FEE", "0.015")),
        },
        Invitations: &InvitationsConfig{
            CreditForVerifiedProperty: ConvertToFloat(getOrElse("CREDIT_FOR_VERIFIED_PROPERTY", "4.0")),
            CreditForInPlatformOrder:  ConvertToFloat(getOrElse("CREDIT_FOR_PLATFORM_ORDER", "2.0")),
        },
    }
}

func getOrElse(envVarName string, defaultValue string) string {
    envVar := os.Getenv(envVarName)
    if envVar == "" {
        return defaultValue
    }
    return envVar
}

func loadInitDBOption(value string) InitDBOption {
    switch strings.ToUpper(value) {
    case string(InitDBOptionMigration):
        return InitDBOptionMigration
    case string(InitDBOptionHard):
        return InitDBOptionHard
    default:
        return InitDBOptionNone
    }
}

func loadEncryptionKey(envVarName string, defaultValue string) string {
    value := getOrElse(envVarName, defaultValue)
    if len(value) != 32 {
        log.Fatal("encryption key must be exactly 32 chars long")
    }
    return value
}
