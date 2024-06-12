package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65535}, []byte{byte(OpConstant), 255, 255}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpMul, []int{}, []byte{byte(OpMul)}},
		{OpDiv, []int{}, []byte{byte(OpDiv)}},
		{OpTrue, []int{}, []byte{byte(OpTrue)}},
		{OpFalse, []int{}, []byte{byte(OpFalse)}},
	  {OpGetLocal, []int{255},[]byte{byte(OpGetLocal), 255}},
    {OpSetLocal, []int{255}, []byte{byte(OpSetLocal), 255}},
  }

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want %d got %d\n", tt.expected, len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want %d got %d", i, b, instruction[i])
			}
		}
	}

}

func TestInstructionString(t *testing.T) {
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpAdd),
		Make(OpPop),
		Make(OpSub),
		Make(OpDiv),
		Make(OpMul),
    Make(OpGetLocal, 1),
	}

	expected := "0000 OpConstant 1\n0003 OpConstant 2\n0006 OpConstant 65535\n0009 OpAdd\n0010 OpPop\n0011 OpSub\n0012 OpDiv\n0013 OpMul\n0014 OpGetLocal 1\n"
	concatted := Instructions{}

	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formatted \nwant %q \ngot %q", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
    {OpSetLocal, []int{255}, 1},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])

		if n != tt.bytesRead {
			t.Fatalf("n wrong. want %d got %d", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want %d got %d", want, operandsRead[i])
			}
		}
	}
}
