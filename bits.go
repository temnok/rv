package rv

func bit(val, i Xint) Xint {
	return (val >> i) & 1
}

func bits(val, i, n Xint) Xint {
	return (val >> i) & (1<<n - 1)
}

func signedBit(val, i Xint) Xint {
	return -bit(val, i)
}
