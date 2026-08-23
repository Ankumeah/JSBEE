package database

import (
	"errors"
	"fmt"
)

var ErrExists = errors.New("Exists")
var ErrExistUser = fmt.Errorf("%w: User exists", ErrExists)

var ErrInvalid = errors.New("Invalid")
var ErrInvalidUser = fmt.Errorf("%w: User invalid", ErrInvalid)
