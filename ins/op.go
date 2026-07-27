package ins

type Op int

func (op Op) code() int { return int(op) }
func (op Op) f5() int   { return op.code() >> 2 & 0x1F }
func (op Op) rd() int   { return op.code() >> 7 & 0x1F }
func (op Op) f3() int   { return op.code() >> 12 & 7 }
func (op Op) rs1() int  { return op.code() >> 15 & 0x1F }
func (op Op) rs2() int  { return op.code() >> 20 & 0x1F }
func (op Op) f7() int   { return op.code() >> 25 & 0x7F }
