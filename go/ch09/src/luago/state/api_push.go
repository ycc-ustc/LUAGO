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
	state.stack.push(newGoClosure(f))
}

func (state *luaState) PushGlobalTable() {
	global := state.registry.get(LUA_RIDX_GLOBALS)
	state.stack.push(global)
}

func (state *luaState) GetGlobal(name string) LuaType {
	t := state.registry.get(LUA_RIDX_GLOBALS)
	return state.getTable(t, name)
}