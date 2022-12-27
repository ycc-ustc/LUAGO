package state

func (state *luaState) PC() int {
	return state.stack.pc
}

func (state *luaState) AddPC(n int) {
	state.stack.pc += n
}

func (state *luaState) Fetch() uint32 {
	i := state.stack.closure.proto.Code[state.stack.pc]
	state.stack.pc++
	return i
}

// 获取常量表中指定索引处常量
func (state *luaState) GetConst(idx int) {
	c := state.stack.closure.proto.Constants[idx]
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

func (state *luaState) RegisterCount() int {
	return int(state.stack.closure.proto.MaxStackSize)
}

func (state *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(state.stack.varargs)
	}
	state.stack.check(n)
	state.stack.pushN(state.stack.varargs, n)
}

func (state *luaState) LoadProto(idx int) {
	proto := state.stack.closure.proto.Protos[idx]
	stack := state.stack
	closure := newLuaClosure(proto)
	stack.push(closure)
	for i, uvInfo := range proto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 { // 捕获的是当前函数的局部变量
			// 如果upValue捕获的外部局部变量还在栈上即open状态，直接引用即可，否则为closed状态
			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}
			if openuv, found := stack.openuvs[uvIdx]; found {
				closure.upvals[i] = openuv
			} else {
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[idx] = closure.upvals[i]
			}
		} else { // 捕获的是更外围函数的局部变量
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}

	}
}

func (state *luaState) CloseUpvalues(a int) {
	for i, openuv := range state.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(state.stack.openuvs, i)
		}
	}
}
