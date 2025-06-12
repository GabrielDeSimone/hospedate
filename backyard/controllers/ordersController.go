package controllers

import (
    "math"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
    "github.com/hospedate/backyard/services"
)

type OrdersController interface {
    Create(newOrderRequest *models.NewOrderRequest) (*models.Order, error)
    GetById(id string) *models.Order
    DeleteById(id string) (int64, error)
    Search(searchParams *models.OrdersSearchParams) []*models.Order
    Edit(searchParams models.OrderEditRequest) (*models.Order, error)
}

type OrdersControllerImpl struct {
    repository           repositories.OrdersRepository
    propertiesRepository repositories.PropertiesRepository
    usersRepository      repositories.UsersRepository
    paymentsRepository   repositories.PaymentsRepository
    paymentsService      services.PaymentsService
    logger               log.Logger
    emailSender          services.EmailNotificationService
    travelerOrderFee     float32
}

func NewOrdersController(
    repository repositories.OrdersRepository,
    propertiesRepository repositories.PropertiesRepository,
    usersRepository repositories.UsersRepository,
    paymentsRepository repositories.PaymentsRepository,
    paymentsService services.PaymentsService,
    emailSender services.EmailNotificationService,
    travelerOrderFee float32,
) OrdersController {
    return &OrdersControllerImpl{
        repository:           repository,
        propertiesRepository: propertiesRepository,
        usersRepository:      usersRepository,
        paymentsRepository:   paymentsRepository,
        paymentsService:      paymentsService,
        logger:               log.GetOrCreateLogger("OrdersController", "INFO"),
        emailSender:          emailSender,
        travelerOrderFee:     travelerOrderFee,
    }
}

func (c *OrdersControllerImpl) totalBilledCents(newOrderRequest *models.NewOrderRequest, property *models.Property) uint {
    subTotal := models.ComputeSubtotalCents(newOrderRequest.DateStart, newOrderRequest.DateEnd, property.Price)
    if newOrderRequest.OrderType == "in_platform" {
        return uint(math.Round(float64(1.0+c.travelerOrderFee) * float64(subTotal)))
    }
    return subTotal
}

func (c *OrdersControllerImpl) Create(newOrderRequest *models.NewOrderRequest) (*models.Order, error) {
    user := c.usersRepository.GetById(newOrderRequest.UserId)
    if user == nil {
        return nil, ErrUserDoesNotExist
    }
    property := c.propertiesRepository.GetById(newOrderRequest.PropertyId)
    if property == nil {
        return nil, ErrPropertyDoesNotExist
    }
    if property.UserId == user.Id {
        return nil, ErrOwnerBookingProperty
    }
    if property.MaxGuests < newOrderRequest.NumberGuests {
        return nil, ErrForTooManyGuests
    }
    if property.Status == "archived" {
        return nil, ErrPropertyArchivedStatus
    }
    var order *models.Order
    var err error
    totalCents := c.totalBilledCents(newOrderRequest, property)
    if newOrderRequest.OrderType == "in_platform" {
        encryptedPk, address, err := c.paymentsService.CreateWallet()
        if err != nil {
            c.logger.Error("Cannot create wallet, will not create order", err.Error())
            return nil, err
        }
        order, err = c.repository.SaveInPlatform(newOrderRequest, property, c.paymentsRepository, address, encryptedPk, totalCents)
        if err == repositories.ErrCollision {
            c.logger.Error("Cannot create in platform order, collision", err.Error())
            return nil, ErrCollision
        } else if err != nil {
            c.logger.Error("Cannot create in platform order", err.Error())
            return nil, err
        }

    } else {
        order, err = c.repository.Save(newOrderRequest, property, totalCents)
        if err != nil {
            c.logger.Error("Cannot create owner_directly order", err.Error())
            return nil, err
        }
        c.emailSender.SendOrderPendingOwnerNotification(order.Id, order.PropertyId)
        c.emailSender.SendOrderPendingTravelerNotification(order.Id, order.UserId)
    }

    return order, nil
}

func (c *OrdersControllerImpl) GetById(id string) *models.Order {
    return c.repository.GetById(id)
}

func (c *OrdersControllerImpl) DeleteById(id string) (int64, error) {
    return c.repository.DeleteById(id)
}

func (c *OrdersControllerImpl) Search(searchParams *models.OrdersSearchParams) []*models.Order {
    return c.repository.Search(searchParams)
}

func (c *OrdersControllerImpl) Edit(orderEditParams models.OrderEditRequest) (*models.Order, error) {
    order, err := c.repository.Edit(orderEditParams)
    if err == nil && orderEditParams.Status != nil && *orderEditParams.Status == "canceled" {
        ownerID := c.propertiesRepository.GetById(c.GetById(orderEditParams.Id).PropertyId).UserId
        travelerID := c.GetById(orderEditParams.Id).UserId
        c.emailSender.SendOrderCanceledNotification(orderEditParams.Id, ownerID)
        c.emailSender.SendOrderCanceledNotification(orderEditParams.Id, travelerID)
    }
    if err == nil && orderEditParams.Status != nil && *orderEditParams.Status == "confirmed" {
        travelerID := c.GetById(orderEditParams.Id).UserId
        c.emailSender.SendOrderConfirmedTravelerNotification(orderEditParams.Id, travelerID)
    }
    return order, err
}
