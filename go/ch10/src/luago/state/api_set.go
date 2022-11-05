package state

import . "luago/api"

func (state *luaState) SetTable(idx int) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	k := state.stack.pop()
	state.setTable(t, k, v)
}

func (state *luaState) setTable(t, k, v luaValue) {
	if tb1, ok := t.(*luaTable); ok {
		tb1.put(k, v)
		return
	}
	panic("not a table")
}

func (state *luaState) SetField(idx int, k string) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, k, v)
}

func (state *luaState) SetI(idx int, i int64) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, i, v)
}

func (state *luaState) SetGlobal(name string) {
	t := state.registry.get(LUA_RIDX_GLOBALS)
	v := state.stack.pop()
	state.setTable(t, name, v)
}

// 注册go函数
func (state *luaState) Register(name string, f GoFunction) {
	state.PushGoFunction(f)
	state.SetGlobal(name)
}