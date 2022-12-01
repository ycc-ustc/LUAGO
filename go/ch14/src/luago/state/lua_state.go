package state

import . "luago/api"

type luaState struct {
	registry *luaTable // 注册表 用户使用 全局变量也通过它实现
	stack    *luaStack
}

func New() *luaState {
	registry := newLuaTable(0, 0)
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
}

func (state *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = state.stack
	state.stack = stack
}

func (state *luaState) popLuaStack() {
	stack := state.stack
	state.stack = stack.prev
	stack.prev = nil
}
