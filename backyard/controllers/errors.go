package controllers

type ContrError string

const ErrDuplicateKey ContrError = "ErrDuplicateKey"

const UnknownError ContrError = "UnknownError"

const ErrCollision ContrError = "ErrCollision"

const ErrUserDoesNotExist ContrError = "ErrUserDoesNotExist"

const ErrPropertyDoesNotExist ContrError = "ErrPropertyDoesNotExist"

const ErrOrderDoesNotExist ContrError = "ErrOrderDoesNotExist"

const ErrOwnerBookingProperty ContrError = "ErrOwnerBookingProperty"

const ErrForTooManyGuests ContrError = "ErrForTooManyGuests"

const ErrBalanceNotEnough ContrError = "ErrBalanceNotEnough"

const ErrInvitationDoesNotExist ContrError = "ErrInvitationDoesNotExist"

const ErrInvitationAlreadyUsed ContrError = "ErrInvitationAlreadyUsed"

const ErrPropertyHasActiveOrders ContrError = "ErrPropertyHasActiveOrders"

const ErrPropertyArchivedStatus ContrError = "ErrPropertyArchivedStatus"

func (e ContrError) Error() string {
    return string(e)
}
