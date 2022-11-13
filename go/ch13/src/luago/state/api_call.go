package state

import (
	. "luago/api"
	"luago/binchunk"
	"luago/vm"
)

// chunk参数给出要加载的字节数组，chunkName指定chunk的名字，mode指定加载模式，'b': 第一个参数必须是文本chunk数据，否则会
// 出错， 't': 第一个参数必须是文本chunk数据，'bt': 第一个参数可以是二进制或者文本chunk数据
//
// 如果加载失败会再栈顶留下一个错误信息，返回一个状态码0表示成功 先忽略状态码 直接返回0
func (state *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	state.stack.push(c)
	if len(proto.Upvalues) > 0 { // 设置_ENV
		env := state.registry.get(LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	return 0
}

func (state *luaState) Call(nArgs, nResults int) {
	val := state.stack.get(-(nArgs + 1))
	c, ok := val.(*closure)
	if !ok {
		if mf := getMetafield(val, "__call", state); mf != nil {
			if c, ok = mf.(*closure); ok {
				state.stack.push(val)
				state.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}
	if ok {
		if c.proto != nil {
			state.callLuaClosure(nArgs, nResults, c)
		} else {
			state.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not a function")
	}
}

func (state *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	newStack := newLuaStack(nRegs+LUA_MINSTACK, state)
	newStack.closure = c

	funcAndArgs := state.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	// 如果传入的参数多于固定参数 则需要把vararg参数记下来
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// 使新的调用帧成为当前帧 然后调用runLuaClosure执行，执行结束抛出调用帧
	state.pushLuaStack(newStack)
	state.runLuaClosure()
	state.popLuaStack()

	// 将返回值从栈顶弹出
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		state.stack.check(len(results))
		state.stack.pushN(results, nResults)
	}
}

func (state *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	newStack := newLuaStack(nArgs+LUA_MINSTACK, state)
	newStack.closure = c

	args := state.stack.popN(nArgs)
	newStack.pushN(args, nArgs)
	state.stack.pop()

	// 使新的调用帧成为当前帧 然后调用goFunc执行，执行结束抛出调用帧
	state.pushLuaStack(newStack)
	r := c.goFunc(state)
	state.popLuaStack()

	// 将返回值从栈顶弹出
	if nResults != 0 {
		results := newStack.popN(r)
		state.stack.check(len(results))
		state.stack.pushN(results, nResults)
	}
}

func (state *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(state.Fetch())
		inst.Execute(state)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

func (state *luaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := state.stack
	status = LUA_ERRUN

	// catch error
	defer func ()  {
		if err := recover(); err != nil{
			// 从调用栈不断弹出调用帧，直到到达发起调用的调用帧
			for state.stack != caller {
				state.popLuaStack()
			}
			state.stack.push(err)
		}
	} ()

	state.Call(nArgs, nResults)
	status = LUA_OK
	return
}