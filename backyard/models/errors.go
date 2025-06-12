package models

type ModelsErr string

const ErrCastPublicKey ModelsErr = "ErrCastPublicKey"
const ErrAttributesNotValid ModelsErr = "ErrAttributesNotValid"
const ErrAttributesNil ModelsErr = "ErrAttributesNil"
const ErrAttributesIncomplete ModelsErr = "ErrAttributesIncomplete"
const ErrAttributesInconsistent ModelsErr = "ErrAttributesInconsistent"
const ErrInvitationNotValid ModelsErr = "ErrInvitationNotValid"
const ErrInvitationKindNotValid ModelsErr = "ErrInvitationKindNotValid"

func (e ModelsErr) Error() string {
    return string(e)
}
