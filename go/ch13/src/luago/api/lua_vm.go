package api

type LuaVM interface {
	LuaState
	PC() int // 返回当前PC 测试用
	AddPC(n int) // 修改PC 用于实现跳转指令
	Fetch() uint32 // 取出当前指令、讲PC指向下一条指令
	GetConst(idx int) // 将指定常量推入栈顶
	GetRK(rk int) // 将指定常量或栈值推入栈顶
	RegisterCount() int // 查询寄存器个数
	LoadVararg(n int) // 加载可变参数
	LoadProto(idx int) // 加载proto文件
	CloseUpvalues(a int) // 闭合Upvalues
}