package ast

type AssignmentExpression struct {
	Identifier
	Expression ASTNode
}

type CallExpression struct {
	BaseNode
	Callee    ASTNode
	Arguments []ASTNode
}

type MemberExpression struct {
	BaseNode
	Object   ASTNode
	Property ASTNode
	Computed bool
}