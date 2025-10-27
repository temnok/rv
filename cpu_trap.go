package rv

func (cpu *CPU) isTrapped() bool {
	return cpu.Update.TrapEnter
}

func (cpu *CPU) trap(cause int) {
	cpu.trapEnter(cause, 0)
}

func (cpu *CPU) trapEnter(cause, tval int) {
	if cpu.isTrapped() {
		panic("double trap")
	}

	mcauseI := cpu.Xmask()
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

	cpu.Update.TrapEnter = true
	cpu.Update.TrapPriv = effectivePriv
	cpu.Update.TrapXepc = cpu.PC
	cpu.Update.TrapXcause = cause
	cpu.Update.TrapXtval = tval

	var tvec int

	switch effectivePriv {
	case PrivM:
		mie := bit(cpu.CSR.Mstatus, MstatusMIE)
		cpu.Update.TrapMstatus = cpu.CSR.Mstatus&^(3<<MstatusMPP|1<<MstatusMPIE|1<<MstatusMIE) |
			cpu.Priv<<MstatusMPP | mie<<MstatusMPIE

		tvec = cpu.CSR.Mtvec

	case PrivS:
		sie := bit(cpu.CSR.Mstatus, MstatusSIE)
		cpu.Update.TrapMstatus = cpu.CSR.Mstatus&^(1<<MstatusSPP|1<<MstatusSPIE|1<<MstatusSIE) |
			cpu.Priv<<MstatusSPP | sie<<MstatusSPIE

		tvec = cpu.CSR.Stvec
	}

	cpu.Update.TrapPC = tvec &^ 3
	if bit(tvec, 0) == 1 && isInterrupt {
		cpu.Update.TrapPC += causeID * 4
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

	cpu.Update.TrapExit = true

	switch retPriv {
	case PrivM:
		cpu.Update.TrapPC = cpu.CSR.Mepc
		cpu.Update.TrapPriv = bits(cpu.CSR.Mstatus, MstatusMPP, 2)

		mie := bit(cpu.CSR.Mstatus, MstatusMPIE)
		cpu.Update.TrapMstatus = cpu.CSR.Mstatus&^(3<<MstatusMPP) |
			(1<<MstatusMPIE | mie<<MstatusMIE)

	case PrivS:
		cpu.Update.TrapPC = cpu.CSR.Sepc
		cpu.Update.TrapPriv = bits(cpu.CSR.Mstatus, MstatusSPP, 1)

		sie := bit(cpu.CSR.Mstatus, MstatusSPIE)
		cpu.Update.TrapMstatus = cpu.CSR.Mstatus&^(1<<MstatusSPP) |
			(1<<MstatusSPIE | sie<<MstatusSIE)
	}

	if cpu.Priv != PrivM {
		cpu.Update.TrapMstatus &^= 1 << MstatusMPRV
	}
}
