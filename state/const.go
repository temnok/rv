package state

const (
	Misa = -1<<63 |
		1<<('i'-'a') | 1<<('m'-'a') | 1<<('a'-'a') | 1<<('c'-'a') |
		1<<('f'-'a') | ('d' - 'a') |
		1<<('u'-'a') | 1<<('s'-'a')

	PageSize = 1 << 12

	PrivU = 0
	PrivS = 1
	PrivM = 3

	AccessFetch = 0
	AccessLoad  = 1
	AccessStore = 3
)
