package models

import (
    "fmt"
    "math/rand"
    "strings"
)

const INVITATION_ID_PREFIX = "HOSP"
const INVITATION_ID_LENGTH = 10
const INVITATION_ID_RND_CHARS = "123456789ABCDEFGHIJKLMNPQRSTUVWXYZ"

type InvitationId struct {
    code string
}

func (i *InvitationId) String() string {
    return fmt.Sprint(INVITATION_ID_PREFIX, i.code)
}

func allCharsValid(s, validSet string) bool {
    for _, char := range s {
        if !strings.ContainsRune(validSet, char) {
            return false
        }
    }
    return true
}

func NewRandomInvitationId() *InvitationId {
    code := RandGivenStringBytes(INVITATION_ID_LENGTH - len(INVITATION_ID_PREFIX), INVITATION_ID_RND_CHARS)
    return newInvitationIdFromCode(code)
}

func newInvitationIdFromCode(code string) *InvitationId {
    return &InvitationId{code}
}

func NewInvitationIdFromStr(str string) (*InvitationId, error) {
    if len(str) != INVITATION_ID_LENGTH || !strings.HasPrefix(str, "HOSP") {
        return nil, ErrInvitationNotValid
    }
    trimmed := str[len(INVITATION_ID_PREFIX):]
    if ! allCharsValid(trimmed, INVITATION_ID_RND_CHARS) {
        return nil, ErrInvitationNotValid
    }

    return newInvitationIdFromCode(trimmed), nil
}

func (i *InvitationId) UnmarshalJSON(data []byte) error {
    // Convert byte slice to string and trim double quotes
    str := strings.Trim(string(data), "\"")

    invitation, err := NewInvitationIdFromStr(str)
    if err != nil {
        return err
    }
    *i = *invitation
    return nil
}

func (i *InvitationId) MarshalJSON() ([]byte, error) {
    return []byte(`"` + i.String() + `"`), nil
}

func RandGivenStringBytes(n int, availableChars string) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = availableChars[rand.Intn(len(availableChars))]
    }
    return string(b)
}
