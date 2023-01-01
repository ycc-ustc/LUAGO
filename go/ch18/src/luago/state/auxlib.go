package state

import (
	"fmt"
	"io/ioutil"
	. "luago/api"
	"luago/stdlib"
)

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

func (state *luaState) Len2(idx int) int64 {
	state.Len(idx)
	i, isNum := state.ToIntegerX(-1)
	if !isNum {
		state.Error2("object length is not a integer")
	}
	state.Pop(1)
	return i
}

func (state *luaState) CheckStack2(sz int, msg string) {
	if !state.CheckStack(sz) {
		if msg != "" {
			state.Error2("stack overflow (%s)", msg)
		} else {
			state.Error2("stack overflow")
		}

	}
}

func (state *luaState) LoadString(s string) int {
	return state.Load([]byte(s), s, "bt")
}

// DoString 加载并使用保护模式执行字符串
func (state *luaState) DoString(str string) bool {
	return state.LoadString(str) == LUA_OK &&
		state.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

// DoFile 加载并使用保护模式执行字符串
func (state *luaState) DoFile(filename string) bool {
	return state.LoadFile(filename) == LUA_OK &&
		state.PCall(0, LUA_MULTRET, 0) == LUA_OK
}

func (state *luaState) LoadFileX(filename, mode string) int {
	if data, err := ioutil.ReadFile(filename); err == nil {
		return state.Load(data, "@"+filename, mode)
	}
	return LUA_ERRFILE
}

func (state *luaState) LoadFile(filename string) int {
	return state.LoadFileX(filename, "bt")
}

// ArgError 用于参数检查方法 检查传递给Go函数的参数 如果非可选参数缺失 或者参数类型和预期不匹配 则使用ArgError抛出错误
func (state *luaState) ArgError(arg int, extraMsg string) int {
	return state.Error2("bad argument #%d (%s)", arg, extraMsg)
}

// ArgCheck 通用的参数检查方法 第一个参数表示检查是否通过 第二个参数表示被检查参数索引 第三个参数是附加信息
func (state *luaState) ArgCheck(cond bool, arg int, extraMsg string) {
	if !cond {
		state.ArgError(arg, extraMsg)
	}
}

// CheckAny 确保某个参数一定存在
func (state *luaState) CheckAny(arg int) {
	if state.Type(arg) == LUA_TNONE {
		state.ArgError(arg, "value expected")
	}
}

// CheckType 确保某一参数属于指定类型
func (state *luaState) CheckType(arg int, t LuaType) {
	if state.Type(arg) != t {
		state.tagError(arg, t)
	}
}

func (state *luaState) tagError(arg int, tag LuaType) {
	state.typeError(arg, state.TypeName(tag))
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

// CheckNumber 检查参数是否是number 是的话返回number值
func (state *luaState) CheckNumber(arg int) float64 {
	f, ok := state.ToNumberX(arg)
	if !ok {
		state.tagError(arg, LUA_TNUMBER)
	}
	return f
}

func (state *luaState) CheckString(arg int) string {
	s, ok := state.ToStringX(arg)
	if !ok {
		state.tagError(arg, LUA_TSTRING)
	}
	return s
}

func (state *luaState) CheckInteger(arg int) int64 {
	i, ok := state.ToIntegerX(arg)
	if !ok {
		state.intError(arg)
	}
	return i
}

func (state *luaState) intError(arg int) {
	if state.IsNumber(arg) {
		state.ArgError(arg, "number has no integer representation")
	} else {
		state.tagError(arg, LUA_TNUMBER)
	}
}

// OptNumber 对可选参数进行检查 如果可选参数有值 进行检查 否则返回默认值
func (state *luaState) OptNumber(arg int, def float64) float64 {
	if state.IsNoneOrNil(arg) {
		return def
	}
	return state.CheckNumber(arg)
}

func (state *luaState) OptString(arg int, def string) string {
	if state.IsNoneOrNil(arg) {
		return def
	}
	return state.CheckString(arg)
}

func (state *luaState) OptInteger(arg int, def int64) int64 {
	if state.IsNoneOrNil(arg) {
		return def
	}
	return state.CheckInteger(arg)
}

func (state *luaState) OpenLibs() {
	libs := map[string]GoFunction{
		"_G": stdlib.OpenBaseLib,
		// TODO
	}
	for name, fun := range libs {
		state.RequireF(name, fun, true)
		state.Pop(1)
	}
}

// RequireF 用于开启单个标准库
func (state *luaState) RequireF(modname string, openf GoFunction, glb bool) {
	state.GetSubTable(LUA_REGISTRYINDEX, "_LOADED")
	state.GetField(-1, modname)
	if !state.ToBoolean(-1) { // package not already loaded?
		state.Pop(1) // remove field
		state.PushGoFunction(openf)
		state.PushString(modname)   // argument to open function
		state.Call(1, 1)            // cannot 'openf' to open module
		state.PushValue(-1)         // make copy of module (call result)
		state.SetField(-3, modname) // _LOADED[modname] = module
	}
	state.Remove(-2) // remove _LOADED table
	if glb {
		state.PushValue(-1)      // copy of module
		state.SetGlobal(modname) // _G[modname] = module
	}
}

func (state *luaState) GetSubTable(idx int, fname string) bool {
	if state.GetField(idx, fname) == LUA_TTABLE {
		return true // table already there
	}
	state.Pop(1) // remove previous result
	idx = state.stack.absIndex(idx)
	state.NewTable()
	state.PushValue(-1)        // copy to be left at top
	state.SetField(idx, fname) // assign new table to field
	return false               // false, because did not find table there
}

func (state *luaState) CallMeta(obj int, event string) bool {
	obj = state.AbsIndex(obj)
	if state.GetMetafield(obj, event) == LUA_TNIL { /* no metafield? */
		return false
	}

	state.PushValue(obj)
	state.Call(1, 1)
	return true
}

func (state *luaState) NewLib(l FuncReg) {
	state.NewLibTable(l)
	state.SetFuncs(l, 0)
}

func (state *luaState) SetFuncs(l FuncReg, nup int) {
	state.CheckStack2(nup, "too many upvalues")
	for name, fun := range l { /* fill the table with given functions */
		for i := 0; i < nup; i++ { /* copy upvalues to the top */
			state.PushValue(-nup)
		}
		// r[-(nup+2)][name]=fun
		state.PushGoClosure(fun, nup) /* closure with those upvalues */
		state.SetField(-(nup + 2), name)
	}
	state.Pop(nup) /* remove upvalues */
}

func (state *luaState) NewLibTable(l FuncReg) {
	state.CreateTable(0, len(l))
}

func (state *luaState) ToString2(idx int) string {
	if state.CallMeta(idx, "__tostring") { /* metafield? */
		if !state.IsString(-1) {
			state.Error2("'__tostring' must return a string")
		}
	} else {
		switch state.Type(idx) {
		case LUA_TNUMBER:
			if state.IsInteger(idx) {
				state.PushString(fmt.Sprintf("%d", state.ToInteger(idx))) // todo
			} else {
				state.PushString(fmt.Sprintf("%g", state.ToNumber(idx))) // todo
			}
		case LUA_TSTRING:
			state.PushValue(idx)
		case LUA_TBOOLEAN:
			if state.ToBoolean(idx) {
				state.PushString("true")
			} else {
				state.PushString("false")
			}
		case LUA_TNIL:
			state.PushString("nil")
		default:
			tt := state.GetMetafield(idx, "__name") /* try name */
			var kind string
			if tt == LUA_TSTRING {
				kind = state.CheckString(-1)
			} else {
				kind = state.TypeName2(idx)
			}

			state.PushString(fmt.Sprintf("%s: %p", kind, state.ToPointer(idx)))
			if tt != LUA_TNIL {
				state.Remove(-2) /* remove '__name' */
			}
		}
	}
	return state.CheckString(-1)
}
