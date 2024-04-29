package core

type Writer interface {
	Execute(data interface{}) error
}
