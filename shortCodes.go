//look at the function name
func countDigits(number int) int {
	count := 0
	for number != 0 {
		number /= 10
		count++
	}
	return count
}

//------------------------------------------------------------
