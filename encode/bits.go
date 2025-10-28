package encode

func bits(val, i, n int) int {
	return (val >> i) & (1<<n - 1)
}

func bit(val, i int) int {
	return (val >> i) & 1
}
