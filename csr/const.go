package csr

const (
	Fflags = 0x001 // Floating-point FLAGS                        https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#fcsr
	Frm    = 0x002 // Floating-point Rounding Mode                https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#norm:dyn_round_enc
	Fcsr   = 0x003 // Floating-point Control and Status Register  https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#fcsr

	Sstatus    = 0x100 // Supervisor STATUS                              https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#sstatus
	Sie        = 0x104 // Supervisor Interrupt Enable                    https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-3-supervisor-interrupt-sip-and-sie-registers
	Stvec      = 0x105 // Supervisor Trap VECtor                         https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-2-supervisor-trap-vector-base-address-stvec-register
	Scounteren = 0x106 // Supervisor COUNTER ENable                      https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-5-counter-enable-scounteren-register
	Sscratch   = 0x140 // Supervisor SCRATCH                             https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-6-supervisor-scratch-sscratch-register
	Sepc       = 0x141 // Supervisor Exception Program Counter           https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-7-supervisor-exception-program-counter-sepc-register
	Scause     = 0x142 // Supervisor Cause                               https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#scause
	Stval      = 0x143 // Supervisor Trap VALue                          https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-9-supervisor-trap-value-stval-register
	Sip        = 0x144 // Supervisor Interrupt Pending                   https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-3-supervisor-interrupt-sip-and-sie-registers
	Stimecmp   = 0x14D // Supervisor TIMEr CoMPare                       https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#stimecmp
	Satp       = 0x180 // Supervisor Address Translation and Protection  https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#satp

	Mstatus       = 0x300 // Machine STATUS                     https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-6-machine-status-mstatus-and-mstatush-registers
	Misa          = 0x301 // Machine ISA                        https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#misa
	Medeleg       = 0x302 // Machine Exception DELEGation       https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-8-machine-trap-delegation-medeleg-and-mideleg-registers
	Mideleg       = 0x303 // Machine Exception DELEGation       https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-8-machine-trap-delegation-medeleg-and-mideleg-registers
	Mie           = 0x304 // Machine Interrupt Enable           https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-9-machine-interrupt-mip-and-mie-registers
	Mtvec         = 0x305 // Machine Trap VECtor                https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-7-machine-trap-vector-base-address-mtvec-register
	Mcounteren    = 0x306 // Machine COUNTER ENable             https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#mcounteren
	Menvcfg       = 0x30A // Machine ENVironment ConFiGuration  https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#sec:menvcfg
	Mcountinhibit = 0x320 // Machine COUNTer INHIBIT            https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-12-machine-counter-inhibit-mcountinhibit-register
	Mscratch      = 0x340 // Machine SCRATCH                    https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-13-machine-scratch-mscratch-register
	Mepc          = 0x341 // Machine Exception Program Counter  https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-14-machine-exception-program-counter-mepc-register
	Mcause        = 0x342 // Machine CAUSE                      https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#mcause
	Mtval         = 0x343 // Machine Trap VALue                 https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-16-machine-trap-value-mtval-register
	Mip           = 0x344 // Machine Interrupt Enable           https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-9-machine-interrupt-mip-and-mie-registers
	Mcycle        = 0xB00 // Machine CYCLE                      https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-10-hardware-performance-monitor
	Minstret      = 0xB02 // Machine INSTructions RETired       https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-10-hardware-performance-monitor

	Cycle   = 0xC00 // https://docs.riscv.org/reference/isa/v20260120/unpriv/counters.html#6-1-1-zicntr-extension-for-base-counters-and-timers
	Time    = 0xC01 // https://docs.riscv.org/reference/isa/v20260120/unpriv/counters.html#6-1-1-zicntr-extension-for-base-counters-and-timers
	Instret = 0xC02 // https://docs.riscv.org/reference/isa/v20260120/unpriv/counters.html#6-1-1-zicntr-extension-for-base-counters-and-timers

	Mvendorid = 0xF11 // Machine VENDOR ID          https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-2-machine-vendor-id-mvendorid-register
	Marchid   = 0xF12 // Machine ARCHitecture ID    https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-3-machine-architecture-id-marchid-register
	Mimpid    = 0xF13 // Machine IMPlementation ID  https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-4-machine-implementation-id-mimpid-register
	Mhartid   = 0xF14 // Machine HART ID            https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-5-hart-id-mhartid-register

	PrivU = 0 // PRIVilege of User        https://docs.riscv.org/reference/isa/v20260120/priv/priv-intro.html#privilege-levels
	PrivS = 1 // PRIVilege of Supervisor  https://docs.riscv.org/reference/isa/v20260120/priv/priv-intro.html#privilege-levels
	PrivM = 3 // PRIVilege of Machine     https://docs.riscv.org/reference/isa/v20260120/priv/priv-intro.html#privilege-levels

	FflagsNX = 0 // iNeXact            https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#bitdef
	FflagsUF = 1 // UnderFlow          https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#bitdef
	FflagsOF = 2 // OverFlow           https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#bitdef
	FflagsDZ = 3 // Divide by Zero     https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#bitdef
	FflagsNV = 4 // iNValid operation  https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#bitdef

	FcsrRM = 5 // Rounding Mode: https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#norm:dyn_round_enc

	McounterenCY = 0 // CYcle    https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#mcounteren
	McounterenTM = 1 // TiMe     https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#mcounteren
	McounterenIR = 2 // InstRet  https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#mcounteren

	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#3-1-1-6-4-memory-privilege-in-mstatus-register
	MstatusSIE  = 1
	MstatusMIE  = 3
	MstatusSPIE = 5
	MstatusMPIE = 7
	MstatusSPP  = 8
	MstatusMPP  = 11 // Machine Previous Privilege
	MstatusFS   = 13
	MstatusMPRV = 17 // Modify PRiVilege
	MstatusSUM  = 18 // permit Supervisor User Memory access
	MstatusMXR  = 19 // Make eXecutable Readable
	MstatusTVM  = 20
	MstatusTSR  = 22
	MstatusUXL  = 32 // https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#xlen-control
	MstatusSXL  = 34
	MstatusSD   = 63 // https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#norm:mstatus_sd_acc

	McauseI = 63

	MipSSIP = 1
	MipMSIP = 3
	MipSTIP = 5
	MipMTIP = 7
	MipSEIP = 9
	MipMEIP = 11

	// https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#sec:menvcfg
	MenvcfgSTCE = 63

	UartTD = 0
	UartTE = 8
	UartTP = 12
	UartRD = 16
	UartRE = 24
	UartRP = 28
)
