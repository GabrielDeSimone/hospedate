package models

type InvitationsSearchParams struct {
    UsedBy      *string
    generatedBy *string
}

func NewInvitationsSearchParams(queryParams map[string][]string) (*InvitationsSearchParams, error) {
    invitationsSearchParams := InvitationsSearchParams{}
    for k, v := range queryParams {
        if k == "generated_by" {
            invitationsSearchParams.generatedBy = &(v[0])
        } else if k == "used_by" {
            invitationsSearchParams.UsedBy = &(v[0])
        }
    }

    return &invitationsSearchParams, nil
}

func (sp *InvitationsSearchParams) HasGeneratedBy() bool {
    return sp.generatedBy != nil
}

func (sp *InvitationsSearchParams) GetGeneratedBy() string {
    return *sp.generatedBy
}

func (sp *InvitationsSearchParams) HasUsedBy() bool {
    return sp.UsedBy != nil
}

func (sp *InvitationsSearchParams) GetUsedBy() string {
    return *sp.UsedBy
}
