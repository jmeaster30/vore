package libvore

import (
	"strconv"
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
		if command != nil {
			commands = append(commands, command)
		}
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
	} else if tokens[token_index].tokenType == EOF {
		return nil, token_index, NoError()
	} else {
		return nil, token_index, NewParseError(*tokens[token_index], "Unexpected token. Expected 'find', 'replace', or 'set'.")
	}
}

func parse_find(tokens []*Token, token_index int) (*AstFind, int, ParseError) {
	all, skipValue, takeValue, lastValue, new_index, amountError := parse_amount(tokens, token_index+1)
	if amountError.isError {
		return nil, new_index, amountError
	}

	findCommand := AstFind{
		all:  all,
		skip: skipValue,
		take: takeValue,
		last: lastValue,
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
	all, skipValue, takeValue, lastValue, new_index, amountError := parse_amount(tokens, token_index+1)
	if amountError.isError {
		return nil, new_index, amountError
	}

	replaceCommand := AstReplace{
		all:  all,
		skip: skipValue,
		take: takeValue,
		last: lastValue,
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

func parse_amount(tokens []*Token, token_index int) (bool, int, int, int, int, ParseError) {
	new_index := consumeIgnoreableTokens(tokens, token_index)

	if tokens[new_index].tokenType == ALL {
		new_index += 1
		return true, 0, 0, 0, new_index, NoError()
	} else if tokens[new_index].tokenType == SKIP {
		new_index = consumeIgnoreableTokens(tokens, new_index+1)
		if tokens[new_index].tokenType == NUMBER {
			skipValue, skipValueError := strconv.Atoi(tokens[new_index].lexeme)
			if skipValueError != nil {
				return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Error converting to int value")
			}

			new_index = consumeIgnoreableTokens(tokens, new_index+1)
			if tokens[new_index].tokenType == TAKE {
				new_index = consumeIgnoreableTokens(tokens, new_index+1)
				if tokens[new_index].tokenType == NUMBER {
					takeValue, takeValueError := strconv.Atoi(tokens[new_index].lexeme)
					if takeValueError != nil {
						return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Error converting to int value")
					}
					new_index++
					return false, skipValue, takeValue, 0, new_index, NoError()
				} else {
					return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
				}
			}
			return true, skipValue, 0, 0, new_index, NoError()
		} else {
			return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
		}
	} else if tokens[new_index].tokenType == TAKE || tokens[new_index].tokenType == TOP {
		new_index = consumeIgnoreableTokens(tokens, new_index+1)
		if tokens[new_index].tokenType == NUMBER {
			takeValue, takeValueError := strconv.Atoi(tokens[new_index].lexeme)
			if takeValueError != nil {
				return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Error converting to int value")
			}
			new_index++
			return false, 0, takeValue, 0, new_index, NoError()
		} else {
			return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
		}
	} else if tokens[new_index].tokenType == LAST {
		new_index = consumeIgnoreableTokens(tokens, new_index+1)
		if tokens[new_index].tokenType == NUMBER {
			lastValue, lastValueError := strconv.Atoi(tokens[new_index].lexeme)
			if lastValueError != nil {
				return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Error converting to int value")
			}
			new_index++
			return true, 0, 0, lastValue, new_index, NoError()
		} else {
			return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
		}
	}
	return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'all', 'skip', or 'take'")
}

func parse_expression(tokens []*Token, token_index int) (AstExpression, int, ParseError) {
	current_token := tokens[token_index]
	if current_token.tokenType == AT {
		return parse_at(tokens, token_index)
	} else if current_token.tokenType == BETWEEN {
		return parse_between(tokens, token_index)
	} else if current_token.tokenType == EXACTLY {
		return parse_exactly(tokens, token_index)
	} else if current_token.tokenType == MAYBE {
		return parse_maybe(tokens, token_index)
	} else if current_token.tokenType == IN {
		return parse_in(tokens, token_index, false)
	} else if current_token.tokenType == OPENCURLY {
		return parse_subroutine(tokens, token_index)
	} else if current_token.tokenType == NOT {
		return parse_not_expression(tokens, token_index)
	} else if current_token.tokenType == STRING || current_token.tokenType == IDENTIFIER ||
		current_token.tokenType == OPENPAREN || current_token.tokenType == ANY ||
		current_token.tokenType == WHITESPACE || current_token.tokenType == DIGIT ||
		current_token.tokenType == UPPER || current_token.tokenType == LOWER ||
		current_token.tokenType == LETTER || current_token.tokenType == LINE ||
		current_token.tokenType == FILE {
		return parse_primary_or_dec(tokens, token_index)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected 'at', 'between', 'exactly', 'maybe', 'in', '<string>', '<identifier>', or a character class ")
}

func parse_at(tokens []*Token, token_index int) (*AstLoop, int, ParseError) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[current_index]

	if current_token.tokenType != LEAST && current_token.tokenType != MOST {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'least' or 'most'.")
	}

	isLeast := current_token.tokenType == LEAST

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	if current_token.tokenType != NUMBER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected a number.")
	}
	value, err := strconv.Atoi(current_token.lexeme)
	if err != nil {
		return nil, current_index, NewParseError(*current_token, "Error converting lexeme to number value")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	literal, next_index, parseError := parse_literal(tokens, current_index)
	if parseError.isError {
		return nil, next_index, parseError
	}

	current_index = consumeIgnoreableTokens(tokens, next_index)
	current_token = tokens[current_index]
	fewest := current_token.tokenType == FEWEST
	if fewest {
		current_index += 1
	}

	var min int
	var max int
	if isLeast {
		min = value
		max = -1
	} else {
		min = 0
		max = value
	}

	atLoop := AstLoop{
		min:    min,
		max:    max,
		fewest: fewest,
		body:   literal,
	}

	return &atLoop, current_index, NoError()
}

func parse_between(tokens []*Token, token_index int) (*AstLoop, int, ParseError) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[current_index]

	if current_token.tokenType != NUMBER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected a number.")
	}
	minValue, err := strconv.Atoi(current_token.lexeme)
	if err != nil {
		return nil, current_index, NewParseError(*current_token, "Error converting lexeme to number value")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	if current_token.tokenType != AND {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'and'.")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	if current_token.tokenType != NUMBER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected a number.")
	}
	maxValue, err := strconv.Atoi(current_token.lexeme)
	if err != nil {
		return nil, current_index, NewParseError(*current_token, "Error converting lexeme to number value")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	literal, next_index, parseError := parse_literal(tokens, current_index)
	if parseError.isError {
		return nil, next_index, parseError
	}

	current_index = consumeIgnoreableTokens(tokens, next_index)
	current_token = tokens[current_index]
	fewest := current_token.tokenType == FEWEST
	if fewest {
		current_index += 1
	}

	between := AstLoop{
		min:    minValue,
		max:    maxValue,
		fewest: fewest,
		body:   literal,
	}

	return &between, current_index, NoError()
}

func parse_exactly(tokens []*Token, token_index int) (*AstLoop, int, ParseError) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[current_index]

	if current_token.tokenType != NUMBER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected a number.")
	}
	value, err := strconv.Atoi(current_token.lexeme)
	if err != nil {
		return nil, current_index, NewParseError(*current_token, "Error converting lexeme to number value")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	literal, next_index, parseError := parse_literal(tokens, current_index)
	if parseError.isError {
		return nil, next_index, parseError
	}

	exactly := AstLoop{
		min:    value,
		max:    value,
		fewest: false,
		body:   literal,
	}

	return &exactly, next_index, NoError()
}

func parse_maybe(tokens []*Token, token_index int) (*AstLoop, int, ParseError) {
	new_index := consumeIgnoreableTokens(tokens, token_index+1)
	literal, next_index, err := parse_literal(tokens, new_index)
	if err.isError {
		return nil, next_index, err
	}

	current_index := consumeIgnoreableTokens(tokens, next_index)
	current_token := tokens[current_index]
	fewest := current_token.tokenType == FEWEST
	if fewest {
		current_index += 1
	}

	maybe := AstLoop{0, 1, fewest, literal}

	return &maybe, current_index, NoError()
}

func parse_not_expression(tokens []*Token, token_index int) (AstExpression, int, ParseError) {
	new_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[new_index]

	if current_token.tokenType == IN {
		return parse_in(tokens, new_index, true)
	} else if current_token.tokenType == STRING {
		ast_str, idx, err := parse_string(tokens, new_index, true)
		if err.isError {
			return nil, idx, err
		}
		return &AstPrimary{literal: ast_str}, idx, err
	} else if current_token.tokenType == ANY ||
		current_token.tokenType == WHITESPACE || current_token.tokenType == DIGIT ||
		current_token.tokenType == UPPER || current_token.tokenType == LOWER ||
		current_token.tokenType == LETTER || current_token.tokenType == LINE ||
		current_token.tokenType == FILE {
		ast_chclass, idx, err := parse_character_class(tokens, new_index, true)
		if err.isError {
			return nil, idx, err
		}
		return &AstPrimary{literal: ast_chclass}, idx, err
	} else {
		return nil, new_index, NewParseError(*current_token, "Unexpected token. Expected 'in', <string>, <character class>")
	}
}

func parse_not_literal(tokens []*Token, token_index int) (AstLiteral, int, ParseError) {
	new_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[new_index]

	if current_token.tokenType == STRING {
		return parse_string(tokens, new_index, true)
	} else if current_token.tokenType == ANY ||
		current_token.tokenType == WHITESPACE || current_token.tokenType == DIGIT ||
		current_token.tokenType == UPPER || current_token.tokenType == LOWER ||
		current_token.tokenType == LETTER || current_token.tokenType == LINE ||
		current_token.tokenType == FILE {
		return parse_character_class(tokens, new_index, true)
	} else {
		return nil, new_index, NewParseError(*current_token, "Unexpected token. Expected 'in', <string>, <character class>")
	}
}

func parse_in(tokens []*Token, token_index int, not bool) (*AstList, int, ParseError) {
	new_index := consumeIgnoreableTokens(tokens, token_index+1)
	contents := []AstListable{}

	listable, next_index, err := parse_listable(tokens, new_index)
	if err.isError {
		return nil, next_index, err
	}
	contents = append(contents, listable)

	current_index := consumeIgnoreableTokens(tokens, next_index)
	current_token := tokens[current_index]
	for current_token.tokenType == COMMA {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		listable, next_index, err := parse_listable(tokens, current_index)
		if err.isError {
			return nil, next_index, err
		}
		contents = append(contents, listable)
		current_index = next_index
		current_token = tokens[current_index]
	}

	inList := AstList{contents: contents, not: not}
	return &inList, current_index, NoError()
}

func isListableClass(t TokenType) bool {
	return t == ANY || t == WHITESPACE || t == DIGIT || t == UPPER || t == LOWER || t == LETTER
}

func parse_listable(tokens []*Token, token_index int) (AstListable, int, ParseError) {
	current_token := tokens[token_index]
	if current_token.tokenType == STRING {
		from, next_index, err := parse_string(tokens, token_index, false)
		if err.isError {
			return nil, next_index, err
		}

		current_index := consumeIgnoreableTokens(tokens, next_index)
		current_token := tokens[current_index]
		if current_token.tokenType != TO {
			return from, current_index, NoError()
		}

		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		to, new_index, terr := parse_string(tokens, current_index, false)
		if terr.isError {
			return nil, new_index, terr
		}

		r := AstRange{
			from: from,
			to:   to,
		}
		return &r, new_index, NoError()

	} else if isListableClass(current_token.tokenType) {
		return parse_character_class(tokens, token_index, false)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected listable literal")
}

func parse_literal(tokens []*Token, token_index int) (AstLiteral, int, ParseError) {
	current_token := tokens[token_index]
	if current_token.tokenType == STRING {
		return parse_string(tokens, token_index, false)
	} else if current_token.tokenType == IDENTIFIER {
		return parse_variable(tokens, token_index)
	} else if current_token.tokenType == OPENPAREN {
		return parse_sub_expression(tokens, token_index)
	} else if current_token.tokenType == NOT {
		return parse_not_literal(tokens, token_index)
	} else if current_token.tokenType == ANY ||
		current_token.tokenType == WHITESPACE || current_token.tokenType == DIGIT ||
		current_token.tokenType == UPPER || current_token.tokenType == LOWER ||
		current_token.tokenType == LETTER || current_token.tokenType == LINE ||
		current_token.tokenType == FILE {
		return parse_character_class(tokens, token_index, false)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected '(', '<string>', '<identifier>', or a character class.")
}

func parse_primary_or_dec(tokens []*Token, token_index int) (AstExpression, int, ParseError) {
	literal, new_index, err := parse_literal(tokens, token_index)
	if err.isError {
		return nil, new_index, err
	}

	current_index := consumeIgnoreableTokens(tokens, new_index)
	current_token := tokens[current_index]
	if current_token.tokenType == EQUAL {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		current_token = tokens[current_index]

		if current_token.tokenType != IDENTIFIER {
			return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected identifier.")
		}

		dec := AstDec{
			name: current_token.lexeme,
			body: literal,
		}
		return &dec, current_index + 1, NoError()
	}

	if current_token.tokenType == OR {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		right_literal, final_index, err := parse_literal(tokens, current_index)
		if err.isError {
			return nil, final_index, err
		}

		branch := AstBranch{
			left:  literal,
			right: right_literal,
		}
		return &branch, final_index, NoError()
	}

	prim := AstPrimary{}
	prim.literal = literal

	return &prim, new_index, NoError()
}

func parse_primary(tokens []*Token, token_index int) (*AstPrimary, int, ParseError) {
	prim := AstPrimary{}

	literal, new_index, err := parse_literal(tokens, token_index)
	if err.isError {
		return nil, new_index, err
	}
	prim.literal = literal

	return &prim, new_index, NoError()
}

func parse_atom(tokens []*Token, token_index int) (AstAtom, int, ParseError) {
	current_token := tokens[token_index]
	if current_token.tokenType == STRING {
		return parse_string(tokens, token_index, false)
	} else if current_token.tokenType == IDENTIFIER {
		return parse_variable(tokens, token_index)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected '<string>' or '<identifier>'.")
}

func parse_string(tokens []*Token, token_index int, not bool) (*AstString, int, ParseError) {
	current_token := tokens[token_index]
	str_literal := AstString{
		not: not,
	}

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

func parse_subroutine(tokens []*Token, token_index int) (*AstSub, int, ParseError) {
	current_token := tokens[token_index+1]
	current_index := token_index + 1
	expr_list := []AstExpression{}

	for current_token.tokenType != CLOSECURLY && current_token.tokenType != FIND && current_token.tokenType != REPLACE && current_token.tokenType != SET && current_token.tokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_expression(tokens, ws_index)
		if parseError.isError {
			return nil, new_index, parseError
		}

		expr_list = append(expr_list, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	if current_token.tokenType != CLOSECURLY {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected '}'")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	if current_token.tokenType != EQUAL {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected '='")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]

	if current_token.tokenType != IDENTIFIER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected identifier.")
	}

	dec := AstSub{
		name: current_token.lexeme,
		body: expr_list,
	}
	return &dec, current_index + 1, NoError()
}

func parse_character_class(tokens []*Token, token_index int, not bool) (*AstCharacterClass, int, ParseError) {
	current_token := tokens[token_index]
	charClass := AstCharacterClass{
		not: not,
	}
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
