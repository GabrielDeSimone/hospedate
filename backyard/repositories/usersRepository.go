package repositories

import (
    "database/sql"
    "time"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    sqlTmp "github.com/hospedate/backyard/sql"
    "github.com/lib/pq"
)

type UsersRepository interface {
    GetByEmailAndPassword(email string, password string) *models.User
    GetById(id string) *models.User
    Save(*models.NewUserRequest) (*models.User, error)
    Edit(userEditRequest models.UserEditRequest) (*models.User, error)
    GetBalanceById(id string) (*models.Balance, error)
    SaveWithdrawal(newWithdrawRequest *models.NewWithdrawalRequest) (*models.Withdrawal, error)
    GetWithdrawalById(id string) *models.Withdrawal
    EditWithdrawal(withdrawalEditRequest models.WithdrawalEditRequest) (*models.Withdrawal, error)
    GetWithdrawals(id string) ([]*models.Withdrawal, error)
}

type pgUserRepository struct {
    db     *sql.DB
    logger log.Logger
}

func NewUsersRepository(db *sql.DB) UsersRepository {
    logger := log.NewLogger("UsersRepository", string(log.INFO_LEVEL))
    repo := &pgUserRepository{db: db, logger: logger}

    return repo
}

func newEmptyBalance(userId string) *models.Balance {
    return &models.Balance{
        UserId:      userId,
        AmountCents: 0,
        CreatedAt:   time.Now().UTC(),
    }
}

func (r *pgUserRepository) GetByEmailAndPassword(email string, password string) *models.User {
    result := r.db.QueryRow(sqlTmp.Users_fetch_user, email, password)
    user, err := UserFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Error("Error when getting an user by email and password: ", err.Error())
        return nil
    }

    return user
}

func (r *pgUserRepository) GetById(id string) *models.User {
    result := r.db.QueryRow(sqlTmp.Users_fetch_by_id, id)
    user, err := UserFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Info("Error when getting a user by id: ", err.Error())
        return nil
    }

    return user
}

func (r *pgUserRepository) GetBalanceById(id string) (*models.Balance, error) {
    result := r.db.QueryRow(sqlTmp.Users_get_balance_by_id, id)
    balance, err := BalanceFromResult(result)

    if err == sql.ErrNoRows {
        return newEmptyBalance(id), nil
    } else if err != nil {
        r.logger.Info("Error when getting user balance: ", err.Error())
        return nil, err
    }

    return balance, nil
}

func (r *pgUserRepository) Save(newUserRequest *models.NewUserRequest) (*models.User, error) {
    // Get the query template
    new_user_sql := sqlTmp.Users_new_user

    // Generate an id
    id := GetRandomId()

    // Execute the query
    _, err := r.db.Exec(new_user_sql, id, newUserRequest.Name, newUserRequest.Email, newUserRequest.Password, newUserRequest.PhoneNumber)

    if err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
            return nil, ErrDuplicateKey
        } else {
            r.logger.Info("Unknown error creating new user: ", err.Error())
            return nil, UnknownError
        }
    }

    r.logger.Info("New user created with id", id)

    return r.GetById(id), nil
}

func (r *pgUserRepository) Edit(userEditRequest models.UserEditRequest) (*models.User, error) {
    // Build the base query template
    query, queryArgs := buildEditQuery(
        sqlTmp.Users_edit_base_query,
        GetParamFieldsPresent(userEditRequest),
        "WHERE id = ?",
        []any{userEditRequest.Id},
    )
    stmt, err := r.db.Prepare(query)

    if err != nil {
        r.logger.Error("Error when creating query template", err.Error())
        return nil, err
    }
    _, err = stmt.Exec(queryArgs...)
    if err != nil {
        r.logger.Error("Error when trying to update user: ", err.Error())
        return nil, err
    } else {
        return r.GetById(userEditRequest.Id), nil
    }
}

func (r *pgUserRepository) SaveWithdrawal(newWithdrawalRequest *models.NewWithdrawalRequest) (*models.Withdrawal, error) {
    // Get the query template
    new_withdraw_sql := sqlTmp.Users_new_withdrawal

    // Generate an id
    id := GetRandomId()

    // Execute the query
    _, err := r.db.Exec(new_withdraw_sql, id, newWithdrawalRequest.UserId,
        newWithdrawalRequest.ReclaimedAmountCents, newWithdrawalRequest.ReclaimedCurrency,
        newWithdrawalRequest.WalletAddress)

    if err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
            return nil, ErrDuplicateKey
        } else {
            r.logger.Error("Unknown error creating new user withdrawal: ", err.Error())
            return nil, UnknownError
        }
    }

    r.logger.Infof("New user withdrawal created with id %v for user id %v", id, newWithdrawalRequest.UserId)

    return r.GetWithdrawalById(id), nil
}

func (r *pgUserRepository) GetWithdrawalById(id string) *models.Withdrawal {
    result := r.db.QueryRow(sqlTmp.Users_get_withdrawal, id)
    withdrawal, err := WithdrawalFromResult(result)
    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Info("Error when getting a user by id: ", err.Error())
        return nil
    }
    return withdrawal
}

func (r *pgUserRepository) GetWithdrawals(id string) ([]*models.Withdrawal, error) {
    // Build the base query template
    rows, err := r.db.Query(sqlTmp.Withdrawals_search_query, id)
    if err == sql.ErrNoRows {
        return []*models.Withdrawal{}, nil
    } else if err != nil {
        r.logger.Error("Error when getting all user withdrawals: ", err.Error())
        return nil, err
    }
    return r.fillWithdrawalSlice(rows), nil
}

func (r *pgUserRepository) EditWithdrawal(withdrawalEditParams models.WithdrawalEditRequest) (*models.Withdrawal, error) {
    // Build the base query template
    query, queryArgs := buildEditQuery(
        sqlTmp.Withdrawals_edit_base_query,
        GetParamFieldsPresent(withdrawalEditParams),
        "WHERE id = ?",
        []any{withdrawalEditParams.Id},
    )
    stmt, err := r.db.Prepare(query)

    if err != nil {
        r.logger.Error("Error when creating query template", err.Error())
        return nil, err
    }
    _, err = stmt.Exec(queryArgs...)
    if err != nil {
        r.logger.Error("Error when trying to update a withdrawal: ", err.Error())
        return nil, err
    } else {
        return r.GetWithdrawalById(withdrawalEditParams.Id), nil
    }
}

func UserFromResult(result RowScanner) (*models.User, error) {
    var user models.User
    err := result.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.PhoneNumber, &user.IsHost)
    if err != nil {
        return nil, err
    } else {
        return &user, nil
    }
}

func BalanceFromResult(result RowScanner) (*models.Balance, error) {
    var balance models.Balance
    err := result.Scan(&balance.UserId, &balance.AmountCents)
    balance.CreatedAt = time.Now().UTC()
    if err != nil {
        return nil, err
    } else {
        return &balance, nil
    }
}

func WithdrawalFromResult(result RowScanner) (*models.Withdrawal, error) {
    var withdrawal models.Withdrawal
    err := result.Scan(
        &withdrawal.Id,
        &withdrawal.UserId,
        &withdrawal.Method,
        &withdrawal.Status,
        &withdrawal.ReclaimedAmountCents,
        &withdrawal.ReclaimedCurrency,
        &withdrawal.CreatedAt,
        &withdrawal.ProcessedAt,
        &withdrawal.WalletAddress)
    if err != nil {
        return nil, err
    } else {
        return &withdrawal, nil
    }
}

func (r *pgUserRepository) fillWithdrawalSlice(rows *sql.Rows) []*models.Withdrawal {
    withdrawals := []*models.Withdrawal{}

    for rows.Next() {
        withdrawal, err := WithdrawalFromResult(rows)

        if err != nil {
            r.logger.Error("Withdrawal Scan failed: ", err.Error())
            return nil
        }
        withdrawals = append(withdrawals, withdrawal)
    }
    return withdrawals
}
