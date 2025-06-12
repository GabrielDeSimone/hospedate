package repositories

import (
    "database/sql"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    sqlTmp "github.com/hospedate/backyard/sql"
    "github.com/lib/pq"
)

type UsersCreditRepository interface {
    GetById(id string) *models.UserCreditInstance
    Save(ownersEarnedRequest *models.NewUserCreditInstanceRequest) (*models.UserCreditInstance, error)
    GetUserCreditInstances(id string) ([]*models.UserCreditInstance, error)
    GetByInvitationId(id models.InvitationId) *models.UserCreditInstance
}

type pgUsersCreditRepository struct {
    db     *sql.DB
    logger log.Logger
}

func NewUsersCreditRepository(db *sql.DB) UsersCreditRepository {
    logger := log.NewLogger("UsersCreditRepository", string(log.INFO_LEVEL))
    repo := &pgUsersCreditRepository{db: db, logger: logger}
    return repo
}

func (r *pgUsersCreditRepository) GetById(id string) *models.UserCreditInstance {
    result := r.db.QueryRow(sqlTmp.Users_credit_fetch_by_id, id)
    userCreditInstance, err := UserCreditInstanceFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Info("Error when getting a user credit instance by id: ", err.Error())
        return nil
    }

    return userCreditInstance
}

func (r *pgUsersCreditRepository) GetByInvitationId(id models.InvitationId) *models.UserCreditInstance {
    result := r.db.QueryRow(sqlTmp.Users_credit_fetch_by_invitationid, id)
    userCreditInstance, err := UserCreditInstanceFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Info("Error when getting a user credit instance by invitation_id: ", err.Error())
        return nil
    }

    return userCreditInstance
}

func (r *pgUsersCreditRepository) Save(userCreditRequest *models.NewUserCreditInstanceRequest) (*models.UserCreditInstance, error) {
    // Get the query template
    new_user_credit_sql := sqlTmp.Users_credit_new_instance

    // Generate an id
    id := GetRandomId()

    // Execute the query
    _, err := r.db.Exec(
        new_user_credit_sql,
        id,
        userCreditRequest.UserId,
        userCreditRequest.InvitationId.String(),
        userCreditRequest.EarnedAmount,
        userCreditRequest.EarnedCurrency,
    )

    if err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
            return nil, ErrDuplicateKey
        } else {
            r.logger.Info("Unknown error creating new user credit instance: ", err.Error())
            return nil, UnknownError
        }
    }

    r.logger.Info("New user credit instance created with id", id)

    return r.GetById(id), nil
}

func (r *pgUsersCreditRepository) GetUserCreditInstances(id string) ([]*models.UserCreditInstance, error) {
    rows, err := r.db.Query(sqlTmp.User_credit_search_query, id)
    if err == sql.ErrNoRows {
        return []*models.UserCreditInstance{}, nil
    } else if err != nil {
        r.logger.Error("Error when getting all user credit instances: ", err.Error())
        return nil, err
    }
    return r.fillUserCreditSlice(rows), nil

}

func (r *pgUsersCreditRepository) fillUserCreditSlice(rows *sql.Rows) []*models.UserCreditInstance {
    earnings := []*models.UserCreditInstance{}

    for rows.Next() {
        earning, err := UserCreditInstanceFromResult(rows)

        if err != nil {
            r.logger.Error("UserCreditInstance Scan failed: ", err.Error())
            return nil
        }
        earnings = append(earnings, earning)
    }
    return earnings
}

func UserCreditInstanceFromResult(result RowScanner) (*models.UserCreditInstance, error) {
    var userCreditInstance models.UserCreditInstance
    var invitationIdStr string
    err := result.Scan(
        &userCreditInstance.Id,
        &userCreditInstance.UserId,
        &invitationIdStr,
        &userCreditInstance.EarnedAmount,
        &userCreditInstance.EarnedCurrency,
        &userCreditInstance.CreatedAt)
    if err != nil {
        return nil, err
    }

    invitationId, err := models.NewInvitationIdFromStr(invitationIdStr)
    if err != nil {
        return nil, err
    }
    userCreditInstance.InvitationId = *invitationId

    return &userCreditInstance, nil

}
