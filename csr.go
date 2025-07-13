package rv

// https://github.com/riscv/riscv-isa-manual/blob/main/src/priv-csrs.adoc#user-content-mcsrnames0

const (
	csrSstatus    = 0x100
	csrSie        = 0x104
	csrStvec      = 0x105
	csrScounteren = 0x106

	csrSscratch = 0x140
	csrSepc     = 0x141
	csrScause   = 0x142
	csrStval    = 0x143
	csrSip      = 0x144

	csrSatp = 0x180

	csrMstatus    = 0x300
	csrMisa       = 0x301
	csrMedeleg    = 0x302
	csrMideleg    = 0x303
	csrMie        = 0x304
	csrMtvec      = 0x305
	csrMcounteren = 0x306

	csrMscratch = 0x340
	csrMepc     = 0x341
	csrMcause   = 0x342
	csrMtval    = 0x343
	csrMip      = 0x344

	csrCycle   = 0xC00
	csrTime    = 0xC01
	csrInstret = 0xC02

	csrCycleh   = 0xC80
	csrTimeh    = 0xC81
	csrInstreth = 0xC82

	csrMvendorid = 0xF11
	csrMarchid   = 0xF12
	csrMimpid    = 0xF13
	csrMhartid   = 0xF14
)

const (
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
	mstatusSXL = 34
	mstatusUXL = 32

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp
	satpMODE32 = 31
	satpMODE64 = 60

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

	mipMSI = 3
	mipSTI = 5
	mipMTI = 7
	mipSEI = 9
)

type CSR struct {
	cycle, cycleh, time, timeh int

	stvec, scounteren, sscratch, sepc, scause, stval, sip, satp int

	mstatus, misa, medeleg, mideleg, mie, mtvec, mcounteren int
	mscratch, mepc, mcause, mtval, mip                      int
	mvendorid, marchid, mimpid, mhartid                     int
}

func (cpu *CPU) csrAddr(i int, write bool) (reg *int, mask int) {
	if write && bits(i, 10, 2) == 3 || cpu.priv < bits(i, 8, 2) {
		return nil, 0
	}

	csr := &cpu.csr
	mask = -1

	switch i {
	case csrSstatus: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sstatus
		reg = &csr.mstatus
		mask = 1<<mstatusSIE | 1<<mstatusSUM | 1<<mstatusMXR | 1<<mstatusSPP
		if !write {
			mask = int(int64(mask) | 1<<mstatusSPIE | 3<<mstatusUXL)
		}

	case csrSie: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_interrupt_sip_and_sie_registers
		reg = &csr.mie // sie
		mask = 1<<mipSEI | 1<<mipSTI

	case csrStvec: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_trap_vector_base_address_stvec_register
		reg = &csr.stvec

	case csrScounteren:
		reg = &csr.scounteren

	case csrSscratch: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_scratch_sscratch_register
		reg = &csr.sscratch

	case csrSepc: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_exception_program_counter_sepc_register
		reg = &csr.sepc

	case csrScause: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#scause
		reg = &csr.scause

	case csrStval: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_trap_value_stval_register
		reg = &csr.stval

	case csrSip: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_interrupt_sip_and_sie_registers
		reg = &csr.mip // sip
		if write {
			mask = 0
		} else {
			mask = 1 << mipSEI
		}

	case csrSatp: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp
		if cpu.priv != PrivS || bit(csr.mstatus, mstatusTVM) == 0 {
			reg = &csr.satp
		}

	case csrMstatus: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
		reg = &csr.mstatus
		if write {
			mask = ^(3<<mstatusSXL | 3<<mstatusUXL)
		}

	case csrMisa: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#misa
		reg = &csr.misa
		if write {
			mask = 0
		}

	case csrMedeleg: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_delegation_medeleg_and_mideleg_registers
		reg = &csr.medeleg

	case csrMideleg: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_delegation_medeleg_and_mideleg_registers
		reg = &csr.mideleg

	case csrMie: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_interrupt_mip_and_mie_registers
		reg = &csr.mie
		// TODO
		//mask = 1<<mipMSI | 1<<mipSTI | 1<<mipMTI | 1<<mipSEI

	case csrMtvec: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_vector_base_address_mtvec_register
		reg = &csr.mtvec

	case csrMcounteren:
		reg = &csr.mcounteren

	case csrMscratch: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_scratch_mscratch_register
		reg = &csr.mscratch

	case csrMepc: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_exception_program_counter_mepc_register
		reg = &csr.mepc

	case csrMcause: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcause
		reg = &csr.mcause

	case csrMtval: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_value_mtval_register
		reg = &csr.mtval

	case csrMip: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_interrupt_mip_and_mie_registers
		reg = &csr.mip
		// TODO
		//mask = 1<<mipMSI | 1<<mipSTI | 1<<mipMTI | 1<<mipSEI

	case csrCycle:
		reg = &csr.cycle

	case csrTime: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_timer_mtime_and_mtimecmp_registers
		reg = &csr.time

	case csrInstret: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_zicntr_extension_for_base_counters_and_timers
		reg = &csr.cycle

	case csrCycleh:
		if !cpu.xlen64() {
			reg = &csr.cycleh
		}

	case csrTimeh: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_timer_mtime_and_mtimecmp_registers
		if !cpu.xlen64() {
			reg = &csr.timeh
		}

	case csrInstreth:
		if !cpu.xlen64() {
			reg = &csr.cycleh
		}

	case csrMvendorid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_vendor_id_mvendorid_register
		reg = &csr.mvendorid

	case csrMarchid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_architecture_id_marchid_register
		reg = &csr.marchid

	case csrMimpid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_implementation_id_mimpid_register
		reg = &csr.mimpid

	case csrMhartid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_hart_id_mhartid_register
		reg = &csr.mhartid
	}

	return
}

func (cpu *CPU) csrRead(i int, val *int) {
	reg, mask := cpu.csrAddr(i, false)

	if reg == nil {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	*val = *reg & mask
}

func (cpu *CPU) csrWrite(i, val int) {
	reg, mask := cpu.csrAddr(i, true)

	if reg == nil {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.updated.csrAddr = reg
	cpu.updated.csrIndex = i
	cpu.updated.csrValue = *reg&^mask | val&mask
}
