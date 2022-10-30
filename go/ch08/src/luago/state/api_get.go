package state

import (
	. "luago/api"
)

func (state *luaState) CreateTable(nArr, nRec int){
	t := newLuaTable(nArr, nRec)
	state.stack.push(t)
}

func (state *luaState) NewTable() {
	state.CreateTable(0, 0)
}

// 方法根据key(从栈顶弹出)从表中(索引由idx取得)取值，然后把值推入栈顶并返回值的类型
func (state *luaState) GetTable(idx int) LuaType {
	t := state.stack.get(idx)
	k := state.stack.pop()
	return state.getTable(t, k)
}

func (state *luaState) getTable(t, k luaValue) LuaType {
	if tb1, ok := t.(*luaTable); ok {
		v := tb1.get(k)
		state.stack.push(v)
		return typeOf(v)
	}
	panic("not a table")
}

// 与GetTable类似，不过key直接给出
func (state *luaState) GetField(idx int, k string) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, k)
}

func (state *luaState) GetI(idx int, k int64) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, k)
}