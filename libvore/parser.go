package libvore

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
	return ParseError{isError: true}
}

func parse(tokens []*Token) ([]*AstCommand, ParseError) {
	commands := []*AstCommand{}

	token_index := 0
	for token_index < len(tokens) {
		token_index = consumeIgnoreableTokens(tokens, token_index)
		command, new_index, e := parse_command(tokens, token_index)
		if e.isError {
			return []*AstCommand{}, e
		}
		token_index = new_index
		commands = append(commands, &command)
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
	return nil, token_index, NoError()
}

func parse_replace(tokens []*Token, token_index int) (*AstReplace, int, ParseError) {
	return nil, token_index, NoError()
}

func parse_set(tokens []*Token, token_index int) (*AstSet, int, ParseError) {
	return nil, token_index, NoError()
}

func parse_amount(tokens []*Token, token_index int) (bool, int, int, int, ParseError) {
	new_index := consumeIgnoreableTokens(tokens, token_index)

	if tokens[new_index].tokenType == ALL {
		new_index += 1
		return true, 0, 0, new_index, NoError()
	} else if tokens[new_index].tokenType == SKIP {

	} else if tokens[new_index].tokenType == TAKE {

	}
	return false, 0, 0, new_index, NewParseError(*tokens[new_index], ":(")
}

func consumeIgnoreableTokens(tokens []*Token, index int) int {
	current_index := index
	for tokens[current_index].tokenType != WS && tokens[current_index].tokenType != COMMENT {
		current_index += 1
	}
	return current_index
}
