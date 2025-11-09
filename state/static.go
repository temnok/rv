package state

type StaticState struct {
	Priv int
	PC   int

	X, F [32]int

	CSR CSR

	Reserved     bool
	ReservedAddr int

	ICache Cache
}
