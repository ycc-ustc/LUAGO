package vm

const MAXARG_Bx = 1<<18 - 1       // 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 131071

/*
 31       22       13       5    0
  +-------+^------+-^-----+-^-----
  |b=9bits |c=9bits |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    bx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |   sbx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    ax=26bits            |op=6|
  +-------+^------+-^-----+-^-----
 31      23      15       7      0
*/

type Instruction uint32

func (instruction Instruction) Opcode() int{
	return int(instruction & 0x3F);
}

func (instruction Instruction) ABC() (a, b, c int){
	a = int(instruction >> 6 & 0xFF)
	b = int(instruction >> 14 & 0x1FF)
	c = int(instruction >> 23 & 0x1FF)
	return
}

func (instruction Instruction) ABx() (a, bx int) {
	a = int(instruction >> 6 & 0xFF)
	bx = int(instruction >> 14)
	return
}

func (instruction Instruction) AsBx() (a, sbx int) {
	a, bx := instruction.ABx()
	return a, bx - MAXARG_sBx
}

func (instruction Instruction) Ax() int {
	return int(instruction >> 6)
}

func (instruction Instruction) OpName() string {
	return opcodes[instruction.Opcode()].name
}

func (instruction Instruction) OpMode() byte {
	return opcodes[instruction.Opcode()].opMode
}

func (instruction Instruction) BMode() byte {
	return opcodes[instruction.Opcode()].argBMode
}

func (instruction Instruction) CMode() byte {
	return opcodes[instruction.Opcode()].argCMode
}