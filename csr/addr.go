package csr

import (
	"github.com/temnok/rv/bi"
	"github.com/temnok/rv/state"
)

func addr(cpu *state.CPU, i int, write bool) (reg *int, mask, shift int) {
	if write && bi.Ts(i, 10, 2) == 3 || cpu.Priv < bi.Ts(i, 8, 2) {
		return
	}

	csr := &cpu.CSR
	mask = -1
	shift = 0

	switch i {
	case Cycle:
		if csr.Mcounteren>>McounterenCY&1 == 1 || cpu.Priv == state.PrivM {
			reg = &csr.Cycle
		}

	case Fcsr: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#fcsr
		reg = &csr.Fcsr
		mask = 0b_111_11111

	case Fflags: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#fcsr
		reg = &csr.Fcsr
		mask = 0b_000_11111

	case Frm: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#fcsr
		reg = &csr.Fcsr
		mask = 0b_111_00000
		shift = 5

	case Instret: // https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#_zicntr_extension_for_base_counters_and_timers
		if csr.Mcounteren>>McounterenIR&1 == 1 || cpu.Priv == state.PrivM {
			reg = &csr.Cycle
		}

	case Marchid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_architecture_id_marchid_register
		reg = &csr.Marchid

	case Mcause: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#mcause
		reg = &csr.Mcause

	case Mcounteren:
		reg = &csr.Mcounteren

	case Medeleg: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_delegation_medeleg_and_mideleg_registers
		reg = &csr.Medeleg

	//case Menvcfg: // https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#sec:menvcfg
	//	reg = &csr.Menvcfg
	//	mask = -1 << 63

	case Mepc: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_exception_program_counter_mepc_register
		reg = &csr.Mepc

	case Mhartid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_hart_id_mhartid_register
		reg = &csr.Mhartid

	case Mideleg: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_delegation_medeleg_and_mideleg_registers
		reg = &csr.Mideleg

	case Mie: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_interrupt_mip_and_mie_registers
		reg = &csr.Mie

	case Mimpid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_implementation_id_mimpid_register
		reg = &csr.Mimpid

	case Mip: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_interrupt_mip_and_mie_registers
		reg = &csr.Mip
		if write {
			mask = 1 << MipSSI
		}

	case Misa: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#misa
		reg = &csr.Misa
		if write {
			mask = 0
		}

	case Mscratch: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_scratch_mscratch_register
		reg = &csr.Mscratch

	case Mstatus: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
		reg = &csr.Mstatus
		if write {
			mask = ^(3<<MstatusSXL | 3<<MstatusUXL)
		}

	case Mtval: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_value_mtval_register
		reg = &csr.Mtval

	case Mtvec: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_trap_vector_base_address_mtvec_register
		reg = &csr.Mtvec

	case Mvendorid: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_vendor_id_mvendorid_register
		reg = &csr.Mvendorid

	case Satp: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#satp
		if cpu.Priv != state.PrivS || bi.T(csr.Mstatus, MstatusTVM) == 0 {
			reg = &csr.Satp
		}

	case Scause: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#scause
		reg = &csr.Scause

	//case Stimecmp: // https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#stimecmp
	//	reg = &csr.Stimecmp

	case Scounteren:
		reg = &csr.Scounteren

	case Sepc: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_exception_program_counter_sepc_register
		reg = &csr.Sepc

	case Sie: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_interrupt_sip_and_sie_registers
		reg = &csr.Mie // sie
		mask = 1<<MipSEI | 1<<MipSTI

	case Sip: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_interrupt_sip_and_sie_registers
		reg = &csr.Mip // sip
		if write {
			mask = 0
		} else {
			mask = 1 << MipSEI
		}

	case Sscratch: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_scratch_sscratch_register
		reg = &csr.Sscratch

	case Sstatus: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#sstatus
		reg = &csr.Mstatus
		mask = 1<<MstatusSIE | 1<<MstatusSUM | 1<<MstatusMXR | 1<<MstatusSPP
		if !write {
			mask |= 1<<MstatusSPIE | 3<<MstatusUXL
		}

	case Stval: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_trap_value_stval_register
		reg = &csr.Stval

	case Stvec: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_supervisor_trap_vector_base_address_stvec_register
		reg = &csr.Stvec

	case Time: // https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_timer_mtime_and_mtimecmp_registers
		if csr.Mcounteren>>McounterenTM&1 == 1 || cpu.Priv == state.PrivM {
			reg = &csr.Time
		}
	}

	return
}
