package service

import "github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/errors"

type errUserAlreadyExists struct {
	message   string
	operation string
}

func NewErrUserAlreadyExists(msg, opName string) errors.ErrInvalidStorageOperation {
	return errUserAlreadyExists{
		message:   msg,
		operation: opName,
	}
}

func (e errUserAlreadyExists) Error() string {
	return e.message
}

func (e errUserAlreadyExists) Operation() string {
	return e.operation
}
