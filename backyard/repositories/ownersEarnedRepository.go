package repositories

import (
    "database/sql"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    sqlTmp "github.com/hospedate/backyard/sql"
    "github.com/lib/pq"
)

type OwnersEarnedRepository interface {
    GetById(id string) *models.OwnerEarnedInstance
    Save(ownersEarnedRequest *models.NewOwnerEarnedInstanceRequest) (*models.OwnerEarnedInstance, error)
    GetOwnerEarnedInstances(id string) ([]*models.OwnerEarnedInstance, error)
}

type pgOwnersEarnedRepository struct {
    db     *sql.DB
    logger log.Logger
}

func NewOwnersEarnedRepository(db *sql.DB) OwnersEarnedRepository {
    logger := log.NewLogger("OwnersEarnedRepository", string(log.INFO_LEVEL))
    repo := &pgOwnersEarnedRepository{db: db, logger: logger}
    return repo
}

func (r *pgOwnersEarnedRepository) GetById(id string) *models.OwnerEarnedInstance {
    result := r.db.QueryRow(sqlTmp.Owners_earned_fetch_by_id, id)
    ownerEarnedInstance, err := OwnerEarnedInstanceFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Info("Error when getting a owner earned instance by id: ", err.Error())
        return nil
    }

    return ownerEarnedInstance
}

func (r *pgOwnersEarnedRepository) Save(ownersEarnedRequest *models.NewOwnerEarnedInstanceRequest) (*models.OwnerEarnedInstance, error) {
    // Get the query template
    new_owner_eaned_sql := sqlTmp.Owners_earned_new_order

    // Generate an id
    id := GetRandomId()

    // Execute the query
    _, err := r.db.Exec(
        new_owner_eaned_sql,
        id,
        ownersEarnedRequest.OrderId,
        ownersEarnedRequest.UserId,
        ownersEarnedRequest.EarnedAmountCents,
        ownersEarnedRequest.EarnedCurrency,
    )

    if err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
            return nil, ErrDuplicateKey
        } else {
            r.logger.Info("Unknown error creating new owner earned instance: ", err.Error())
            return nil, UnknownError
        }
    }

    r.logger.Info("New owner earned instance created with id", id)

    return r.GetById(id), nil
}

func (r *pgOwnersEarnedRepository) GetOwnerEarnedInstances(id string) ([]*models.OwnerEarnedInstance, error) {
    rows, err := r.db.Query(sqlTmp.Owners_earned_search_query, id)
    if err == sql.ErrNoRows {
        return []*models.OwnerEarnedInstance{}, nil
    } else if err != nil {
        r.logger.Error("Error when getting all user owners earned instances: ", err.Error())
        return nil, err
    }
    return r.fillOwnersEarnedSlice(rows), nil

}

func (r *pgOwnersEarnedRepository) fillOwnersEarnedSlice(rows *sql.Rows) []*models.OwnerEarnedInstance {
    earnings := []*models.OwnerEarnedInstance{}

    for rows.Next() {
        earning, err := OwnerEarnedInstanceFromResult(rows)

        if err != nil {
            r.logger.Error("OwnerEarnedInstance Scan failed: ", err.Error())
            return nil
        }
        earnings = append(earnings, earning)
    }
    return earnings
}

func OwnerEarnedInstanceFromResult(result RowScanner) (*models.OwnerEarnedInstance, error) {
    var ownerEarnedInstance models.OwnerEarnedInstance
    err := result.Scan(
        &ownerEarnedInstance.Id,
        &ownerEarnedInstance.OrderId,
        &ownerEarnedInstance.UserId,
        &ownerEarnedInstance.EarnedAmountCents,
        &ownerEarnedInstance.EarnedCurrency,
        &ownerEarnedInstance.CreatedAt)
    if err != nil {
        return nil, err
    } else {
        return &ownerEarnedInstance, nil
    }
}
