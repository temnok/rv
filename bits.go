package rv

func bit(val, i int) int {
	return (val >> i) & 1
}

func bits(val, i, n int) int {
	return (val >> i) & (1<<n - 1)
}

func signedBit(val, i int) int {
	return -bit(val, i)
}
