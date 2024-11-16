package interpreter

import (
	"fmt"
	"runtime/debug"

	lx "github.com/devnazir/gosh-script/internal/lexer"
	"github.com/devnazir/gosh-script/pkg/ast"
	"github.com/devnazir/gosh-script/pkg/oops"
	"github.com/devnazir/gosh-script/pkg/semantics"
)

func NewInterpreter() *Interpreter {
	symbolTable := semantics.NewSymbolTable()
	scopeResolver := semantics.NewScopeResolver(symbolTable)

	return &Interpreter{
		symbolTable:   symbolTable,
		scopeResolver: scopeResolver,
	}
}

func (i *Interpreter) Interpret(p ast.ASTNode) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			debug.PrintStack()
		}

		// utils.PrintJson(i.symbolTable.Scopes)
	}()

	program := p.(*ast.Program)
	entryPoint := program.EntryPoint

	for _, nodeItem := range program.Body {
		i.InterpretNode(nodeItem, entryPoint)
	}
}

func (i *Interpreter) InterpretNode(nodeItem ast.ASTNode, entryPoint string) {
	switch (nodeItem).(type) {
	case ast.VariableDeclaration:
		i.InterpretVariableDeclaration((nodeItem).(ast.VariableDeclaration))

	case ast.ShellExpression:
		i.InterpretShellExpression(InterpretShellExpression{
			expression:    (nodeItem).(ast.ShellExpression),
			captureOutput: false,
		})

	case ast.SubShell:
		res := i.InterpretSubShell((nodeItem).(ast.SubShell).Arguments.(string))
		fmt.Printf("%v", res)

	case ast.AssignmentExpression:
		i.InterpretAssigmentExpression((nodeItem).(ast.AssignmentExpression))

	case ast.SourceDeclaration:
		i.InterpretSourceDeclaration((nodeItem).(ast.SourceDeclaration).Sources, entryPoint)

	case ast.FunctionDeclaration:
		name := (nodeItem).(ast.FunctionDeclaration).Identifier.Name

		if name == "init" && len((nodeItem).(ast.FunctionDeclaration).Parameters) > 0 {
			oops.InitFunctionCannotHaveParametersError((nodeItem).(ast.FunctionDeclaration))
		}

		i.symbolTable.Insert((nodeItem).(ast.FunctionDeclaration).Identifier.Name, semantics.SymbolInfo{
			Kind:  lx.KeywordFunc,
			Value: (nodeItem).(ast.FunctionDeclaration),
			Line:  (nodeItem).(ast.FunctionDeclaration).Line,
		})

	case ast.CallExpression:
		info := i.scopeResolver.ResolveScope((nodeItem).(ast.CallExpression).Callee.(ast.Identifier).Name)
		arguments := (nodeItem).(ast.CallExpression).Arguments
		i.InterpretBodyFunction(info.Value.(ast.FunctionDeclaration), arguments)
	}

	// utils.PrintJson(i.symbolTable.Scopes)
}
