package models

import (
    "time"
)

type NewPropertyRequest struct {
    MaxGuests    uint   `json:"max_guests" binding:"required"`
    AirbnbRoomId string `json:"airbnb_room_id" binding:"required"`
    Price        uint   `json:"price" binding:"required"`
    User_id      string `json:"user_id" binding:"required"`
    City         string `json:"city" binding:"required"`
}

type PropertyEditRequest struct {
    Id             string   `db:"id" json:"id"`
    Title          *string  `db:"title" json:"title"`
    Description    *string  `db:"description" json:"description"`
    Price          *uint    `db:"price" json:"price"`
    Images         []string `db:"images" json:"images"`
    Status         *string  `db:"status" json:"status"`
    IsVerified     *bool    `db:"is_verified" json:"is_verified"`
    BookingOptions *string  `db:"booking_options" json:"booking_options"`
    MaxGuests      *uint    `db:"max_guests" json:"max_guests"`
    ExAirbnbRoomId *string  `db:"ex_airbnb_room_id"`
    Accommodation  *string  `db:"accommodation" json:"accommodation"`
    Location       *string  `db:"location" json:"location"`
    Wifi           *string  `db:"wifi" json:"wifi"`
    TV             *string  `db:"tv" json:"tv"`
    Microwave      *string  `db:"microwave" json:"microwave"`
    Oven           *string  `db:"oven" json:"oven"`
    Kettle         *string  `db:"kettle" json:"kettle"`
    Toaster        *string  `db:"toaster" json:"toaster"`
    CoffeeMachine  *string  `db:"coffee_machine" json:"coffee_machine"`
    AC             *string  `db:"air_conditioning" json:"air_conditioning"`
    Heating        *string  `db:"heating" json:"heating"`
    Parking        *string  `db:"parking" json:"parking"`
    Pool           *string  `db:"pool" json:"pool"`
    Gym            *string  `db:"gym" json:"gym"`
    HalfBathrooms  *uint    `db:"half_bathrooms" json:"half_bathrooms"`
    Bedrooms       *uint    `db:"bedrooms" json:"bedrooms"`
}

type Property struct {
    Id             string    `json:"id"`
    Title          string    `json:"title"`
    Description    string    `json:"description"`
    MaxGuests      uint      `json:"max_guests"`
    AirbnbRoomId   *string   `json:"airbnb_room_id"`
    ExAirbnbRoomId *string   `json:"ex_airbnb_room_id"`
    Price          uint      `json:"price"`
    UserId         string    `json:"user_id"`
    City           string    `json:"city"`
    Status         string    `json:"status"`
    IsVerified     bool      `json:"is_verified"`
    Images         []string  `json:"images"`
    BookingOptions *string   `json:"booking_options"`
    CreatedAt      time.Time `json:"created_at"`
    Accommodation  *string   `json:"accommodation"`
    Location       *string   `json:"location"`
    Wifi           *string   `json:"wifi"`
    TV             *string   `json:"tv"`
    Microwave      *string   `json:"microwave"`
    Oven           *string   `json:"oven"`
    Kettle         *string   `json:"kettle"`
    Toaster        *string   `json:"toaster"`
    CoffeeMachine  *string   `json:"coffee_machine"`
    AC             *string   `json:"air_conditioning"`
    Heating        *string   `json:"heating"`
    Parking        *string   `json:"parking"`
    Pool           *string   `json:"pool"`
    Gym            *string   `json:"gym"`
    HalfBathrooms  *uint     `json:"half_bathrooms"`
    Bedrooms       *uint     `json:"bedrooms"`
}
