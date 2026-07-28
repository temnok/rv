package csr

type Registers struct {
	Fcsr int // Floating-point Control and Status Register  https://docs.riscv.org/reference/isa/v20260120/unpriv/f-st-ext.html#fcsr

	Mcause     int // Machine CAUSE                      https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#mcause
	Mcounteren int // Machine COUNTER ENable             https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#mcounteren
	Mcycle     int // Machine CYCLE                      https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-10-hardware-performance-monitor
	Medeleg    int // Machine Exception DELEGation       https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-8-machine-trap-delegation-medeleg-and-mideleg-registers
	Menvcfg    int // Machine ENVironment ConFiGuration  https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#sec:menvcfg
	Mepc       int // Machine Exception Program Counter  https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-14-machine-exception-program-counter-mepc-register
	Mideleg    int // Machine Interrupt DELEGation       https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-8-machine-trap-delegation-medeleg-and-mideleg-registers
	Mie        int // Machine Interrupt Enable           https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-9-machine-interrupt-mip-and-mie-registers
	Mip        int // Machine Interrupt Pending          https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-9-machine-interrupt-mip-and-mie-registers
	Mscratch   int // Machine SCRATCH                    https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-13-machine-scratch-mscratch-register
	Mstatus    int // Machine STATUS                     https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-6-machine-status-mstatus-and-mstatush-registers
	Mtval      int // Machine Trap VALue                 https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-16-machine-trap-value-mtval-register
	Mtvec      int // Machine Trap VECtor                https://docs.riscv.org/reference/isa/v20260120/priv/machine.html#2-1-1-7-machine-trap-vector-base-address-mtvec-register

	Priv int // current PRIVilege: hidden non-addressable register

	Satp       int // Supervisor Address Translation and Protection  https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#satp
	Scause     int // Supervisor Cause                               https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#scause
	Scounteren int // Supervisor COUNTER ENable                      https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-5-counter-enable-scounteren-register
	Sepc       int // Supervisor Exception Program Counter           https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-7-supervisor-exception-program-counter-sepc-register
	Sscratch   int // Supervisor SCRATCH                             https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-6-supervisor-scratch-sscratch-register
	Stimecmp   int // Supervisor TIMEr CoMPare                       https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#stimecmp
	Stval      int // Supervisor Trap VALue                          https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-9-supervisor-trap-value-stval-register
	Stvec      int // Supervisor Trap VECtor                         https://docs.riscv.org/reference/isa/v20260120/priv/supervisor.html#11-1-1-2-supervisor-trap-vector-base-address-stvec-register

	Uart int // UART state: custom hidden non-addressable register
}
