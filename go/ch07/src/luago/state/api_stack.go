package state

func (state *luaState) GetTop() int {
	return state.stack.top
}

func (state *luaState) AbsIndex(idx int) int {
	return state.stack.absIndex(idx)
}

// 保证栈的容量不小于n 如果小于n就用nil填满
func (state *luaState) CheckStack(n int) bool {
	state.stack.check(n)
	return true
}

func (state *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		state.stack.pop()
	}
}

func (state *luaState) Copy(fromIdx, ToIdx int) {
	val := state.stack.get(fromIdx)
	state.stack.set(ToIdx, val)
}

// 把指定索引处的值推到栈顶
func (state *luaState) PushValue(idx int) {
	val := state.stack.get(idx)
	state.stack.push(val)
}

// 将指定索引处的值置换为栈顶元素
func (state *luaState) Replace(idx int) {
	val := state.stack.pop()
	state.stack.set(idx, val)
}

func (state *luaState) Insert(idx int) {
	state.Rotate(idx, 1)
}

func (state *luaState) Remove(idx int) {
	state.Rotate(idx, -1)
	state.Pop(1)
}

func (state *luaState) Rotate(idx, n int) {
	t := state.stack.top - 1
	p := state.stack.absIndex(idx) - 1
	var m int
	if n >= 0 { //n的符号代表旋转的方向，>0向栈顶 < 0向栈底
		m = t - n
	} else {
		m = p - n - 1
	}
	state.stack.Reverse(p, m)
	state.stack.Reverse(m+1, t)
	state.stack.Reverse(p, t)
}

func (state *luaState) SetTop(idx int) {
	newTop := state.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow")
	}

	n := state.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			state.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			state.stack.push(nil)
		}
	}
}
