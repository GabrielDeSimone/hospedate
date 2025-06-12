package models

type OrdersSearchParams struct {
    UserId  *string
    OwnerId *string
    Status  *string
}

func NewOrdersSearchParams(queryParams map[string][]string) (*OrdersSearchParams, error) {
    ordersSearchParams := OrdersSearchParams{}
    for k, v := range queryParams {
        if k == "user_id" {
            ordersSearchParams.UserId = &(v[0])
        } else if k == "owner_id" {
            ordersSearchParams.OwnerId = &(v[0])
        } else if k == "status" {
            ordersSearchParams.Status = &(v[0])
        }
    }

    return &ordersSearchParams, nil
}

func (sp *OrdersSearchParams) HasUserId() bool {
    return sp.UserId != nil
}

func (sp *OrdersSearchParams) HasOwnerId() bool {
    return sp.OwnerId != nil
}

func (sp *OrdersSearchParams) HasStatus() bool {
    return sp.Status != nil
}

func (sp *OrdersSearchParams) GetUserId() string {
    return *sp.UserId
}

func (sp *OrdersSearchParams) GetOwnerId() string {
    return *sp.OwnerId
}

func (sp *OrdersSearchParams) GetStatus() string {
    return *sp.Status
}
