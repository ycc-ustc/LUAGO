package ast

type Exp interface{}

type NilExp struct {
	Line int
}

type TrueExp struct {
	Line int
}

type FalseExp struct {
	Line int
}

type VarargExp struct {
	Line int
}

type IntegerExp struct {
	Line int
	Val  int64
}

type FloatExp struct {
	Line int
	Val  float64
}

type StringExp struct {
	Line int
	Str  string
}

type NameExp struct {
	Line int
	Name string
}

// UnopExp 一元运算符表达式
type UnopExp struct {
	Line int    // line of operator
	Op   string // operator
	Exp  Exp
}

// BinopExp 二元运算符表达式
type BinopExp struct {
	Line int    // line of operator
	Op   string // operator
	Exp1 Exp
	Exp2 Exp
}

// ConcatExp 拼接运算符表达式
type ConcatExp struct {
	Line int // line of last ...
	Exps []Exp
}

// TableConstructorExp 表构造表达式
type TableConstructorExp struct {
	Line     int // line of '{'
	LastLine int // line of '}'
	KeyExps  []Exp
	ValExps  []Exp
}

// FuncDefExp 函数定义表达式
type FuncDefExp struct {
	Line     int // line of '{'
	LastLine int // line of '}'
	ParList  []string
	IsVararg bool
	Block    *Block
}

// ParentsExp 圆括号表达式
type ParentsExp struct {
	Exp Exp
}

// TableAccessExp 表访问表达式
type TableAccessExp struct {
	LastLine  int // line of ']'
	PrefixExp Exp
	KeyExp    Exp
}

// FuncCallExp 函数调用表达式
type FuncCallExp struct {
	Line      int // line of '('
	LastLine  int // line of ')'
	PrefixExp Exp
	NameExp   *StringExp
	Args      []Exp
}
