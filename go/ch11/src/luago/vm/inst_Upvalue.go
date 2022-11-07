package vm

import . "luago/api"


// iABC()模式
//
// 把当前闭包的某个Upvalue值(B伪索引指定)拷贝到目标寄存器(A指定)中，C没用
func getUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(LuaUpvalueIndex(b), a)
}

// iABC()模式
//
// 把当前闭包的某个Upvalue值(B伪索引指定)赋值为目标寄存器(A指定)中的值，C没用
func setUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(a, LuaUpvalueIndex(b))
}

// iABC()模式
//
// 如果upvalue(B指定)是table 可以根据键值(C指定)从该表中取值，放入目标寄存器(A指定)中
func getTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(LuaUpvalueIndex(b))
	vm.Replace(a)
}

// iABC()模式
//
// 如果upvalue(A指定)是table 可以根据键值(B指定)更新值(C指定)
func setTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(LuaUpvalueIndex(a))
}