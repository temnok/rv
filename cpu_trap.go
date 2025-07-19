package rv

func (cpu *CPU) isTrapped() bool {
	return cpu.Updated.TrapEnter
}

func (cpu *CPU) trap(cause int) {
	cpu.trapEnter(cause, 0)
}

func (cpu *CPU) trapEnter(cause, tval int) {
	if cpu.isTrapped() {
		panic("double trap")
	}

	mcauseI := cpu.XLen - 1
	isInterrupt := bit(cause, mcauseI) == 1
	causeID := bits(cause, 0, 5)

	deleg := cpu.CSR.Medeleg
	if isInterrupt {
		deleg = cpu.CSR.Mideleg
	}

	effectivePriv := PrivM
	if cpu.Priv <= PrivS && bit(deleg, causeID) == 1 {
		effectivePriv = PrivS
	}

	cpu.Updated.TrapEnter = true
	cpu.Updated.TrapPriv = effectivePriv
	cpu.Updated.TrapXepc = cpu.PC
	cpu.Updated.TrapXcause = cause
	cpu.Updated.TrapXtval = tval

	var tvec int

	switch effectivePriv {
	case PrivM:
		mie := bit(cpu.CSR.Mstatus, MstatusMIE)
		cpu.Updated.TrapMstatus = cpu.CSR.Mstatus&^(3<<MstatusMPP|1<<MstatusMPIE|1<<MstatusMIE) |
			cpu.Priv<<MstatusMPP | mie<<MstatusMPIE

		tvec = cpu.CSR.Mtvec

	case PrivS:
		sie := bit(cpu.CSR.Mstatus, MstatusSIE)
		cpu.Updated.TrapMstatus = cpu.CSR.Mstatus&^(1<<MstatusSPP|1<<MstatusSPIE|1<<MstatusSIE) |
			cpu.Priv<<MstatusSPP | sie<<MstatusSPIE

		tvec = cpu.CSR.Stvec
	}

	cpu.Updated.TrapPC = tvec &^ 3
	if bit(tvec, 0) == 1 && isInterrupt {
		cpu.Updated.TrapPC += causeID * 4
	}
}

// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#otherpriv
// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#privstack
func (cpu *CPU) trapExit(retPriv int) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#virt-control
	trap := cpu.Priv == PrivS && bit(cpu.CSR.Mstatus, MstatusTSR) == 1

	if trap || retPriv > cpu.Priv {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.TrapExit = true

	switch retPriv {
	case PrivM:
		cpu.Updated.TrapPC = cpu.CSR.Mepc
		cpu.Updated.TrapPriv = bits(cpu.CSR.Mstatus, MstatusMPP, 2)

		mie := bit(cpu.CSR.Mstatus, MstatusMPIE)
		cpu.Updated.TrapMstatus = cpu.CSR.Mstatus&^(3<<MstatusMPP) |
			(1<<MstatusMPIE | mie<<MstatusMIE)

	case PrivS:
		cpu.Updated.TrapPC = cpu.CSR.Sepc
		cpu.Updated.TrapPriv = bits(cpu.CSR.Mstatus, MstatusSPP, 1)

		sie := bit(cpu.CSR.Mstatus, MstatusSPIE)
		cpu.Updated.TrapMstatus = cpu.CSR.Mstatus&^(1<<MstatusSPP) |
			(1<<MstatusSPIE | sie<<MstatusSIE)
	}

	if cpu.Priv != PrivM {
		cpu.Updated.TrapMstatus &^= 1 << MstatusMPRV
	}
}
