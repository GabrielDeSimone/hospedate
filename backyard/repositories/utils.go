package repositories

import (
    "database/sql"
    "fmt"
    "math/rand"
    "reflect"
    "regexp"
    "strings"

    "github.com/hospedate/backyard/models"
)

func GetRandomId() string {
    return RandStringBytes(10)
}

func RandStringBytes(n int) string {
    const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func ReplaceArgsTemplate(query string) string {
    placeholder := regexp.MustCompile(`\?`)
    i := 0
    replacer := func(string) string {
        i++
        return fmt.Sprintf("$%d", i)
    }
    return placeholder.ReplaceAllStringFunc(query, replacer)
}

type RowScanner interface {
    Scan(dest ...interface{}) error
}

type DbOrTx interface {
    Exec(query string, args ...any) (sql.Result, error)
    Prepare(query string) (*sql.Stmt, error)
    QueryRow(query string, args ...any) *sql.Row
}

func GetDbOrTx(db *sql.DB, tx *sql.Tx) DbOrTx {
    if tx != nil {
        return tx
    } else {
        return db
    }
}

func GetParamFieldsPresent(editParams interface{}) []models.ParamField {
    /*
     * The argument editParams should be a struct type.
     * It is intended to be used with structs that represent a set of fields
     * to be edited like OrderEditRequest or PaymentEditRequest.
     *
     * This function omits the "Id" field if present and returns a
     * list of ParamField objects for each attribute that is non nil.
     */
    reflectedType := reflect.TypeOf(editParams)
    reflectedValue := reflect.ValueOf(editParams)
    var result []models.ParamField
    for i := 0; i < reflectedType.NumField(); i++ {
        field := reflectedType.Field(i)
        if field.Name == "Id" {
            continue
        }
        fieldValue := reflectedValue.FieldByName(field.Name)
        if !fieldValue.IsNil() {
            result = append(
                result,
                models.ParamField{Value: fieldValue.Interface(), DbField: field.Tag.Get("db")},
            )
        }
    }
    return result
}

func buildEditQuery(baseQuery string, editParams []models.ParamField, whereClause string, extraQueryArgs []any) (string, []any) {
    editSets := []string{}
    queryArgs := []any{}

    for _, paramField := range editParams {
        editSets = append(editSets, paramField.DbField+" = ?")
        queryArgs = append(queryArgs, paramField.Value)
    }
    queryArgs = append(queryArgs, extraQueryArgs...)

    finalQuery := fmt.Sprintf(
        "%v %v %v ",
        baseQuery,
        strings.Join(editSets, ", "),
        whereClause,
    )
    return ReplaceArgsTemplate(finalQuery), queryArgs
}
