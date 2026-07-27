package bit

//go:fix inline
func GetN(val, i, n int) int {
	return (val >> i) &^ (-1 << n)
}

//go:fix inline
func Get(val, i int) int {
	return (val >> i) & 1
}

//go:fix inline
func IsSet(val, i int) bool {
	return Get(val, i) == 1
}

//go:fix inline
func IsNotSet(val, i int) bool {
	return Get(val, i) == 0
}
