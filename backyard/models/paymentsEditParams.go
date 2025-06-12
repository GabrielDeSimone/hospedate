package models

import (
    "github.com/gin-gonic/gin"
)

type PaymentEditRequest struct {
    Id                  string  `json:"id" db:"id"`
    ReceivedAmountCents *uint   `json:"received_amount_cents" db:"received_amount_cents"`
    ReceivedCurrency    *string `json:"received_currency" db:"received_currency"`
    RevertedAmountCents *uint   `json:"reverted_amount_cents" db:"reverted_amount_cents"`
    RevertedCurrency    *string `json:"reverted_currency" db:"reverted_currency"`
    Status              *string `json:"status" db:"status"`
}

func NewPaymentEditRequest(ctx *gin.Context) (*PaymentEditRequest, error) {
    var paymentEditRequest PaymentEditRequest

    err := ctx.BindJSON(&paymentEditRequest)
    if err != nil {
        return nil, ErrAttributesNotValid
    }
    isValid, err := isValidPaymentEditRequest(paymentEditRequest)
    if !isValid {
        return nil, err
    }
    return &paymentEditRequest, nil

}

func isValidPaymentEditRequest(request PaymentEditRequest) (bool, error) {
    // At least one field must not be null
    fields := []interface{}{
        request.ReceivedAmountCents,
        request.ReceivedCurrency,
        request.RevertedAmountCents,
        request.RevertedCurrency,
        request.Status,
    }
    validFields := 0
    for _, field := range fields {
        if field != nil {
            validFields++
        }
    }
    if validFields == 0 {
        return false, ErrAttributesNil
    }

    // If Received_currency is provided, Received_amount must be provided and vice versa
    if (request.ReceivedAmountCents == nil) != (request.ReceivedCurrency == nil) {
        return false, ErrAttributesIncomplete
    }

    // If Reverted_amount is provided, Reverted_currency must be provided and vice versa
    if (request.RevertedAmountCents == nil) != (request.RevertedCurrency == nil) {
        return false, ErrAttributesIncomplete
    }

    //  If Reverted_amount and Received_amount are provided, Reverted_amount cannot be greater than Received_amount
    if request.RevertedAmountCents != nil && request.ReceivedAmountCents != nil {
        revertedAmount := *request.RevertedAmountCents
        receivedAmount := *request.ReceivedAmountCents
        if revertedAmount > uint(receivedAmount) {
            return false, ErrAttributesInconsistent
        }

        // Received_currency must be equal to Reverted_currency
        if *request.ReceivedCurrency != *request.RevertedCurrency {
            return false, ErrAttributesInconsistent
        }
    }

    return true, nil
}
