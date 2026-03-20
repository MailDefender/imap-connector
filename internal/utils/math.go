package utils

func MinOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if i < min {
			min = i
		}
	}

	return min
}
