package ins

import "github.com/temnok/rv/bi"

type Op int

func (op Op) Code() int { return int(op) }
func (op Op) F5() int   { return bi.Ts(op.Code(), 2, 5) }
func (op Op) Rd() int   { return bi.Ts(op.Code(), 7, 5) }
func (op Op) F3() int   { return bi.Ts(op.Code(), 12, 3) }
func (op Op) Rs1() int  { return bi.Ts(op.Code(), 15, 5) }
func (op Op) Rs2() int  { return bi.Ts(op.Code(), 20, 5) }
func (op Op) F7() int   { return bi.Ts(op.Code(), 25, 7) }
