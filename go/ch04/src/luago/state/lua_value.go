package state

import . "luago/api"

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMBER
	case float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	default:
		panic("todo")
	}
}

func ConvertToBoolean(val luaValue) bool {
	switch b := val.(type) {
	case nil:
		return false
	case bool:
		return b
	default:
		return true
	}
}