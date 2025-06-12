package models

type PaymentsSearchParams struct {
    userId      *string
    ownerId     *string
    orderId     *string
}


func NewPaymentsSearchParams(queryParams map[string][]string) (*PaymentsSearchParams, error) {
    paymentsSearchParams := PaymentsSearchParams{}
    for k, v := range queryParams {
        if k == "user_id" {
            paymentsSearchParams.userId = &(v[0])
        } else if k == "owner_id" {
            paymentsSearchParams.ownerId = &(v[0])
        } else if k == "order_id" {
            paymentsSearchParams.orderId = &(v[0])
        }
    }

    return &paymentsSearchParams, nil
}

func (sp *PaymentsSearchParams) HasUserId() bool {
    return sp.userId != nil
}

func (sp *PaymentsSearchParams) HasOwnerId() bool {
    return sp.ownerId != nil
}

func (sp *PaymentsSearchParams) HasOrderId() bool {
    return sp.orderId != nil
}

func (sp *PaymentsSearchParams) GetUserId() string {
    return *sp.userId
}

func (sp *PaymentsSearchParams) GetOwnerId() string {
    return *sp.ownerId
}

func (sp *PaymentsSearchParams) GetOrderId() string {
    return *sp.orderId
}
