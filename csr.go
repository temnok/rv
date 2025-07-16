package rv

// https://github.com/riscv/riscv-isa-manual/blob/main/src/priv-csrs.adoc#user-content-mcsrnames0

const (
	Fflags = 0x001
	Frm    = 0x002
	Fcsr   = 0x003

	Sstatus    = 0x100
	Sie        = 0x104
	Stvec      = 0x105
	Scounteren = 0x106

	Sscratch = 0x140
	Sepc     = 0x141
	Scause   = 0x142
	Stval    = 0x143
	Sip      = 0x144

	Satp = 0x180

	Mstatus    = 0x300
	Misa       = 0x301
	Medeleg    = 0x302
	Mideleg    = 0x303
	Mie        = 0x304
	Mtvec      = 0x305
	Mcounteren = 0x306

	Mscratch = 0x340
	Mepc     = 0x341
	Mcause   = 0x342
	Mtval    = 0x343
	Mip      = 0x344

	Cycle   = 0xC00
	Time    = 0xC01
	Instret = 0xC02

	Cycleh   = 0xC80
	Timeh    = 0xC81
	Instreth = 0xC82

	Mvendorid = 0xF11
	Marchid   = 0xF12
	Mimpid    = 0xF13
	Mhartid   = 0xF14
)

const (
	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
	MstatusFS   = 13
	MstatusUXL  = 32
	MstatusSXL  = 34
	MstatusSD32 = 31
	MstatusSD64 = 63

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp
	SatpMODE32 = 31
	SatpMODE64 = 60

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
	MstatusSIE  = 1
	MstatusMIE  = 3
	MstatusSPIE = 5
	MstatusMPIE = 7
	MstatusSPP  = 8
	MstatusMPP  = 11
	MstatusMPRV = 17
	MstatusSUM  = 18
	MstatusMXR  = 19
	MstatusTVM  = 20
	MstatusTSR  = 22

	MipMSI = 3
	MipSTI = 5
	MipMTI = 7
	MipSEI = 9

	FSoff     = 0b_00
	FSinitial = 0b_01
	FSclean   = 0b_10
	FSdirty   = 0b_11
)

type CSR struct {
	Fcsr, Cycle, Cycleh, Time, Timeh int

	Stvec, Scounteren, Sscratch, Sepc, Scause, Stval, Sip, Satp int

	Mstatus, Misa, Medeleg, Mideleg, Mie, Mtvec, Mcounteren int
	Mscratch, Mepc, Mcause, Mtval, Mip                      int
	Mvendorid, Marchid, Mimpid, Mhartid                     int
}

func (cpu *CPU) csrAddr(i int, write bool) (reg *int, mask, shift int) {
	if write && bits(i, 10, 2) == 3 || cpu.Priv < bits(i, 8, 2) {
		return
	}

	csr := &cpu.CSR
	mask = -1
	shift = 0

	switch i {
	case Fflags: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#fcsr
		reg = &csr.Fcsr
		mask = 0b_000_11111

	case Frm: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#fcsr
		reg = &csr.Fcsr
		mask = 0b_111_00000
		shift = 5

	case Fcsr: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#fcsr
		reg = &csr.Fcsr
		mask = 0b_111_11111

	case Sstatus: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sstatus
		reg = &csr.Mstatus
		mask = 1<<MstatusSIE | 1<<MstatusSUM | 1<<MstatusMXR | 1<<MstatusSPP
		if !write {
			mask = int(int64(mask) | 1<<MstatusSPIE | 3<<MstatusUXL)
		}

	case Sie: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_interrupt_sip_and_sie_registers
		reg = &csr.Mie // sie
		mask = 1<<MipSEI | 1<<MipSTI

	case Stvec: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_trap_vector_base_address_stvec_register
		reg = &csr.Stvec

	case Scounteren:
		reg = &csr.Scounteren

	case Sscratch: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_scratch_sscratch_register
		reg = &csr.Sscratch

	case Sepc: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_exception_program_counter_sepc_register
		reg = &csr.Sepc

	case Scause: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#scause
		reg = &csr.Scause

	case Stval: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_trap_value_stval_register
		reg = &csr.Stval

	case Sip: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_interrupt_sip_and_sie_registers
		reg = &csr.Mip // sip
		if write {
			mask = 0
		} else {
			mask = 1 << MipSEI
		}

	case Satp: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp
		if cpu.Priv != PrivS || bit(csr.Mstatus, MstatusTVM) == 0 {
			reg = &csr.Satp
		}

	case Mstatus: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
		reg = &csr.Mstatus
		if write {
			mask = ^(3<<MstatusSXL | 3<<MstatusUXL)
		}

	case Misa: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#misa
		reg = &csr.Misa
		if write {
			mask = 0
		}

	case Medeleg: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_delegation_medeleg_and_mideleg_registers
		reg = &csr.Medeleg

	case Mideleg: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_delegation_medeleg_and_mideleg_registers
		reg = &csr.Mideleg

	case Mie: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_interrupt_mip_and_mie_registers
		reg = &csr.Mie

	case Mtvec: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_vector_base_address_mtvec_register
		reg = &csr.Mtvec

	case Mcounteren:
		reg = &csr.Mcounteren

	case Mscratch: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_scratch_mscratch_register
		reg = &csr.Mscratch

	case Mepc: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_exception_program_counter_mepc_register
		reg = &csr.Mepc

	case Mcause: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcause
		reg = &csr.Mcause

	case Mtval: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_value_mtval_register
		reg = &csr.Mtval

	case Mip: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_interrupt_mip_and_mie_registers
		reg = &csr.Mip

	case Cycle:
		reg = &csr.Cycle

	case Time: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_timer_mtime_and_mtimecmp_registers
		reg = &csr.Time

	case Instret: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_zicntr_extension_for_base_counters_and_timers
		reg = &csr.Cycle

	case Cycleh:
		if !cpu.xlen64() {
			reg = &csr.Cycleh
		}

	case Timeh: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_timer_mtime_and_mtimecmp_registers
		if !cpu.xlen64() {
			reg = &csr.Timeh
		}

	case Instreth:
		if !cpu.xlen64() {
			reg = &csr.Cycleh
		}

	case Mvendorid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_vendor_id_mvendorid_register
		reg = &csr.Mvendorid

	case Marchid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_architecture_id_marchid_register
		reg = &csr.Marchid

	case Mimpid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_implementation_id_mimpid_register
		reg = &csr.Mimpid

	case Mhartid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_hart_id_mhartid_register
		reg = &csr.Mhartid
	}

	return
}

func (cpu *CPU) csrRead(i int, val *int) {
	reg, mask, shift := cpu.csrAddr(i, false)

	if reg == nil {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	*val = (*reg & mask) >> shift
}

func (cpu *CPU) csrWrite(i, val int) {
	reg, mask, shift := cpu.csrAddr(i, true)

	if reg == nil {
		cpu.trap(ExceptionIllegalIstruction)
		return
	}

	cpu.Updated.CSRAddr = reg
	cpu.Updated.CSRIndex = i
	cpu.Updated.CSRValue = *reg&^mask | (val<<shift)&mask
}
