package state

const (
	PageSize = 1 << 12

	PrivU = 0
	PrivS = 1
	PrivM = 3

	AccessFetch = 0
	AccessLoad  = 1
	AccessStore = 3
)
