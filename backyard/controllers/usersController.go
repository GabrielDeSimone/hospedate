package controllers

import (
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

type UsersController interface {
    Get(email string, password string) *models.User
    GetById(id string) *models.User
    Create(newUserRequest *models.NewUserRequest, invitationsController InvitationsController) (*models.User, error)
    Edit(userEditParams models.UserEditRequest) (*models.User, error)
    GetBalance(id string) (*models.Balance, error)
    CreateWithdrawal(newWithdrawalRequest *models.NewWithdrawalRequest) (*models.Withdrawal, error)
    GetWithdrawalById(id string) *models.Withdrawal
    EditWithdrawal(withdrawalEditRequest models.WithdrawalEditRequest) (*models.Withdrawal, error)
    GetWithdrawals(id string) ([]*models.Withdrawal, error)
    GetOwnerEarnedInstances(id string) ([]*models.OwnerEarnedInstance, error)
    GetUserCreditInstances(id string) ([]*models.UserCreditInstance, error)
    SendCreditToUser(userCreditRequest *models.NewUserCreditInstanceRequest) (*models.UserCreditInstance, error)
}

type UsersControllerImpl struct {
    repository       repositories.UsersRepository
    ownersEarnedRepo repositories.OwnersEarnedRepository
    userCreditRepo   repositories.UsersCreditRepository
}

func NewUsersController(repository repositories.UsersRepository,
    ownersEarnedRepo repositories.OwnersEarnedRepository,
    userCreditRepo repositories.UsersCreditRepository) UsersController {
    return &UsersControllerImpl{repository: repository,
        ownersEarnedRepo: ownersEarnedRepo,
        userCreditRepo:   userCreditRepo}
}

func (c *UsersControllerImpl) GetById(id string) *models.User {
    return c.repository.GetById(id)
}

func (c *UsersControllerImpl) Get(email string, password string) *models.User {
    return c.repository.GetByEmailAndPassword(email, password)
}

func (c *UsersControllerImpl) Create(newUserRequest *models.NewUserRequest, invitationsController InvitationsController) (*models.User, error) {
    if newUserRequest.InvitationId != nil {
        isValid, err := invitationsController.IsAvailableInvitationId(newUserRequest.InvitationId)
        if !isValid {
            return nil, err
        }
    }
    user, err := c.repository.Save(newUserRequest)
    if err == repositories.ErrDuplicateKey {
        return nil, ErrDuplicateKey
    } else if err != nil {
        return nil, UnknownError
    }
    if newUserRequest.InvitationId != nil {
        invitation, err := invitationsController.Edit(models.InvitationEditRequest{Id: *newUserRequest.InvitationId, UsedBy: &user.Id})
        if err != nil {
            return nil, UnknownError
        }

        // If user was invited by an Owners invitation, make them host
        if invitation.Kind == models.ForOwner {
            trueV := true
            request := models.UserEditRequest{
                Id:     user.Id,
                IsHost: &(trueV),
            }
            user, err = c.repository.Edit(request)
            if err != nil {
                return nil, err
            }
        }
    }
    return user, nil
}

func (c *UsersControllerImpl) Edit(userEditParams models.UserEditRequest) (*models.User, error) {
    return c.repository.Edit(userEditParams)
}

func (c *UsersControllerImpl) GetBalance(id string) (*models.Balance, error) {
    user := c.GetById(id)
    if user == nil {
        return nil, ErrUserDoesNotExist
    }
    return c.repository.GetBalanceById(id)
}

func (c *UsersControllerImpl) GetWithdrawals(id string) ([]*models.Withdrawal, error) {
    user := c.GetById(id)
    if user == nil {
        return nil, ErrUserDoesNotExist
    }
    return c.repository.GetWithdrawals(id)
}

func (c *UsersControllerImpl) CreateWithdrawal(newWithdrawalRequest *models.NewWithdrawalRequest) (*models.Withdrawal, error) {
    balance, err := c.GetBalance(newWithdrawalRequest.UserId)
    if err != nil {
        return nil, err
    }

    if balance.AmountCents < newWithdrawalRequest.ReclaimedAmountCents {
        return nil, ErrBalanceNotEnough
    }
    withdraw, err := c.repository.SaveWithdrawal(newWithdrawalRequest)

    if err != nil {
        return nil, UnknownError
    }

    return withdraw, nil
}

func (c *UsersControllerImpl) GetWithdrawalById(id string) *models.Withdrawal {
    return c.repository.GetWithdrawalById(id)
}

func (c *UsersControllerImpl) EditWithdrawal(withdrawalEditRequest models.WithdrawalEditRequest) (*models.Withdrawal, error) {
    return c.repository.EditWithdrawal(withdrawalEditRequest)
}

func (c *UsersControllerImpl) GetOwnerEarnedInstances(id string) ([]*models.OwnerEarnedInstance, error) {
    user := c.GetById(id)
    if user == nil {
        return nil, ErrUserDoesNotExist
    }
    return c.ownersEarnedRepo.GetOwnerEarnedInstances(id)
}

func (c *UsersControllerImpl) GetUserCreditInstances(id string) ([]*models.UserCreditInstance, error) {
    user := c.GetById(id)
    if user == nil {
        return nil, ErrUserDoesNotExist
    }
    return c.userCreditRepo.GetUserCreditInstances(id)
}

func (c *UsersControllerImpl) SendCreditToUser(userCreditRequest *models.NewUserCreditInstanceRequest) (*models.UserCreditInstance, error) {
    user := c.GetById(userCreditRequest.UserId)
    if user == nil {
        return nil, ErrUserDoesNotExist
    }
    return c.userCreditRepo.Save(userCreditRequest)
}
