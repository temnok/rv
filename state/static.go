package state

type StaticState struct {
	Priv int
	PC   int

	X, F [32]int

	CSR CSR

	Reservation int
}
