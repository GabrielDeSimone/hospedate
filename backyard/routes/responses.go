package routes

import "fmt"

type ResponseOk struct {
    Data interface{} `json:"data"`
}

type ResponseErr struct {
    Error ErrorCode `json:"error"`
}

type NumberRowsDeleted int64

type DeletedRows struct {
    Nrows NumberRowsDeleted `json:"deleted_rows"`
}

type ErrorCode string

const ErrNotFound ErrorCode = "NotFound"
const ErrInternal ErrorCode = "InternalServerError"
const ErrEmailOrPhoneAlreadyTaken ErrorCode = "ErrEmailOrPhoneAlreadyTaken"
const ErrPropertyAlreadyTaken ErrorCode = "ErrPropertyAlreadyTaken"
const ErrBadRequest ErrorCode = "ErrBadRequest"
const ErrCollision ErrorCode = "ErrCollision"
const ErrNoValuesToUpdate ErrorCode = "ErrNoValuesToUpdate"
const ErrOwnerBookingProperty ErrorCode = "ErrOwnerBookingProperty"
const ErrGuestsNumberExceeded ErrorCode = "ErrGuestsNumberExceeded"
const ErrBalanceNotEnough ErrorCode = "ErrBalanceNotEnough"
const ErrInvitationNotValid ErrorCode = "ErrInvitationNotValid"
const ErrInvitationAlreadyUsed ErrorCode = "ErrInvitationAlreadyUsed"
const ErrInvitationDoesNotExist ErrorCode = "ErrInvitationDoesNotExist"
const ErrUserNotFound ErrorCode = "ErrUserNotFound"
const ErrPropertyArchivedStatus ErrorCode = "ErrPropertyArchivedStatus"
const ErrPropertyHasActiveOrders ErrorCode = "ErrPropertyHasActiveOrders"

func NewErrBadRequest(reason string) ErrorCode {
    return ErrorCode(fmt.Sprintf("%s: %s", ErrBadRequest, reason))
}
