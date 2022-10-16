package state

import (
	. "luago/api"
	"luago/number"
	"math"
)

type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc func(float64, float64) float64
}


var (
	iadd = func(a, b int64) int64 { return a + b }
	fadd = func(a, b float64) float64 { return a + b }
	isub = func(a, b int64) int64 { return a - b }
	fsub = func(a, b float64) float64 { return a - b }
	imul = func(a, b int64) int64 { return a * b }
	fmul = func(a, b float64) float64 { return a * b }
	imod = number.IMod
	fmod = number.FMod
	pow = math.Pow
	div = func(a, b float64) float64 { return a/ b }
	iidiv = number.IFloorDiv
	fidiv = number.FFloorDiv
	band = func(a, b int64) int64 { return a & b }
	bor = func(a, b int64) int64 { return a | b }
	bxor = func(a, b int64) int64 { return a ^ b }
	shl = number.ShfitLeft
	shr = number.ShfitRight
	iunm = func(a, _ int64) int64 { return -a }
	funm = func(a, _ float64) float64 { return -a }
	bnot = func(a, _ int64) int64 { return ^a }
)

var operators = []operator{
	operator{iadd, fadd},
	operator{isub, fsub},
	operator{imul, fmul},
	operator{imod, fmod},
	operator{nil, pow},
	operator{nil, div},
	operator{iidiv, fidiv},
	operator{band, nil},
	operator{bor, nil},
	operator{bxor, nil},
	operator{shl, nil},
	operator{shr, nil},
	operator{iunm, funm},
	operator{bnot, nil},
}

func (state *luaState) Arith(op ArithOp) {
	var a, b luaValue
	b = state.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT{
		a = state.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		state.stack.push(result)
	} else {
		panic("arithmetic error!")
	}
}

func _arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil {  //bitwise
		if x, ok := ConvertToInteger(a); ok {
			if y, ok := ConvertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else {
		if op.integerFunc != nil {
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
			if x, ok := ConvertToFloat(a); ok {
				if y, ok := ConvertToFloat(b); ok {
					return op.floatFunc(x, y)
				}
			}
		}
	}	
	return nil
}