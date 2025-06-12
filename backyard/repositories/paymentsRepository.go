package repositories

import (
    "database/sql"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    sqlTmp "github.com/hospedate/backyard/sql"
    "github.com/lib/pq"
)

type PaymentsRepository interface {
    Save(newPaymentRequest *models.NewPaymentRequest) (*models.Payment, error)
    saveTx(newPaymentRequest *models.NewPaymentRequest, tx *sql.Tx) (*models.Payment, error)
    GetById(id string) *models.Payment
    getByIdTx(id string, tx *sql.Tx) *models.Payment
    DeleteById(id string) (int64, error)
    DeleteByOrderId(id string) (int64, error)
    Search(queryParams *models.PaymentsSearchParams) []*models.Payment
    Edit(paymentEditParams models.PaymentEditRequest) (*models.Payment, error)
    EditByOrderId(paymentEditParams models.PaymentEditRequestByOrderId) (*models.Payment, error)
    editByOrderIdTx(paymentEditParams models.PaymentEditRequestByOrderId, tx *sql.Tx) (*models.Payment, error)
    GetByOrderId(order_id string) *models.Payment
    getByOrderIdTx(order_id string, tx *sql.Tx) *models.Payment
    GetPkByOrderId(order_id string) (string, error)
}

type pgPaymentsRepository struct {
    db     *sql.DB
    logger log.Logger
}

func NewPaymentsRepository(db *sql.DB) PaymentsRepository {
    logger := log.NewLogger("PaymentsRepository", string(log.INFO_LEVEL))
    repo := &pgPaymentsRepository{db: db, logger: logger}

    return repo
}

func (r *pgPaymentsRepository) fillPaymentSlice(rows *sql.Rows) []*models.Payment {
    payments := []*models.Payment{}

    for rows.Next() {
        payment, err := PaymentFromResult(rows)

        if err != nil {
            r.logger.Error("Payment Scan failed: ", err.Error())
            return nil
        }
        payments = append(payments, payment)
    }
    return payments
}

func (r *pgPaymentsRepository) Save(newPaymentRequest *models.NewPaymentRequest) (*models.Payment, error) {
    return r.saveTx(newPaymentRequest, nil)
}

func (r *pgPaymentsRepository) saveTx(newPaymentRequest *models.NewPaymentRequest, tx *sql.Tx) (*models.Payment, error) {
    // Get DbOrTx object
    dbOrTx := GetDbOrTx(r.db, tx)

    // Get the query template
    new_payment_sql := sqlTmp.Payments_new_payment

    // Generate an id
    id := GetRandomId()

    // Execute the query
    _, err := dbOrTx.Exec(
        new_payment_sql,
        id,
        newPaymentRequest.OrderId,
        newPaymentRequest.TravelerAmountCents,
        newPaymentRequest.TravelerCurrency,
    )

    if err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
            return nil, ErrDuplicateKey // el pago de la orden ya existe
        } else {
            r.logger.Error("Unknown error creating new payment: ", err.Error())
            return nil, UnknownError
        }
    }

    r.logger.Info("New payment created with id", id)

    return r.getByIdTx(id, tx), nil
}

func (r *pgPaymentsRepository) GetById(id string) *models.Payment {
    return r.getByIdTx(id, nil)
}

func (r *pgPaymentsRepository) getByIdTx(id string, tx *sql.Tx) *models.Payment {
    dbOrTx := GetDbOrTx(r.db, tx)
    result := dbOrTx.QueryRow(sqlTmp.Payments_fetch_by_id, id)
    payment, err := PaymentFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Error("Error when fetching payment: ", err.Error())
        return nil
    }

    return payment
}

func (r *pgPaymentsRepository) GetByOrderId(orderId string) *models.Payment {
    return r.getByOrderIdTx(orderId, nil)
}

func (r *pgPaymentsRepository) getByOrderIdTx(orderId string, tx *sql.Tx) *models.Payment {
    dbOrTx := GetDbOrTx(r.db, tx)
    result := dbOrTx.QueryRow(sqlTmp.Payments_fetch_by_order_id, orderId)
    payment, err := PaymentFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Error("Error when fetching payment: ", err.Error())
        return nil
    }

    return payment
}

func (r *pgPaymentsRepository) GetPkByOrderId(orderId string) (string, error) {
    result := r.db.QueryRow(sqlTmp.Payments_fetch_PK_by_order_id, orderId)

    var PK string

    err := result.Scan(
        &PK,
    )

    if err == sql.ErrNoRows {
        return "", err
    } else if err != nil {
        r.logger.Error("Error when fetching payment PK from orderid: ", err.Error())
        return "", err
    }

    return PK, nil
}

func (r *pgPaymentsRepository) DeleteById(id string) (int64, error) {
    result, err := r.db.Exec(sqlTmp.Payments_delete_by_id, id)
    if err != nil {
        r.logger.Error("Error when trying to delete a payment: ", err.Error())
        return 0, err
    } else {
        return result.RowsAffected()
    }
}

func (r *pgPaymentsRepository) DeleteByOrderId(id string) (int64, error) {
    result, err := r.db.Exec(sqlTmp.Payments_delete_by_order_id, id)
    if err != nil {
        r.logger.Error("Error when trying to delete a payment: ", err.Error())
        return 0, err
    } else {
        return result.RowsAffected()
    }
}

func (r *pgPaymentsRepository) Search(queryParams *models.PaymentsSearchParams) []*models.Payment {
    // Build the base query template
    queryTmp := sqlTmp.Payments_search_base_query
    queryArgs := []any{}

    // Add filters if necessary
    if queryParams.HasUserId() {
        queryTmp += " AND order_id IN (SELECT id FROM orders WHERE user_id = ?)"
        queryArgs = append(queryArgs, queryParams.GetUserId())
    }
    if queryParams.HasOwnerId() {
        queryTmp += " AND order_id IN (SELECT id FROM orders WHERE (property_id IN (SELECT id FROM properties WHERE user_id = ?)))"
        queryArgs = append(queryArgs, queryParams.GetOwnerId())
    }
    if queryParams.HasOrderId() {
        queryTmp += " AND order_id = ? "
        queryArgs = append(queryArgs, queryParams.GetOrderId())
    }

    // change template symbols from ? to $1, etc
    queryTmp = ReplaceArgsTemplate(queryTmp)

    // Exec the query and get the results
    stmt, err := r.db.Prepare(queryTmp)

    if err != nil {
        r.logger.Error("Error when creating query template", err.Error())
        return nil
    }

    rows, err := stmt.Query(queryArgs...)
    if err == sql.ErrNoRows {
        return []*models.Payment{}
    } else if err != nil {
        r.logger.Error("Error when searching payments: ", err.Error())
        return nil
    }

    defer rows.Close()

    return r.fillPaymentSlice(rows)
}

func (r *pgPaymentsRepository) Edit(paymentEditParams models.PaymentEditRequest) (*models.Payment, error) {
    // Build the base query template
    query, queryArgs := buildEditQuery(
        sqlTmp.Payments_edit_base_query,
        GetParamFieldsPresent(paymentEditParams),
        "WHERE id = ?",
        []any{paymentEditParams.Id},
    )
    stmt, err := r.db.Prepare(query)

    if err != nil {
        r.logger.Error("Error when creating query template", err.Error())
        return nil, err
    }
    _, err = stmt.Exec(queryArgs...)
    if err != nil {
        r.logger.Error("Error when trying to update a payment: ", err.Error())
        return nil, err
    } else {
        return r.getByIdTx(paymentEditParams.Id, nil), nil
    }
}

func (r *pgPaymentsRepository) EditByOrderId(paymentEditParams models.PaymentEditRequestByOrderId) (*models.Payment, error) {
    return r.editByOrderIdTx(paymentEditParams, nil)
}

func (r *pgPaymentsRepository) editByOrderIdTx(paymentEditParams models.PaymentEditRequestByOrderId, tx *sql.Tx) (*models.Payment, error) {
    // Get payment by order id
    payment := r.getByOrderIdTx(paymentEditParams.Id, tx)
    if payment == nil {
        r.logger.Error("Error when getting payment by order id ")
        return nil, sql.ErrNoRows
    }

    // Get DbOrTx object and build query
    dbOrTx := GetDbOrTx(r.db, tx)
    query, queryArgs := buildEditQuery(
        sqlTmp.Payments_edit_base_query,
        GetParamFieldsPresent(paymentEditParams),
        "WHERE order_id = ?",
        []any{paymentEditParams.Id},
    )
    stmt, err := dbOrTx.Prepare(query)

    if err != nil {
        r.logger.Error("Error when creating query template", err.Error())
        return nil, err
    }
    _, err = stmt.Exec(queryArgs...)
    if err != nil {
        r.logger.Error("Error when trying to update a payment: ", err.Error())
        return nil, err
    } else {
        return r.getByOrderIdTx(paymentEditParams.Id, tx), nil
    }

}

func PaymentFromResult(result RowScanner) (*models.Payment, error) {
    payment := models.Payment{}

    err := result.Scan(
        &payment.Id,
        &payment.OrderId,
        &payment.Method,
        &payment.Status,
        &payment.TravelerAmountCents,
        &payment.TravelerCurrency,
        &payment.ReceivedAmountCents,
        &payment.ReceivedCurrency,
        &payment.RevertedAmountCents,
        &payment.RevertedCurrency,
        &payment.CreatedAt,
    )

    if err != nil {
        return nil, err
    } else {
        return &payment, nil
    }
}
