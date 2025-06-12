package models

import (
    "encoding/json"
)

const (
    ForTraveler InvitationKind = "for_traveler"
    ForOwner    InvitationKind = "for_owner"
)

type InvitationKind string

func (i *InvitationKind) UnmarshalJSON(data []byte) error {
    var role string
    if err := json.Unmarshal(data, &role); err != nil {
        return err
    }
    invKind := InvitationKind(role)
    if invKind != ForTraveler && invKind != ForOwner {
        return ErrInvitationKindNotValid
    }

    *i = invKind
    return nil
}

func (i *InvitationKind) MarshalJSON() ([]byte, error) {
    return json.Marshal(string(*i))
}

func (i InvitationKind) String() string {
    return string(i)
}

func NewInvitationKind(role string) (*InvitationKind, error) {
    invKind := InvitationKind(role)
    if invKind != ForTraveler && invKind != ForOwner {
        return nil, ErrInvitationKindNotValid
    }
    return &invKind, nil
}
