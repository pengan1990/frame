package core

type CheckPointer interface {
	Save(data interface{}) error
}
