package model

type UnauthorizedError struct{}

func (e *UnauthorizedError) Error() string {
	return "unauthorized"
}
