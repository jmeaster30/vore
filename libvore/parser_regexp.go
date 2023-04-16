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

	if next_index >= len(regexp) {
		return &AstPrimary{start}, next_index, nil
	}

	op := regexp[next_index]
	if op == '*' {
		fewest := next_index+1 < len(regexp) && regexp[next_index+1] == '?'
		exp := &AstLoop{0, -1, fewest, &AstPrimary{start}, ""}
		if fewest {
			return exp, next_index + 2, nil
		} else {
			return exp, next_index + 1, nil
		}
	} else if op == '+' {
		fewest := next_index+1 < len(regexp) && regexp[next_index+1] == '?'
		exp := &AstLoop{1, -1, fewest, &AstPrimary{start}, ""}
		if fewest {
			return exp, next_index + 2, nil
		} else {
			return exp, next_index + 1, nil
		}
	} else if op == '?' {
		fewest := next_index+1 < len(regexp) && regexp[next_index+1] == '?'
		exp := &AstLoop{0, 1, fewest, &AstPrimary{start}, ""}
		if fewest {
			return exp, next_index + 2, nil
		} else {
			return exp, next_index + 1, nil
		}
	} else if op == '{' {
		from, idx, err := parse_regexp_number(regexp_token, regexp, next_index+1)
		if err != nil {
			return nil, idx, err
		}
		comma_or_brace := regexp[idx]
		var end_idx int
		var exp *AstLoop
		if comma_or_brace == ',' {
			if regexp[idx+1] == '}' {
				exp = &AstLoop{from, -1, false, &AstPrimary{start}, ""}
				end_idx = idx + 1
			} else {
				to, idx2, err := parse_regexp_number(regexp_token, regexp, idx+1)
				if err != nil {
					return nil, idx, err
				}
				brace := regexp[idx2]
				if brace != '}' {
					return nil, idx2, NewParseError(*regexp_token, "Unexpected character. Expected '}'")
				}

				exp = &AstLoop{from, to, false, &AstPrimary{start}, ""}
				end_idx = idx2 + 1
			}
		} else if comma_or_brace == '}' {
			exp = &AstLoop{from, from, false, &AstPrimary{start}, ""}
			end_idx = idx + 1
		}
		exp.fewest = end_idx+1 < len(regexp) && regexp[end_idx+1] == '?'
		if exp.fewest {
			return exp, end_idx + 2, nil
		} else {
			return exp, end_idx + 1, nil
		}
	} else if op == '|' {
		end, idx, err := parse_regexp_pattern(regexp_token, regexp, next_index+1)
		return &AstBranch{start, end}, idx, err
	}

	return &AstPrimary{start}, next_index, nil
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

func parse_regexp_literal(regexp_token *Token, regexp string, index int) (AstLiteral, int, error) {
	c := regexp[index]
	var start AstLiteral
	next_index := index
	if c == '^' {
		start = &AstCharacterClass{false, ClassLineStart}
		next_index += 1
	} else if c == '$' {
		start = &AstCharacterClass{false, ClassLineEnd}
		next_index += 1
	} else if c == '\\' {
		exp, idx, err := parse_regexp_escape_characters(regexp_token, regexp, index+1)
		if err != nil {
			return nil, index, err
		}
		start = exp
		next_index = idx
	} else if c == '(' {
		exp, idx, err := parse_regexp_groups(regexp_token, regexp, index+1)
		if err != nil {
			return nil, index, err
		}
		start = exp
		next_index = idx
	} else if c == '.' {
		start = &AstString{true, "\n", false}
		next_index += 1
	} else {
		start = &AstString{false, string(c), false}
		next_index += 1
	}
	return start, next_index, nil
}

func parse_regexp_escape_characters(regexp_token *Token, regexp string, index int) (AstLiteral, int, error) {
	return nil, index, nil
}

func parse_regexp_groups(regexp_token *Token, regexp string, index int) (AstLiteral, int, error) {
	return nil, index, nil
}
