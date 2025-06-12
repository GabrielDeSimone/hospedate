package repositories

import (
    "database/sql"
    "time"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    sqlTmp "github.com/hospedate/backyard/sql"
)

type OrdersRepository interface {
    Save(newOrderRequest *models.NewOrderRequest, property *models.Property, subTotalCents uint) (*models.Order, error)
    saveTx(newOrderRequest *models.NewOrderRequest, property *models.Property, tx *sql.Tx, totalCents uint) (*models.Order, error)
    SaveInPlatform(newOrderRequest *models.NewOrderRequest, property *models.Property, paymentsRepo PaymentsRepository, addressWallet *models.Address, encryptedPk string, totalBilledCents uint) (*models.Order, error)
    GetById(id string) *models.Order
    getByIdTx(id string, tx *sql.Tx) *models.Order
    DeleteById(id string) (int64, error)
    Search(queryParams *models.OrdersSearchParams) []*models.Order
    Edit(orderEditParams models.OrderEditRequest) (*models.Order, error)
    editTx(orderEditParams models.OrderEditRequest, tx *sql.Tx) (*models.Order, error)
    GetEphemeralOrders() []*models.Order
}

type pgOrdersRepository struct {
    db     *sql.DB
    logger log.Logger
}

func NewOrdersRepository(db *sql.DB) OrdersRepository {
    logger := log.NewLogger("OrdersRepository", string(log.INFO_LEVEL))
    repo := &pgOrdersRepository{db: db, logger: logger}

    return repo
}

func (r *pgOrdersRepository) fillOrderSlice(rows *sql.Rows) []*models.Order {
    orders := []*models.Order{}

    for rows.Next() {
        order, err := OrderFromResult(rows)

        if err != nil {
            r.logger.Error("Order Scan failed: ", err.Error())
            return nil
        }
        orders = append(orders, order)
    }
    return orders
}

func (r *pgOrdersRepository) SaveInPlatform(
    newOrderRequest *models.NewOrderRequest,
    property *models.Property,
    paymentsRepo PaymentsRepository,
    addressWallet *models.Address,
    encryptedPk string,
    totalBilledCents uint,
) (*models.Order, error) {

    tx, err := r.db.Begin()
    if err != nil {
        r.logger.Error("Cannot create transaction", err.Error())
        return nil, UnknownError
    }

    order, err := r.saveTx(newOrderRequest, property, tx, totalBilledCents)
    if err == ErrCollision {
        r.logger.Error("Cannot create order due to a collision, doing rollback", err.Error())
        tx.Rollback()
        return nil, ErrCollision
    } else if err != nil {
        r.logger.Error("Cannot create order, doing rollback", err.Error())
        tx.Rollback()
        return nil, UnknownError
    }

    _, err = paymentsRepo.saveTx(&models.NewPaymentRequest{
        OrderId:             order.Id,
        Method:              "crypto_wallet",
        TravelerAmountCents: order.TotalBilledCents,
        TravelerCurrency:    order.PriceCurrency,
    }, tx)

    if err != nil {
        r.logger.Error("Cannot create payment, doing rollback", err.Error())
        tx.Rollback()
        return nil, UnknownError
    }

    _, err = paymentsRepo.editByOrderIdTx(models.PaymentEditRequestByOrderId{
        Id:         order.Id,
        PrivateKey: &encryptedPk,
    }, tx)
    if err != nil {
        r.logger.Error("Error when saving privateKey, doing rollback", err.Error())
        tx.Rollback()
        return nil, UnknownError
    }

    ephemeral := "ephemeral"
    order, err = r.editTx(models.OrderEditRequest{
        Id:            order.Id,
        WalletAddress: &(addressWallet.Address),
        Status:        &ephemeral,
    }, tx)

    if err != nil {
        r.logger.Error("Cannot edit order, doing rollback", err.Error())
        tx.Rollback()
        return nil, UnknownError
    }

    err = tx.Commit()
    if err != nil {
        r.logger.Error("Cannot commit transaction to create order, doing rollback", err.Error())
        tx.Rollback()
        return nil, UnknownError
    }

    return order, nil
}

func (r *pgOrdersRepository) Save(newOrderRequest *models.NewOrderRequest, property *models.Property, subTotalCents uint) (*models.Order, error) {
    return r.saveTx(newOrderRequest, property, nil, subTotalCents)
}

func (r *pgOrdersRepository) saveTx(newOrderRequest *models.NewOrderRequest, property *models.Property, tx *sql.Tx, totalCents uint) (*models.Order, error) {
    // Get the query template
    check_colision_sql := sqlTmp.Block_check_colision
    new_order_sql := sqlTmp.Order_new_order
    dbOrTx := GetDbOrTx(r.db, tx)

    // Generate an id
    id := GetRandomId()

    // Execute the query
    result := r.db.QueryRow(
        check_colision_sql,
        newOrderRequest.DateStart.String(),
        newOrderRequest.DateEnd.String(),
        newOrderRequest.PropertyId,
    )
    var collision_number int
    err := result.Scan(&collision_number)
    if err != nil {
        r.logger.Error("Unknown error checking block colision: ", err.Error())
        return nil, UnknownError
    } else if collision_number == 0 {
        _, err := dbOrTx.Exec(
            new_order_sql,
            id,
            newOrderRequest.UserId,
            newOrderRequest.PropertyId,
            newOrderRequest.DateStart.String(),
            newOrderRequest.DateEnd.String(),
            newOrderRequest.NumberGuests,
            property.Price,
            totalCents,
            newOrderRequest.OrderType,
        )
        if err != nil {
            r.logger.Error("Unknown error creating new order: ", err.Error())
            return nil, UnknownError
        } else {
            r.logger.Info("New order created with id", id)
            order := r.getByIdTx(id, tx)

            if order == nil {
                r.logger.Error("Could not find a order that was just created with id", id)
                return nil, UnknownError
            }

            return order, nil
        }
    } else {
        return nil, ErrCollision
    }
}

func (r *pgOrdersRepository) GetById(id string) *models.Order {
    return r.getByIdTx(id, nil)
}

func (r *pgOrdersRepository) getByIdTx(id string, tx *sql.Tx) *models.Order {
    dbOrTx := GetDbOrTx(r.db, tx)
    result := dbOrTx.QueryRow(sqlTmp.Orders_fetch_by_id, id)
    order, err := OrderFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Error("Error when fetching an order: ", err.Error())
        return nil
    }

    return order
}

func (r *pgOrdersRepository) GetEphemeralOrders() []*models.Order {
    // Build the base query template
    rows, err := r.db.Query(sqlTmp.Orders_fetch_ephemeral)
    if err == sql.ErrNoRows {
        return []*models.Order{}
    } else if err != nil {
        r.logger.Error("Error when getting ephemeral orders: ", err.Error())
        return nil
    }

    defer rows.Close()

    return r.fillOrderSlice(rows)
}

func (r *pgOrdersRepository) DeleteById(id string) (int64, error) {
    result, err := r.db.Exec(sqlTmp.Orders_delete_by_id, id)
    if err != nil {
        r.logger.Error("Error when trying to delete an order: ", err.Error())
        return 0, err
    } else {
        return result.RowsAffected()
    }
}

func (r *pgOrdersRepository) Search(queryParams *models.OrdersSearchParams) []*models.Order {
    // Build the base query template
    queryTmp := sqlTmp.Orders_search_base_query
    queryArgs := []any{}

    // Add filters if necessary
    if queryParams.HasUserId() {
        queryTmp += " AND user_id = ?"
        queryArgs = append(queryArgs, queryParams.GetUserId())
    }
    if queryParams.HasOwnerId() {
        queryTmp += " AND property_id IN (SELECT id FROM properties WHERE user_id = ?)"
        queryArgs = append(queryArgs, queryParams.GetOwnerId())
    }
    if queryParams.HasStatus() {
        queryTmp += " AND status = ?"
        queryArgs = append(queryArgs, queryParams.GetStatus())
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
        return []*models.Order{}
    } else if err != nil {
        r.logger.Error("Error when searching orders: ", err.Error())
        return nil
    }

    defer rows.Close()

    return r.fillOrderSlice(rows)
}

func (r *pgOrdersRepository) Edit(orderEditParams models.OrderEditRequest) (*models.Order, error) {
    return r.editTx(orderEditParams, nil)
}

func (r *pgOrdersRepository) editTx(orderEditParams models.OrderEditRequest, tx *sql.Tx) (*models.Order, error) {
    // Get DbOrTx object and build query
    dbOrTx := GetDbOrTx(r.db, tx)
    query, queryArgs := buildEditQuery(
        sqlTmp.Orders_edit_base_query,
        GetParamFieldsPresent(orderEditParams),
        "WHERE id = ?",
        []any{orderEditParams.Id},
    )
    stmt, err := dbOrTx.Prepare(query)

    if err != nil {
        r.logger.Error("Error when creating query template", err.Error())
        return nil, UnknownError
    }
    _, err = stmt.Exec(queryArgs...)
    if err != nil {
        r.logger.Error("Error when trying to update an order: ", err.Error())
        return nil, UnknownError
    } else {
        return r.getByIdTx(orderEditParams.Id, tx), nil
    }
}

func OrderFromResult(result RowScanner) (*models.Order, error) {
    order := models.Order{}

    // Create time.Time variables to get date_start and date_end fields
    var date_start_time time.Time
    var date_end_time time.Time

    err := result.Scan(
        &order.Id,
        &order.UserId,
        &order.Status,
        &order.PropertyId,
        &date_start_time,
        &date_end_time,
        &order.NumberGuests,
        &order.Price,
        &order.PriceCurrency,
        &order.TotalBilledCents,
        &order.CanceledBy,
        &order.OrderType,
        &order.CreatedAt,
        &order.WalletAddress,
    )

    // convert time.Time variables to Date
    order.DateStart = *models.NewDateFromTime(date_start_time)
    order.DateEnd = *models.NewDateFromTime(date_end_time)

    if err != nil {
        return nil, err
    } else {
        return &order, nil
    }
}
