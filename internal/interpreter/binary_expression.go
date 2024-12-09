package interpreter

import (
	"fmt"
	"reflect"

	"github.com/devnazir/gosh-script/pkg/ast"
	"github.com/devnazir/gosh-script/pkg/oops"
)

func (i *Interpreter) InterpretBinaryExpr(b ast.ASTNode) interface{} {

	switch b.GetType() {
	case ast.StringLiteralTree:
		return b.(ast.StringLiteral).Raw
	case ast.NumberLiteralTree:
		return b.(ast.NumberLiteral).Value
	case ast.BooleanLiteralTree:
		return b.(ast.BooleanLiteral).Value
	case ast.IdentifierTree:
		name := b.(ast.Identifier).Name
		info := i.scopeResolver.ResolveScope(name)

		return info.Value

	case ast.SubShellTree:
		value := i.InterpretSubShell(b.(ast.SubShell).Arguments)
		return value

	case ast.MemberExpressionTree:
		memberExpr := b.(ast.MemberExpression)
		value := i.EvaluateMemberExpr(memberExpr, memberExpr.Computed)
		return value

	case ast.BinaryExpressionTree:
		leftValue := i.InterpretBinaryExpr(b.(ast.BinaryExpression).Left)
		rightValue := i.InterpretBinaryExpr(b.(ast.BinaryExpression).Right)
		operator := b.(ast.BinaryExpression).Operator
		isConcat := false

		var leftFloat, rightFloat float64
		var isLeftFloat, isRightFloat bool

		boolType := reflect.TypeOf(true)
		stringType := reflect.TypeOf("")
		leftType := reflect.TypeOf(leftValue)
		rightType := reflect.TypeOf(rightValue)

		if leftType == stringType || rightType == stringType {
			isConcat = true
		}

		if leftType == boolType || rightType == boolType {
			panic(oops.SyntaxError(b.(ast.BinaryExpression).Right, "Invalid operation"))
		}

		switch v := leftValue.(type) {
		case int:
			leftFloat = float64(v)
		case float64:
			leftFloat = v
			isLeftFloat = true
		case string:
			leftValue = string(v)
		default:
			return v
		}

		switch v := rightValue.(type) {
		case int:
			rightFloat = float64(v)
		case float64:
			rightFloat = v
			isRightFloat = true
		case string:
			rightValue = string(v)
		default:
			return v
		}

		if isConcat {
			if operator != "+" {
				panic(oops.SyntaxError(b.(ast.BinaryExpression).Right, "Invalid operation"))
			}

			return fmt.Sprintf("%v", leftValue) + fmt.Sprintf("%v", rightValue)
		}

		var result interface{}
		switch operator {
		case "+":
			result = leftFloat + rightFloat
		case "-":
			result = leftFloat - rightFloat
		case "*":
			result = leftFloat * rightFloat
		case "/":
			if rightFloat == 0 {
				return 0
			}
			result = leftFloat / rightFloat
		}

		if isLeftFloat || isRightFloat {
			return result
		}
		return int(result.(float64))

	default:
		return 0
	}
}
