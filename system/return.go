package systemutil

func MustReturn[T any](t T, _ error) T {
	return t
}

func ReturnRef[T any](t T) *T {
	return &t
}

func MustReturnRef[T any](t T, _ error) *T {
	return &t
}

func ReturnUnref[T any](t *T) T {
	return *t
}

func MustReturnUnref[T any](t *T, _ error) T {
	return *t
}
