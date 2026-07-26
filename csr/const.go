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
	Sscratch   = 0x140
	Sepc       = 0x141
	Scause     = 0x142
	Stval      = 0x143
	Sip        = 0x144
	Stimecmp   = 0x14D
	Satp       = 0x180

	Mstatus       = 0x300
	Misa          = 0x301
	Medeleg       = 0x302
	Mideleg       = 0x303
	Mie           = 0x304
	Mtvec         = 0x305
	Mcounteren    = 0x306
	Menvcfg       = 0x30A
	Mcountinhibit = 0x320
	Mscratch      = 0x340
	Mepc          = 0x341
	Mcause        = 0x342
	Mtval         = 0x343
	Mip           = 0x344
	Mcycle        = 0xB00
	Minstret      = 0xB02

	Cycle   = 0xC00
	Time    = 0xC01
	Instret = 0xC02

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

	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#mcounteren
	McounterenCY = 0
	McounterenTM = 1
	McounterenIR = 2

	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#3-1-1-12-machine-counter-inhibit-mcountinhibit-register
	McountinhibitCY = 0
	McountinhibitIR = 2

	// https://riscv.github.io/riscv-isa-manual/snapshot/privileged/#_machine_status_mstatus_and_mstatush_registers
	MstatusFS = 13
	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#xlen-control
	MstatusUXL = 32
	MstatusSXL = 34
	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#norm:mstatus_sd_acc
	MstatusSD = 63

	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#3-1-1-6-4-memory-privilege-in-mstatus-register
	MstatusSIE  = 1
	MstatusMIE  = 3
	MstatusSPIE = 5
	MstatusMPIE = 7
	MstatusSPP  = 8
	MstatusMPP  = 11 // Machine Previous Privilege
	MstatusMPRV = 17 // Modify PRiVilege
	MstatusSUM  = 18 // permit Supervisor User Memory access
	MstatusMXR  = 19 // Make eXecutable Readable
	MstatusTVM  = 20
	MstatusTSR  = 22

	McauseI = 63

	MipSSIP = 1
	MipMSIP = 3
	MipSTIP = 5
	MipMTIP = 7
	MipSEIP = 9
	MipMEIP = 11

	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#sec:menvcfg
	MenvcfgSTCE = 63

	FSoff     = 0b_00
	FSinitial = 0b_01
	FSclean   = 0b_10
	FSdirty   = 0b_11

	UartTX = 0
	UartRX = 8
	UartIE = 16
	UartIP = 24
)
