package repositories

type RepoError string

const ErrDuplicateKey RepoError = "ErrDuplicateKey"

const UnknownError RepoError = "UnknownError"

const ErrCollision RepoError = "ErrCollision"

const ErrUserDoesNotExist RepoError = "ErrUserDoesNotExist"

const ErrWrongMigrationChain RepoError = "ErrWrongMigrationChain"

func (e RepoError) Error() string {
    return string(e)
}
