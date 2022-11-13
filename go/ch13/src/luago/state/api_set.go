package state

import . "luago/api"

func (state *luaState) SetTable(idx int) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	k := state.stack.pop()
	state.setTable(t, k, v, false)
}

func (state *luaState) RawSet(idx int) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	k := state.stack.pop()
	state.setTable(t, k, v, true)
}


// raw为true表示忽略元方法 如果t是表并且键已经在表中或者需要忽略元方法，或者表中没用__newindex元方法，直接取值即可
func (state *luaState) setTable(t, k, v luaValue, raw bool) { 
	if tb1, ok := t.(*luaTable); ok {
		if raw || tb1.get(k) != nil || !tb1.hasMetaField("__newindex"){
			tb1.put(k, v)
			return
		}
	}
	if !raw {
		if mf := getMetafield(t, "__newindex", state); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				state.setTable(x, k, v, false)
				return
			case *closure:
				state.stack.push(mf)
				state.stack.push(t)
				state.stack.push(k)
				state.stack.push(v)
				state.Call(3, 0)
				return
			}
		}
	}
	panic("index error")
}

func (state *luaState) SetField(idx int, k string) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, k, v, false)
}

func (state *luaState) SetI(idx int, i int64) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, i, v, false)
}

func (state *luaState) RawSetI(idx int, i int64) {
	t := state.stack.get(idx)
	v := state.stack.pop()
	state.setTable(t, i, v, true)
}

func (state *luaState) SetGlobal(name string) {
	t := state.registry.get(LUA_RIDX_GLOBALS)
	v := state.stack.pop()
	state.setTable(t, name, v, false)
}

func (state *luaState) SetMetatable (idx int) {
	val := state.stack.get(idx)
	mtVal := state.stack.pop()

	if mtVal == nil {
		setMetatable(val, nil, state)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, state)
	} else {
		panic("table expected")
	}
}

// 注册go函数
func (state *luaState) Register(name string, f GoFunction) {
	state.PushGoFunction(f)
	state.SetGlobal(name)
}