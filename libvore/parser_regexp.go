package libvore

import (
	"strconv"
)

// This is mostly built off of the ECMAScript spec for regular expressions
// https://262.ecma-international.org/13.0/#sec-patterns
// this was mostly because it was the easiest to find the grammar

/*

PATTERN :: DISJUNCTION

DISJUNCTION :: ALTERNATIVE
			:: ALTERNATIVE | DISJUNCTION

ALTERNATIVE :: [empty]
			:: ALTERNATIVE TERM

TERM :: ASSERTION
	 :: ATOM
	 :: ATOM QUANTIFIER

ASSERTION :: ^
		  :: $
		  :: \ b
		  :: \ B
		  :: ( ? = DISJUNCTION )
		  :: ( ? ! DISJUNCTION )
		  :: ( ? < = DISJUNCTION )
		  :: ( ? < ! DISJUNCTION )

QUANTIFIER :: QUANTIFIERPREFIX
		   :: QUANTIFIERPREFIX ?

QUANTIFIERPREFIX :: *
			     :: +
				 :: ?
				 :: { DECIMALDIGITS }
				 :: { DECIMALDIGITS , }
				 :: { DECIMALDIGITS , DECIMALDIGITS }

ATOM :: PATTERNCHARACTER
	 :: .
	 :: \ ATOMESCAPE
	 :: CHARACTERCLASS
	 :: ( GROUPSPECIFIER DISJUNCTION )
	 :: ( ? : DISJUNCTION )

SYNTAXCHARACTER :: one of ^ $ \ . * + ? ( ) [ ] { }

PATTERNCHARACTER :: any single character except for SYNTAXCHARACTER

ATOMESCAPE :: unimplemented

CHARACTERCLASS :: unimplemented

GROUPSPECIFIER :: [empty]
			   :: ? GROUPNAME

GROUPNAME :: < IDENTIFIER >     -- will break ECMAScript standard here and just do the same style identifiers in vore

*/

func parse_regexp(tokens []*Token, token_index int) (AstExpression, int, error) {
	regexp_token := tokens[token_index]
	regexp := regexp_token.Lexeme

	index := 0
	results := []AstExpression{}
	for index < len(regexp) {
		exp, next_index, err := parse_regexp_pattern(regexp_token, regexp, index)
		if err != nil {
			return nil, next_index, err
		}
		results = append(results, exp)
		index = next_index
	}

	s := &AstPrimary{
		&AstSubExpr{results},
	}

	return s, token_index + 1, nil
}

func parse_regexp_pattern(regexp_token *Token, regexp string, index int) (AstExpression, int, error) {

	start, next_index, err := parse_regexp_literal(regexp_token, regexp, index)
	if err != nil {
		return nil, next_index, err
	}

	if next_index < len(regexp) {
		if regexp[next_index] == '|' {
			end, idx, err := parse_regexp_pattern(regexp_token, regexp, next_index+1)
			return &AstBranch{&AstSubExpr{[]AstExpression{start}}, end}, idx, err
		} else {
			return start, next_index, nil
		}
	} else {
		return start, next_index, nil
	}
}

func parse_regexp_number(regexp_token *Token, regexp string, index int) (int, int, error) {
	c := regexp[index]
	result := ""
	idx := index
	for c >= '0' && c <= '9' {
		result += string(c)
		idx += 1
		c = regexp[idx]
	}
	if result == "" {
		return -1, index, NewParseError(*regexp_token, "Unexpected Token. Expected number")
	}
	value, err := strconv.Atoi(result)
	if err != nil {
		return value, idx, NewParseError(*regexp_token, "Error converting string to number")
	}
	return value, idx, nil
}

func parse_regexp_literal(regexp_token *Token, regexp string, index int) (AstExpression, int, error) {
	c := regexp[index]
	var start AstLiteral
	next_index := index
	if c == '^' {
		start = &AstCharacterClass{false, ClassLineStart}
		next_index += 1
		return &AstPrimary{start}, next_index, nil
	} else if c == '$' {
		start = &AstCharacterClass{false, ClassLineEnd}
		next_index += 1
		return &AstPrimary{start}, next_index, nil
	} else if c == '\\' {
		exp, idx, err := parse_regexp_escape_characters(regexp_token, regexp, index+1)
		if err != nil {
			return nil, index, err
		}
		return &AstPrimary{exp}, idx, nil
	} else if c == '(' {
		start, next_index, err := parse_regexp_groups(regexp_token, regexp, index+1)
		if err != nil {
			return nil, next_index, err
		}
		exp, idx, err := parse_regexp_quantifier(regexp_token, regexp, next_index)
		if err != nil {
			return nil, idx, err
		}
		if exp == nil {
			return &AstPrimary{start}, idx, nil
		}
		exp.body = &AstPrimary{start}
		return exp, idx, nil
	} else if c == '.' {
		start = &AstString{true, "\n", false}
		next_index += 1
		exp, idx, err := parse_regexp_quantifier(regexp_token, regexp, next_index)
		if err != nil {
			return nil, idx, err
		}
		if exp == nil {
			return &AstPrimary{start}, idx, nil
		}
		exp.body = &AstPrimary{start}
		return exp, idx, nil
	} else {
		start = &AstString{false, string(c), false}
		next_index += 1
		exp, idx, err := parse_regexp_quantifier(regexp_token, regexp, next_index)
		if err != nil {
			return nil, idx, err
		}
		if exp == nil {
			return &AstPrimary{start}, idx, nil
		}
		exp.body = &AstPrimary{start}
		return exp, idx, nil
	}
}

func parse_regexp_quantifier(regexp_token *Token, regexp string, index int) (*AstLoop, int, error) {
	if index >= len(regexp) {
		return nil, index, nil
	}
	op := regexp[index]
	var end_idx int
	var exp *AstLoop
	if op == '*' {
		exp = &AstLoop{0, -1, false, nil, ""}
		end_idx = index + 1
	} else if op == '+' {
		exp = &AstLoop{1, -1, false, nil, ""}
		end_idx = index + 1
	} else if op == '?' {
		exp = &AstLoop{0, 1, false, nil, ""}
		end_idx = index + 1
	} else if op == '{' {
		from, idx, err := parse_regexp_number(regexp_token, regexp, index+1)
		if err != nil {
			return nil, idx, err
		}
		comma_or_brace := regexp[idx]

		if comma_or_brace == ',' {
			if regexp[idx+1] == '}' {
				exp = &AstLoop{from, -1, false, nil, ""}
				end_idx = idx + 2
			} else {
				to, idx2, err := parse_regexp_number(regexp_token, regexp, idx+1)
				if err != nil {
					return nil, idx, err
				}
				brace := regexp[idx2]
				if brace != '}' {
					return nil, idx2, NewParseError(*regexp_token, "Unexpected character. Expected '}'")
				}

				exp = &AstLoop{from, to, false, nil, ""}
				end_idx = idx2 + 1
			}
		} else if comma_or_brace == '}' {
			exp = &AstLoop{from, from, false, nil, ""}
			end_idx = idx + 1
		}
	} else {
		exp = nil
		end_idx = index
	}
	exp.fewest = end_idx < len(regexp) && regexp[end_idx] == '?'
	if exp.fewest {
		end_idx += 1
	}
	return exp, end_idx, nil
}

func parse_regexp_escape_characters(regexp_token *Token, regexp string, index int) (AstLiteral, int, error) {
	return nil, index, nil
}

func parse_regexp_groups(regexp_token *Token, regexp string, index int) (AstLiteral, int, error) {
	return nil, index, nil
}
