package ast

type Stat interface{}
type EmptyStat struct{} // ';'
type BreakStat struct { // break
	Line int
}

type LabelStat struct { // '::' Name '::'
	Name string
}

type GotoStat struct { // goto Name
	Name string
}

type DoStat struct { // do block end
	Block *Block
}

type FuncCallStat = FuncCallExp // function call

type WhileStat struct {
	Exp   Exp
	Block *Block
}

type RepeatStat struct {
	Exp   Exp
	Block *Block
}

type IfStat struct {
	Exps   []Exp
	Blocks []*Block
}

type ForNumStat struct {
	LineOfFor int
	LineOfDo  int
	VarName   string
	InitExp   Exp
	LimitExp  Exp
	Step      Exp
	Block     *Block
}

type ForInStat struct {
	LineOfDo int
	NameList []string
	ExpList  []Exp
	Block    *Block
}

type LocalVarDeclStat struct {
	LastLine int
	NameList []string
	ExpList  []Exp
}

type AssignStat struct {
	LastLine int
	VarList  []string
	ExpList  []Exp
}

type LocalFuncDefStat struct {
	Name string
	Exp  *FuncDefExp
}
