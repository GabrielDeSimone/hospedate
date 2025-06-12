package repositories

import (
    "database/sql"
    "encoding/json"
    "time"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    sqlTmp "github.com/hospedate/backyard/sql"
    "github.com/lib/pq"
)

type PropertiesRepository interface {
    Save(*models.NewPropertyRequest) (*models.Property, error)
    GetById(id string) *models.Property
    Search(queryParams *models.PropertiesSearchParams) []*models.Property
    GetBlockById(id_property string, id_block string) *models.Block
    SaveBlock(newBlockRequest *models.NewBlockRequest) (*models.Block, error)
    GetBlocksByPropertyId(id_property string) []*models.Block
    DeleteBlockById(id_property string, id_block string) (int64, error)
    Edit(propertyEditParams models.PropertyEditRequest) (*models.Property, error)
    SavePropertyArchived(property *models.Property, propertyEditParams *models.PropertyEditRequest) (*models.Property, error)
    setAirbnbRoomIdNullTx(id string, tx *sql.Tx) (*models.Property, error)
    editTx(propertyEditParams models.PropertyEditRequest, tx *sql.Tx) (*models.Property, error)
    getByIdTx(id string, tx *sql.Tx) *models.Property
}

type pgPropertiesRepository struct {
    db     *sql.DB
    logger log.Logger
}

func NewPropertiesRepository(db *sql.DB) PropertiesRepository {
    logger := log.NewLogger("PropertiesRepository", string(log.INFO_LEVEL))
    repo := &pgPropertiesRepository{db: db, logger: logger}

    return repo
}

func (r *pgPropertiesRepository) fillPropertySlice(rows *sql.Rows) []*models.Property {
    properties := []*models.Property{}

    for rows.Next() {
        property, err := PropertyFromResult(rows, r.logger)

        if err != nil {
            r.logger.Error("Property Scan failed: ", err.Error())
            return nil
        }
        properties = append(properties, property)
    }
    return properties
}

func (r *pgPropertiesRepository) fillBlockSlice(rows *sql.Rows) []*models.Block {
    var blocks = []*models.Block{}

    for rows.Next() {
        var block *models.Block
        block, err := BlockFromResult(rows)
        if err != nil {
            r.logger.Error("Block Scan failed: ", err.Error())
            return nil
        }
        blocks = append(blocks, block)
    }
    return blocks
}

func (r *pgPropertiesRepository) Save(newPropertyRequest *models.NewPropertyRequest) (*models.Property, error) {
    // Get the query template
    new_property_sql := sqlTmp.Properties_new_property

    // Generate an id
    id := GetRandomId()

    // Execute the query
    _, err := r.db.Exec(
        new_property_sql,
        id,
        newPropertyRequest.MaxGuests,
        newPropertyRequest.AirbnbRoomId,
        newPropertyRequest.Price,
        newPropertyRequest.User_id,
        newPropertyRequest.City,
    )

    if err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
            return nil, ErrDuplicateKey
        } else if ok && err.Code.Name() == "foreign_key_violation" {
            return nil, ErrUserDoesNotExist
        } else {
            r.logger.Error("Unknown error creating new property: ", err.Error())
            return nil, UnknownError
        }
    }

    r.logger.Info("New property created with id", id)

    return r.GetById(id), nil
}

func (r *pgPropertiesRepository) GetById(id string) *models.Property {
    return r.getByIdTx(id, nil)
}

func (r *pgPropertiesRepository) getByIdTx(id string, tx *sql.Tx) *models.Property {
    dbOrTx := GetDbOrTx(r.db, tx)
    result := dbOrTx.QueryRow(sqlTmp.Properties_fetch_by_id, id)
    property, err := PropertyFromResult(result, r.logger)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Error("Error when fetching a property: ", err.Error())
        return nil
    }

    return property
}

func (r *pgPropertiesRepository) setAirbnbRoomIdNullTx(id string, tx *sql.Tx) (*models.Property, error) {
    _, err := tx.Exec(sqlTmp.Properties_set_airbnb_room_id_NULL, id)
    if err != nil {
        r.logger.Errorf("Error when trying to set airbnb_room_id null for id %v: %v ", id, err.Error())
        return nil, UnknownError
    } else {
        return r.getByIdTx(id, tx), nil
    }
}

func (r *pgPropertiesRepository) SavePropertyArchived(
    property *models.Property,
    propertyEditParams *models.PropertyEditRequest,
) (*models.Property, error) {

    tx, err := r.db.Begin()
    if err != nil {
        r.logger.Error("Cannot create transaction", err.Error())
        return nil, UnknownError
    }

    _, err = r.setAirbnbRoomIdNullTx(property.Id, tx)
    if err != nil {
        r.logger.Error("Cannot set null airbnb room, doing rollback", err.Error())
        tx.Rollback()
        return nil, UnknownError
    }
    propertyEditParams.ExAirbnbRoomId = property.AirbnbRoomId

    propertyEdited, err := r.editTx(*propertyEditParams, tx)
    if err != nil {
        r.logger.Error("Error when editing property, doing rollback", err.Error())
        tx.Rollback()
        return nil, UnknownError
    }

    err = tx.Commit()
    if err != nil {
        r.logger.Error("Cannot commit transaction to set null airbnb room, doing rollback", err.Error())
        tx.Rollback()
        return nil, UnknownError
    }

    return propertyEdited, nil

}

func (r *pgPropertiesRepository) Search(searchParams *models.PropertiesSearchParams) []*models.Property {
    // Build the base query template
    query_tmp := sqlTmp.Properties_search_base_query
    query_args := []any{}

    // Add filters if necessary
    if searchParams.HasCity() {
        query_tmp += " AND city = ?"
        query_args = append(query_args, searchParams.GetCity())
    }

    if searchParams.HasUserId() {
        query_tmp += " AND user_id = ?"
        query_args = append(query_args, searchParams.GetUserId())
    }

    if searchParams.HasStatus() {
        query_tmp += " AND status = ?"
        query_args = append(query_args, searchParams.GetStatus())
    }

    if searchParams.HasDates() {
        query_tmp += " AND id NOT IN " + sqlTmp.Properties_get_blocked_properties
        query_args = append(query_args, searchParams.GetDateStart().String(), searchParams.GetDateEnd().String())
    }

    if searchParams.HasGuests() {
        query_tmp += " AND max_guests >= ?"
        query_args = append(query_args, searchParams.GetGuests())
    }

    if searchParams.HasBookingOptions() {
        query_tmp += " AND booking_options = ?"
        query_args = append(query_args, searchParams.GetBookingOptions())
    }

    if searchParams.HasAccommodation() {
        query_tmp += " AND accommodation = ?"
        query_args = append(query_args, searchParams.GetAccommodation())
    }

    if searchParams.HasLocation() {
        query_tmp += " AND location = ?"
        query_args = append(query_args, searchParams.GetLocation())
    }

    if searchParams.HasWifi() {
        query_tmp += " AND wifi = ?"
        query_args = append(query_args, searchParams.GetWifi())
    }

    if searchParams.HasTV() {
        query_tmp += " AND tv = ?"
        query_args = append(query_args, searchParams.GetTV())
    }

    if searchParams.HasMicrowave() {
        query_tmp += " AND microwave = ?"
        query_args = append(query_args, searchParams.GetMicrowave())
    }

    if searchParams.HasOven() {
        query_tmp += " AND oven = ?"
        query_args = append(query_args, searchParams.GetOven())
    }

    if searchParams.HasKettle() {
        query_tmp += " AND kettle = ?"
        query_args = append(query_args, searchParams.GetKettle())
    }

    if searchParams.HasToaster() {
        query_tmp += " AND toaster = ?"
        query_args = append(query_args, searchParams.GetToaster())
    }

    if searchParams.HasCoffeeMachine() {
        query_tmp += " AND coffee_machine = ?"
        query_args = append(query_args, searchParams.GetCoffeeMachine())
    }

    if searchParams.HasAC() {
        query_tmp += " AND air_conditioning = ?"
        query_args = append(query_args, searchParams.GetAC())
    }

    if searchParams.HasHeating() {
        query_tmp += " AND heating = ?"
        query_args = append(query_args, searchParams.GetHeating())
    }

    if searchParams.HasParking() {
        query_tmp += " AND parking = ?"
        query_args = append(query_args, searchParams.GetParking())
    }

    if searchParams.HasPool() {
        query_tmp += " AND pool = ?"
        query_args = append(query_args, searchParams.GetPool())
    }

    if searchParams.HasGym() {
        query_tmp += " AND gym = ?"
        query_args = append(query_args, searchParams.GetGym())
    }

    if searchParams.HasBathrooms() {
        query_tmp += " AND half_bathrooms = CAST (? AS DECIMAL)"
        query_args = append(query_args, searchParams.GetBathrooms())
    }

    if searchParams.HasBedrooms() {
        query_tmp += " AND bedrooms = CAST(? AS INTEGER)"
        query_args = append(query_args, searchParams.GetBedrooms())
    }

    // change template symbols from ? to $1, etc
    query_tmp = ReplaceArgsTemplate(query_tmp)

    // Exec the query and get the results
    stmt, err := r.db.Prepare(query_tmp)

    if err != nil {
        r.logger.Error("Error when creating query template", err.Error())
    }

    rows, err := stmt.Query(query_args...)
    if err == sql.ErrNoRows {
        return []*models.Property{}
    } else if err != nil {
        if err, ok := err.(*pq.Error); ok && err.Code.Name() == "invalid_text_representation" {
            // this error happens when status doesn't match
            // the values from the enum.
            return []*models.Property{}
        } else {
            r.logger.Error("Error when searching properties: ", err.Error())
            return nil
        }
    }

    defer rows.Close()

    return r.fillPropertySlice(rows)
}

func (r *pgPropertiesRepository) Edit(propertyEditParams models.PropertyEditRequest) (*models.Property, error) {
    return r.editTx(propertyEditParams, nil)
}

func (r *pgPropertiesRepository) editTx(propertyEditParams models.PropertyEditRequest, tx *sql.Tx) (*models.Property, error) {
    // Get DbOrTx object and build query
    dbOrTx := GetDbOrTx(r.db, tx)
    paramFields := GetParamFieldsPresent(propertyEditParams)
    if propertyEditParams.Images != nil {
        // we need to marshall the slice of strings
        result, err := marshalSlice(propertyEditParams.Images)
        if err != nil {
            r.logger.Error("Error when marshalling slice of images, omitting this column", err.Error())
        } else {
            updatedParamFields := []models.ParamField{}
            for _, paramField := range paramFields {
                if paramField.DbField == "images" {
                    updatedParamFields = append(updatedParamFields,
                        models.ParamField{Value: result, DbField: paramField.DbField},
                    )
                } else {
                    updatedParamFields = append(updatedParamFields, paramField)
                }
            }
            paramFields = updatedParamFields
        }
    }

    query, queryArgs := buildEditQuery(
        sqlTmp.Properties_edit_base_query,
        paramFields,
        "WHERE id = ?",
        []any{propertyEditParams.Id},
    )
    stmt, err := dbOrTx.Prepare(query)

    if err != nil {
        r.logger.Error("Error when creating query template", err.Error())
        return nil, UnknownError
    }
    _, err = stmt.Exec(queryArgs...)
    if err != nil {
        r.logger.Error("Error when trying to update a property: ", err.Error())
        return nil, UnknownError
    } else {
        return r.getByIdTx(propertyEditParams.Id, tx), nil
    }
}

func (r *pgPropertiesRepository) SaveBlock(newBlockRequest *models.NewBlockRequest) (*models.Block, error) {
    // Get the query template
    check_colision_sql := sqlTmp.Block_check_colision
    new_block_sql := sqlTmp.Block_new_block

    // Generate an id
    id := GetRandomId()

    // Execute the query
    result := r.db.QueryRow(
        check_colision_sql,
        newBlockRequest.DateStart.String(),
        newBlockRequest.DateEnd.String(),
        newBlockRequest.PropertyId,
    )
    var collision_number int
    err := result.Scan(&collision_number)
    if err != nil {
        r.logger.Error("Unknown error checking block colision: ", err.Error())
        return nil, UnknownError
    } else if collision_number == 0 {
        _, err := r.db.Exec(
            new_block_sql,
            id,
            newBlockRequest.PropertyId,
            newBlockRequest.DateStart.String(),
            newBlockRequest.DateEnd.String(),
        )
        if err != nil {
            if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
                return nil, ErrDuplicateKey
            } else {
                r.logger.Error("Unknown error creating new block: ", err.Error())
                return nil, UnknownError
            }
        } else {
            r.logger.Info("New block created with id", id)
            block := r.GetBlockById(newBlockRequest.PropertyId, id)

            if block == nil {
                r.logger.Error("Could not find a block that was just created with id", id)
                return nil, UnknownError
            }

            return block, nil
        }
    } else {
        return nil, ErrCollision
    }
}

func (r *pgPropertiesRepository) GetBlockById(id_property string, id_block string) *models.Block {
    result := r.db.QueryRow(sqlTmp.Properties_fetch_block_by_id, id_property, id_block)
    block, err := BlockFromResult(result)

    if err == sql.ErrNoRows {
        return nil
    } else if err != nil {
        r.logger.Error("Error when fetching a block: ", err.Error())
        return nil
    }

    return block
}

func (r *pgPropertiesRepository) GetBlocksByPropertyId(id_property string) []*models.Block {
    rows, err := r.db.Query(sqlTmp.Properties_fetch_blocks_by_property_id, id_property)
    if err == sql.ErrNoRows {
        return []*models.Block{}
    } else if err != nil {
        r.logger.Error("Error fetching blocks by property: ", err.Error())
        return nil
    }
    defer rows.Close()

    return r.fillBlockSlice(rows)
}

func (r *pgPropertiesRepository) DeleteBlockById(id_property string, id_block string) (int64, error) {

    result, err := r.db.Exec(sqlTmp.Properties_delete_block_by_id, id_property, id_block)
    if err != nil {
        r.logger.Error("Error when trying to delete a block: ", err.Error())
        return 0, err
    } else {
        return result.RowsAffected()
    }
}

func PropertyFromResult(result RowScanner, logger log.Logger) (*models.Property, error) {
    property := models.Property{}
    var defaultSlice []uint8

    err := result.Scan(
        &property.Id,
        &property.Title,
        &property.Description,
        &property.MaxGuests,
        &property.AirbnbRoomId,
        &property.ExAirbnbRoomId,
        &property.Price,
        &property.UserId,
        &property.City,
        &property.Status,
        &property.IsVerified,
        &defaultSlice,
        &property.BookingOptions,
        &property.CreatedAt,
        &property.Accommodation,
        &property.Location,
        &property.Wifi,
        &property.TV,
        &property.Microwave,
        &property.Oven,
        &property.Kettle,
        &property.Toaster,
        &property.CoffeeMachine,
        &property.AC,
        &property.Heating,
        &property.Parking,
        &property.Pool,
        &property.Gym,
        &property.HalfBathrooms,
        &property.Bedrooms,
    )

    if len(defaultSlice) == 0 {
        property.Images = []string{}
    } else {
        result, err := unMarshalSlice(defaultSlice)
        if err != nil {
            logger.Error("Error unmarshalling slice", err.Error())
            property.Images = []string{}
        } else {
            property.Images = result
        }
    }

    if err != nil {
        return nil, err
    } else {
        return &property, nil
    }
}

func BlockFromResult(result RowScanner) (*models.Block, error) {
    block := models.Block{}

    // Create time.Time variables to get date_start and date_end fields
    var date_start_time time.Time
    var date_end_time time.Time

    err := result.Scan(
        &block.Id,
        &date_start_time,
        &date_end_time,
        &block.PropertyId,
        &block.CreatedAt,
    )

    // convert time.Time variables to Date
    block.DateStart = *models.NewDateFromTime(date_start_time)
    block.DateEnd = *models.NewDateFromTime(date_end_time)

    if err != nil {
        return nil, err
    } else {
        return &block, nil
    }
}

func marshalSlice(inputSlice []string) ([]byte, error) {
    result, err := json.Marshal(inputSlice)
    return result, err
}

func unMarshalSlice(byteSlice []uint8) ([]string, error) {
    var stringSlice []string

    err := json.Unmarshal(byteSlice, &stringSlice)
    if err != nil {
        return nil, err
    }
    return stringSlice, nil
}
