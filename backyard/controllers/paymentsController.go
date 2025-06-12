package controllers

import (
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

type PaymentsController interface {
    Create(newPaymentRequest *models.NewPaymentRequest) (*models.Payment, error)
    GetById(id string) *models.Payment
    DeleteById(id string) (int64, error)
    Search(searchParams *models.PaymentsSearchParams) []*models.Payment
    Edit(searchParams models.PaymentEditRequest) (*models.Payment, error)
}

type PaymentsControllerImpl struct {
    repository       repositories.PaymentsRepository
    ordersRepository repositories.OrdersRepository
}

func NewPaymentsController(
    repository repositories.PaymentsRepository,
    ordersRepository repositories.OrdersRepository,
) PaymentsController {
    return &PaymentsControllerImpl{
        repository:       repository,
        ordersRepository: ordersRepository,
    }
}

func (c *PaymentsControllerImpl) Create(newPaymentRequest *models.NewPaymentRequest) (*models.Payment, error) {
    order := c.ordersRepository.GetById(newPaymentRequest.OrderId)
    if order == nil {
        return nil, ErrOrderDoesNotExist
    }

    payment, err := c.repository.Save(newPaymentRequest)

    if err != nil {
        return nil, err
    }

    return payment, nil
}

func (c *PaymentsControllerImpl) GetById(id string) *models.Payment {
    return c.repository.GetById(id)
}

func (c *PaymentsControllerImpl) DeleteById(id string) (int64, error) {
    return c.repository.DeleteById(id)
}

func (c *PaymentsControllerImpl) Search(searchParams *models.PaymentsSearchParams) []*models.Payment {
    return c.repository.Search(searchParams)
}

func (c *PaymentsControllerImpl) Edit(paymentEditParams models.PaymentEditRequest) (*models.Payment, error) {
    return c.repository.Edit(paymentEditParams)
}
