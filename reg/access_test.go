package reg

import "testing"

func TestAccess(t *testing.T) {
	const ref = 0x7766554433221100

	tests := []struct {
		inReg, inData, outReg, outData, offset, width int
		isWrite                                       bool
	}{
		{
			inReg: ref, inData: 0, outReg: ref, outData: ref, offset: 0, width: 8,
		},
		{
			inReg: ref, inData: 0, outReg: ref, outData: 0x33221100, offset: 0, width: 4,
		},
		{
			inReg: ref, inData: 0, outReg: ref, outData: 0x77665544, offset: 4, width: 4,
		},
		{
			inReg: ref, inData: 0, outReg: ref, outData: 0x3322, offset: 2, width: 2,
		},
		{
			inReg: ref, inData: 0, outReg: ref, outData: 0x00, offset: 0, width: 1,
		},
		{
			inReg: ref, inData: 0, outReg: ref, outData: 0x11, offset: 1, width: 1,
		},
		{
			inReg: ref, inData: 0, outReg: ref, outData: 0x77, offset: 7, width: 1,
		},

		{
			inReg: 0, inData: ref, outReg: ref, outData: ref, offset: 0, width: 8,
			isWrite: true,
		},
		{
			inReg: -1, inData: ref, outReg: -1<<32 | 0x33221100, outData: ref, offset: 0, width: 4,
			isWrite: true,
		},
		{
			inReg: -1, inData: ref, outReg: 0x33221100FFFFFFFF, outData: ref, offset: 4, width: 4,
			isWrite: true,
		},
		{
			inReg: -1, inData: ref, outReg: -1<<16 | 0x1100, outData: ref, offset: 0, width: 2,
			isWrite: true,
		},
		{
			inReg: -1, inData: ref, outReg: -1<<8 | 0x00, outData: ref, offset: 0, width: 1,
			isWrite: true,
		},
		{
			inReg: -1, inData: ref, outReg: -1<<16 | 0x00FF, outData: ref, offset: 1, width: 1,
			isWrite: true,
		},
		{
			inReg: -1, inData: ref, outReg: 0x00FFFFFF_FFFFFFFF, outData: ref, offset: 7, width: 1,
			isWrite: true,
		},
	}

	for _, test := range tests {
		reg, data := test.inReg, test.inData

		Access(&reg, &data, test.offset, test.width, test.isWrite)

		if reg != test.outReg || data != test.outData {
			t.Fatalf("Access(%x,%x,%v,%v,%v):\nwant reg=%x, data=%x\n got reg=%x, data=%x",
				test.inReg, test.inData, test.offset, test.width, test.isWrite, test.outReg, test.outData, reg, data)
		}
	}
}
