package utils

func SliceSome[T comparable](s []T, fn func(v T, i int) bool) bool {
	for i, val := range s {
		if fn(val, i) {
			return true
		}
	}
	return false
}

// 查找项
func SliceFind[T comparable](slice []T, fn func(v T, i int) bool) T {
	for i, val := range slice {
		if fn(val, i) {
			return val
		}
	}
	return *new(T)
}
