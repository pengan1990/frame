package core

type Reader interface {
	IsDone() bool
	Execute() (interface{}, error)
}
