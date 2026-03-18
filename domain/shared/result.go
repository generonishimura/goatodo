package shared

type Result[T any] struct {
	value T
	err   string
	ok    bool
}

func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, ok: true}
}

func Err[T any](err string) Result[T] {
	return Result[T]{err: err, ok: false}
}

func (r Result[T]) IsOk() bool {
	return r.ok
}

func (r Result[T]) IsErr() bool {
	return !r.ok
}

func (r Result[T]) Value() T {
	return r.value
}

func (r Result[T]) Error() string {
	return r.err
}
