package services

import (
    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

type UsersCreditService interface {
    HandleCreditTravelerEvent(userId string)
    HandleCreditOwnerEvent(userId string)
}

type UsersCreditServiceImp struct {
    logger                    log.Logger
    usersCreditRepo           repositories.UsersCreditRepository
    invitationsRepo           repositories.InvitationsRepository
    usersRepo                 repositories.UsersRepository
    emailSender               EmailNotificationService
    CreditForInPlatformOrder  float32
    CreditForVerifiedProperty float32
}

func NewUsersCreditService(
    usersCreditRepo repositories.UsersCreditRepository,
    invitationsRepo repositories.InvitationsRepository,
    usersRepo repositories.UsersRepository,
    emailSender EmailNotificationService,
    CreditForInPlatformOrder float32,
    CreditForVerifiedProperty float32,
) UsersCreditService {
    logger := log.GetOrCreateLogger("UsersCreditService", "INFO")
    return &UsersCreditServiceImp{
        logger:                    logger,
        usersCreditRepo:           usersCreditRepo,
        invitationsRepo:           invitationsRepo,
        usersRepo:                 usersRepo,
        emailSender:               emailSender,
        CreditForInPlatformOrder:  CreditForInPlatformOrder,
        CreditForVerifiedProperty: CreditForVerifiedProperty,
    }
}

func (ucs *UsersCreditServiceImp) HandleCreditTravelerEvent(userId string) {
    ucs.handleCreditEvent(userId, models.ForTraveler, ucs.CreditForInPlatformOrder, ucs.emailSender.SendUserCreditNotificationSourceTraveler)
}

func (ucs *UsersCreditServiceImp) HandleCreditOwnerEvent(userId string) {
    ucs.handleCreditEvent(userId, models.ForOwner, ucs.CreditForVerifiedProperty, ucs.emailSender.SendUserCreditNotificationSourceOwner)
}

func (ucs *UsersCreditServiceImp) handleCreditEvent(userId string, invitaionKind models.InvitationKind, creditAmount float32, sender func(amount float32, recipientID string, userInvitedName string) error) {
    invitationsSearchParams := models.InvitationsSearchParams{UsedBy: &userId}
    invitation := ucs.invitationsRepo.Search(&invitationsSearchParams)
    if len(invitation) == 1 && invitation[0].Kind == invitaionKind {
        creditInstance := ucs.usersCreditRepo.GetByInvitationId(invitation[0].Id)
        if creditInstance == nil {
            newUserCreditInstanceRequest := models.NewUserCreditInstanceRequest{UserId: invitation[0].GeneratedBy,
                InvitationId:   invitation[0].Id,
                EarnedAmount:   creditAmount,
                EarnedCurrency: "USDT"}
            _, err := ucs.usersCreditRepo.Save(&newUserCreditInstanceRequest)
            if err != nil {
                ucs.logger.Errorf("Error sending credit to user %v: %v", invitation[0].GeneratedBy, err.Error())
            } else {
                userToNotify := invitation[0].GeneratedBy
                userInvited := ucs.usersRepo.GetById(*invitation[0].UsedBy)
                sender(newUserCreditInstanceRequest.EarnedAmount, userToNotify, userInvited.Name)
            }

        }
    }
}
