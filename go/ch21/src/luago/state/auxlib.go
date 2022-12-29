package state

import (
	"fmt"
	. "luago/api"
)

func (state *luaState) CheckType(arg int, t LuaType) {
	if state.Type(arg) != t {
		state.tagError(arg, t)
	}
}

func (state *luaState) tagError(arg int, tag LuaType) {
	state.typeError(arg, state.TypeName(LuaType(tag)))
}

func (state *luaState) typeError(arg int, tname string) int {
	var typeArg string /* name for the type of the actual argument */
	if state.GetMetafield(arg, "__name") == LUA_TSTRING {
		typeArg = state.ToString(-1) /* use the given type name */
	} else if state.Type(arg) == LUA_TLIGHTUSERDATA {
		typeArg = "light userdata" /* special name for messages */
	} else {
		typeArg = state.TypeName2(arg) /* standard name */
	}
	msg := tname + " expected, got " + typeArg
	state.PushString(msg)
	return state.ArgError(arg, msg)
}

func (state *luaState) TypeName2(idx int) string {
	return state.TypeName(state.Type(idx))
}

func (state *luaState) PushFString(fmtStr string, a ...interface{}) {
	str := fmt.Sprintf(fmtStr, a...)
	state.stack.push(str)
}

func (state *luaState) Error2(fmt string, a ...interface{}) int {
	state.PushFString(fmt, a...) // todo
	return state.Error()
}

func (state *luaState) ArgError(arg int, extraMsg string) int {
	// bad argument #arg to 'funcname' (extramsg)
	return state.Error2("bad argument #%d (%s)", arg, extraMsg) // todo
}

func (state *luaState) GetMetafield(obj int, event string) LuaType {
	if !state.GetMetatable(obj) { /* no metatable? */
		return LUA_TNIL
	}

	state.PushString(event)
	tt := state.RawGet(-2)
	if tt == LUA_TNIL { /* is metafield nil? */
		state.Pop(2) /* remove metatable and metafield */
	} else {
		state.Remove(-2) /* remove only metatable */
	}
	return tt /* return metafield type */
}

func (state *luaState) ArgCheck(cond bool, arg int, extraMsg string) {
	if !cond {
		state.ArgError(arg, extraMsg)
	}
}
