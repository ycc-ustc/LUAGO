package api

type FuncReg map[string]GoFunction

// auxiliary library
type AuxLib interface {
	/* Error-report functions */
	Error2(fmt string, a ...interface{}) int
	CheckType(arg int, t LuaType)
	ArgError(arg int, extraMsg string) int
	TypeName2(idx int) string
	ArgCheck(cond bool, arg int, extraMsg string)
}
