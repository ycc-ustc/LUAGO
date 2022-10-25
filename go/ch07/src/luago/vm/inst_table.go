package vm

import (
	. "luago/api"
)

const LFIELD_PER_FLUSH = 50

func newTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

// iABC模式
//
// 键由C指定(RK) 表由B索引 放入A指定的索引处
func getTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// iABC模式
//
// 键值分别由由B C指定(RK) 表由A索引
func setTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

// iABC模式
//
// C如果大于0则表示实际的批次数+1，否则真正的批次数保存在下一个EXTRAARG指令中
func setList(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}

	idx := int64(c * LFIELD_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		vm.PushValue(a + j)
		vm.SetI(a, idx)
	}
}
