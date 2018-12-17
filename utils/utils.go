package utils

// Map function to array of ints
func Map(vs []string, f func(string) int) []int {
	vsm := make([]int, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// MapIntToStr map function to array of ints and returns array of strings
func MapIntToStr(vs []int, f func(int) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}