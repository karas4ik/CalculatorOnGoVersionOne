package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var precedence = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func applyOperation(a, b float64, op rune) float64 {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		if b == 0 {
			panic("division by zero")
		}
		return a / b
	}
	return 0
}

func Calc(expression string) (float64, error) {
	var values []float64
	var ops []rune

	expression = strings.ReplaceAll(expression, " ", "")

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if unicode.IsDigit(char) || char == '.' {
			start := i
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				i++
			}
			numStr := expression[start:i]
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, errors.New("invalid number: " + numStr)
			}
			values = append(values, num)
			i--
		} else if char == '(' {
			ops = append(ops, char)
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if len(values) < 2 {
					return 0, errors.New("invalid expression: not enough values")
				}
				v2 := values[len(values)-1]
				values = values[:len(values)-1]
				v1 := values[len(values)-1]
				values = values[:len(values)-1]
				operator := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				values = append(values, applyOperation(v1, v2, operator))
			}
			if len(ops) == 0 {
				return 0, errors.New("mismatched parentheses")
			}
			ops = ops[:len(ops)-1] // Удаляем '('
		} else if op, ok := precedence[char]; ok {
			for len(ops) > 0 && precedence[ops[len(ops)-1]] >= op {
				if len(values) < 2 {
					return 0, errors.New("invalid expression: not enough values")
				}
				v2 := values[len(values)-1]
				values = values[:len(values)-1]
				v1 := values[len(values)-1]
				values = values[:len(values)-1]
				operator := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				values = append(values, applyOperation(v1, v2, operator))
			}
			ops = append(ops, char)
		} else {
			return 0, errors.New("invalid character: " + string(char))
		}
	}

	for len(ops) > 0 {
		if len(values) < 2 {
			return 0, errors.New("invalid expression: not enough values")
		}
		v2 := values[len(values)-1]
		values = values[:len(values)-1]
		v1 := values[len(values)-1]
		values = values[:len(values)-1]
		operator := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		values = append(values, applyOperation(v1, v2, operator))
	}

	if len(values) != 1 {
		return 0, errors.New("invalid expression: unmatched values and operators")
	}
	return values[0], nil
}
