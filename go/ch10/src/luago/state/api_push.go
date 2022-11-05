package state

import . "luago/api"

func (state *luaState) PushNil() {
	state.stack.push(nil)
}

func (state *luaState) PushBoolean(b bool) {
	state.stack.push(b)
}

func (state *luaState) PushInteger(n int64) {
	state.stack.push(n)
}

func (state *luaState) PushNumber(n float64) {
	state.stack.push(n)
}

func (state *luaState) PushString(s string) {
	state.stack.push(s)
}

func (state *luaState) PushGoFunction(f GoFunction) {
	state.stack.push(newGoClosure(f, 0))
}

func (state *luaState) PushGlobalTable() {
	global := state.registry.get(LUA_RIDX_GLOBALS)
	state.stack.push(global)
}

func (state *luaState) GetGlobal(name string) LuaType {
	t := state.registry.get(LUA_RIDX_GLOBALS)
	return state.getTable(t, name)
}


// 先创建一个go闭包，然后讲upval放进去 最后推入栈中
func (state *luaState) PushGoClosure(f GoFunction, n int) {
	closure := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := state.stack.pop()
		closure.upvals[n-1] = &upvalue{&val}
	}
	state.stack.push(closure)
}