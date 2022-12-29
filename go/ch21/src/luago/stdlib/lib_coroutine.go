package stdlib

import . "luago/api"

var coFunc = map[string]GoFunction{
	"create":      coCreate,    // coroutine.create
	"resume":      coResume,    // coroutine.resume
	"yield":       coYield,     // coroutine.yield
	"status":      coStatus,    // coroutine.status
	"isyieldable": coYieldable, // coroutine.isyieldable
	"running":     coRunning,   // coroutine.running
	"wrap":        coWrap,      //coroutine.wrap
}

func coCreate(ls LuaState) int {
	ls.CheckType(1, LUA_TFUNCTION)
	ls2 := ls.NewThread()
	ls.PushValue(1)  /* move function to top */
	ls.XMove(ls2, 1) /* move function from ls to ls2 */
	return 1
}

func coStatus(ls LuaState) int {
	co := ls.ToThread(1)
	ls.ArgCheck(co != nil, 1, "thread expected")
	if ls == co {
		ls.PushString("running")
	} else {
		switch co.Status() {
		case LUA_YIELD:
			ls.PushString("suspended")
		case LUA_OK:
			if co.GetStack() {
				ls.PushString("normal")
			} else if co.GetTop() == 0 {
				ls.PushString("dead")
			} else {
				ls.PushString("suspended")
			}
		default:
			ls.PushString("dead")
		}

	}
	return 1
}

func coYield(ls LuaState) int {
	return ls.Yield(ls.GetTop())
}

func coResume(ls LuaState) int {
	co := ls.ToThread(1)
	ls.ArgCheck(co != nil, 1, "thread expected")

	if r := _auxResume(ls, co, ls.GetTop()-1); r < 0 {
		ls.PushBoolean(false)
		ls.Insert(-2)
		return 2 /* return false + error message */
	} else {
		ls.PushBoolean(true)
		ls.Insert(-(r + 1))
		return r + 1 /* return true + 'resume' returns*/
	}
}

func _auxResume(ls LuaState, co LuaState, nArgs int) int {
	if !ls.CheckStack(nArgs) {
		ls.PushString("too many arguments to resume")
		return -1 /* error flag */
	}
	if co.Status() == LUA_OK && co.GetTop() == 0 {
		ls.PushString("cannot resume dead coroutine")
		return -1
	}
	ls.XMove(co, nArgs)
	status := co.Resume(ls, nArgs)
	if status == LUA_OK || status == LUA_YIELD {
		nRes := co.GetTop()
		if !ls.CheckStack(nRes + 1) {
			co.Pop(nRes) /* remove results anyway */
			ls.PushString("too many results to resume")
			return -1 /* error flag */
		}
		co.XMove(ls, nRes) /* move yielded values */
		return nRes
	} else {
		co.XMove(ls, 1) /* move error message */
		return -1
	}
}
