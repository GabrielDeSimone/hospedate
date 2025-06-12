package services

import (
    "math"
    "time"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

type OrderUpdateService interface {
    Start()
}

type statusUpdaterFunc func(order *models.Order) error

type OrderUpdateServiceImpl struct {
    ordersRepo         repositories.OrdersRepository
    ownersEarnedRepo   repositories.OwnersEarnedRepository
    propertiesRepo     repositories.PropertiesRepository
    logger             log.Logger
    sleepCheckUpdate   time.Duration
    sleepCheckCancel   time.Duration
    timeToCancel       time.Duration
    timeToInProgress   time.Duration
    timeToCompleted    time.Duration
    OwnerOrderFee      float32
    emailSender        EmailNotificationService
    usersCreditService UsersCreditService
}

func NewOrderUpdateService(
    ordersRepo repositories.OrdersRepository,
    ownersEarnedRepo repositories.OwnersEarnedRepository,
    propertiesRepo repositories.PropertiesRepository,
    sleepCheckUpdate time.Duration,
    sleepCheckCancel time.Duration,
    timeToCancel time.Duration,
    timeToInProgress time.Duration,
    timeToCompleted time.Duration,
    OwnerOrderFee float32,
    emailSender EmailNotificationService,
    usersCreditService UsersCreditService,
) OrderUpdateService {
    logger := log.GetOrCreateLogger("OrderUpdateService", "INFO")
    return &OrderUpdateServiceImpl{
        ordersRepo:         ordersRepo,
        ownersEarnedRepo:   ownersEarnedRepo,
        propertiesRepo:     propertiesRepo,
        logger:             logger,
        sleepCheckUpdate:   sleepCheckUpdate,
        sleepCheckCancel:   sleepCheckCancel,
        timeToCancel:       timeToCancel,
        timeToInProgress:   timeToInProgress,
        timeToCompleted:    timeToCompleted,
        OwnerOrderFee:      OwnerOrderFee,
        emailSender:        emailSender,
        usersCreditService: usersCreditService,
    }
}

func (s *OrderUpdateServiceImpl) Start() {
    go s.newOrderStatus()
    go s.cancelPendingOrder()
    s.logger.Info("Order Update service initialized")
}

func (s *OrderUpdateServiceImpl) EarnedAmount(order *models.Order, property *models.Property) uint {
    subTotalCents := models.ComputeSubtotalCents(order.DateStart, order.DateEnd, property.Price)
    ownerFeeDeduction := uint(math.Round(float64(subTotalCents) * float64(s.OwnerOrderFee)))

    return subTotalCents - ownerFeeDeduction
}

func (s *OrderUpdateServiceImpl) sendToOwnersEarned(order *models.Order) error {
    property := s.propertiesRepo.GetById(order.PropertyId)
    newOwnerEarnedInstanceRequest := models.NewOwnerEarnedInstanceRequest{OrderId: order.Id, UserId: property.UserId,
        EarnedAmountCents: s.EarnedAmount(order, property),
        EarnedCurrency:    "USDT"}
    _, err := s.ownersEarnedRepo.Save(&newOwnerEarnedInstanceRequest)
    if err != nil {
        s.logger.Errorf("Error saving order %v to owners earned: %v", order.Id, err.Error())
        return err
    }
    userToNotify := s.propertiesRepo.GetById(order.PropertyId).UserId
    s.emailSender.SendOwnerFundsNotification(order.Id, newOwnerEarnedInstanceRequest.EarnedAmountCents, userToNotify)
    return nil

}

func (s *OrderUpdateServiceImpl) processOrdersWithStatus(status string, updater statusUpdaterFunc) {
    searchParams := models.OrdersSearchParams{Status: &status}
    orders := s.ordersRepo.Search(&searchParams)

    if len(orders) > 0 {
        s.logger.Infof("Got %v %v orders from the DB", len(orders), status)
    } else {
        s.logger.Debugf("Got %v %v orders from the DB", len(orders), status)
    }

    for _, order := range orders {
        err := updater(order)
        if err != nil {
            s.logger.Errorf("Error updating order %v : %v", order.Id, err.Error())
        }
    }
}

func (s *OrderUpdateServiceImpl) checkUpdateToInProgress(order *models.Order) error {
    if time.Now().UTC().Sub(order.DateStart.AsTime()) >= s.timeToInProgress {
        s.logger.Infof("Updating to in_progress order %v created at %v", order.Id, order.CreatedAt)
        err := updateOrderStatus(s.ordersRepo, s.logger, order.Id, "in_progress")
        if err != nil {
            return err
        } else if order.OrderType == "in_platform" {
            s.usersCreditService.HandleCreditTravelerEvent(order.UserId)
            return s.sendToOwnersEarned(order)
        }
    }
    return nil
}

func (s *OrderUpdateServiceImpl) checkUpdateToCompleted(order *models.Order) error {
    if time.Now().UTC().Sub(order.DateEnd.AsTime()) >= s.timeToCompleted {
        s.logger.Infof("Updating to completed order %v created at %v", order.Id, order.CreatedAt)
        err := updateOrderStatus(s.ordersRepo, s.logger, order.Id, "completed")
        if err != nil {
            return err
        }
    }
    return nil
}

func orderInPendingStatus(orderCreated time.Time) time.Duration {
    // This function returns an estimation of how much time
    // a given order has been with the "pending" status
    // based on the time the order was created.
    // Prerequisite: the order must be in "pending" status
    // when this function is called.

    estimation := time.Now().UTC().Sub(orderCreated) - (1 * time.Hour)

    if estimation < 0 {
        return 0
    } else {
        return estimation
    }
}

func (s *OrderUpdateServiceImpl) UpdateToCanceled(order *models.Order) error {
    orderHasExpired := orderInPendingStatus(order.CreatedAt) >= s.timeToCancel

    if orderHasExpired {
        s.logger.Infof("Updating to canceled order %v created at %v", order.Id, order.CreatedAt)
        canceled := "canceled"
        canceledBy := "owner"
        orderEditRequest := models.OrderEditRequest{
            Id:         order.Id,
            Status:     &canceled,
            CanceledBy: &canceledBy,
        }
        _, err := s.ordersRepo.Edit(orderEditRequest)
        if err != nil {
            s.logger.Error("Error updating order status:", err.Error())
        } else {
            s.emailSender.SendOrderCanceledNotification(order.Id, order.UserId)
        }
        return err
    }
    return nil
}

func (s *OrderUpdateServiceImpl) newOrderStatus() {
    for {
        s.processOrdersWithStatus("confirmed", s.checkUpdateToInProgress)
        s.processOrdersWithStatus("in_progress", s.checkUpdateToCompleted)
        time.Sleep(s.sleepCheckUpdate)
    }
}

func (s *OrderUpdateServiceImpl) cancelPendingOrder() {
    for {
        s.processOrdersWithStatus("pending", s.UpdateToCanceled)
        time.Sleep(s.sleepCheckCancel)
    }
}
