package libvore

import (
	"strconv"
)

func parse(tokens []*Token) ([]AstCommand, error) {
	commands := []AstCommand{}

	token_index := 0
	for token_index < len(tokens)-1 {
		ws_index := consumeIgnoreableTokens(tokens, token_index)
		command, new_index, e := parse_command(tokens, ws_index)
		if e != nil {
			return []AstCommand{}, e
		}
		token_index = new_index
		if command != nil {
			commands = append(commands, command)
		}
	}

	return commands, nil
}

func parse_command(tokens []*Token, token_index int) (AstCommand, int, error) {
	if tokens[token_index].TokenType == FIND {
		return parse_find(tokens, token_index)
	} else if tokens[token_index].TokenType == REPLACE {
		return parse_replace(tokens, token_index)
	} else if tokens[token_index].TokenType == SET {
		return parse_set(tokens, token_index)
	} else if tokens[token_index].TokenType == EOF {
		return nil, token_index, nil
	} else {
		return nil, token_index, NewParseError(*tokens[token_index], "Unexpected token. Expected 'find', 'replace', or 'set'.")
	}
}

func parse_find(tokens []*Token, token_index int) (*AstFind, int, error) {
	all, skipValue, takeValue, lastValue, new_index, amountError := parse_amount(tokens, token_index+1)
	if amountError != nil {
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
	for current_token.TokenType != FIND && current_token.TokenType != REPLACE && current_token.TokenType != SET && current_token.TokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_expression(tokens, ws_index)
		if parseError != nil {
			return nil, new_index, parseError
		}

		findCommand.body = append(findCommand.body, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	return &findCommand, current_index, nil
}

func parse_replace(tokens []*Token, token_index int) (*AstReplace, int, error) {
	all, skipValue, takeValue, lastValue, new_index, amountError := parse_amount(tokens, token_index+1)
	if amountError != nil {
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
	for current_token.TokenType != WITH && current_token.TokenType != FIND && current_token.TokenType != REPLACE && current_token.TokenType != SET && current_token.TokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_expression(tokens, ws_index)
		if parseError != nil {
			return nil, new_index, parseError
		}

		replaceCommand.body = append(replaceCommand.body, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	if current_token.TokenType != WITH {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'with'.")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	for current_token.TokenType != FIND && current_token.TokenType != REPLACE && current_token.TokenType != SET && current_token.TokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_atom(tokens, ws_index)
		if parseError != nil {
			return nil, new_index, parseError
		}

		replaceCommand.result = append(replaceCommand.result, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	return &replaceCommand, current_index, nil
}

func parse_set(tokens []*Token, token_index int) (*AstSet, int, error) {
	var current_index = consumeIgnoreableTokens(tokens, token_index+1)
	var current_token = tokens[current_index]

	if current_token.TokenType != IDENTIFIER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected identifier")
	}

	name := current_token.Lexeme

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]

	if current_token.TokenType != TO {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'to'")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	var body AstSetBody
	if current_token.TokenType == PATTERN {
		expr, next_index, err := parse_set_pattern(tokens, current_index)
		if err != nil {
			return nil, next_index, err
		}
		body = expr
		current_index = next_index
	} else if current_token.TokenType == MATCHES {
		expr, next_index, err := parse_set_matches(tokens, current_index)
		if err != nil {
			return nil, next_index, err
		}
		body = expr
		current_index = next_index
	} else if current_token.TokenType == TRANSFORM {
		expr, next_index, err := parse_set_transform(tokens, current_index)
		if err != nil {
			return nil, next_index, err
		}
		body = expr
		current_index = next_index
	} else {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'pattern', 'transform', or 'matches'")
	}

	setCommand := AstSet{
		id:   name,
		body: body,
	}
	return &setCommand, current_index, nil
}

func parse_set_transform(tokens []*Token, token_index int) (AstSetBody, int, error) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)

	if tokens[current_index].TokenType == BEGIN {
		current_index += 1
	}

	statements, next_index, err := parse_process_statements(tokens, current_index)
	if err != nil {
		return nil, next_index, err
	}

	if tokens[next_index].TokenType != END {
		return nil, next_index, NewParseError(*tokens[next_index], "Unexpected token. Expected 'end'.")
	}

	return &AstSetTransform{statements}, next_index + 1, err
}

func parse_set_pattern(tokens []*Token, token_index int) (AstSetBody, int, error) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)
	expr, next_index, err := parse_expression(tokens, current_index)
	if err != nil {
		return nil, next_index, err
	}

	current_index = consumeIgnoreableTokens(tokens, next_index)
	if tokens[current_index].TokenType != BEGIN {
		return &AstSetPattern{expr, []AstProcessStatement{}}, current_index, err
	}

	statements, next_index, err := parse_process_statements(tokens, current_index+1)
	if err != nil {
		return nil, next_index, err
	}

	if tokens[next_index].TokenType != END {
		return nil, next_index, NewParseError(*tokens[next_index], "Unexpected token. Expected 'end'.")
	}

	return &AstSetPattern{expr, statements}, next_index + 1, err
}

func parse_set_matches(tokens []*Token, token_index int) (AstSetBody, int, error) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)
	command, next_index, err := parse_command(tokens, current_index)
	if err != nil {
		return nil, next_index, err
	}
	return &AstSetMatches{command}, next_index, err
}

func parse_amount(tokens []*Token, token_index int) (bool, int, int, int, int, error) {
	new_index := consumeIgnoreableTokens(tokens, token_index)

	if tokens[new_index].TokenType == ALL {
		new_index += 1
		return true, 0, 0, 0, new_index, nil
	} else if tokens[new_index].TokenType == SKIP {
		new_index = consumeIgnoreableTokens(tokens, new_index+1)
		if tokens[new_index].TokenType == NUMBER {
			skipValue, skipValueError := strconv.Atoi(tokens[new_index].Lexeme)
			if skipValueError != nil {
				return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Error converting to int value")
			}

			new_index = consumeIgnoreableTokens(tokens, new_index+1)
			if tokens[new_index].TokenType == TAKE {
				new_index = consumeIgnoreableTokens(tokens, new_index+1)
				if tokens[new_index].TokenType == NUMBER {
					takeValue, takeValueError := strconv.Atoi(tokens[new_index].Lexeme)
					if takeValueError != nil {
						return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Error converting to int value")
					}
					new_index++
					return false, skipValue, takeValue, 0, new_index, nil
				} else {
					return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
				}
			}
			return true, skipValue, 0, 0, new_index, nil
		} else {
			return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
		}
	} else if tokens[new_index].TokenType == TAKE || tokens[new_index].TokenType == TOP {
		new_index = consumeIgnoreableTokens(tokens, new_index+1)
		if tokens[new_index].TokenType == NUMBER {
			takeValue, takeValueError := strconv.Atoi(tokens[new_index].Lexeme)
			if takeValueError != nil {
				return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Error converting to int value")
			}
			new_index++
			return false, 0, takeValue, 0, new_index, nil
		} else {
			return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
		}
	} else if tokens[new_index].TokenType == LAST {
		new_index = consumeIgnoreableTokens(tokens, new_index+1)
		if tokens[new_index].TokenType == NUMBER {
			lastValue, lastValueError := strconv.Atoi(tokens[new_index].Lexeme)
			if lastValueError != nil {
				return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Error converting to int value")
			}
			new_index++
			return true, 0, 0, lastValue, new_index, nil
		} else {
			return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected a number")
		}
	}
	return false, 0, 0, 0, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'all', 'skip', or 'take'")
}

func parse_expression(tokens []*Token, token_index int) (AstExpression, int, error) {
	current_token := tokens[token_index]
	if current_token.TokenType == AT {
		return parse_at(tokens, token_index)
	} else if current_token.TokenType == BETWEEN {
		return parse_between(tokens, token_index)
	} else if current_token.TokenType == EXACTLY {
		return parse_exactly(tokens, token_index)
	} else if current_token.TokenType == MAYBE {
		return parse_maybe(tokens, token_index)
	} else if current_token.TokenType == IN {
		return parse_in(tokens, token_index, false)
	} else if current_token.TokenType == OPENCURLY {
		return parse_subroutine(tokens, token_index)
	} else if current_token.TokenType == NOT {
		return parse_not_expression(tokens, token_index)
	} else if current_token.TokenType == STRING || current_token.TokenType == IDENTIFIER ||
		current_token.TokenType == OPENPAREN || current_token.TokenType == ANY ||
		current_token.TokenType == WHITESPACE || current_token.TokenType == DIGIT ||
		current_token.TokenType == UPPER || current_token.TokenType == LOWER ||
		current_token.TokenType == LETTER || current_token.TokenType == LINE ||
		current_token.TokenType == FILE || current_token.TokenType == WORD ||
		current_token.TokenType == WHOLE {
		return parse_primary_or_dec(tokens, token_index)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected 'at', 'between', 'exactly', 'maybe', 'in', '<string>', '<identifier>', or a character class ")
}

func parse_at(tokens []*Token, token_index int) (*AstLoop, int, error) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[current_index]

	if current_token.TokenType != LEAST && current_token.TokenType != MOST {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'least' or 'most'.")
	}

	isLeast := current_token.TokenType == LEAST

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	if current_token.TokenType != NUMBER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected a number.")
	}
	value, err := strconv.Atoi(current_token.Lexeme)
	if err != nil {
		return nil, current_index, NewParseError(*current_token, "Error converting lexeme to number value")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	expr, next_index, parseError := parse_expression(tokens, current_index)
	if parseError != nil {
		return nil, next_index, parseError
	}

	current_index = consumeIgnoreableTokens(tokens, next_index)
	current_token = tokens[current_index]
	fewest := current_token.TokenType == FEWEST
	if fewest {
		current_index += 1
	}

	loopName := ""
	current_index = consumeIgnoreableTokens(tokens, current_index)
	current_token = tokens[current_index]
	if current_token.TokenType == NAMED {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		nameToken := tokens[current_index]
		if nameToken.TokenType == IDENTIFIER || nameToken.TokenType == STRING {
			loopName = nameToken.Lexeme
			current_index += 1
		} else {
			return nil, current_index, parseError
		}
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
		body:   expr,
		name:   loopName,
	}

	return &atLoop, current_index, nil
}

func parse_between(tokens []*Token, token_index int) (*AstLoop, int, error) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[current_index]

	if current_token.TokenType != NUMBER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected a number.")
	}
	minValue, err := strconv.Atoi(current_token.Lexeme)
	if err != nil {
		return nil, current_index, NewParseError(*current_token, "Error converting lexeme to number value")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	if current_token.TokenType != AND {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'and'.")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	if current_token.TokenType != NUMBER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected a number.")
	}
	maxValue, err := strconv.Atoi(current_token.Lexeme)
	if err != nil {
		return nil, current_index, NewParseError(*current_token, "Error converting lexeme to number value")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	expr, next_index, parseError := parse_expression(tokens, current_index)
	if parseError != nil {
		return nil, next_index, parseError
	}

	current_index = consumeIgnoreableTokens(tokens, next_index)
	current_token = tokens[current_index]
	fewest := current_token.TokenType == FEWEST
	if fewest {
		current_index += 1
	}

	loopName := ""
	current_index = consumeIgnoreableTokens(tokens, current_index)
	current_token = tokens[current_index]
	if current_token.TokenType == NAMED {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		nameToken := tokens[current_index]
		if nameToken.TokenType == IDENTIFIER || nameToken.TokenType == STRING {
			loopName = nameToken.Lexeme
			current_index += 1
		} else {
			return nil, current_index, NewParseError(*current_token, "Expected identifier following keyword 'named'")
		}
	}

	between := AstLoop{
		min:    minValue,
		max:    maxValue,
		fewest: fewest,
		body:   expr,
		name:   loopName,
	}

	return &between, current_index, nil
}

func parse_exactly(tokens []*Token, token_index int) (*AstLoop, int, error) {
	current_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[current_index]

	if current_token.TokenType != NUMBER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected a number.")
	}
	value, err := strconv.Atoi(current_token.Lexeme)
	if err != nil {
		return nil, current_index, NewParseError(*current_token, "Error converting lexeme to number value")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	expr, next_index, parseError := parse_expression(tokens, current_index)
	if parseError != nil {
		return nil, next_index, parseError
	}

	loopName := ""
	current_index = consumeIgnoreableTokens(tokens, current_index)
	current_token = tokens[current_index]
	if current_token.TokenType == NAMED {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		nameToken := tokens[current_index]
		if nameToken.TokenType == IDENTIFIER || nameToken.TokenType == STRING {
			loopName = nameToken.Lexeme
			current_index += 1
		} else {
			return nil, current_index, parseError
		}
	}

	exactly := AstLoop{
		min:    value,
		max:    value,
		fewest: false,
		body:   expr,
		name:   loopName,
	}

	return &exactly, next_index, nil
}

func parse_maybe(tokens []*Token, token_index int) (*AstLoop, int, error) {
	new_index := consumeIgnoreableTokens(tokens, token_index+1)
	expr, next_index, err := parse_expression(tokens, new_index)
	if err != nil {
		return nil, next_index, err
	}

	current_index := consumeIgnoreableTokens(tokens, next_index)
	current_token := tokens[current_index]
	fewest := current_token.TokenType == FEWEST
	if fewest {
		current_index += 1
	}

	maybe := AstLoop{0, 1, fewest, expr, ""}

	return &maybe, current_index, nil
}

func parse_not_expression(tokens []*Token, token_index int) (AstExpression, int, error) {
	new_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[new_index]

	if current_token.TokenType == IN {
		return parse_in(tokens, new_index, true)

		// TODO I think these 2 cases are unneccessary we can just call not_literal
		// TODO calling not_literal will help resolve the bug with not parsing this part correctly
		// I like adding comments :)
	} else if current_token.TokenType == STRING {
		ast_str, idx, err := parse_string(tokens, new_index, true)
		if err != nil {
			return nil, idx, err
		}
		return &AstPrimary{literal: ast_str}, idx, err
	} else if current_token.TokenType == ANY ||
		current_token.TokenType == WHITESPACE || current_token.TokenType == DIGIT ||
		current_token.TokenType == UPPER || current_token.TokenType == LOWER ||
		current_token.TokenType == LETTER || current_token.TokenType == LINE ||
		current_token.TokenType == FILE || current_token.TokenType == WORD ||
		current_token.TokenType == WHOLE {
		ast_chclass, idx, err := parse_character_class(tokens, new_index, true)
		if err != nil {
			return nil, idx, err
		}
		return &AstPrimary{literal: ast_chclass}, idx, err
	} else {
		return nil, new_index, NewParseError(*current_token, "Unexpected token. Expected 'in', <string>, <character class>")
	}
}

func parse_not_literal(tokens []*Token, token_index int) (AstLiteral, int, error) {
	new_index := consumeIgnoreableTokens(tokens, token_index+1)
	current_token := tokens[new_index]

	if current_token.TokenType == STRING {
		return parse_string(tokens, new_index, true)
	} else if current_token.TokenType == ANY ||
		current_token.TokenType == WHITESPACE || current_token.TokenType == DIGIT ||
		current_token.TokenType == UPPER || current_token.TokenType == LOWER ||
		current_token.TokenType == LETTER || current_token.TokenType == LINE ||
		current_token.TokenType == FILE || current_token.TokenType == WORD ||
		current_token.TokenType == WHOLE {
		return parse_character_class(tokens, new_index, true)
	} else {
		return nil, new_index, NewParseError(*current_token, "Unexpected token. Expected 'in', <string>, <character class>")
	}
}

func parse_in(tokens []*Token, token_index int, not bool) (*AstList, int, error) {
	new_index := consumeIgnoreableTokens(tokens, token_index+1)
	contents := []AstListable{}

	listable, next_index, err := parse_listable(tokens, new_index)
	if err != nil {
		return nil, next_index, err
	}
	contents = append(contents, listable)
	current_index := consumeIgnoreableTokens(tokens, next_index)
	current_token := tokens[current_index]
	for current_token.TokenType == COMMA {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		listable, next_index, err := parse_listable(tokens, current_index)
		if err != nil {
			return nil, next_index, err
		}
		contents = append(contents, listable)
		current_index = next_index
		current_token = tokens[current_index]
	}
	inList := AstList{contents: contents, not: not}
	return &inList, current_index, nil
}

func isListableClass(t TokenType) bool {
	return t == ANY || t == WHITESPACE || t == DIGIT || t == UPPER || t == LOWER || t == LETTER
}

func parse_listable(tokens []*Token, token_index int) (AstListable, int, error) {
	current_token := tokens[token_index]
	if current_token.TokenType == STRING {
		from, next_index, err := parse_string(tokens, token_index, false)
		if err != nil {
			return nil, next_index, err
		}

		current_index := consumeIgnoreableTokens(tokens, next_index)
		current_token := tokens[current_index]
		if current_token.TokenType != TO {
			return from, current_index, nil
		}

		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		to, new_index, terr := parse_string(tokens, current_index, false)
		if terr != nil {
			return nil, new_index, terr
		}

		r := AstRange{
			from: from,
			to:   to,
		}
		return &r, new_index, nil

	} else if isListableClass(current_token.TokenType) {
		return parse_character_class(tokens, token_index, false)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected listable literal")
}

func parse_literal(tokens []*Token, token_index int) (AstLiteral, int, error) {
	current_token := tokens[token_index]
	if current_token.TokenType == STRING {
		return parse_string(tokens, token_index, false)
	} else if current_token.TokenType == IDENTIFIER {
		return parse_variable(tokens, token_index)
	} else if current_token.TokenType == OPENPAREN {
		return parse_sub_expression(tokens, token_index)
	} else if current_token.TokenType == NOT {
		return parse_not_literal(tokens, token_index)
	} else if current_token.TokenType == ANY ||
		current_token.TokenType == WHITESPACE || current_token.TokenType == DIGIT ||
		current_token.TokenType == UPPER || current_token.TokenType == LOWER ||
		current_token.TokenType == LETTER || current_token.TokenType == LINE ||
		current_token.TokenType == FILE || current_token.TokenType == WORD ||
		current_token.TokenType == WHOLE {
		return parse_character_class(tokens, token_index, false)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected '(', '<string>', '<identifier>', or a character class.")
}

func parse_primary_or_dec(tokens []*Token, token_index int) (AstExpression, int, error) {
	literal, new_index, err := parse_literal(tokens, token_index)
	if err != nil {
		return nil, new_index, err
	}

	current_index := consumeIgnoreableTokens(tokens, new_index)
	current_token := tokens[current_index]
	if current_token.TokenType == EQUAL {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		current_token = tokens[current_index]

		if current_token.TokenType != IDENTIFIER {
			return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected identifier.")
		}

		dec := AstDec{
			name: current_token.Lexeme,
			body: literal,
		}
		return &dec, current_index + 1, nil
	}

	if current_token.TokenType == OR {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		right_expression, final_index, err := parse_primary_or_or(tokens, current_index)
		if err != nil {
			return nil, final_index, err
		}

		branch := AstBranch{
			left:  literal,
			right: right_expression,
		}
		return &branch, final_index, nil
	}

	prim := AstPrimary{}
	prim.literal = literal

	return &prim, new_index, nil
}

func parse_primary_or_or(tokens []*Token, token_index int) (AstExpression, int, error) {
	literal, new_index, err := parse_literal(tokens, token_index)
	if err != nil {
		return nil, new_index, err
	}

	current_index := consumeIgnoreableTokens(tokens, new_index)
	current_token := tokens[current_index]
	if current_token.TokenType == OR {
		current_index = consumeIgnoreableTokens(tokens, current_index+1)
		right_expression, final_index, err := parse_primary_or_or(tokens, current_index)
		if err != nil {
			return nil, final_index, err
		}

		branch := AstBranch{
			left:  literal,
			right: right_expression,
		}
		return &branch, final_index, nil
	}

	prim := AstPrimary{}
	prim.literal = literal

	return &prim, new_index, nil
}

func parse_atom(tokens []*Token, token_index int) (AstAtom, int, error) {
	current_token := tokens[token_index]
	if current_token.TokenType == STRING {
		return parse_string(tokens, token_index, false)
	} else if current_token.TokenType == IDENTIFIER {
		return parse_variable(tokens, token_index)
	}
	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected '<string>' or '<identifier>'.")
}

func parse_string(tokens []*Token, token_index int, not bool) (*AstString, int, error) {
	current_token := tokens[token_index]
	str_literal := AstString{
		not: not,
	}

	if current_token.TokenType == STRING {
		str_literal.value = current_token.Lexeme
		return &str_literal, token_index + 1, nil
	}

	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected a string")
}

func parse_variable(tokens []*Token, token_index int) (*AstVariable, int, error) {
	current_token := tokens[token_index]
	var_literal := AstVariable{}

	if current_token.TokenType == IDENTIFIER {
		var_literal.name = current_token.Lexeme
		return &var_literal, token_index + 1, nil
	}

	return nil, token_index, NewParseError(*current_token, "Unexpected token. Expected a variable")
}

func parse_sub_expression(tokens []*Token, token_index int) (*AstSubExpr, int, error) {
	current_token := tokens[token_index+1]
	current_index := token_index + 1
	expr_list := []AstExpression{}

	for current_token.TokenType != CLOSEPAREN && current_token.TokenType != FIND && current_token.TokenType != REPLACE && current_token.TokenType != SET && current_token.TokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_expression(tokens, ws_index)
		if parseError != nil {
			return nil, new_index, parseError
		}

		expr_list = append(expr_list, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	if current_token.TokenType != CLOSEPAREN {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected ')'")
	}

	sub_expr := AstSubExpr{body: expr_list}
	return &sub_expr, current_index + 1, nil
}

func parse_subroutine(tokens []*Token, token_index int) (*AstSub, int, error) {
	current_token := tokens[token_index+1]
	current_index := token_index + 1
	expr_list := []AstExpression{}

	for current_token.TokenType != CLOSECURLY && current_token.TokenType != FIND && current_token.TokenType != REPLACE && current_token.TokenType != SET && current_token.TokenType != EOF {
		ws_index := consumeIgnoreableTokens(tokens, current_index)
		expr, new_index, parseError := parse_expression(tokens, ws_index)
		if parseError != nil {
			return nil, new_index, parseError
		}

		expr_list = append(expr_list, expr)
		current_index = consumeIgnoreableTokens(tokens, new_index)
		current_token = tokens[current_index]
	}

	if current_token.TokenType != CLOSECURLY {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected '}'")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]
	if current_token.TokenType != EQUAL {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected '='")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]

	if current_token.TokenType != IDENTIFIER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected identifier.")
	}

	dec := AstSub{
		name: current_token.Lexeme,
		body: expr_list,
	}
	return &dec, current_index + 1, nil
}

func parse_character_class(tokens []*Token, token_index int, not bool) (*AstCharacterClass, int, error) {
	current_token := tokens[token_index]
	charClass := AstCharacterClass{
		not: not,
	}
	if current_token.TokenType == ANY {
		charClass.classType = ClassAny
		return &charClass, token_index + 1, nil
	} else if current_token.TokenType == WHITESPACE {
		charClass.classType = ClassWhitespace
		return &charClass, token_index + 1, nil
	} else if current_token.TokenType == DIGIT {
		charClass.classType = ClassDigit
		return &charClass, token_index + 1, nil
	} else if current_token.TokenType == UPPER {
		charClass.classType = ClassUpper
		return &charClass, token_index + 1, nil
	} else if current_token.TokenType == LOWER {
		charClass.classType = ClassLower
		return &charClass, token_index + 1, nil
	} else if current_token.TokenType == LETTER {
		charClass.classType = ClassLetter
		return &charClass, token_index + 1, nil
	} else if current_token.TokenType == LINE {
		new_index := consumeIgnoreableTokens(tokens, token_index+1)
		if tokens[new_index].TokenType == START {
			charClass.classType = ClassLineStart
			return &charClass, new_index + 1, nil
		} else if tokens[new_index].TokenType == END {
			charClass.classType = ClassLineEnd
			return &charClass, new_index + 1, nil
		}
		return nil, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'start' or 'end'")
	} else if current_token.TokenType == FILE {
		new_index := consumeIgnoreableTokens(tokens, token_index+1)
		if tokens[new_index].TokenType == START {
			charClass.classType = ClassFileStart
			return &charClass, new_index + 1, nil
		} else if tokens[new_index].TokenType == END {
			charClass.classType = ClassFileEnd
			return &charClass, new_index + 1, nil
		}
		return nil, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'start' or 'end'")
	} else if current_token.TokenType == WORD {
		new_index := consumeIgnoreableTokens(tokens, token_index+1)
		if tokens[new_index].TokenType == START {
			charClass.classType = ClassWordStart
			return &charClass, new_index + 1, nil
		} else if tokens[new_index].TokenType == END {
			charClass.classType = ClassWordEnd
			return &charClass, new_index + 1, nil
		}
		return nil, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'start' or 'end'")
	} else if current_token.TokenType == WHOLE {
		new_index := consumeIgnoreableTokens(tokens, token_index+1)
		if tokens[new_index].TokenType == LINE {
			charClass.classType = ClassWholeLine
			return &charClass, new_index + 1, nil
		} else if tokens[new_index].TokenType == FILE {
			charClass.classType = ClassWholeFile
			return &charClass, new_index + 1, nil
		} else if tokens[new_index].TokenType == WORD {
			charClass.classType = ClassWholeWord
			return &charClass, new_index + 1, nil
		}
		return nil, new_index, NewParseError(*tokens[new_index], "Unexpected token. Expected 'file', 'line', or 'word'")
	}
	return nil, token_index, NewParseError(*tokens[token_index], "Unexpected token. Expected a character class: 'any', 'whitespace', 'digit', 'upper', 'lower', 'letter', 'whole word', 'whole line', 'whole file', 'word start', 'word end', 'line start', 'line end', 'file start', or 'file end'.")
}

func parse_process_statements(tokens []*Token, index int) ([]AstProcessStatement, int, error) {
	statements := []AstProcessStatement{}
	token_index := index
	for token_index < len(tokens)-1 {
		ws_index := consumeIgnoreableTokens(tokens, token_index)
		command, new_index, e := parse_process_statement(tokens, ws_index)
		if e != nil {
			return nil, new_index, e
		}
		token_index = new_index
		if command != nil {
			statements = append(statements, command)
		} else {
			break
		}
	}

	return statements, token_index, nil
}

func parse_process_statement(tokens []*Token, index int) (AstProcessStatement, int, error) {
	if tokens[index].TokenType == SET {
		return parse_process_set(tokens, index)
	} else if tokens[index].TokenType == IF {
		return parse_process_if(tokens, index)
	} else if tokens[index].TokenType == RETURN {
		return parse_process_return(tokens, index)
	} else if tokens[index].TokenType == DEBUG {
		return parse_process_debug(tokens, index)
	} else if tokens[index].TokenType == LOOP {
		return parse_process_loop(tokens, index)
	} else if tokens[index].TokenType == BREAK {
		return AstProcessBreak{}, index + 1, nil
	} else if tokens[index].TokenType == CONTINUE {
		return AstProcessContinue{}, index + 1, nil
	} else if tokens[index].TokenType == END {
		return nil, index, nil
	} else if tokens[index].TokenType == ELSE {
		return nil, index, nil
	} else {
		return nil, index, NewParseError(*tokens[index], "Unexpected token. Expected 'set', 'if', 'return', 'debug', 'loop', 'else', or 'end'.")
	}
}

func parse_process_set(tokens []*Token, index int) (AstProcessStatement, int, error) {
	var current_index = consumeIgnoreableTokens(tokens, index+1)
	var current_token = tokens[current_index]

	if current_token.TokenType != IDENTIFIER {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected identifier")
	}

	name := current_token.Lexeme

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	current_token = tokens[current_index]

	if current_token.TokenType != TO {
		return nil, current_index, NewParseError(*current_token, "Unexpected token. Expected 'to'")
	}

	current_index = consumeIgnoreableTokens(tokens, current_index+1)
	expr, next_index, err := parse_process_expression(tokens, current_index)
	if err != nil {
		return nil, next_index, err
	}

	setStatement := AstProcessSet{
		name: name,
		expr: expr,
	}
	return &setStatement, next_index, nil
}

func parse_process_if(tokens []*Token, index int) (AstProcessStatement, int, error) {
	current_index := consumeIgnoreableTokens(tokens, index+1)
	expr, next_index, err := parse_process_expression(tokens, current_index)
	if err != nil {
		return nil, next_index, err
	}

	next_index = consumeIgnoreableTokens(tokens, next_index)
	if tokens[next_index].TokenType != THEN {
		return nil, next_index, NewParseError(*tokens[next_index], "Unexpected token. Expected 'then'.")
	}

	next_index = consumeIgnoreableTokens(tokens, next_index+1)
	trueBody, follow_index, err := parse_process_statements(tokens, next_index)
	if err != nil {
		return nil, follow_index, err
	}

	falseBody := []AstProcessStatement{}
	if tokens[follow_index].TokenType == ELSE {
		falseBody, follow_index, err = parse_process_statements(tokens, follow_index+1)
		if err != nil {
			return nil, follow_index, err
		}
	}

	follow_index = consumeIgnoreableTokens(tokens, follow_index)
	if tokens[follow_index].TokenType != END {
		return nil, follow_index, NewParseError(*tokens[follow_index], "Unexpected token. Expected 'end'.")
	}

	return &AstProcessIf{expr, trueBody, falseBody}, follow_index + 1, nil
}

func parse_process_return(tokens []*Token, index int) (AstProcessStatement, int, error) {
	current_index := consumeIgnoreableTokens(tokens, index+1)
	expr, next_index, err := parse_process_expression(tokens, current_index)
	if err != nil {
		return nil, next_index, err
	}

	return &AstProcessReturn{expr}, next_index, err
}

func parse_process_debug(tokens []*Token, index int) (AstProcessStatement, int, error) {
	current_index := consumeIgnoreableTokens(tokens, index+1)
	expr, next_index, err := parse_process_expression(tokens, current_index)
	if err != nil {
		return nil, next_index, err
	}

	return &AstProcessDebug{expr}, next_index, err
}

func parse_process_loop(tokens []*Token, index int) (AstProcessStatement, int, error) {
	current_index := consumeIgnoreableTokens(tokens, index+1)

	body, next_index, err := parse_process_statements(tokens, current_index)
	if err != nil {
		return nil, next_index, err
	}

	next_index = consumeIgnoreableTokens(tokens, next_index)
	if tokens[next_index].TokenType != END {
		return nil, next_index, NewParseError(*tokens[next_index], "Unexpected token. Expected 'end'.")
	}

	return &AstProcessLoop{body}, next_index + 1, nil
}

func parse_process_expression(tokens []*Token, index int) (AstProcessExpression, int, error) {
	exprTokens, next_index := getProcessExpressionTokens(tokens, index)
	expr, fail_index, err := parse_expr_pratt(exprTokens, 0, 0)
	if err != nil {
		return nil, index + fail_index, err
	}
	return expr, next_index, nil
}

func parse_expr_pratt(tokens []*Token, index int, minPrecedence int) (AstProcessExpression, int, error) {
	token_index := index + 1
	var lhs AstProcessExpression
	if tokens[index].TokenType == STRING {
		lhs = AstProcessString{tokens[index].Lexeme}
	} else if tokens[index].TokenType == TRUE {
		lhs = AstProcessBoolean{true}
	} else if tokens[index].TokenType == FALSE {
		lhs = AstProcessBoolean{false}
	} else if tokens[index].TokenType == NUMBER {
		intval, err := strconv.Atoi(tokens[index].Lexeme)
		if err != nil {
			intval = 0
		}
		lhs = AstProcessNumber{intval}
	} else if tokens[index].TokenType == IDENTIFIER {
		lhs = AstProcessVariable{tokens[index].Lexeme}
	} else if tokens[index].TokenType == OPENPAREN {
		subexpr, next_index, err := parse_expr_pratt(tokens, index+1, 0)
		if err != nil {
			return nil, next_index, err
		}
		if tokens[next_index].TokenType != CLOSEPAREN {
			return nil, next_index, err
		}
		token_index = next_index + 1
		lhs = subexpr
	} else if isPrefixOp(tokens[index].TokenType) {
		rprec := prefixPrecedence(tokens[index].TokenType)
		rhs, next_index, err := parse_expr_pratt(tokens, index+1, rprec)
		if err != nil {
			return nil, next_index, err
		}
		lhs = AstProcessUnaryExpression{tokens[index].TokenType, rhs}
		token_index = next_index
	} else {
		return nil, index, NewParseError(*tokens[index], "Unexpected token. Expected string, number, variable, or unary operator")
	}

	for token_index < len(tokens) {
		if tokens[token_index].TokenType == CLOSEPAREN {
			break
		}
		if !isBinaryOp(tokens[token_index].TokenType) {
			return nil, token_index, NewParseError(*tokens[token_index], "Unexpected token. Expected binary operator.")
		}
		op := tokens[token_index].TokenType
		lprec, rprec := infixPrecedence(tokens[token_index].TokenType)
		if lprec < minPrecedence {
			break
		}
		rhs, next_index, err := parse_expr_pratt(tokens, token_index+1, rprec)
		if err != nil {
			return nil, next_index, err
		}
		token_index = next_index
		lhs = AstProcessBinaryExpression{op, lhs, rhs}
	}

	return lhs, token_index, nil
}

func getProcessExpressionTokens(tokens []*Token, index int) ([]*Token, int) {
	exprTokens := []*Token{}
	token_index := index
	for token_index < len(tokens) {
		if isProcessExprEnd(tokens[token_index].TokenType) {
			break
		} else if tokens[token_index].TokenType == WS {
			token_index += 1
		} else {
			exprTokens = append(exprTokens, tokens[token_index])
			token_index += 1
		}
	}
	return exprTokens, token_index
}

func isProcessExprEnd(tokenType TokenType) bool {
	return tokenType == SET || tokenType == THEN || tokenType == IF || tokenType == ELSE || tokenType == END || tokenType == DEBUG || tokenType == RETURN || tokenType == LOOP
}

func isPrefixOp(tokenType TokenType) bool {
	return tokenType == NOT || tokenType == HEAD || tokenType == TAIL
}

func prefixPrecedence(tokenType TokenType) int {
	if tokenType == NOT {
		return 11
	} else if tokenType == HEAD || tokenType == TAIL {
		return 12
	}
	return -1
}

func isBinaryOp(tokenType TokenType) bool {
	return tokenType == AND || tokenType == OR || tokenType == PLUS || tokenType == MINUS || tokenType == MOD || tokenType == MULT || tokenType == DIV || tokenType == LESS || tokenType == GREATER || tokenType == LESSEQ || tokenType == GREATEREQ || tokenType == DEQUAL || tokenType == NEQUAL
}

func infixPrecedence(tokenType TokenType) (int, int) {
	if tokenType == AND || tokenType == OR {
		return 1, 2
	} else if tokenType == DEQUAL || tokenType == NEQUAL {
		return 3, 4
	} else if tokenType == LESS || tokenType == GREATER || tokenType == LESSEQ || tokenType == GREATEREQ {
		return 5, 6
	} else if tokenType == PLUS || tokenType == MINUS {
		return 7, 8
	} else if tokenType == MULT || tokenType == DIV || tokenType == MOD {
		return 9, 10
	}
	return -1, -1
}

func consumeIgnoreableTokens(tokens []*Token, index int) int {
	current_index := index
	for tokens[current_index].TokenType == WS || tokens[current_index].TokenType == COMMENT {
		current_index += 1
	}
	return current_index
}
