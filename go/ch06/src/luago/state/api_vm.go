package state

func (state *luaState) PC() int {
	return state.pc
}

func (state *luaState) AddPC(n int) {
	state.pc += n
}

func (state *luaState) Fetch() uint32 {
	i := state.proto.Code[state.pc]
	state.pc++
	return i
}

// 获取常量表中指定索引处常量
func (state *luaState) GetConst(idx int) {
	c := state.proto.Constants[idx]
	state.stack.push(c)
}

// 将常量值或寄存器中的值推入栈中
func (state *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant 一共占9位 最高位是1代表常量去掉最高位为实际值 否则表示寄存器索引值
		state.GetConst(rk & 0xFF)
	} else { // register lua指令操作数寄存器索引是从0开始的 而lua api里栈是1开始索引 所以要 + 1
		state.PushValue(rk + 1)
	}
}
