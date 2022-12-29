package state

import . "luago/api"

func (state *luaState) NewThread() LuaState {
	t := &luaState{registry: state.registry}
	t.pushLuaStack(newLuaStack(LUA_MINSTACK, t))
	state.stack.push(t)
	return t
}

func (state *luaState) Resume(from LuaState, nArgs int) int {
	lsFrom := from.(*luaState)
	if lsFrom.coChan == nil {
		lsFrom.coChan = make(chan int)
	}

	if state.coChan == nil { // start coroutine
		state.coChan = make(chan int)
		state.coCaller = lsFrom
		go func() {
			state.coStatus = state.PCall(nArgs, -1, 0)
			lsFrom.coChan <- 1
		}()
	} else { // resume coroutine
		state.coStatus = LUA_OK
		state.coChan <- 1
	}
	<-lsFrom.coChan // wait coroutine to finish or yield
	return state.coStatus
}

func (state *luaState) Yield(nResult int) int {
	state.coStatus = LUA_YIELD
	state.coCaller.coChan <- 1
	<-state.coChan
	return state.GetTop()
}

func (state *luaState) Status() int {
	return state.coStatus
}

func (state *luaState) GetStack() bool {
	return state.stack.prev != nil
}
