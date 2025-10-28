package bi

func Ts(val, i, n int) int {
	return (val >> i) & (1<<n - 1)
}

func T(val, i int) int {
	return (val >> i) & 1
}
