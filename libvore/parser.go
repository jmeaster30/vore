package libvore

import (
	"strconv"
)

type ParseState int

const (
	P_START ParseState = iota
	P_COMMAND
	P_FIND
	P_REPLACE
	P_SET
	P_COMMAND_AMOUNT
)

type ParseError struct {
	isError bool
	token   Token
	message string
}

func NewParseError(t Token, message string) ParseError {
	return ParseError{
		isError: true,
		token:   t,
		message: message,
	}
}

func NoError() ParseError {
	return ParseError{isError: false}
}

func parse(tokens []*Token) ([]AstCommand, ParseError) {
	commands := []AstCommand{}

	token_index := 0
	for token_index < len(tokens)-1 {
		ws_index := consumeIgnoreableTokens(tokens, token_index)
		command, new_index, e := parse_command(tokens, ws_index)
		if e.isError {
			return []AstCommand{}, e
		}
		token_index = new_index
		commands = append(commands, command)
	}

	return commands, NoError()
}

func parse_command(tokens []*Token, token_index int) (AstCommand, int, ParseError) {
	if tokens[token_index].tokenType == FIND {
		return parse_find(tokens, token_index)
	} else if tokens[token_index].tokenType == REPLACE {
		return parse_replace(tokens, token_index)
	} else if tokens[token_index].tokenType == SET {
		return parse_set(tokens, token_index)
	} else {
		return nil, token_index, NewParseError(*tokens[token_index], "Unexpected token. Expected 'find', 'replace', or 'set'.")
	}
}

func parse_find(tokens []*Token, token_index int) (*AstFind, int, ParseError) {
	all, skipValue, takeValue, new_index, amountError := parse_amount(tokens, token_index+1)
	if amountError.isError {
		return nil, new_index, amountError
	}

	findCommand := AstFind{
		all:  all,
		skip: skipValue,
		take: takeValue,
		body: []AstExpression{},
	}

	var current_token = tokens[new_index]
	var current_index = new_index
	for current_token.tokenType != FIND && current_token.tokenType != REPLACE && current_token.tokenType != SET && current_token.tokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_expression(tokens, ws_index)
		if parseError.isError {
			return nil, new_index, parseError
		}

		findCommand.body = append(findCommand.body, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	return &findCommand, current_index, NoError()
}

func parse_replace(tokens []*Token, token_index int) (*AstReplace, int, ParseError) {
	all, skipValue, takeValue, new_index, amountError := parse_amount(tokens, token_index+1)
	if amountError.isError {
		return nil, new_index, amountError
	}

	replaceCommand := AstReplace{
		all:  all,
		skip: skipValue,
		take: takeValue,
		body: []AstExpression{},
	}

	var current_token = tokens[new_index]
	var current_index = new_index
	for current_token.tokenType != WITH && current_token.tokenType != FIND && current_token.tokenType != REPLACE && current_token.tokenType != SET && current_token.tokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_expression(tokens, ws_index)
		if parseError.isError {
			return nil, new_index, parseError
		}

		replaceCommand.body = append(replaceCommand.body, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	if current_token.tokenType != WITH {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'with'.")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	for current_token.tokenType != FIND && current_token.tokenType != REPLACE && current_token.tokenType != SET && current_token.tokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_atom(tokens, ws_index)
		if parseError.isError {
			return nil, new_index, parseError
		}

		replaceCommand.result = append(replaceCommand.result, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	return &replaceCommand, current_index, NoError()
}

func parse_set(tokens []*Token, token_index int) (*AstSet, int, ParseError) {
	var current_index = consumeIgnoreableTokens(tokens, token_index+1)
	var current_token = tokens[current_index]

	if current_token.tokenType != IDENTIFIER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected identifier")
	}

	name := current_token.lexeme

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]

	if current_token.tokenType != TO {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'to'")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	expr, new_index, err := parse_expression(tokens, current_index)
	if err.isError {
		return nil, new_index, err
	}
	current_index = new_index

	setCommand := AstSet{
		id:   name,
		expr: expr,
	}
	return &setCommand, current_index, NoError()
}

func parse_amount(tokens []*Token, token_index int) (bool, int, int, int, ParseError) {
	new_index := consumeIgnoreableTokens(tokens, token_index)

	if tokens[new_index].tokenType == ALL {
		new_index += 1
		return true, 0, 0, new_index, NoError()
	} else if tokens[new_index].tokenType == SKIP {
		new_index = consumeIgnoreableTokens(tokens, new_index+1)
		if tokens[new_index].tokenType == NUMBER {
			skipValue, skipValueError := strconv.Atoi(tokens[new_index].lexeme)
			if skipValueError != nil {
				return false, 0, 0, 0, NewParseError(*tokens[new_index], "Error converting to int value")
			}

			new_index = consumeIgnoreableTokens(tokens, new_index+1)
			if tokens[new_index].tokenType == TAKE {
				new_index = consumeIgnoreableTokens(tokens, new_index+1)
				if tokens[new_index].tokenType == NUMBER {
					takeValue, takeValueError := strconv.Atoi(tokens[new_index].lexeme)
					if takeValueError != nil {
						return false, 0, 0, 0, NewParseError(*tokens[new_index], "Error converting to int value")
					}
					new_index++
					return false, skipValue, takeValue, new_index, NoError()
				} else {
					return false, 0, 0, 0, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
				}
			}
			return true, skipValue, 0, new_index, NoError()
		} else {
			return false, 0, 0, 0, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
		}
	} else if tokens[new_index].tokenType == TAKE {
		new_index = consumeIgnoreableTokens(tokens, new_index+1)
		if tokens[new_index].tokenType == NUMBER {
			takeValue, takeValueError := strconv.Atoi(tokens[new_index].lexeme)
			if takeValueError != nil {
				return false, 0, 0, 0, NewParseError(*tokens[new_index], "Error converting to int value")
			}
			new_index++
			return false, 0, takeValue, new_index, NoError()
		} else {
			return false, 0, 0, 0, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
		}
	}
	return false, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'all', 'skip', or 'take'")
}

func parse_expression(tokens []*Token, token_index int) (AstExpression, int, ParseError) {
	current_token := tokens[token_index]
	if current_token.tokenType == AT {
	} else if current_token.tokenType == BETWEEN {
	} else if current_token.tokenType == EXACTLY {
	} else if current_token.tokenType == MAYBE {
	} else if current_token.tokenType == IN {
	} else if current_token.tokenType == STRING || current_token.tokenType == IDENTIFIER ||
		current_token.tokenType == OPENPAREN || current_token.tokenType == ANY ||
		current_token.tokenType == WHITESPACE || current_token.tokenType == DIGIT ||
		current_token.tokenType == UPPER || current_token.tokenType == LOWER ||
		current_token.tokenType == LETTER || current_token.tokenType == LINE ||
		current_token.tokenType == FILE {
		return parse_primary(tokens, token_index)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected 'at', 'between', 'exactly', 'maybe', 'in', '<string>', '<identifier>', or a character class ")
}

func parse_primary(tokens []*Token, token_index int) (*AstPrimary, int, ParseError) {
	current_token := tokens[token_index]
	prim := AstPrimary{}
	if current_token.tokenType == STRING {
		str_literal, new_index, err := parse_string(tokens, token_index)
		if err.isError {
			return nil, new_index, err
		}
		prim.literal = str_literal
		return &prim, new_index, NoError()
	} else if current_token.tokenType == IDENTIFIER {
		variable, new_index, err := parse_variable(tokens, token_index)
		if err.isError {
			return nil, new_index, err
		}
		prim.literal = variable
		return &prim, new_index, NoError()
	} else if current_token.tokenType == OPENPAREN {
		sub_expr, new_index, err := parse_sub_expression(tokens, token_index)
		if err.isError {
			return nil, new_index, err
		}
		prim.literal = sub_expr
		return &prim, new_index, NoError()
	} else if current_token.tokenType == ANY ||
		current_token.tokenType == WHITESPACE || current_token.tokenType == DIGIT ||
		current_token.tokenType == UPPER || current_token.tokenType == LOWER ||
		current_token.tokenType == LETTER || current_token.tokenType == LINE ||
		current_token.tokenType == FILE {
		cc_literal, new_index, err := parse_character_class(tokens, token_index)
		if err.isError {
			return nil, new_index, err
		}
		prim.literal = cc_literal

		return &prim, new_index, NoError()
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected '(', '<string>', '<identifier>', or a character class.")
}

func parse_atom(tokens []*Token, token_index int) (AstAtom, int, ParseError) {
	current_token := tokens[token_index]
	if current_token.tokenType == STRING {
		return parse_string(tokens, token_index)
	} else if current_token.tokenType == IDENTIFIER {
		return parse_variable(tokens, token_index)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected '<string>' or '<identifier>'.")
}

func parse_string(tokens []*Token, token_index int) (*AstString, int, ParseError) {
	current_token := tokens[token_index]
	str_literal := AstString{}

	if current_token.tokenType == STRING {
		str_literal.value = current_token.lexeme
		return &str_literal, token_index + 1, NoError()
	}

	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected a string")
}

func parse_variable(tokens []*Token, token_index int) (*AstVariable, int, ParseError) {
	current_token := tokens[token_index]
	var_literal := AstVariable{}

	if current_token.tokenType == IDENTIFIER {
		var_literal.name = current_token.lexeme
		return &var_literal, token_index + 1, NoError()
	}

	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected a variable")
}

func parse_sub_expression(tokens []*Token, token_index int) (*AstSubExpr, int, ParseError) {
	current_token := tokens[token_index+1]
	current_index := token_index + 1
	expr_list := []AstExpression{}

	for current_token.tokenType != CLOSEPAREN && current_token.tokenType != FIND && current_token.tokenType != REPLACE && current_token.tokenType != SET && current_token.tokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_expression(tokens, ws_index)
		if parseError.isError {
			return nil, new_index, parseError
		}

		expr_list = append(expr_list, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	if current_token.tokenType != CLOSEPAREN {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected ')'")
	}

	sub_expr := AstSubExpr{body: expr_list}
	return &sub_expr, current_index + 1, NoError()
}

func parse_character_class(tokens []*Token, token_index int) (*AstCharacterClass, int, ParseError) {
	current_token := tokens[token_index]
	charClass := AstCharacterClass{}
	if current_token.tokenType == ANY {
		charClass.classType = ClassAny
		return &charClass, token_index + 1, NoError()
	} else if current_token.tokenType == WHITESPACE {
		charClass.classType = ClassWhitespace
		return &charClass, token_index + 1, NoError()
	} else if current_token.tokenType == DIGIT {
		charClass.classType = ClassDigit
		return &charClass, token_index + 1, NoError()
	} else if current_token.tokenType == UPPER {
		charClass.classType = ClassUpper
		return &charClass, token_index + 1, NoError()
	} else if current_token.tokenType == LOWER {
		charClass.classType = ClassLower
		return &charClass, token_index + 1, NoError()
	} else if current_token.tokenType == LETTER {
		charClass.classType = ClassLetter
		return &charClass, token_index + 1, NoError()
	} else if current_token.tokenType == LINE {
		new_index := consumeIgnoreableTokens(tokens, token_index+1)
		if tokens[new_index].tokenType == START {
			charClass.classType = ClassLineStart
			return &charClass, new_index + 1, NoError()
		} else if tokens[new_index].tokenType == END {
			charClass.classType = ClassLineEnd
			return &charClass, new_index + 1, NoError()
		}
		return nil, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'start' or 'end'")
	} else if current_token.tokenType == FILE {
		new_index := consumeIgnoreableTokens(tokens, token_index+1)
		if tokens[new_index].tokenType == START {
			charClass.classType = ClassFileStart
			return &charClass, new_index + 1, NoError()
		} else if tokens[new_index].tokenType == END {
			charClass.classType = ClassFileEnd
			return &charClass, new_index + 1, NoError()
		}
		return nil, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'start' or 'end'")
	}
	return nil, token_index, NewParseError(*tokens[token_index], "Unexpected token. Expected a character class: 'any', 'whitespace', 'digit', 'upper', 'lower', 'letter', 'line start', 'line end', 'file start', or 'file end'.")
}

func consumeIgnoreableTokens(tokens []*Token, index int) int {
	current_index := index
	for tokens[current_index].tokenType == WS || tokens[current_index].tokenType == COMMENT {
		current_index += 1
	}
	return current_index
}
