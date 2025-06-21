package rv

func bit(val, i int32) int32 {
	return (val >> i) & 1
}

func bits(val, i, n int32) int32 {
	return (val >> i) & (1<<n - 1)
}

func signedBit(val, i int32) int32 {
	return -bit(val, i)
}
