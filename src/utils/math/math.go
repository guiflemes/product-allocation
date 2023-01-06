package math

func Sum(slice []int) int {
	result := 0
	for _, v := range slice {
		result += v
	}
	return result
}
