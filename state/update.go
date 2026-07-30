package state

import "github.com/temnok/rv/csr"

const (
	UpdatePC          = 1 << 0
	UpdatePriv        = 1 << 1
	UpdateReservation = 1 << 2

	UpdateXreg = 1 << 3
	UpdateFreg = 1 << 4
	UpdateCreg = 1 << 5

	UpdateFcsr    = 1 << 6
	UpdateMstatus = 1 << 7
	UpdateMepc    = 1 << 8
	UpdateSepc    = 1 << 9
	UpdateMcause  = 1 << 10
	UpdateScause  = 1 << 11
	UpdateMtval   = 1 << 12
	UpdateStval   = 1 << 13
	UpdateMcycle  = 1 << 14
	UpdateMip     = 1 << 15
	UpdateUart    = 1 << 16

	UpdateRAM = 1 << 17
)

type Updates struct {
	Targets int

	PC          int
	Priv        int
	Reservation int

	Xreg, Xval int
	Freg, Fval int
	Creg, Cval int

	Fcsr    int
	Mstatus int
	Xepc    int
	Xcause  int
	Xtval   int
	Mcycle  int
	Mip     int
	Uart    int

	RAMPos, RAMVal int
}

func Update(cpu *CPU) {
	up := &cpu.Update

	if up.Targets&UpdatePC != 0 {
		cpu.PC = up.PC
	}

	if up.Targets&UpdatePriv != 0 {
		cpu.CSR.Priv = up.Priv
	}

	if up.Targets&UpdateReservation != 0 {
		cpu.Reservation = up.Reservation
	}

	if up.Targets&UpdateFcsr != 0 {
		cpu.CSR.Fcsr = up.Fcsr
	}

	if up.Targets&UpdateMstatus != 0 {
		cpu.CSR.Mstatus = up.Mstatus
	}

	if up.Targets&UpdateMepc != 0 {
		cpu.CSR.Mepc = up.Xepc
	}

	if up.Targets&UpdateSepc != 0 {
		cpu.CSR.Sepc = up.Xepc
	}

	if up.Targets&UpdateMcause != 0 {
		cpu.CSR.Mcause = up.Xcause
	}

	if up.Targets&UpdateScause != 0 {
		cpu.CSR.Scause = up.Xcause
	}

	if up.Targets&UpdateMtval != 0 {
		cpu.CSR.Mtval = up.Xtval
	}

	if up.Targets&UpdateStval != 0 {
		cpu.CSR.Stval = up.Xtval
	}

	if up.Targets&UpdateMcycle != 0 {
		cpu.CSR.Mcycle = up.Mcycle
	}

	if up.Targets&UpdateMip != 0 {
		cpu.CSR.Mip = up.Mip
	}

	if up.Targets&UpdateUart != 0 {
		cpu.CSR.Uart = up.Uart
	}

	if up.Targets&UpdateXreg != 0 && up.Xreg != 0 {
		cpu.X[up.Xreg] = up.Xval
	}

	if up.Targets&UpdateFreg != 0 {
		cpu.F[up.Freg] = up.Fval
	}

	if up.Targets&UpdateCreg != 0 {
		csr.Write(&cpu.CSR, up.Creg, up.Cval)
	}

	if up.Targets&UpdateRAM != 0 {
		cpu.RAM[up.RAMPos] = up.RAMVal
	}
}
