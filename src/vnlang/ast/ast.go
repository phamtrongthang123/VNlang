package ast

import (
	"bytes"
	"math/big"
	"strings"
	"vnlang/scanner"
	"vnlang/token"
)

// The base Node interface
type Node interface {
	Position() scanner.Position
	TokenLiteral() string
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Position() scanner.Position {
	if len(p.Statements) > 0 {
		return p.Statements[0].Position()
	} else {
		return scanner.Position{}
	}
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Statements
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()             {}
func (ls *LetStatement) Position() scanner.Position { return ls.Token.Pos }
func (ls *LetStatement) TokenLiteral() string       { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()             {}
func (rs *ReturnStatement) Position() scanner.Position { return rs.Token.Pos }
func (rs *ReturnStatement) TokenLiteral() string       { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type BreakStatement struct {
	Token token.Token // the 'ngắt' token
}

func (bs *BreakStatement) statementNode()             {}
func (bs *BreakStatement) Position() scanner.Position { return bs.Token.Pos }
func (bs *BreakStatement) TokenLiteral() string       { return bs.Token.Literal }
func (bs *BreakStatement) String() string             { return bs.Token.Literal }

type ContinueStatement struct {
	Token token.Token // the 'tiếp' token
}

func (cs *ContinueStatement) statementNode()             {}
func (cs *ContinueStatement) Position() scanner.Position { return cs.Token.Pos }
func (cs *ContinueStatement) TokenLiteral() string       { return cs.Token.Literal }
func (cs *ContinueStatement) String() string             { return cs.Token.Literal }

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()             {}
func (es *ExpressionStatement) Position() scanner.Position { return es.Token.Pos }
func (es *ExpressionStatement) TokenLiteral() string       { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()             {}
func (bs *BlockStatement) Position() scanner.Position { return bs.Token.Pos }
func (bs *BlockStatement) TokenLiteral() string       { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String() + " ")
	}

	return out.String()
}

// Expressions
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()            {}
func (i *Identifier) Position() scanner.Position { return i.Token.Pos }
func (i *Identifier) TokenLiteral() string       { return i.Token.Literal }
func (i *Identifier) String() string             { return i.Value }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()            {}
func (b *Boolean) Position() scanner.Position { return b.Token.Pos }
func (b *Boolean) TokenLiteral() string       { return b.Token.Literal }
func (b *Boolean) String() string             { return b.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value *big.Int
}

func (il *IntegerLiteral) expressionNode()            {}
func (il *IntegerLiteral) Position() scanner.Position { return il.Token.Pos }
func (il *IntegerLiteral) TokenLiteral() string       { return il.Token.Literal }
func (il *IntegerLiteral) String() string             { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()            {}
func (fl *FloatLiteral) Position() scanner.Position { return fl.Token.Pos }
func (fl *FloatLiteral) TokenLiteral() string       { return fl.Token.Literal }
func (fl *FloatLiteral) String() string             { return fl.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()            {}
func (pe *PrefixExpression) Position() scanner.Position { return pe.Token.Pos }
func (pe *PrefixExpression) TokenLiteral() string       { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()            {}
func (oe *InfixExpression) Position() scanner.Position { return oe.Token.Pos }
func (oe *InfixExpression) TokenLiteral() string       { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   []Expression
	Consequence []*BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()            {}
func (ie *IfExpression) Position() scanner.Position { return ie.Token.Pos }
func (ie *IfExpression) TokenLiteral() string       { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("nếu ")
	out.WriteString(ie.Condition[0].String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence[0].String())

	for i := range ie.Condition[1:] {
		out.WriteString("còn_nếu ")
		out.WriteString(ie.Condition[i].String())
		out.WriteString(" ")
		out.WriteString(ie.Consequence[i].String())
	}

	if ie.Alternative != nil {
		out.WriteString("ngược_lại ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type WhileExpression struct {
	Token     token.Token // The 'khi' token
	Condition Expression
	Body      *BlockStatement
}

func (ie *WhileExpression) expressionNode()            {}
func (ie *WhileExpression) Position() scanner.Position { return ie.Token.Pos }
func (ie *WhileExpression) TokenLiteral() string       { return ie.Token.Literal }
func (ie *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("khi ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Body.String())

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()            {}
func (fl *FunctionLiteral) Position() scanner.Position { return fl.Token.Pos }
func (fl *FunctionLiteral) TokenLiteral() string       { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()            {}
func (ce *CallExpression) Position() scanner.Position { return ce.Token.Pos }
func (ce *CallExpression) TokenLiteral() string       { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()            {}
func (sl *StringLiteral) Position() scanner.Position { return sl.Token.Pos }
func (sl *StringLiteral) TokenLiteral() string       { return "\"" + sl.Token.Literal + "\"" }
func (sl *StringLiteral) String() string             { return "\"" + sl.Token.Literal + "\"" }

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()            {}
func (al *ArrayLiteral) Position() scanner.Position { return al.Token.Pos }
func (al *ArrayLiteral) TokenLiteral() string       { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()            {}
func (ie *IndexExpression) Position() scanner.Position { return ie.Token.Pos }
func (ie *IndexExpression) TokenLiteral() string       { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()            {}
func (hl *HashLiteral) Position() scanner.Position { return hl.Token.Pos }
func (hl *HashLiteral) TokenLiteral() string       { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
