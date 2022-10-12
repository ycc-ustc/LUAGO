package state

import . "luago/api"

type luaValue interface{}
type luaType int

func typeOf(val luaValue) luaType {
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
