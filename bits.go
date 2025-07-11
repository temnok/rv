package rv

func bits(val, i, n int) int {
	return (val >> i) & (1<<n - 1)
}

func bit(val, i int) int {
	return (val >> i) & 1
}

func signBit(val, i int) int {
	return -bit(val, i)
}
