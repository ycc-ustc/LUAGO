package state

import (
	"luago/number"
	"math"
)

type luaTable struct {
	arr  []luaValue
	_map map[luaValue]luaValue
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 { // 表可能是当作数组使用的 先创建数组
		t.arr = make([]luaValue, nArr)
	}
	if nRec > 0 { // 创建map
		t._map = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (table *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(table.arr)) {
			return table.arr[idx-1]
		}
	}
	return table._map[key]
}

func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (table *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN")
	}
	key = _floatToInteger(key)

	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(table.arr))
		if idx <= arrLen {
			table.arr[idx-1] = val
			if idx == arrLen && val == nil {
				table._shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(table._map, key) // key可能已经存在了map中
			if val != nil {
				table.arr = append(table.arr, val)
				table._expandArray()
			}
			return
		}
	}

	if val != nil {
		if table._map == nil {
			table._map = make(map[luaValue]luaValue, 8) // 默认容量为8
		}
		table._map[key] = val
	} else {
		delete(table._map, key)
	}
}


func (table *luaTable) _shrinkArray() {
	for i := len(table.arr) - 1; i >= 0; i-- {
		if table.arr[i] == nil {
			table.arr = table.arr[0:i]
		}
	}
}

func (table *luaTable) _expandArray() {
	for idx := int64(len(table.arr)) + 1; true; idx++ {
		if val, found := table._map[idx]; found {
			delete(table._map, idx)
			table.arr[idx] = val
		} else {
			break
		}
	}
}

func (table *luaTable) len() int {
	return len(table.arr)
}

