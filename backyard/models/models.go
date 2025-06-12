package models

import (
    "math"
    "time"
)

type User struct {
    Id          string    `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    CreatedAt   time.Time `json:"created_at"`
    PhoneNumber string    `json:"phone_number"`
    IsHost      bool      `json:"is_host"`
}

type NewUserRequest struct {
    Name         string        `json:"name" binding:"required"`
    Email        string        `json:"email" binding:"required,email" `
    Password     string        `json:"password" binding:"required,min=8"`
    PhoneNumber  string        `json:"phone_number" binding:"required"`
    InvitationId *InvitationId `json:"invitation_id"`
}

type UserEditRequest struct {
    Id     string `json:"id" db:"id"`
    IsHost *bool  `json:"is_host" db:"is_host"`
}

type Block struct {
    Id         string    `json:"id"`
    DateStart  Date      `json:"date_start"`
    DateEnd    Date      `json:"date_end"`
    PropertyId string    `json:"property_id"`
    CreatedAt  time.Time `json:"created_at"`
}

type NewBlockRequest struct {
    DateStart  Date   `json:"date_start" binding:"required"`
    DateEnd    Date   `json:"date_end" binding:"required"`
    PropertyId string `json:"property_id"`
}

type Order struct {
    Id               string    `json:"id"`
    UserId           string    `json:"user_id"`
    PropertyId       string    `json:"property_id"`
    Status           string    `json:"status"`
    DateStart        Date      `json:"date_start"`
    DateEnd          Date      `json:"date_end"`
    NumberGuests     uint      `json:"number_guests"`
    Price            uint      `json:"price"`
    PriceCurrency    string    `json:"price_currency"`
    TotalBilledCents uint      `json:"total_billed_cents"`
    CanceledBy       *string   `json:"canceled_by"`
    OrderType        string    `json:"order_type"`
    CreatedAt        time.Time `json:"created_at"`
    WalletAddress    string    `json:"wallet_address"`
}

type NewOrderRequest struct {
    UserId       string `json:"user_id" binding:"required" `
    PropertyId   string `json:"property_id" binding:"required" `
    DateStart    Date   `json:"date_start" binding:"required"`
    DateEnd      Date   `json:"date_end" binding:"required"`
    NumberGuests uint   `json:"number_guests" binding:"required" `
    OrderType    string `json:"order_type" binding:"required" `
}

type OrderEditRequest struct {
    Id            string  `json:"id" db:"id"`
    CanceledBy    *string `json:"canceled_by" db:"canceled_by"`
    Status        *string `json:"status" db:"status"`
    WalletAddress *string `json:"wallet_address" db:"wallet_address"`
}

type Payment struct {
    Id                  string    `json:"id"`
    OrderId             string    `json:"order_id"`
    Method              string    `json:"method"`
    Status              string    `json:"status"`
    TravelerAmountCents uint      `json:"traveler_amount_cents"`
    TravelerCurrency    string    `json:"traveler_currency"`
    ReceivedAmountCents *uint     `json:"received_amount_cents"`
    ReceivedCurrency    *string   `json:"received_currency"`
    RevertedAmountCents *uint     `json:"reverted_amount_cents"`
    RevertedCurrency    *string   `json:"reverted_currency"`
    CreatedAt           time.Time `json:"created_at"`
}

type NewPaymentRequest struct {
    OrderId             string `json:"order_id" binding:"required" `
    Method              string `json:"method" binding:"required"`
    TravelerAmountCents uint   `json:"traveler_amount_cents" binding:"required" `
    TravelerCurrency    string `json:"traveler_currency" binding:"required"`
}

type AirbnbFetcherStatus string

type NewAirbnbFetcherStatus struct {
    Status AirbnbFetcherStatus `json:"status"`
}

type PaymentEditRequestByOrderId struct {
    Id                  string  `json:"id" db:"id"`
    ReceivedAmountCents *uint   `json:"received_amount_cents" db:"received_amount_cents"`
    ReceivedCurrency    *string `json:"received_currency" db:"received_currency"`
    RevertedAmountCents *uint   `json:"reverted_amount_cents" db:"reverted_amount_cents"`
    RevertedCurrency    *string `json:"reverted_currency" db:"reverted_currency"`
    Status              *string `json:"status" db:"status"`
    PrivateKey          *string `json:"private_key" db:"private_key"`
}

type ParamField struct {
    Value   interface{}
    DbField string
}

type Withdrawal struct {
    Id                   string     `json:"id"`
    UserId               string     `json:"user_id"`
    Method               string     `json:"method"`
    Status               string     `json:"status"`
    ReclaimedAmountCents uint       `json:"reclaimed_amount_cents"`
    ReclaimedCurrency    string     `json:"reclaimed_currency"`
    CreatedAt            time.Time  `json:"created_at"`
    ProcessedAt          *time.Time `json:"processed_at"`
    WalletAddress        string     `json:"wallet_address"`
}

type NewWithdrawalRequest struct {
    UserId               string `json:"user_id" `
    ReclaimedAmountCents uint   `json:"reclaimed_amount_cents" binding:"required" `
    ReclaimedCurrency    string `json:"reclaimed_currency" binding:"required"`
    WalletAddress        string `json:"wallet_address" binding:"required"`
}

type Balance struct {
    UserId      string    `json:"user_id"`
    AmountCents uint      `json:"amount_cents"`
    CreatedAt   time.Time `json:"created_at"`
}

type WithdrawalEditRequest struct {
    Id     string  `json:"id" db:"id"`
    Status *string `json:"status" db:"status"`
}

type OwnerEarnedInstance struct {
    Id                string    `json:"id"`
    UserId            string    `json:"user_id"`
    OrderId           string    `json:"order_id"`
    EarnedAmountCents uint      `json:"earned_amount_cents"`
    EarnedCurrency    string    `json:"earned_currency"`
    CreatedAt         time.Time `json:"created_at"`
}

type NewOwnerEarnedInstanceRequest struct {
    UserId            string `json:"user_id" binding:"required" `
    OrderId           string `json:"order_id" binding:"required"`
    EarnedAmountCents uint   `json:"earned_amount_cents" binding:"required" `
    EarnedCurrency    string `json:"earned_currency" binding:"required"`
}

type Invitation struct {
    Id          InvitationId   `json:"id"`
    UsedBy      *string        `json:"used_by"`
    Kind        InvitationKind `json:"kind"`
    GeneratedBy string         `json:"generated_by"`
    CreatedAt   time.Time      `json:"created_at"`
}

type NewInvitationRequest struct {
    Kind        InvitationKind `json:"kind" binding:"required"`
    GeneratedBy string         `json:"generated_by" binding:"required"`
}

type InvitationEditRequest struct {
    Id     InvitationId `json:"id" db:"id"`
    UsedBy *string      `json:"used_by" db:"used_by"`
}

type UserHostApplication struct {
    UserId string `json:"user_id" binding:"required"`
}

type ExternalInvitationRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Body  string `json:"body" binding:"required"`
}

type UserCreditInstance struct {
    Id             string       `json:"id"`
    UserId         string       `json:"user_id"`
    InvitationId   InvitationId `json:"invitation_id"`
    EarnedAmount   float32      `json:"earned_amount"`
    EarnedCurrency string       `json:"earned_currency"`
    CreatedAt      time.Time    `json:"created_at"`
}

type NewUserCreditInstanceRequest struct {
    UserId         string       `json:"user_id"`
    InvitationId   InvitationId `json:"invitation_id" binding:"required"`
    EarnedAmount   float32      `json:"earned_amount" binding:"required" `
    EarnedCurrency string       `json:"earned_currency" binding:"required"`
}

func ComputeSubtotalCents(DateStart Date, DateEnd Date, Price uint) uint {
    days := DateEnd.Sub(&DateStart)
    subtotal := days * float64(Price)
    return uint(math.Round(subtotal * 100))
}
