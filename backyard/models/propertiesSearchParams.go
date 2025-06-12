package models

import (
    "errors"
    "strconv"
)

type PropertiesSearchParams struct {
    dateStart      *Date
    dateEnd        *Date
    city           *string
    guests         *uint
    userId         *string
    status         *string
    accommodation  *string
    location       *string
    wifi           *string
    tv             *string
    microwave      *string
    oven           *string
    kettle         *string
    toaster        *string
    coffeeMachine  *string
    ac             *string
    heating        *string
    parking        *string
    pool           *string
    gym            *string
    bookingOptions *string
    halfBathrooms  *uint
    bedrooms       *uint
}

func NewPropSearchParams(queryParams map[string][]string) (*PropertiesSearchParams, error) {
    propSearchParams := PropertiesSearchParams{}

    for k, v := range queryParams {
        if k == "city" {
            propSearchParams.city = &(v[0])
        } else if k == "date_start" {
            dateStart, err := NewDateFromStr(v[0])
            if err != nil {
                return nil, err
            }
            propSearchParams.dateStart = dateStart
        } else if k == "date_end" {
            dateEnd, err := NewDateFromStr(v[0])
            if err != nil {
                return nil, err
            }
            propSearchParams.dateEnd = dateEnd
        } else if k == "guests" {
            numGuests, err := strconv.Atoi(v[0])
            if err != nil {
                return nil, err
            } else if numGuests < 0 {
                return nil, errors.New("guests cannot be a negative number")
            }
            numGuestsUint := uint(numGuests)
            propSearchParams.guests = &numGuestsUint
        } else if k == "user_id" {
            propSearchParams.userId = &(v[0])
        } else if k == "status" {
            propSearchParams.status = &(v[0])
        } else if k == "accommodation" {
            propSearchParams.accommodation = &(v[0])
        } else if k == "location" {
            propSearchParams.location = &(v[0])
        } else if k == "wifi" {
            propSearchParams.wifi = &(v[0])
        } else if k == "tv" {
            propSearchParams.tv = &(v[0])
        } else if k == "microwave" {
            propSearchParams.microwave = &(v[0])
        } else if k == "oven" {
            propSearchParams.oven = &(v[0])
        } else if k == "kettle" {
            propSearchParams.kettle = &(v[0])
        } else if k == "toaster" {
            propSearchParams.toaster = &(v[0])
        } else if k == "coffee_machine" {
            propSearchParams.coffeeMachine = &(v[0])
        } else if k == "ac" {
            propSearchParams.ac = &(v[0])
        } else if k == "heating" {
            propSearchParams.heating = &(v[0])
        } else if k == "parking" {
            propSearchParams.parking = &(v[0])
        } else if k == "pool" {
            propSearchParams.pool = &(v[0])
        } else if k == "gym" {
            propSearchParams.gym = &(v[0])
        } else if k == "booking_options" {
            propSearchParams.bookingOptions = &(v[0])
        } else if k == "half_bathrooms" {
            halfBathrooms, err := strconv.Atoi(v[0])
            if err != nil {
                return nil, err
            }
            halfBathroomsUint := uint(halfBathrooms)
            propSearchParams.halfBathrooms = &halfBathroomsUint
        } else if k == "bedrooms" {
            bedrooms, err := strconv.Atoi(v[0])
            if err != nil {
                return nil, err
            }
            bedroomsUint := uint(bedrooms)
            propSearchParams.bedrooms = &bedroomsUint
        }
    }

    return &propSearchParams, nil
}

func (sp *PropertiesSearchParams) HasDates() bool {
    return sp.HasDateStart() && sp.HasDateEnd()
}

func (sp *PropertiesSearchParams) HasOnlyOneDate() bool {
    return (sp.HasDateStart() && !sp.HasDateEnd()) || (!sp.HasDateStart() && sp.HasDateEnd())
}

func (sp *PropertiesSearchParams) HasDateStart() bool {
    return sp.dateStart != nil
}

func (sp *PropertiesSearchParams) HasDateEnd() bool {
    return sp.dateEnd != nil
}

func (sp *PropertiesSearchParams) HasCity() bool {
    return sp.city != nil
}

func (sp *PropertiesSearchParams) HasUserId() bool {
    return sp.userId != nil
}

func (sp *PropertiesSearchParams) HasStatus() bool {
    return sp.status != nil
}

func (sp *PropertiesSearchParams) HasGuests() bool {
    return sp.guests != nil
}

func (sp *PropertiesSearchParams) GetDateStart() *Date {
    return sp.dateStart
}

func (sp *PropertiesSearchParams) GetDateEnd() *Date {
    return sp.dateEnd
}

func (sp *PropertiesSearchParams) GetGuests() uint {
    return *sp.guests
}

func (sp *PropertiesSearchParams) GetCity() string {
    return *sp.city
}

func (sp *PropertiesSearchParams) GetUserId() string {
    return *sp.userId
}

func (sp *PropertiesSearchParams) GetStatus() string {
    return *sp.status
}

func (sp *PropertiesSearchParams) HasAccommodation() bool {
    return sp.accommodation != nil
}

func (sp *PropertiesSearchParams) GetAccommodation() string {
    return *sp.accommodation
}

func (sp *PropertiesSearchParams) HasLocation() bool {
    return sp.location != nil
}

func (sp *PropertiesSearchParams) GetLocation() string {
    return *sp.location
}

func (sp *PropertiesSearchParams) HasWifi() bool {
    return sp.wifi != nil
}

func (sp *PropertiesSearchParams) GetWifi() string {
    return *sp.wifi
}

func (sp *PropertiesSearchParams) HasTV() bool {
    return sp.tv != nil
}

func (sp *PropertiesSearchParams) GetTV() string {
    return *sp.tv
}

func (sp *PropertiesSearchParams) HasMicrowave() bool {
    return sp.microwave != nil
}

func (sp *PropertiesSearchParams) GetMicrowave() string {
    return *sp.microwave
}

func (sp *PropertiesSearchParams) HasOven() bool {
    return sp.oven != nil
}

func (sp *PropertiesSearchParams) GetOven() string {
    return *sp.oven
}

func (sp *PropertiesSearchParams) HasKettle() bool {
    return sp.kettle != nil
}

func (sp *PropertiesSearchParams) GetKettle() string {
    return *sp.kettle
}

func (sp *PropertiesSearchParams) HasToaster() bool {
    return sp.toaster != nil
}

func (sp *PropertiesSearchParams) GetToaster() string {
    return *sp.toaster
}

func (sp *PropertiesSearchParams) HasCoffeeMachine() bool {
    return sp.coffeeMachine != nil
}

func (sp *PropertiesSearchParams) GetCoffeeMachine() string {
    return *sp.coffeeMachine
}

func (sp *PropertiesSearchParams) HasAC() bool {
    return sp.ac != nil
}

func (sp *PropertiesSearchParams) GetAC() string {
    return *sp.ac
}

func (sp *PropertiesSearchParams) HasHeating() bool {
    return sp.heating != nil
}

func (sp *PropertiesSearchParams) GetHeating() string {
    return *sp.heating
}

func (sp *PropertiesSearchParams) HasParking() bool {
    return sp.parking != nil
}

func (sp *PropertiesSearchParams) GetParking() string {
    return *sp.parking
}

func (sp *PropertiesSearchParams) HasPool() bool {
    return sp.pool != nil
}

func (sp *PropertiesSearchParams) GetPool() string {
    return *sp.pool
}

func (sp *PropertiesSearchParams) HasGym() bool {
    return sp.gym != nil
}

func (sp *PropertiesSearchParams) GetGym() string {
    return *sp.gym
}

func (sp *PropertiesSearchParams) HasBookingOptions() bool {
    return sp.bookingOptions != nil
}

func (sp *PropertiesSearchParams) GetBookingOptions() string {
    return *sp.bookingOptions
}

func (sp *PropertiesSearchParams) HasBathrooms() bool {
    return sp.halfBathrooms != nil
}

func (sp *PropertiesSearchParams) GetBathrooms() uint {
    return *sp.halfBathrooms
}

func (sp *PropertiesSearchParams) HasBedrooms() bool {
    return sp.bedrooms != nil
}

func (sp *PropertiesSearchParams) GetBedrooms() uint {
    return *sp.bedrooms
}
