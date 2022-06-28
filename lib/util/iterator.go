package util

type Iterator[E interface{}] interface {
	HasNext() bool
	Next() E
	Remove()
}
