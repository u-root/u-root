// Copyright 2026 The HuaTuo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/net/bpf"
)

type extendedTokenKind uint8

const (
	extendedAtom extendedTokenKind = iota
	extendedOperator
	extendedLeftParen
	extendedRightParen
	extendedLeftBracket
	extendedRightBracket
	extendedColon
)

type extendedToken struct {
	kind  extendedTokenKind
	text  string
	start int
	end   int
}

type extendedParser struct {
	source        string
	tokens        []extendedToken
	pos           int
	depth         int
	lastPrimitive *primitive
}

const (
	maxFilterExpressionLength = 64 * 1024
	maxFilterTokens           = 8192
	maxFilterNesting          = 128
)

func parseFilterExpression(expression string) (Filter, error) {
	if len(expression) > maxFilterExpressionLength {
		return nil, fmt.Errorf("filter expression exceeds %d bytes", maxFilterExpressionLength)
	}
	tokens, err := scanExtended(expression)
	if err != nil {
		return nil, err
	}
	if err := validateExtendedLimits(tokens); err != nil {
		return nil, err
	}
	if feature := unsupportedKeyword(tokens); feature != "" {
		return nil, unsupportedFeature(feature)
	}
	if !needsExtendedParser(tokens) {
		parsed := NewExpression(expression)
		if parsed == nil {
			return nil, ErrEmptyFilter
		}
		return parsed.Compile(), nil
	}
	parser := extendedParser{source: expression, tokens: tokens}
	filter, err := parser.parseOr()
	if err != nil {
		return nil, err
	}
	if parser.pos != len(tokens) {
		return nil, fmt.Errorf("unexpected token %q", tokens[parser.pos].text)
	}
	return filter, nil
}

func unsupportedKeyword(tokens []extendedToken) string {
	for _, token := range tokens {
		if token.kind != extendedAtom {
			continue
		}
		switch strings.ToLower(token.text) {
		case "broadcast", "inbound", "outbound", "ifindex", "protochain", "less", "greater":
			return strings.ToLower(token.text)
		}
	}
	return ""
}

func scanExtended(expression string) ([]extendedToken, error) {
	tokens := make([]extendedToken, 0, len(expression)/2)
	for index := 0; index < len(expression); {
		if unicode.IsSpace(rune(expression[index])) {
			index++
			continue
		}
		start := index
		if index+1 < len(expression) {
			pair := expression[index : index+2]
			switch pair {
			case "&&", "||", "==", "!=", ">=", "<=", "<<", ">>":
				tokens = append(tokens, extendedToken{
					kind: extendedOperator, text: pair, start: start, end: index + 2,
				})
				index += 2
				continue
			}
		}
		switch expression[index] {
		case '(':
			tokens = append(tokens, extendedToken{kind: extendedLeftParen, text: "(", start: start, end: start + 1})
			index++
		case ')':
			tokens = append(tokens, extendedToken{kind: extendedRightParen, text: ")", start: start, end: start + 1})
			index++
		case '[':
			tokens = append(tokens, extendedToken{kind: extendedLeftBracket, text: "[", start: start, end: start + 1})
			index++
		case ']':
			tokens = append(tokens, extendedToken{kind: extendedRightBracket, text: "]", start: start, end: start + 1})
			index++
		case ':':
			tokens = append(tokens, extendedToken{kind: extendedColon, text: ":", start: start, end: start + 1})
			index++
		case '+', '-', '*', '/', '%', '&', '|', '^', '!', '>', '<', '=':
			tokens = append(tokens, extendedToken{
				kind: extendedOperator, text: expression[index : index+1], start: start, end: start + 1,
			})
			index++
		default:
			for index < len(expression) && isExtendedAtomByte(expression[index]) {
				index++
			}
			if start == index {
				return nil, fmt.Errorf("illegal character %q", expression[index])
			}
			tokens = append(tokens, extendedToken{
				kind: extendedAtom, text: expression[start:index], start: start, end: index,
			})
		}
		if len(tokens) > maxFilterTokens {
			return nil, fmt.Errorf("filter expression exceeds %d tokens", maxFilterTokens)
		}
	}
	return tokens, nil
}

func validateExtendedLimits(tokens []extendedToken) error {
	depth := 0
	for _, token := range tokens {
		switch token.kind {
		case extendedLeftParen, extendedLeftBracket:
			depth++
			if depth > maxFilterNesting {
				return fmt.Errorf("filter expression exceeds nesting limit %d", maxFilterNesting)
			}
		case extendedRightParen, extendedRightBracket:
			depth--
			if depth < 0 {
				return fmt.Errorf("unexpected closing delimiter %q", token.text)
			}
		}
	}
	if depth != 0 {
		return fmt.Errorf("unclosed filter delimiter")
	}
	return nil
}

func isExtendedAtomByte(value byte) bool {
	return value >= 'a' && value <= 'z' ||
		value >= 'A' && value <= 'Z' ||
		value >= '0' && value <= '9' ||
		value == '_' || value == '.' || value == '-' || value == '\\'
}

func needsExtendedParser(tokens []extendedToken) bool {
	for _, token := range tokens {
		if token.kind == extendedLeftBracket {
			return true
		}
		if token.kind == extendedOperator && isComparisonOperator(token.text) {
			return true
		}
		if token.kind != extendedAtom {
			continue
		}
		switch strings.ToLower(token.text) {
		case "vlan", "mpls", "len":
			return true
		}
	}
	return false
}

func (p *extendedParser) parseOr() (Filter, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	filters := Filters{left}
	for p.matchBoolean("or", "||") {
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		filters = append(filters, right)
	}
	if len(filters) == 1 {
		return left, nil
	}
	return composite{filters: filters, and: false}, nil
}

func (p *extendedParser) parseAnd() (Filter, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	filters := Filters{left}
	for p.matchBoolean("and", "&&") {
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		filters = append(filters, right)
	}
	if len(filters) == 1 {
		return left, nil
	}
	return composite{filters: filters, and: true}, nil
}

func (p *extendedParser) parseUnary() (Filter, error) {
	p.depth++
	defer func() { p.depth-- }()
	if p.depth > maxFilterNesting {
		return nil, fmt.Errorf("filter expression exceeds nesting limit %d", maxFilterNesting)
	}
	if p.matchBoolean("not", "!") {
		inner, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return negated{inner: inner}, nil
	}
	if p.matchKind(extendedLeftParen) {
		inner, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if !p.matchKind(extendedRightParen) {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		return inner, nil
	}
	return p.parseClause()
}

func (p *extendedParser) parseClause() (Filter, error) {
	start := p.pos
	end := p.clauseEnd()
	if start == end {
		return nil, fmt.Errorf("missing filter clause")
	}
	clause := p.tokens[start:end]
	p.pos = end

	if clause[0].kind == extendedAtom {
		switch strings.ToLower(clause[0].text) {
		case "vlan":
			return parseTransition(clause, filterKindVLAN)
		case "mpls":
			return parseTransition(clause, filterKindMPLS)
		}
	}
	if comparison := comparisonIndex(clause); comparison >= 0 {
		return parseArithmeticComparison(clause, comparison)
	}
	if isIncompleteArithmeticClause(clause) {
		return nil, fmt.Errorf("packet arithmetic requires a comparison")
	}

	text := strings.TrimSpace(p.source[clause[0].start:clause[len(clause)-1].end])
	expression := NewExpression(text)
	if expression == nil {
		return nil, ErrInvalidFilter
	}
	element := expression.Next()
	primitiveValue, ok := element.(primitive)
	if !ok || expression.Next() != nil {
		return nil, ErrInvalidFilter
	}
	setPrimitiveDefaults(&primitiveValue, p.lastPrimitive)
	p.lastPrimitive = &primitiveValue
	if primitiveValue.negator {
		return negated{inner: primitiveValue}, nil
	}
	return primitiveValue, nil
}

func isIncompleteArithmeticClause(tokens []extendedToken) bool {
	for _, token := range tokens {
		if token.kind == extendedLeftBracket || token.kind == extendedRightBracket {
			return true
		}
	}
	return len(tokens) == 1 && tokens[0].kind == extendedAtom &&
		strings.EqualFold(tokens[0].text, "len")
}

func (p *extendedParser) clauseEnd() int {
	brackets := 0
	parentheses := 0
	for index := p.pos; index < len(p.tokens); index++ {
		token := p.tokens[index]
		switch token.kind {
		case extendedLeftBracket:
			brackets++
		case extendedRightBracket:
			brackets--
		case extendedLeftParen:
			parentheses++
		case extendedRightParen:
			if brackets == 0 && parentheses == 0 {
				return index
			}
			parentheses--
		}
		if brackets != 0 || parentheses != 0 || !isBooleanToken(token, "and", "or", "&&", "||") {
			continue
		}
		if isDirectionJoin(p.tokens, p.pos, index) {
			continue
		}
		return index
	}
	return len(p.tokens)
}

func isDirectionJoin(tokens []extendedToken, start, index int) bool {
	if index != start+1 || start >= len(tokens) || index+1 >= len(tokens) {
		return false
	}
	return strings.EqualFold(tokens[start].text, "src") && strings.EqualFold(tokens[index+1].text, "dst")
}

func parseTransition(tokens []extendedToken, kind filterKind) (Filter, error) {
	if len(tokens) > 2 {
		return nil, fmt.Errorf("invalid transition %q", tokens[0].text)
	}
	value := ""
	if len(tokens) == 2 {
		if tokens[1].kind != extendedAtom {
			return nil, fmt.Errorf("invalid transition value %q", tokens[1].text)
		}
		value = tokens[1].text
	}
	return primitive{
		kind:      kind,
		direction: filterDirectionSrcOrDst,
		id:        value,
	}, nil
}

func comparisonIndex(tokens []extendedToken) int {
	brackets := 0
	parentheses := 0
	for index, token := range tokens {
		switch token.kind {
		case extendedLeftBracket:
			brackets++
		case extendedRightBracket:
			brackets--
		case extendedLeftParen:
			parentheses++
		case extendedRightParen:
			parentheses--
		case extendedOperator:
			if brackets == 0 && parentheses == 0 && isComparisonOperator(token.text) {
				return index
			}
		}
	}
	return -1
}

func parseArithmeticComparison(tokens []extendedToken, comparison int) (Filter, error) {
	if comparison == 0 || comparison == len(tokens)-1 {
		return nil, fmt.Errorf("comparison %q is missing an operand", tokens[comparison].text)
	}
	left, err := parseArithmetic(tokens[:comparison])
	if err != nil {
		return nil, err
	}
	right, err := parseArithmetic(tokens[comparison+1:])
	if err != nil {
		return nil, err
	}
	test, err := comparisonTest(tokens[comparison].text)
	if err != nil {
		return nil, err
	}
	return arithmeticComparison{left: left, test: test, right: right}, nil
}

func comparisonTest(operator string) (bpf.JumpTest, error) {
	switch operator {
	case "=", "==":
		return bpf.JumpEqual, nil
	case "!=":
		return bpf.JumpNotEqual, nil
	case ">":
		return bpf.JumpGreaterThan, nil
	case "<":
		return bpf.JumpLessThan, nil
	case ">=":
		return bpf.JumpGreaterOrEqual, nil
	case "<=":
		return bpf.JumpLessOrEqual, nil
	default:
		return 0, fmt.Errorf("unsupported comparison operator %q", operator)
	}
}

func isComparisonOperator(operator string) bool {
	switch operator {
	case "=", "==", "!=", ">", "<", ">=", "<=":
		return true
	default:
		return false
	}
}

func (p *extendedParser) matchBoolean(values ...string) bool {
	if p.pos >= len(p.tokens) || !isBooleanToken(p.tokens[p.pos], values...) {
		return false
	}
	p.pos++
	return true
}

func isBooleanToken(token extendedToken, values ...string) bool {
	for _, value := range values {
		if strings.EqualFold(token.text, value) {
			return true
		}
	}
	return false
}

func (p *extendedParser) matchKind(kind extendedTokenKind) bool {
	if p.pos >= len(p.tokens) || p.tokens[p.pos].kind != kind {
		return false
	}
	p.pos++
	return true
}

type arithmeticParser struct {
	tokens []extendedToken
	pos    int
	depth  int
}

func parseArithmetic(tokens []extendedToken) (arithmeticExpression, error) {
	parser := arithmeticParser{tokens: tokens}
	expression, err := parser.parseBinary(1)
	if err != nil {
		return nil, err
	}
	if parser.pos != len(tokens) {
		return nil, fmt.Errorf("unexpected arithmetic token %q", tokens[parser.pos].text)
	}
	return expression, nil
}

func (p *arithmeticParser) parseBinary(minimumPrecedence int) (arithmeticExpression, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for p.pos < len(p.tokens) {
		operator, precedence, ok := arithmeticOperator(p.tokens[p.pos].text)
		if !ok || precedence < minimumPrecedence {
			break
		}
		p.pos++
		right, err := p.parseBinary(precedence + 1)
		if err != nil {
			return nil, err
		}
		left = arithmeticBinary{left: left, op: operator, right: right}
	}
	return left, nil
}

func (p *arithmeticParser) parsePrimary() (arithmeticExpression, error) {
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("missing arithmetic operand")
	}
	token := p.tokens[p.pos]
	if token.kind == extendedLeftParen {
		p.depth++
		defer func() { p.depth-- }()
		if p.depth > maxFilterNesting {
			return nil, fmt.Errorf("arithmetic expression exceeds nesting limit %d", maxFilterNesting)
		}
		p.pos++
		expression, err := p.parseBinary(1)
		if err != nil {
			return nil, err
		}
		if p.pos >= len(p.tokens) || p.tokens[p.pos].kind != extendedRightParen {
			return nil, fmt.Errorf("missing arithmetic closing parenthesis")
		}
		p.pos++
		return expression, nil
	}
	if token.kind != extendedAtom {
		return nil, fmt.Errorf("unexpected arithmetic token %q", token.text)
	}
	p.pos++
	name := strings.ToLower(token.text)
	if value, ok := arithmeticNames[name]; ok {
		return arithmeticConstant{value: value}, nil
	}
	if name == "len" {
		return packetLength{}, nil
	}
	if protocol, ok := packetProtocols[name]; ok && p.pos < len(p.tokens) && p.tokens[p.pos].kind == extendedLeftBracket {
		return p.parsePacketAccess(protocol)
	}
	value, err := strconv.ParseUint(token.text, 0, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid arithmetic value %q", token.text)
	}
	return arithmeticConstant{value: uint32(value)}, nil
}

func (p *arithmeticParser) parsePacketAccess(protocol packetProtocol) (arithmeticExpression, error) {
	p.pos++
	offsetStart := p.pos
	depth := 0
	separator := -1
	end := -1
	for index := p.pos; index < len(p.tokens); index++ {
		switch p.tokens[index].kind {
		case extendedLeftParen, extendedLeftBracket:
			depth++
		case extendedRightParen:
			depth--
		case extendedColon:
			if depth == 0 {
				separator = index
			}
		case extendedRightBracket:
			if depth == 0 {
				end = index
				index = len(p.tokens)
			} else {
				depth--
			}
		}
	}
	if end < 0 || offsetStart == end {
		return nil, fmt.Errorf("invalid packet access")
	}
	offsetEnd := end
	size := 1
	if separator >= 0 {
		offsetEnd = separator
		if separator+2 != end {
			return nil, fmt.Errorf("packet access width must be 1, 2, or 4")
		}
		width, err := strconv.Atoi(p.tokens[separator+1].text)
		if err != nil || width != 1 && width != 2 && width != 4 {
			return nil, fmt.Errorf("packet access width must be 1, 2, or 4")
		}
		size = width
	}
	offset, err := parseArithmetic(p.tokens[offsetStart:offsetEnd])
	if err != nil {
		return nil, fmt.Errorf("parsing packet offset: %w", err)
	}
	p.pos = end + 1
	return packetAccess{protocol: protocol, offset: offset, size: size}, nil
}

func arithmeticOperator(value string) (bpf.ALUOp, int, bool) {
	switch value {
	case "|":
		return bpf.ALUOpOr, 1, true
	case "^":
		return bpf.ALUOpXor, 2, true
	case "&":
		return bpf.ALUOpAnd, 3, true
	case "<<":
		return bpf.ALUOpShiftLeft, 4, true
	case ">>":
		return bpf.ALUOpShiftRight, 4, true
	case "+":
		return bpf.ALUOpAdd, 5, true
	case "-":
		return bpf.ALUOpSub, 5, true
	case "*":
		return bpf.ALUOpMul, 6, true
	case "/":
		return bpf.ALUOpDiv, 6, true
	case "%":
		return bpf.ALUOpMod, 6, true
	default:
		return 0, 0, false
	}
}

var packetProtocols = map[string]packetProtocol{
	"ether": packetProtocolEther,
	"ip":    packetProtocolIP,
	"ip6":   packetProtocolIP6,
	"tcp":   packetProtocolTCP,
	"udp":   packetProtocolUDP,
}

var arithmeticNames = map[string]uint32{
	"tcpflags":   13,
	"tcp-fin":    0x01,
	"tcp-syn":    0x02,
	"tcp-rst":    0x04,
	"tcp-push":   0x08,
	"tcp-ack":    0x10,
	"tcp-urg":    0x20,
	"tcp-ece":    0x40,
	"tcp-cwr":    0x80,
	"icmp-type":  0,
	"icmp-code":  1,
	"icmp6-type": 0,
	"icmp6-code": 1,
}
