package main

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
