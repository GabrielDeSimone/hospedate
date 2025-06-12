package services

type ServiceErr string

const ErrKeyNotExistInJson ServiceErr = "ErrKeyNotExistInJson"
const ErrHttpRequestStatus ServiceErr = "ErrHttpRequestStatus"
const ErrOrderNotFound ServiceErr = "ErrOrderNotFound"
const UnknownError ServiceErr = "UnknownError"
const ErrUserNotFound ServiceErr = "ErrUserNotFound"

func (e ServiceErr) Error() string {
    return string(e)
}
