package vm

import . "luago/api"

// iABC模式
//
// 用于给连续n个局部变量设置初始值（即nil）
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

// iABC模式
//
// 给单个寄存器设置布尔值，A为要赋值的寄存器索引，B表示bool值，C非0表示跳过下条指令
func loadBool(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}
}

// iABx模式
//
// 将常量表中的某个常量加载到指定寄存器中，寄存器索引由A指定，常量表索引由Bx指定
func loadK(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.GetConst(bx)
	vm.Replace(a)
}

// iABx模式
//
// 需要和EXTRAARG指令(iAx模式)搭配使用，用后者的Ax操作数指定常量索引 Ax有26bit
func loadKx(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a += 1
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(a)
}
