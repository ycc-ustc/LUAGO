package state

import (
	"fmt"
	. "luago/api"
)

func (state *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

func (state *luaState) Type(idx int) LuaType {
	if state.stack.isValid(idx) {
		val := state.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

func (state *luaState) IsNone(idx int) bool {
	return state.Type(idx) == LUA_TNONE
}

func (state *luaState) IsNil(idx int) bool {
	return state.Type(idx) == LUA_TNIL
}

func (state *luaState) IsNoneOrNil(idx int) bool {
	return state.Type(idx) <= LUA_TNIL
}

func (state *luaState) IsBoolean(idx int) bool {
	return state.Type(idx) == LUA_TBOOLEAN
}

func (state *luaState) IsString(idx int) bool {
	t := state.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

func (state *luaState) IsNumber(idx int) bool {
	_, ok := state.ToNumberX(idx)
	return ok
}

func (state *luaState) IsInteger(idx int) bool {
	val := state.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (state *luaState) ToBoolean(idx int) bool {
	val := state.stack.get(idx)
	return ConvertToBoolean(val)
}

func (state *luaState) ToNumberX(idx int) (float64, bool) {
	val := state.stack.get(idx)
	return ConvertToFloat(val)
}

func (state *luaState) ToNumber(idx int) float64 {
	n, _ := state.ToNumberX(idx)
	return n
}

func (state *luaState) ToIntegerX(idx int) (int64, bool) {
	val := state.stack.get(idx)
	return ConvertToInteger(val)
}

func (state *luaState) ToInteger(idx int) int64 {
	i, _ := state.ToIntegerX(idx)
	return i
}

func (state *luaState) ToStringX(idx int) (string, bool) {
	val := state.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		state.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

func (state *luaState) ToString(idx int) string {
	s, _ := state.ToStringX(idx)
	return s
}
