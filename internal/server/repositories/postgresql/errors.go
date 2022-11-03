package postgresql

import "errors"

var (
	ErrDBNoDBConn                    = errors.New("no database connection")
	ErrDBUsernameIsAlreadyTaken      = errors.New("username is already taken")
	ErrDBInvalidUsernamePasswordPair = errors.New("invalid username/password pair")
)
