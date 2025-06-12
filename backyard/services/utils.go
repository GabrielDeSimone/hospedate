package services

import (
    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

func updateOrderStatus(ordersRepo repositories.OrdersRepository, logger log.Logger, orderId string, status string) error {
    orderEditRequest := models.OrderEditRequest{
        Id:     orderId,
        Status: &status,
    }
    _, err := ordersRepo.Edit(orderEditRequest)
    if err != nil {
        logger.Error("Error updating order status:", err.Error())
    }
    return err
}

