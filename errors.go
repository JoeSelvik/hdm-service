package main

// UserError contains error information presented to the users of our API, probably via json.
type UserError struct {
	Msg  string `json:"message"`
	Code int    `json:"error"`
}

// ApplicationError contains information about errors that arise while accessing resources.
type ApplicationError struct {
	Msg  string
	Err  error
	Code int
}

// Error returns a human-readable representation of a ApplicationError.
func (err *ApplicationError) Error() string {
	return err.Msg
}

// UserError return a user-facing error
func (err *ApplicationError) UserError() *UserError {
	return &UserError{Msg: err.Msg, Code: err.Code}
}
