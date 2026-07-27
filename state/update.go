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
	UpdateEpc     = 1 << 8
	UpdateCause   = 1 << 9
	UpdateTval    = 1 << 10
	UpdateMcycle  = 1 << 11
	UpdateMip     = 1 << 12
	UpdateUart    = 1 << 13

	UpdateRAM = 1 << 14
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
	Epc     int
	Cause   int
	Tval    int
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

	if up.Targets&UpdateEpc != 0 {
		if up.Priv == csr.PrivM {
			cpu.CSR.Mepc = up.Epc
		} else {
			cpu.CSR.Sepc = up.Epc
		}
	}

	if up.Targets&UpdateCause != 0 {
		if up.Priv == csr.PrivM {
			cpu.CSR.Mcause = up.Cause
		} else {
			cpu.CSR.Scause = up.Cause
		}
	}

	if up.Targets&UpdateTval != 0 {
		if up.Priv == csr.PrivM {
			cpu.CSR.Mtval = up.Tval
		} else {
			cpu.CSR.Stval = up.Tval
		}
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
