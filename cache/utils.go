package cache

func PointerTo[T any](v T) *T {
	return &v
}
