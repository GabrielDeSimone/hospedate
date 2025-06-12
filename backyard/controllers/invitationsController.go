package controllers

import (
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

type InvitationsController interface {
    GetById(id *models.InvitationId) *models.Invitation
    Create(request *models.NewInvitationRequest) (*models.Invitation, error)
    IsAvailableInvitationId(id *models.InvitationId) (bool, error)
    Edit(invitationEditRequest models.InvitationEditRequest) (*models.Invitation, error)
    Search(searchParams *models.InvitationsSearchParams) []*models.Invitation
}

type InvitationsControllerImpl struct {
    repository repositories.InvitationsRepository
}

func NewInvitationsController(repository repositories.InvitationsRepository) InvitationsController {
    return &InvitationsControllerImpl{repository: repository}
}

func (c InvitationsControllerImpl) GetById(id *models.InvitationId) *models.Invitation {
    return c.repository.GetById(id)
}

func (c InvitationsControllerImpl) Create(request *models.NewInvitationRequest) (*models.Invitation, error) {
    invitation, err := c.repository.Save(request)

    if err != nil {
        return nil, UnknownError
    }

    return invitation, nil
}

func (c InvitationsControllerImpl) IsAvailableInvitationId(id *models.InvitationId) (bool, error) {
    invitation := c.repository.GetById(id)
    if invitation == nil {
        return false, ErrInvitationDoesNotExist
    }
    if invitation.UsedBy != nil {
        return false, ErrInvitationAlreadyUsed
    }
    return true, nil
}

func (c InvitationsControllerImpl) Edit(invitationEditRequest models.InvitationEditRequest) (*models.Invitation, error) {
    return c.repository.Edit(invitationEditRequest)
}

func (c InvitationsControllerImpl) Search(searchParams *models.InvitationsSearchParams) []*models.Invitation {
    return c.repository.Search(searchParams)
}
