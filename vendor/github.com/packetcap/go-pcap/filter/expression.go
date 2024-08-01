package filter

import (
	"bufio"
	"bytes"
	"strings"
)

type ExpressionToken int

const (
	eof rune = 0
)
const (
	tokenAnd ExpressionToken = iota
	tokenOr
	tokenNot
	tokenLeft
	tokenRight
	tokenWhitespace
	tokenWord
	tokenEOF
	tokenIllegal
	tokenSrc
	tokenDst
	tokenGateway
	tokenProto
	tokenIP4
	tokenIP6
	tokenNet
	tokenTCP
	tokenUDP
	tokenID
	tokenHost
	tokenPort
	tokenPortRange
	tokenEther
)

var lexerTokens = map[string]ExpressionToken{
	"and":       tokenAnd,
	"or":        tokenOr,
	"not":       tokenNot,
	"gateway":   tokenGateway,
	"proto":     tokenProto,
	"ether":     tokenEther,
	"src":       tokenSrc,
	"dst":       tokenDst,
	"net":       tokenNet,
	"port":      tokenPort,
	"host":      tokenHost,
	"portrange": tokenPortRange,
	"ip":        tokenIP4,
	"ip4":       tokenIP4,
	"ip6":       tokenIP6,
	"tcp":       tokenTCP,
	"udp":       tokenUDP,
}

type buffer struct {
	token ExpressionToken
	word  string
}
type Expression struct {
	raw    string
	lexer  expressionLexer
	buffer buffer
}

type expressionLexer struct {
	reader *bufio.Reader
}

func NewExpression(s string) *Expression {
	if s == "" {
		return nil
	}
	e := &Expression{
		raw: s,
		lexer: expressionLexer{
			reader: bufio.NewReader(strings.NewReader(s)),
		},
	}
	// initialize the buffer
	e.scan()
	return e
}

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

// isAlpha returns true if the rune is a letter.
func isAlpha(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}

// isValidWord returns true if the rune is part of a valid word, which is broader
// than just alphanumeric, e.g. 10.100.100.100/24 or fe200::
func isValidWord(ch rune) bool {
	return isAlpha(ch) || ch == '/' || ch == '.' || ch == ':' || ch == '-'
}

// scanWhitespace scan past all of the next whitespace
func (e *expressionLexer) scanWhitespace() ExpressionToken {
	// read until we either do not have whitespace or have an EOF
	// be sure to put back the whitespace for the next read
	for {
		if ch := e.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			e.unread()
			break
		}
	}
	return tokenWhitespace
}

func (e *expressionLexer) unread() {
	_ = e.reader.UnreadRune()
}

func (e *expressionLexer) read() rune {
	ch, _, err := e.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// scanWord consumes the current rune and all contiguous ident runes until one of:
// - eof
// - whitespace
// - non-ascii
func (e *expressionLexer) scanWord() (ExpressionToken, string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// Read every subsequent character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
forloop:
	for {
		ch := e.read()
		switch {
		case ch == eof:
			break forloop
		case !isValidWord(ch):
			e.unread()
			break forloop
		default:
			buf.WriteRune(ch)
		}
	}

	// we now have a full word
	word := buf.String()
	if val, ok := lexerTokens[word]; ok {
		return val, word
	}

	// Otherwise return as a regular identifier.
	return tokenID, word
}

// Scan read the next element from the expression and convert it into a token
// It might return a primitive, a composite or a joiner.
func (e *expressionLexer) Scan() (ExpressionToken, string) {
	ch := e.read()
	if ch == eof {
		return tokenEOF, ""
	}

	// handle whitespace by reading it all as a single one
	switch {
	case isWhitespace(ch):
		e.unread()
		return e.scanWhitespace(), ""
	case ch == '(':
		return tokenLeft, string(ch)
	case ch == ')':
		return tokenRight, string(ch)
	case isAlpha(ch):
		e.unread()
		return e.scanWord()
	}
	return tokenIllegal, ""
}

// Compile build an abstract syntax tree of the expression, implemented in
// a Filter.
func (e *Expression) Compile() Filter {
	// create a root element, which should be a composite. If it ends up having
	// just one member, we will return just that at the end.
	var combo composite

	for {
		var fe Element
		if fe = e.Next(); fe == nil {
			break
		}
		switch fe.Type() {
		case Primitive:
			p := fe.(primitive)
			setPrimitiveDefaults(&p, combo.LastPrimitive())
			combo.filters = append(combo.filters, p)
		case Composite:
			c := fe.(composite)
			combo.filters = append(combo.filters, c)
		case Joiner:
			// it is not a primitive, so it is a joiner
			isAnd := fe.(*and)
			combo.and = bool(*isAnd)
		}
	}
	return combo.Distill()
}

func (e *Expression) scan() (ExpressionToken, string) {
	tok, word := e.buffer.token, e.buffer.word
	e.buffer.token, e.buffer.word = e.lexer.Scan()
	return tok, word
}

func (e *Expression) scanPastWhitespace() (ExpressionToken, string) {
	var (
		tok  ExpressionToken
		word string
	)
	for {
		tok, word = e.buffer.token, e.buffer.word
		e.buffer.token, e.buffer.word = e.lexer.Scan()
		if tok != tokenWhitespace {
			break
		}
	}
	return tok, word
}

func (e *Expression) peek() (ExpressionToken, string) {
	return e.buffer.token, e.buffer.word
}

func (e *Expression) peekPastWhitespace() (ExpressionToken, string) {
	var (
		tok  ExpressionToken
		word string
	)
	for {
		tok, word = e.buffer.token, e.buffer.word
		if tok != tokenWhitespace {
			break
		}
		e.buffer.token, e.buffer.word = e.lexer.Scan()
	}

	return tok, word
}

func (e *Expression) HasNext() bool {
	tok, _ := e.peek()
	return tok != tokenEOF
}

// Next get the next element. If none left, return nil.
// It might return a primitive, a composite or a joiner.
func (e *Expression) Next() Element {
	if !e.HasNext() {
		return nil
	}

	var inElement bool

	p := primitive{
		direction: filterDirectionUnset,
		kind:      filterKindUnset,
		protocol:  filterProtocolUnset,
	}

tokens:
	for {
		tok, _ := e.peekPastWhitespace()
		// handle the case where the next element will be the end of us
		if inElement && (tok == tokenAnd || tok == tokenOr || tok == tokenEOF) {
			// we hit "and" or "or". If we already have started building a primitive,
			// return the started one. Else return a joiner.
			// We account for the special case of "src and dst" or "src or dst" below.
			return p
		}

		tok, word := e.scanPastWhitespace()

		// indicate we are in an element
		inElement = true

		// first look for sub-element, negator, and other special cases
		// if we hit the indication of the end of an element - eof, and, or, left, right -
		// we are starting a new element. If we were in the middle of an element, we would
		// have handled it in the "if inElement { }" statement above
		switch tok {
		case tokenEOF:
			return nil
		case tokenAnd:
			j := and(true)
			return &j
		case tokenOr:
			j := and(false)
			return &j
		case tokenLeft:
			// start a new sub-element
			return e.tokenBrace()
		case tokenRight:
			// end a sub-element
			return p
		case tokenNot:
			p.negator = true
			continue tokens
		case tokenGateway:
			// this really needs to use the composite of two primitives
			p.protocol = filterProtocolEther
			p.kind = filterKindHost
			continue tokens
		case tokenProto:
			// the next word is the sub-protocol
			tok, word := e.scanPastWhitespace()
			if tok == tokenEOF {
				continue tokens
			}
			// we will accept the protocol as "name" or "\name", because some get escaped
			protoName := strings.TrimLeft(word, "\\")
			if sub, ok := subProtocols[protoName]; ok {
				p.subProtocol = sub
			} else {
				p.subProtocol = filterSubProtocolUnknown
				p.id = protoName
			}
			continue tokens
		case tokenSrc:
			// handle the "src or dst"/"src and dst" case
			if nToken, _ := e.peekPastWhitespace(); nToken == tokenAnd || nToken == tokenOr {
				// get that next token
				andor, _ := e.scanPastWhitespace()
				// get the one after that
				isDst, _ := e.scanPastWhitespace()
				switch {
				case andor == tokenAnd && isDst == tokenDst:
					p.direction = filterDirectionSrcAndDst
				case andor == tokenOr && isDst == tokenDst:
					p.direction = filterDirectionSrcOrDst
				default:
					return nil
				}
			} else {
				p.direction = filterDirectionSrc
			}
		case tokenDst:
			p.direction = filterDirectionDst
		}
		// it must be a primitive word, so find it
		if kind, ok := kinds2[tok]; ok {
			p.kind = kind
		} else if protocol, ok := protocols[word]; ok {
			p.protocol = protocol
		} else if subprotocol, ok := subProtocols[word]; ok {
			p.subProtocol = subprotocol
		} else {
			p.id = word
		}
	}
}

// tokenBrace process the innards of a "( ... )"
func (e *Expression) tokenBrace() Filter {
	return e.Compile()
}

// setPrimitiveDefaults set defaults on expressions
func setPrimitiveDefaults(p, lastPrimitive *primitive) {
	// if nothing was set, do not try to fix it
	if p.direction == filterDirectionUnset && p.protocol == filterProtocolUnset && p.kind == filterKindUnset && p.subProtocol == filterSubProtocolUnset {
		if lastPrimitive == nil {
			return
		}

		// we only copy over the previous ones if everything else is identical, per the manpage:
		/*
			To save typing, identical qualifier lists can be omitted. E.g., `tcp dst port ftp or ftp-data or domain' is exactly the same as `tcp dst port ftp or tcp dst port ftp-data or tcp dst port domain'
		*/
		p.direction = lastPrimitive.direction
		p.kind = lastPrimitive.kind
		p.protocol = lastPrimitive.protocol
		p.subProtocol = lastPrimitive.subProtocol
	}
	// special cases
	//if (p.subProtocol == filterSubProtocolUDP || p.subProtocol == filterSubProtocolTCP || p.subProtocol == filterSubProtocolIcmp) && p.protocol == filterProtocolUnset {
	//p.protocol = filterProtocolIP
	//}

	if p.kind == filterKindUnset && p.direction != filterDirectionUnset && (p.protocol == filterProtocolEther || p.protocol == filterProtocolIP || p.protocol == filterProtocolIP6 || p.protocol == filterProtocolArp || p.protocol == filterProtocolRarp) {
		p.kind = filterKindHost
	}
	if p.direction == filterDirectionUnset {
		p.direction = filterDirectionSrcOrDst
	}
	if p.kind == filterKindUnset && p.protocol == filterProtocolUnset && p.subProtocol == filterSubProtocolUnset {
		p.kind = filterKindHost
	}
}
