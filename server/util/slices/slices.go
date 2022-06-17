package slices

// Remove removes the first occurrence of the given item from the slice.
// If the item is not found, the slice is unchanged.
func Remove[T comparable](slice []T, item T) []T {
	for i, v := range slice {
		if v == item {
			slice[i] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}

	return slice
}

// RemoveAt removes the element at the given index from the slice.
// If the index is out of bounds, RemoveAt panics.
func RemoveAt[T comparable](slice []T, index int) []T {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

// Filter returns a new slice containing only the elements of the original slice that satisfy the predicate.
// The original slice is not modified.
func Filter[T comparable](slice []T, f func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map returns a new slice containing the results of applying the function to each element of the original slice.
// The original slice is not modified.
func Map[T any, U any](slice []T, f func(T) U) []U {
	var result []U
	for _, v := range slice {
		result = append(result, f(v))
	}
	return result
}

// Contains returns true if the given item is in the slice.
// The original slice is not modified.
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first occurrence of the given item in the slice.
// If the item is not found, -1 is returned.
func IndexOf[T comparable](slice []T, item T) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

type numeric interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

// Max returns the largest element in the slice or 0 if the slice is empty.
func Max[T numeric](slice []T) T {
	if len(slice) == 0 {
		return 0
	}

	max := slice[0]
	for _, v := range slice {
		if v > max {
			max = v
		}
	}
	return max
}

// Min returns the smallest element in the slice or 0 if the slice is empty.
func Min[T numeric](slice []T) T {
	if len(slice) == 0 {
		return 0
	}

	min := slice[0]
	for _, v := range slice {
		if v < min {
			min = v
		}
	}
	return min
}

// Sum returns the sum of all elements in the slice.
func Sum[T numeric](slice []T) T {
	var sum T
	for _, v := range slice {
		sum += v
	}
	return sum
}

// Average returns the average of all elements in the slice.
func Average[T numeric](slice []T) T {
	return Sum(slice) / T(len(slice))
}

// Unique returns a new slice containing the unique elements of the original slice.
// Only the first occurrence of each unique element is kept.
func Unique[T comparable](slice []T) []T {
	var have map[T]struct{}
	for _, v := range slice {
		have[v] = struct{}{}
	}
	var result []T
	for k := range have {
		result = append(result, k)
	}
	return result
}

// UniqueBy returns a new slice containing the unique elements of the original slice,
// where the uniqueness is determined by the return value of a function applied to each element.
// Only the first occurrence of each unique element is kept.
func UniqueBy[T any, U comparable](slice []T, f func(T) U) []T {
	var have map[U]T
	for _, v := range slice {
		unique := f(v)
		if _, ok := have[unique]; !ok {
			have[unique] = v
		}
	}
	var result []T
	for _, v := range have {
		result = append(result, v)
	}
	return result
}

// Reduce returns the result of applying a function to each element of the slice and accumulating the result.
func Reduce[T any, U any](slice []T, f func(U, T) U, initial U) U {
	result := initial
	for _, v := range slice {
		result = f(result, v)
	}
	return result
}

// Some returns true if any element of the slice satisfies the predicate.
func Some[T comparable](slice []T, f func(T) bool) bool {
	for _, v := range slice {
		if f(v) {
			return true
		}
	}
	return false
}

// Every returns true if all elements of the slice satisfy the predicate.
func Every[T comparable](slice []T, f func(T) bool) bool {
	for _, v := range slice {
		if !f(v) {
			return false
		}
	}
	return true
}

// Associate returns a map that maps each value of a slice to an element of another.
func Associate[K comparable, V any](keys []K, values []V) map[K]V {
	result := make(map[K]V)
	for i, v := range keys {
		if i >= len(values) {
			break
		}
		result[v] = values[i]
	}
	return result
}

// AssociateBy returns a map that maps each element of the slice to the result of applying the function to it.
func AssociateBy[K comparable, V any](keys []K, valueGenerator func(K) V) map[K]V {
	result := make(map[K]V)
	for _, v := range keys {
		result[v] = valueGenerator(v)
	}
	return result
}

// AssociateReverseBy returns a map that maps each result of applying the function to an element to the original element.
func AssociateReverseBy[V any, K comparable](values []V, keyGenerator func(V) K) map[K]V {
	result := make(map[K]V)
	for _, v := range values {
		result[keyGenerator(v)] = v
	}
	return result
}

// Intersperse separates elements of a slice with values and returns the result.
// If the slice is empty, Intersperse returns an empty slice.
func Intersperse[T any](values []T, separator T) []T {
	if len(values) == 0 {
		return []T{}
	}

	result := []T{values[0]}
	for _, v := range values[1:] {
		result = append(result, separator, v)
	}
	return result
}

// IntersperseBy separates elements of a slice with the result of a function and returns the result.
// The function is given the original index of the element before it in the slice.
// If the slice is empty, IntersperseBy returns an empty slice.
func IntersperseBy[T any](values []T, separatorGenerator func(int) T) []T {
	if len(values) == 0 {
		return []T{}
	}

	result := []T{values[0]}
	for i, v := range values[1:] {
		result = append(result, separatorGenerator(i), v)
	}
	return result
}
