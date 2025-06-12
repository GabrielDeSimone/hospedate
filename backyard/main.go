package main

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/hospedate/backyard/config"
    "github.com/hospedate/backyard/controllers"
    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/middlewares"
    "github.com/hospedate/backyard/repositories"
    "github.com/hospedate/backyard/routes"
    "github.com/hospedate/backyard/services"

    _ "github.com/lib/pq"
)

var logger log.Logger

func connectWithDB(dbConfig *config.DBConfig) *sql.DB {
    db_url := fmt.Sprintf("postgres://%s:%s/%s?sslmode=disable&user=%s&password=%s", dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.User, dbConfig.Password)
    logger.Infof("Connecting to database at %v:%v", dbConfig.Host, dbConfig.Port)
    attempts := 0

    var db *sql.DB
    var err error

    // While there are errors connecting to the database, keep trying
    for {
        db, err = sql.Open("postgres", db_url)
        if err != nil {
            attempts++
            logger.Errorf("%v - Attempt %v - Retrying in 5 seconds...", err.Error(), attempts)
            time.Sleep(5 * time.Second)
            continue
        }

        err = db.Ping()
        if err != nil {
            attempts++
            logger.Errorf("%v - Attempt %v - Retrying in 5 seconds...", err.Error(), attempts)
            time.Sleep(5 * time.Second)
            continue
        }
        break
    }

    return db
}

func initDB(globalRepository repositories.GlobalRepository, dbConfig *config.DBConfig) {
    switch dbConfig.Init_db {
    case config.InitDBOptionMigration:
        globalRepository.Migration()
    case config.InitDBOptionHard:
        globalRepository.HardInitDB()
    default:
        logger.Info("Skipping DB initialization")
    }
}

func main() {
    // load config
    config_v := config.GetConfig()
    // Initialize the logger
    logger = log.NewLogger("main", config_v.LogLevel)
    logger.Infof("Backyard version %v | Product version %v", config_v.BackyardVersion, config_v.ProductVersion)

    db := connectWithDB(config_v.Db)
    logger.Infof("Connected with the DB at %v:%v", config_v.Db.Host, config_v.Db.Port)
    defer db.Close()

    // Initialize global repository and init DB
    globalRepository := repositories.NewGlobalRepository(db)
    initDB(globalRepository, config_v.Db)

    // Initialize repositories
    usersRepository := repositories.NewUsersRepository(db)
    invitationsRepository := repositories.NewInvitationsRepository(db)
    propertiesRepository := repositories.NewPropertiesRepository(db)
    ordersRepository := repositories.NewOrdersRepository(db)
    paymentsRepository := repositories.NewPaymentsRepository(db)
    ownersEarnedRepository := repositories.NewOwnersEarnedRepository(db)
    usersCreditRepository := repositories.NewUsersCreditRepository(db)
    logger.Debug("Repositories initialized")

    // Initialize services
    emailSender := services.NewEmailNotificationService(usersRepository,
        propertiesRepository,
        config_v.EmailService.SenderAddress,
        config_v.EmailService.EMAILS_CHANNEL_CAPACITY,
        config_v.EmailService.TimeOutEmailSend,
        config_v.EmailService.Disabled,
        config_v.EmailService.PathPrefix,
    )
    emailSender.Start()
    airbnbFetcher := services.NewAirbnbFetcher(
        propertiesRepository,
        config_v.AirbnbFetcherService.Disabled,
        config_v.AirbnbFetcherService.PathPrefix,
    )
    airbnbFetcher.Start()
    blockchainService := services.NewBlockchainService(config_v.BlockchainService.USDTContract,
        config_v.BlockchainService.MethodBalanceOf,
        config_v.BlockchainService.Url,
        config_v.BlockchainService.Api_key_tron)
    usersCreditService := services.NewUsersCreditService(usersCreditRepository,
        invitationsRepository,
        usersRepository,
        emailSender,
        config_v.Invitations.CreditForInPlatformOrder,
        config_v.Invitations.CreditForVerifiedProperty)
    paymentsService := services.NewPaymentService(
        ordersRepository,
        paymentsRepository,
        blockchainService,
        config_v.PaymentsService.EncryptionKey,
        config_v.PaymentsService.Disabled,
        emailSender,
    )
    paymentsService.Start()

    OrderUpdateService := services.NewOrderUpdateService(
        ordersRepository,
        ownersEarnedRepository,
        propertiesRepository,
        config_v.OrderUpdateService.SleepCheckUpdate,
        config_v.OrderUpdateService.SleepCheckCancel,
        config_v.OrderUpdateService.TimeToCancel,
        config_v.OrderUpdateService.TimeToInProgress,
        config_v.OrderUpdateService.TimeToCompleted,
        config_v.InPlatformOrderFees.OwnerOrderFee,
        emailSender,
        usersCreditService,
    )
    OrderUpdateService.Start()
    // Initialize controllers
    usersController := controllers.NewUsersController(usersRepository, ownersEarnedRepository, usersCreditRepository)
    invitationsController := controllers.NewInvitationsController(invitationsRepository)
    propertiesController := controllers.NewPropertiesController(
        propertiesRepository,
        ordersRepository,
        airbnbFetcher,
        usersCreditService,
    )
    ordersController := controllers.NewOrdersController(
        ordersRepository,
        propertiesRepository,
        usersRepository,
        paymentsRepository,
        paymentsService,
        emailSender,
        config_v.InPlatformOrderFees.TravelerOrderFee)
    paymentsController := controllers.NewPaymentsController(paymentsRepository, ordersRepository)
    logger.Debug("Controllers initialized")

    // Initialize the router
    router := gin.New()
    router.Use(middlewares.GinLogger())
    router.Use(gin.Recovery())

    // Initialize the routes
    routes.UsersRoutes(router, usersController, invitationsController)
    routes.InvitationsRoutes(router, invitationsController)
    routes.PropertiesRoutes(router, propertiesController)
    routes.OrdersRoutes(router, ordersController)
    routes.PaymentsRoutes(router, paymentsController)
    routes.AirbnbFetcherRoutes(router, airbnbFetcher)
    routes.NotificationServiceRoutes(router, emailSender, usersController)
    logger.Debug("Routes initialized")

    // Start the server
    router.Run(fmt.Sprintf("%v:%v", config_v.Server.Host, config_v.Server.Port))
}
