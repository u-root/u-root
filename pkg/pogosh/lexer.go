// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

import (
	"fmt"
	"strings"
)

type token struct {
	value string // TODO: make []byte
	ttype tokenType
}

type tokenType uint8

// Token types
const (
	ttError = iota // TODO: is this used?
	ttEOF
	ttWord
	ttAssignmentWord
	ttName
	ttNewLine
	ttIONumber
	ttAndIf     // &&
	ttOrIf      // ||
	ttDSemi     // ;;
	ttDLess     // <<
	ttDGreat    // >>
	ttLessAnd   // <&
	ttGreatAnd  // >&
	ttLessGreat // <>
	ttDLessDash // <<-
	ttClobber   // >|
	ttIf        // if
	ttThen      // then
	ttElse      // else
	ttElif      // elif
	ttFi        // fi
	ttDo        // do
	ttDone      // done
	ttCase      // case
	ttEsac      // esac
	ttWhile     // while
	ttUntil     // until
	ttFor       // for
	ttLBrace    // {
	ttRBrace    // }
	ttBang      // !
	ttIn        // in
)

var operators = map[string]tokenType{
	"&&":  ttAndIf,
	"||":  ttOrIf,
	";;":  ttDSemi,
	"<<":  ttDLess,
	">>":  ttDGreat,
	"<&":  ttLessAnd,
	">&":  ttGreatAnd,
	"<>":  ttLessGreat,
	"<<-": ttDLessDash,
	">|":  ttClobber,
}

var reservedWords = map[string]tokenType{
	"if":    ttIf,
	"then":  ttThen,
	"else":  ttElse,
	"elif":  ttElif,
	"fi":    ttFi,
	"do":    ttDo,
	"done":  ttDone,
	"case":  ttCase,
	"esac":  ttEsac,
	"while": ttWhile,
	"until": ttUntil,
	"for":   ttFor,
	"{":     ttLBrace,
	"}":     ttRBrace,
	"!":     ttBang,
	"in":    ttIn,
}

var portableCharSet = "\x00\a\b\t\n\v\f\r !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxy{|}~"

// tokenize splits the input into an array of tokens.
// TODO: memoize?
func tokenize(script string) ([]token, error) {
	ts := []token{}
	b := 0 // Beginning of current token
	i := 0 // Index of current character

	// Tokenizer states
	const (
		sStart = iota
		sEscape
		sOperator
		sWord
		sWordEscape
		sSingleQuote
		sDoubleQuote
		sDoubleQuoteEscape
		sLineComment
		sLineCommentEscape
	)
	state := sStart

	// Iterate over each character + an imaginary blank character.
	for {
		// Current character being processed
		var c byte

		// Check for EOF
		if i == len(script) {
			switch state {
			case sStart, sOperator, sWord, sLineComment:
				// Use an imaginary blank character to delimit the last token.
				c = ' '
			default:
				return ts, fmt.Errorf("INCOMPLETE") // TODO
			}
		} else {
			c = script[i]
		}

		// The scanner is implemented with a DFA:
		// * outer switch -- selects state
		// * inner switch -- selects transition
		switch state {

		// Start state
		case sStart:
			b = i
			// TODO: \r
			switch c {
			default:
				state = sWord
			case ' ', '\t':
				state = sStart
			case '\n':
				ts = append(ts, token{script[i : i+1], ttNewLine})
			case '\\':
				state = sEscape
			case '\'':
				state = sSingleQuote
			case '"':
				state = sDoubleQuote
			case '#':
				state = sLineComment
			case '&', '|', ';', '<', '>':
				state = sOperator
			}

		// Escape
		case sEscape:
			switch c {
			case '\n':
				state = sStart
			default:
				state = sWord
			}

		// Words
		case sWord:
			switch c {
			case ' ', '\t', '\n', '#', '&', '|', ';', '<', '>':
				// The token may contain a line escape. This is cleaned up
				// during variable expansion.
				ts = append(ts, token{script[b:i], ttWord})
				state = sStart
				i--
			case '\\':
				state = sWordEscape
			}
		case sWordEscape:
			state = sWord

		// Single quotes
		case sSingleQuote:
			// This optimization iterates quicker.
			for script[i] != '\'' {
				i++
			}
			state = sWord

		// Double quotes
		case sDoubleQuote:
			switch c {
			case '"':
				state = sWord
			case '\\':
				state = sDoubleQuoteEscape
			}
		case sDoubleQuoteEscape:
			state = sDoubleQuote

		// Line comment
		case sLineComment:
			switch c {
			case '\n':
				ts = append(ts, token{script[i : i+1], ttNewLine})
				state = sStart
			case '\\':
				state = sLineCommentEscape
			}
		case sLineCommentEscape:
			state = sLineComment

		// Operators
		case sOperator:
			switch c {
			case '&', '|', ';', '<', '>', '-':
				_, ok := operators[script[b:i+1]]
				if ok {
					break
				}
				fallthrough
			default:
				word := strings.ReplaceAll(script[b:i], "\\\n", "")
				op, ok := operators[word]
				if ok {
					ts = append(ts, token{word, op})
				} else {
					ts = append(ts, token{word, ttWord})
				}
				state = sStart
				i--
			}
		}

		if i == len(script) {
			break
		}
		i++
	}

	return ts, nil
}
