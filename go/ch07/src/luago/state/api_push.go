package state

func (state *luaState) PushNil() {
	state.stack.push(nil)
}

func (state *luaState) PushBoolean(b bool) {
	state.stack.push(b)
}

func (state *luaState) PushInteger(n int64) {
	state.stack.push(n)
}

func (state *luaState) PushNumber(n float64) {
	state.stack.push(n)
}

func (state *luaState) PushString(s string) {
	state.stack.push(s)
}
