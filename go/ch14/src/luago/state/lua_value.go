package state

import (
	"fmt"
	. "luago/api"
	"luago/number"
)

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
	case *luaTable:
		return LUA_TTABLE
	case *closure:
		return LUA_TFUNCTION
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

func ConvertToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}

func ConvertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		return _stringToInteger(x)
	default:
		return 0, false
	}
}

func _stringToInteger(s string) (int64, bool) {
	if i, ok := number.ParseInteger(s); ok {
		return i, ok
	}
	if f, ok := number.ParseFloat(s); ok {
		return number.FloatToInteger(f)
	}
	return 0, false
}

func setMetatable(val luaValue, mt *luaTable, ls *luaState) {
	if t, ok := val.(*luaTable); ok {
		t.metatable = mt
		return
	}
	// 按照约定下划线开头后跟大写字母的字段名是保留给lua实现使用的，所以我们使用_MT1这种字段名，
	// 以免与用户放在注册表中的数据冲突
	key := fmt.Sprintf("_MT%d", typeOf(val))
	ls.registry.put(key, mt)
}

// 先判断值是否是表，是的话直接返回metatable字段，否则要从注册表中取出来
func getMetatable(val luaValue, ls *luaState) *luaTable {
	if t, ok := val.(*luaTable); ok {
		return t.metatable
	}
	key := fmt.Sprintf("_MT%d", typeOf(val))
	if mt := ls.registry.get(key); mt != nil {
		return mt.(*luaTable)
	}
	return nil
}
