package ins

import "github.com/temnok/rv/bi"

type Op int

func (op Op) code() int { return int(op) }
func (op Op) f5() int   { return bi.Ts(op.code(), 2, 5) }
func (op Op) rd() int   { return bi.Ts(op.code(), 7, 5) }
func (op Op) f3() int   { return bi.Ts(op.code(), 12, 3) }
func (op Op) rs1() int  { return bi.Ts(op.code(), 15, 5) }
func (op Op) rs2() int  { return bi.Ts(op.code(), 20, 5) }
func (op Op) f7() int   { return bi.Ts(op.code(), 25, 7) }
