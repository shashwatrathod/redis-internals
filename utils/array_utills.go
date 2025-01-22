package utils

// Map applies the provided function f to each element of the slice arr
// and returns a new slice containing the transformed elements.
func Map(arr []interface{}, f func(interface{}) interface{}) []interface{} {
	result := make([]interface{}, len(arr))

	for i, v := range arr {
		result[i] = f(v)
	}

	return result
}
