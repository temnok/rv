package rv

type CLINT struct {
	cpu                       *CPU
	mswi, mtimecmp, mtimecmph int32
}

func (c *CLINT) init(cpu *CPU) {
	*c = CLINT{
		cpu: cpu,
	}
}

func (c *CLINT) access(addr int32, data *int32, isWrite bool) {

}
