package interpreter

import (
	"errors"
	"fmt"

	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/parser/ast"
	"github.com/rpromyshlennikov/lox_tree_walk_interpretator/pkg/scanner"
)

type RuntimeError struct {
	token   scanner.Token
	message string
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf(`Runtime error: "%s" at token: %v`, re.message, re.token)
}

func NewRuntimeError(token scanner.Token, message string) *RuntimeError {
	return &RuntimeError{
		token:   token,
		message: message,
	}
}

type Interpreter struct{}

func (i Interpreter) Interpret(expr ast.Expr) (result string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			rErr, ok := recovered.(*RuntimeError)
			if !ok {
				err = errors.New("unknown error in runtime during interpret")
			}
			err = rErr
		}
	}()
	if expr == nil {
		return "", NewRuntimeError(
			scanner.NewToken(scanner.EOF, "", nil, 0),
			"no expression given",
		)
	}
	value := i.evaluate(expr)
	return i.stringify(value), nil
}

func (i Interpreter) VisitUnary(unary *ast.Unary) any {
	right := i.evaluate(unary.Right)
	switch unary.Operator.Kind() {
	case scanner.BANG:
		return !i.isTruthy(right)
	case scanner.MINUS:
		i.checkNumberOperands(unary.Operator, right)
		return -right.(float64)
	}
	// Unreachable.
	return nil
}

// TODO: gocyclo considers this function too difficult,
// refactor to several private operator methods after book will be completed.
//
//gocyclo:ignore
func (i Interpreter) VisitBinary(binary *ast.Binary) any {
	left := i.evaluate(binary.Left)
	right := i.evaluate(binary.Right)
	switch binary.Operator.Kind() {
	case scanner.BANGEQUAL:
		return !i.isEqual(left, right)
	case scanner.EQUALEQUAL:
		return i.isEqual(left, right)

	case scanner.GREATER:
		i.checkNumberOperands(binary.Operator, left, right)
		return left.(float64) > right.(float64)
	case scanner.GREATEREQUAL:
		i.checkNumberOperands(binary.Operator, left, right)
		return left.(float64) >= right.(float64)
	case scanner.LESS:
		i.checkNumberOperands(binary.Operator, left, right)
		return left.(float64) < right.(float64)
	case scanner.LESSEQUAL:
		i.checkNumberOperands(binary.Operator, left, right)
		return left.(float64) <= right.(float64)

	case scanner.MINUS:
		i.checkNumberOperands(binary.Operator, left, right)
		return left.(float64) - right.(float64)
	case scanner.PLUS:
		const floatType = "float"
		const stringType = "string"
		var leftType string
		var rightType string
		switch left.(type) {
		case float64:
			leftType = floatType
		case string:
			leftType = stringType
		}
		switch right.(type) {
		case float64:
			rightType = floatType
		case string:
			rightType = stringType
		}
		if leftType == floatType && rightType == floatType {
			return left.(float64) + right.(float64)
		}
		if leftType == stringType && rightType == stringType {
			return left.(string) + right.(string)
		}
		err := NewRuntimeError(
			binary.Operator,
			"invalid type for operator '+' given, must be numbers or strings.",
		)
		panic(err)
	case scanner.SLASH:
		i.checkNumberOperands(binary.Operator, left, right)
		return left.(float64) / right.(float64)
	case scanner.STAR:
		i.checkNumberOperands(binary.Operator, left, right)
		return left.(float64) * right.(float64)
	}
	// Unreachable.
	return nil
}

func (i Interpreter) VisitLiteral(literal *ast.Literal) any {
	return literal.Value
}

func (i Interpreter) VisitGrouping(grouping *ast.Grouping) any {
	return i.evaluate(grouping.Expression)
}

func (i Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i Interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if value, ok := value.(bool); ok {
		return value
	}
	return true
}

func (i Interpreter) isEqual(left any, right any) bool {
	return left == right
}

func (i Interpreter) checkNumberOperands(operator scanner.Token, operands ...any) {
	for _, operand := range operands {
		switch operand.(type) {
		case float64:
			continue
		}
		err := NewRuntimeError(
			operator,
			fmt.Sprintf(
				"invalid type for operator %s given, must be number.",
				operator.Kind(),
			),
		)
		panic(err)
	}
}

func (i Interpreter) stringify(value any) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprint(value)
}
