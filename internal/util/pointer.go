package util

func ToPointer[T any](t T) *T {
	return &t
}
