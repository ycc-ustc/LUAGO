package state

import (
	. "luago/api"
)

func (state *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	state.stack.push(t)
}

func (state *luaState) NewTable() {
	state.CreateTable(0, 0)
}

func (state *luaState) RawGet(idx int) LuaType {
	t := state.stack.get(idx)
	k := state.stack.pop()
	return state.getTable(t, k, true)
}

// GetTable 方法根据key(从栈顶弹出)从表中(索引由idx取得)取值，然后把值推入栈顶并返回值的类型
func (state *luaState) GetTable(idx int) LuaType {
	t := state.stack.get(idx)
	k := state.stack.pop()
	return state.getTable(t, k, false)
}

func (state *luaState) GetGlobal(name string) LuaType {
	t := state.registry.get(LUA_RIDX_GLOBALS)
	return state.getTable(t, name, false)
}

// raw为true表示忽略元方法 如果t是表并且键已经在表中或者需要忽略元方法，或者表中没用__index元方法，直接取值即可
func (state *luaState) getTable(t, k luaValue, raw bool) LuaType {
	if tb1, ok := t.(*luaTable); ok {
		v := tb1.get(k)
		if raw || v != nil || !tb1.hasMetaField("__index") {
			state.stack.push(v)
			return typeOf(v)
		}
	}
	if !raw {
		if mf := getMetafield(t, "__index", state); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return state.getTable(x, k, false)
			case *closure:
				state.stack.push(mf)
				state.stack.push(t)
				state.stack.push(k)
				state.Call(2, 1)
				v := state.stack.get(-1)
				return typeOf(v)
			}
		}
	}
	panic("index error")
}

// 与GetTable类似，不过key直接给出
func (state *luaState) GetField(idx int, k string) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, k, false)
}

func (state *luaState) GetI(idx int, k int64) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, k, false)
}

func (state *luaState) RawGetI(idx int, k int64) LuaType {
	t := state.stack.get(idx)
	return state.getTable(t, k, true)
}

func (state *luaState) GetMetatable(idx int) bool {
	val := state.stack.get(idx)

	if mt := getMetatable(val, state); mt != nil {
		state.stack.push(mt)
		return true
	} else {
		return false
	}
}
