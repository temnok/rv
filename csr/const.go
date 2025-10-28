package csr

const (
	// https://github.com/riscv/riscv-isa-manual/blob/main/src/priv-csrs.adoc#user-content-mcsrnames0
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
	Mstatush   = 0x310
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

	// https://riscv.github.io/riscv-isa-manual/snapshot/unprivileged/#bitdef
	FflagsNX = 0
	FflagsUF = 1
	FflagsOF = 2
	FflagsDZ = 3
	FflagsNV = 4
	FcsrRM   = 5

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
