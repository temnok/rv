package rv

// https://github.com/riscv/riscv-isa-manual/blob/main/src/priv-csrs.adoc#user-content-mcsrnames0

const (
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#misabase
	misaMXL = Xlen - 2

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
	mstatusSXL = 34
	mstatusUXL = 32

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp
	satpMODE = Xlen - 1 - (Xlen/64)*3

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
	mstatusSIE  = 1
	mstatusMIE  = 3
	mstatusSPIE = 5
	mstatusMPIE = 7
	mstatusSPP  = 8
	mstatusMPP  = 11
	mstatusMPRV = 17
	mstatusSUM  = 18
	mstatusMXR  = 19
	mstatusTVM  = 20
	mstatusTSR  = 22

	mcauseI = Xlen - 1

	mipMSI = 3
	mipSTI = 5
	mipMTI = 7
	mipSEI = 9
)

type CSR struct {
	cycle, cycleh, mtime, mtimeh int

	stvec, scounteren, sscratch, sepc, scause, stval, sip, satp int

	mstatus, misa, medeleg, mideleg, mie, mtvec, mcounteren int
	mscratch, mepc, mcause, mtval, mip                      int
	mvendorid, marchid, mimpid, mhartid                     int
}

func (cpu *CPU) csrAccess(i int, val *int, write bool) {
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_csr_address_mapping_conventions
	if write && bits(i, 10, 2) == 3 || cpu.priv < bits(i, 8, 2) {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	csr := &cpu.csr
	mask := -1

	var reg *int

	switch i {
	case 0x100: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sstatus
		reg = &csr.mstatus
		mask = 1<<mstatusSIE | 1<<mstatusSUM | 1<<mstatusMXR | 1<<mstatusSPP
		if !write {
			mask = int(int64(mask) | 1<<mstatusSPIE | 3<<mstatusUXL)
		}

	case 0x104: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_interrupt_sip_and_sie_registers
		reg = &csr.mie
		mask = 1<<mipSEI | 1<<mipSTI

	case 0x105: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_trap_vector_base_address_stvec_register
		reg = &csr.stvec

	case 0x106:
		reg = &csr.scounteren

	case 0x140: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_scratch_sscratch_register
		reg = &csr.sscratch

	case 0x141: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_exception_program_counter_sepc_register
		reg = &csr.sepc

	case 0x142: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#scause
		reg = &csr.scause

	case 0x143: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_trap_value_stval_register
		reg = &csr.stval

	case 0x144: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_interrupt_sip_and_sie_registers
		reg = &csr.mip
		if write {
			mask = 0
		} else {
			mask = 1 << mipSEI
		}

	case 0x180: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp
		reg = &csr.satp
		if cpu.priv == PrivS && bit(csr.mstatus, mstatusTVM) != 0 {
			cpu.trap(ExceptionIllegalIstruction)
			return
		}

	case 0x300: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
		reg = &csr.mstatus
		if write {
			tmp := ^int64(0b_11<<mstatusSXL | 0b_11<<mstatusUXL)
			mask = int(tmp)
		}

	case 0x301: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#misa
		reg = &csr.misa
		if write {
			mask = 0
		}

	case 0x302: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_delegation_medeleg_and_mideleg_registers
		reg = &csr.medeleg

	case 0x303: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_delegation_medeleg_and_mideleg_registers
		reg = &csr.mideleg

	case 0x304: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_interrupt_mip_and_mie_registers
		reg = &csr.mie
		// TODO
		//mask = 1<<mipMSI | 1<<mipSTI | 1<<mipMTI | 1<<mipSEI

	case 0x305: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_vector_base_address_mtvec_register
		reg = &csr.mtvec

	case 0x306:
		reg = &csr.mcounteren

	case 0x340: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_scratch_mscratch_register
		reg = &csr.mscratch

	case 0x341: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_exception_program_counter_mepc_register
		reg = &csr.mepc

	case 0x342: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcause
		reg = &csr.mcause

	case 0x343: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_value_mtval_register
		reg = &csr.mtval

	case 0x344: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_interrupt_mip_and_mie_registers
		reg = &csr.mip
		// TODO
		//mask = 1<<mipMSI | 1<<mipSTI | 1<<mipMTI | 1<<mipSEI

	case 0xC00:
		reg = &csr.cycle

	case 0xC01: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_timer_mtime_and_mtimecmp_registers
		reg = &csr.mtime

	case 0xC02:
		reg = &csr.cycle

	case 0xC80:
		if Xlen32 {
			reg = &csr.cycleh
		}

	case 0xC81: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_timer_mtime_and_mtimecmp_registers
		if Xlen32 {
			reg = &csr.mtimeh
		}

	case 0xC82:
		if Xlen32 {
			reg = &csr.cycleh
		}

	case 0xF11: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_vendor_id_mvendorid_register
		reg = &csr.mvendorid

	case 0xF12: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_architecture_id_marchid_register
		reg = &csr.marchid

	case 0xF13: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_implementation_id_mimpid_register
		reg = &csr.mimpid

	case 0xF14: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_hart_id_mhartid_register
		reg = &csr.mhartid
	}

	if reg == nil {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	if write {
		*reg = *reg&^mask | *val&mask
	} else {
		*val = *reg & mask
	}
}

func (cpu *CPU) csrRead(i int, val *int) {
	cpu.csrAccess(i, val, false)
}

func (cpu *CPU) csrWrite(i, val int) {
	cpu.csrAccess(i, &val, true)
}
