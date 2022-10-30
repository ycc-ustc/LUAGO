package state

type luaStack struct {
	slots []luaValue
	top   int
	prev *luaStack
	closure *closure
	varargs []luaValue
	pc int
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

func (stack *luaStack) check(n int) {
	free := len(stack.slots) - stack.top
	for i := free; i < n; i++ {
		stack.slots = append(stack.slots, nil)
	}
}

func (stack *luaStack) push(val luaValue) {
	if stack.top == len(stack.slots) {
		panic(`stack overflow`)
	}
	stack.slots[stack.top] = val
	stack.top++
}

func (stack *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}
	for i := 0; i < n; i++ {
		if i < nVals {
			stack.push(vals[i])
		} else {
			stack.push(nil)
		}
	}
}
func (stack *luaStack) pop() luaValue {
	if stack.top < 1 {
		panic(`stack underflow`)
	}
	stack.top--
	val := stack.slots[stack.top]
	stack.slots[stack.top] = nil
	return val
}

func (stack *luaStack) popN(n int) []luaValue {
	vals :=  make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		vals[i] = stack.pop()
	}
	return vals
}
func (stack *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return idx + stack.top + 1
}

func (stack *luaStack) isValid(idx int) bool {
	absIdx := stack.absIndex(idx)
	return absIdx > 0 && absIdx <= stack.top
}

func (stack *luaStack) get(idx int) luaValue {
	absIdx := stack.absIndex(idx)
	if stack.isValid(absIdx) {
		return stack.slots[absIdx-1]
	}
	return nil
}

func (stack *luaStack) set(idx int, val luaValue) {
	absIdx := stack.absIndex(idx)
	if stack.isValid(absIdx) {
		stack.slots[absIdx-1] = val
		return
	}
	panic("invalid index")
}

func (stack *luaStack) Reverse (from , to int) {
	slots := stack.slots;
	for from < to {
		slots[from], slots[to] = slots[to], slots[from];
		from++;
		to--;
	}
}