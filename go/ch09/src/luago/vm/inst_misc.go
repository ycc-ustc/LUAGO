package vm

import . "luago/api"

// iABC模式

// 将源寄存器的值copy到目的寄存器
func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(b, a)
}

// iAsBx模式
//
// pc跳转
func jmp(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		panic("to do")
	}
}
