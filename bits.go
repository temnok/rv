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

func setBits(addr *int, i, n, val int) {
	mask := (1<<n - 1) << i
	*addr = *addr&^mask | val<<i
}

func orBit(val int) int {
	if val == 0 {
		return 0
	} else {
		return 1
	}
}
