package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrInputTooShort           = errors.New("input too short")
	ErrInvalidCharacter        = errors.New("invalid character")
	ErrDirectiveClosingMissing = errors.New("directive closing missing")
	ErrUnexpectedGroupClosing  = errors.New("unexpected group closing found")
	ErrUnexpectedDirective     = errors.New("unexpected directive found")
	ErrDuplicateID             = errors.New("duplicate id found")
)

const (
	plus               = '+'
	dive               = '>'
	ascend             = '^'
	openingBracket     = '['
	closingBracket     = ']'
	openingParenthesis = '('
	closingParenthesis = ')'
	openingBrace       = '{'
	closingBrace       = '}'
	atSign             = '@'
	dollarSign         = '$'
	hashSign           = '#'
	equalSign          = '='
	colon              = ':'
	dash               = '-'
	dotSign            = '.'
	space              = ' '
	quote              = '"'
	star               = '*'
)

type Action int32

const (
	Add        Action = plus
	Dive       Action = dive
	Ascend     Action = ascend
	CloseGroup Action = closingParenthesis
)

func allowedClassName(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '-' || r >= '0' && r <= '9'
}

func allowedXMLTagName(r rune) bool {
	return allowedHTMLTagName(r) || r == colon
}

func allowedHTMLTagName(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
}

func allowedNumbers(r rune) bool {
	return r >= '0' && r <= '9'
}

func allowedText(r rune) bool {
	return r != closingBrace
}

func allowedQuoteContent(r rune) bool {
	return r != '"'
}

func allowedUnquotedAttribute(r rune) bool {
	return r != ' ' && r != closingBracket
}

type Lexer struct {
	mode Mode
}

func NewLexer(mode Mode) *Lexer {
	return &Lexer{
		mode: mode,
	}
}

func (l *Lexer) FindTokenValue(runes []rune, allowed func(r rune) bool) (string, int) {
	var (
		builder strings.Builder
		length  int
	)

	for _, r := range runes { // nolint: varnamelen
		if !allowed(r) {
			break
		}

		length++

		builder.WriteRune(r)
	}

	tokenValue := builder.String()

	return tokenValue, length
}

func (l *Lexer) FindNumbering(runes []rune) (int, bool, string, int, error) {
	pos := 0

	for _, r := range runes {
		if r != dollarSign {
			break
		}

		pos++
	}

	numbering := string(runes[:pos])

	if len(runes) <= pos || runes[pos] != atSign {
		return 1, false, numbering, pos, nil
	}

	pos++

	// atSign can't be the last character in the numbering directive
	if len(runes[pos:]) == 0 {
		return 0, false, "", pos, ErrInputTooShort
	}

	var (
		reverse = false
		start   = 1
	)

	if runes[pos] == dash {
		pos++

		reverse = true
	}

	str, length := l.FindTokenValue(runes[pos:], allowedNumbers)
	if length == 0 {
		// @ should not be defined on its own
		if !reverse && numbering == "" {
			return 0, false, "", pos, ErrInputTooShort
		}

		return start, reverse, numbering, pos, nil
	}

	num, _ := strconv.Atoi(str)

	start = num
	pos += length

	return start, reverse, numbering, pos, nil
}

func (l *Lexer) FindClassOrIDToken(runes []rune) (*AttrValue, int, error) {
	if len(runes) < 2 { // nolint: gomnd
		return nil, 0, ErrInputTooShort
	}

	if runes[0] != dotSign && runes[0] != hashSign {
		return nil, 0, ErrInvalidCharacter
	}

	value, length := l.FindTokenValue(runes[1:], allowedClassName)

	if length == 0 {
		return nil, 1, ErrInputTooShort
	}

	var token *AttrValue

	if runes[0] == hashSign {
		token = NewID(value)
	} else {
		token = NewClass(value)
	}

	pos := length + 1

	if pos < len(runes) && (runes[pos] == dollarSign || runes[pos] == atSign) {
		start, reverse, numbering, numLength, err := l.FindNumbering(runes[pos:])
		if err != nil {
			return nil, pos + numLength, err
		}

		token.SetNumbering(numbering).SetStart(start).SetReverse(reverse)

		pos += numLength
	}

	return token, pos, nil
}

func (l *Lexer) FindAttribute(runes []rune) (string, string, bool, int, error) {
	name, length := l.FindTokenValue(runes, allowedClassName)

	// Attribute name can't be empty
	if length == 0 {
		return "", "", false, 0, ErrInputTooShort
	}

	// No equal sign after attribute name
	if len(runes[length:]) == 0 || runes[length] == ' ' || runes[length] == closingBracket {
		return name, "", false, length, nil
	}

	// Equal sign is expected after attribute name for attributes with a Value
	if runes[length] != equalSign {
		return "", "", false, length, ErrInvalidCharacter
	}

	pos := length + 1

	if len(runes[pos:]) == 0 || runes[pos] == space {
		return name, "", true, pos, nil
	}

	if runes[pos] == quote {
		value, valueLength := l.FindTokenValue(runes[pos+1:], allowedQuoteContent)

		return name, value, true, pos + valueLength + 2, nil // nolint:gomnd
	}

	value, valueLength := l.FindTokenValue(runes[pos:], allowedUnquotedAttribute)

	return name, value, true, pos + valueLength, nil
}

func (l *Lexer) FindAttributeTokens(runes []rune) (AttrList, int, error) {
	if len(runes) < 2 { // nolint:gomnd
		return nil, 0, ErrInputTooShort
	}

	if runes[0] != openingBracket {
		return nil, 0, ErrInvalidCharacter
	}

	var attributes AttrList

	pos := 1

	for {
		if runes[pos] == closingBracket {
			return attributes, pos + 1, nil
		}

		name, value, hasEqualSign, length, err := l.FindAttribute(runes[pos:])
		pos += length

		if err != nil {
			return nil, pos, err
		}

		attributes = append(attributes, NewAttr(name, value))

		if !hasEqualSign {
			attributes[len(attributes)-1].HasNoEqualSign()
		}

		for _, ch := range runes[pos:] {
			if ch != ' ' {
				break
			}

			pos++
		}

		if len(runes[pos:]) == 0 {
			break
		}
	}

	return nil, 0, ErrDirectiveClosingMissing
}

func (l *Lexer) FindClassToken(runes []rune) (*AttrValue, int, error) {
	classToken, length, err := l.FindClassOrIDToken(runes)
	if err != nil {
		return nil, length, err
	}

	return classToken, length, nil
}

func (l *Lexer) FindIDToken(runes []rune) (*AttrValue, int, error) {
	idToken, length, err := l.FindClassOrIDToken(runes)
	if err != nil {
		return nil, length, err
	}

	return idToken, length, nil
}

func (l *Lexer) FindAllAttributeTokens(token *TagToken, runes []rune) (int, error) {
	if len(runes) == 0 {
		return 0, nil
	}

	pos := 0

	for {
		if len(runes) <= pos {
			break
		}

		switch runes[pos] {
		case openingBracket:
			currentTokens, attrLength, err := l.FindAttributeTokens(runes[pos:])
			if err != nil {
				return pos, err
			}

			token.Attributes = append(token.Attributes, currentTokens...)
			pos += attrLength

		case dotSign:
			currentToken, attrLength, err := l.FindClassToken(runes[pos:])
			if err != nil {
				return pos, err
			}

			token.Classes = append(token.Classes, currentToken)
			pos += attrLength

		case hashSign:
			if token.ID != nil {
				return pos, ErrDuplicateID
			}

			currentToken, attrLength, err := l.FindIDToken(runes[pos:])
			if err != nil {
				return pos, err
			}

			token.ID = currentToken
			pos += attrLength
		default:
			return pos, nil
		}
	}

	return pos, nil
}

func (l *Lexer) FindRepeat(runes []rune) (int, int, error) {
	if len(runes) == 0 || runes[0] != star {
		return 1, 0, nil
	}

	repeatStr, numLength := l.FindTokenValue(runes[1:], allowedNumbers)
	if numLength == 0 {
		return 0, 0, ErrInputTooShort
	}

	repeat, err := strconv.Atoi(repeatStr)
	if err != nil {
		return 0, 0, ErrInvalidCharacter
	}

	if repeat < 1 {
		repeat = 1
	}

	return repeat, numLength + 1, nil
}

func (l *Lexer) NextTextToken(runes []rune) (*Text, int, error) {
	if len(runes) == 0 || runes[0] != openingBrace {
		return nil, 0, nil
	}

	value, length := l.FindTokenValue(runes[1:], allowedText)
	if length == 0 {
		if runes[1] == closingBrace {
			return nil, 2, nil // nolint: gomnd
		}

		panic("logic error")
	}

	pos := length + 1

	if len(runes[pos:]) == 0 || runes[pos] != closingBrace {
		return nil, 0, ErrDirectiveClosingMissing
	}

	return NewText(value), length + 2, nil // nolint: gomnd
}

func (l *Lexer) NextTagToken(runes []rune) (*TagToken, int, error) {
	if len(runes) == 0 {
		return nil, 0, ErrInputTooShort
	}

	f := allowedXMLTagName

	value, length := l.FindTokenValue(runes, f)

	if length == 0 {
		return nil, 0, ErrInputTooShort
	}

	token := NewTagToken(value, 1)
	pos := length

	classLength, err := l.FindAllAttributeTokens(token, runes[pos:])
	if err != nil {
		return nil, pos + classLength, err
	}

	pos += classLength

	repeat, repeatLength, err := l.FindRepeat(runes[pos:])
	if err != nil {
		return nil, pos + repeatLength, err
	}

	if repeatLength > 0 {
		pos += repeatLength
		token.Repeat = repeat
	}

	textToken, textLength, err := l.NextTextToken(runes[pos:])
	if err != nil {
		return nil, pos + textLength, err
	}

	token.Text = textToken

	return token, pos + textLength, nil
}

// nolint: ireturn
func (l *Lexer) NextSubjectToken(runes []rune) (Token, int, error) {
	if len(runes) == 0 {
		return nil, 0, ErrInputTooShort
	}

	if runes[0] == openingParenthesis {
		tokens, pos, err := l.Tokenize(runes[1:], true)
		if err != nil {
			return nil, pos, errors.Wrap(err, "failed to tokenize the remaining runes")
		}

		repeat, numLength := 1, 0
		if len(runes) > pos {
			repeat, numLength, err = l.FindRepeat(runes[pos+1:])
			if err != nil {
				return nil, pos + numLength, errors.Wrap(err, "failed to find the repeat")
			}
		}

		return NewGroupToken(repeat, tokens...), pos + numLength + 1, nil
	}

	return l.NextTagToken(runes)
}

// nolint: cyclop
func (l *Lexer) NextDirectiveToken(runes []rune) (*DirectiveToken, int, error) {
	if len(runes) == 0 {
		return nil, 0, nil
	}

	var token *DirectiveToken

	for _, curRune := range runes {
		switch {
		case curRune == plus || curRune == dive:
			if token != nil {
				return nil, 1, ErrInvalidCharacter
			}

			token = NewDirectiveToken(Action(curRune), 1)

		case curRune == ascend:
			if token == nil {
				token = NewDirectiveToken(Action(curRune), 0)
			}

			if token.Name != Action(curRune) {
				return nil, token.Repeat, ErrInvalidCharacter
			}

			token.Repeat++

		case curRune == closingParenthesis:
			if token != nil {
				return nil, 1, ErrInvalidCharacter
			}

			return NewDirectiveToken(CloseGroup, 1), 1, nil

		default:
			goto breakout
		}
	}

breakout:
	if token == nil {
		return nil, 0, ErrInvalidCharacter
	}

	return token, token.Repeat, nil
}

func (l *Lexer) Tokenize(runes []rune, inGroup bool) ([]Token, int, error) {
	var (
		tokens    []Token
		lastToken Token
		directive *DirectiveToken
	)

	pos := 0

	subject, length, err := l.NextSubjectToken(runes[pos:])
	pos += length
	lastToken = subject

	if err != nil {
		return nil, pos, fmt.Errorf("invalid subject token, error at %d, err: %w", pos, err)
	}

	tokens = append(tokens, subject)

	for pos < len(runes) {
		if pos == len(runes) {
			break
		}

		directive, length, err = l.NextDirectiveToken(runes[pos:])
		pos += length

		if err != nil {
			return nil, pos, fmt.Errorf("invalid directive token, error at %d, err: %w", pos, err)
		}

		if directive.Name == CloseGroup && inGroup {
			return tokens, pos, nil
		}

		subject, length, err = l.NextSubjectToken(runes[pos:])
		pos += length

		if err != nil {
			return nil, pos, fmt.Errorf("invalid subject token, error at %d, err: %w", pos, err)
		}

		tokens = l.act(directive, subject, tokens, lastToken)

		lastToken = subject
	}

	return tokens, pos, nil
}

func (l *Lexer) act(directive *DirectiveToken, subject Token, tokens []Token, lastToken Token) []Token {
	switch directive.Name {
	case Add:
		parent := lastToken.GetParent()
		if parent == nil {
			tokens = append(tokens, subject)
		} else {
			parent.AddChildren(subject)
		}

	case Dive:
		lastToken.AddChildren(subject)

	case Ascend:
		parent := lastToken

		for i := 0; i < directive.Repeat+1; i++ {
			parent = parent.GetParent()

			// Traversing too high is ignored, in line with Emmet
			if parent == nil {
				break
			}
		}

		if parent == nil {
			tokens = append(tokens, subject)
		} else {
			parent.AddChildren(subject)
		}

	case CloseGroup:
		panic("logic error")

	default:
		panic("logic error")
	}

	return tokens
}
