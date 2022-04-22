package async

type Result[T any] interface {
	Error() error
	Data() T
}

type result[T any] struct {
	data T
	err  error
}

func (o result[T]) Error() error {
	return o.err
}

func (o result[T]) Data() T {
	return o.data
}
