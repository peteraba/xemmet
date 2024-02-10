package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_FindTokenValue(t *testing.T) {
	t.Parallel()

	type args struct {
		runes   []rune
		allowed func(r rune) bool
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantParsed string
		wantLength int
	}{
		{
			name: "default",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("foo bar baz quix"),
				allowed: func(r rune) bool {
					return r != 'q'
				},
			},
			wantParsed: "foo bar baz ",
			wantLength: 12,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotParsed, gotLength := tt.sut.FindTokenValue(tt.args.runes, tt.args.allowed)

			assert.Equal(t, tt.wantParsed, gotParsed)
			assert.Equal(t, tt.wantLength, gotLength)
		})
	}
}

func TestLexer_FindNumbering(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name          string
		sut           *Lexer
		args          args
		wantStart     int
		wantReverse   bool
		wantNumbering string
		wantLength    int
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "empty",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(""),
			},
			wantStart:     1,
			wantReverse:   false,
			wantNumbering: "",
			wantLength:    0,
			wantErr:       assert.NoError,
		},
		{
			name: "numbering only",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("$$$"),
			},
			wantStart:     1,
			wantReverse:   false,
			wantNumbering: "$$$",
			wantLength:    3,
			wantErr:       assert.NoError,
		},
		{
			name: "@ only",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("@"),
			},
			wantStart:     0,
			wantReverse:   false,
			wantNumbering: "",
			wantLength:    1,
			wantErr:       assert.Error,
		},
		{
			name: "valid start",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("@3"),
			},
			wantStart:     3,
			wantReverse:   false,
			wantNumbering: "",
			wantLength:    2,
			wantErr:       assert.NoError,
		},
		{
			name: "valid start, continued with @",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("@3@"),
			},
			wantStart:     3,
			wantReverse:   false,
			wantNumbering: "",
			wantLength:    2,
			wantErr:       assert.NoError,
		},
		{
			name: "valid start, continued",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("@3$"),
			},
			wantStart:     3,
			wantReverse:   false,
			wantNumbering: "",
			wantLength:    2,
			wantErr:       assert.NoError,
		},
		{
			name: "numbering and start",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("$$@3"),
			},
			wantStart:     3,
			wantReverse:   false,
			wantNumbering: "$$",
			wantLength:    4,
			wantErr:       assert.NoError,
		},
		{
			name: "numbering and start, continued",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("$$@3*3"),
			},
			wantStart:     3,
			wantReverse:   false,
			wantNumbering: "$$",
			wantLength:    4,
			wantErr:       assert.NoError,
		},
		{
			name: "@ with minus only",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("@-"),
			},
			wantStart:     1,
			wantReverse:   true,
			wantNumbering: "",
			wantLength:    2,
			wantErr:       assert.NoError,
		},
		{
			name: "numbering with start, in reverse",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("$$@-3"),
			},
			wantStart:     3,
			wantReverse:   true,
			wantNumbering: "$$",
			wantLength:    5,
			wantErr:       assert.NoError,
		},
		{
			name: "numbering with start, in reverse, continued",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune("$$@-3+foo"),
			},
			wantStart:     3,
			wantReverse:   true,
			wantNumbering: "$$",
			wantLength:    5,
			wantErr:       assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotStart, gotReverse, gotNumbering, gotLength, err := tt.sut.FindNumbering(tt.args.runes)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantStart, gotStart)
			assert.Equal(t, tt.wantReverse, gotReverse)
			assert.Equal(t, tt.wantNumbering, gotNumbering)
			assert.Equal(t, tt.wantLength, gotLength)
		})
	}
}

func TestLexer_FindClassOrIDToken(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantToken  *AttrValue
		wantLength int
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "simple id",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("#div")},
			wantToken:  NewID("div"),
			wantLength: 4,
			wantErr:    assert.NoError,
		},
		{
			name:       "longer id",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("#div+p")},
			wantToken:  NewID("div"),
			wantLength: 4,
			wantErr:    assert.NoError,
		},
		{
			name:       "simple class",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(".div")},
			wantToken:  NewClass("div"),
			wantLength: 4,
			wantErr:    assert.NoError,
		},
		{
			name:       "longer class",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(".div+p")},
			wantToken:  NewClass("div"),
			wantLength: 4,
			wantErr:    assert.NoError,
		},
		{
			name:       "dash supported",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(".div-p")},
			wantToken:  NewClass("div-p"),
			wantLength: 6,
			wantErr:    assert.NoError,
		},
		{
			name:       "underscore supported",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(".div_p")},
			wantToken:  NewClass("div_p"),
			wantLength: 6,
			wantErr:    assert.NoError,
		},
		{
			name:       "dot not supported",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(".div_p.")},
			wantToken:  NewClass("div_p"),
			wantLength: 6,
			wantErr:    assert.NoError,
		},
		{
			name:       "strange character",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(".divőp")},
			wantToken:  NewClass("div"),
			wantLength: 4,
			wantErr:    assert.NoError,
		},
		{
			name:       "empty",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune{}},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid character",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("!p")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "simple hash",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("#")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "hash after hash",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("##")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "simple dot",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(".")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "dot after dot",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("..")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "hash after dot",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(".#")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "dot after hash",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("#.")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "empty",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune{}},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid @",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("#doo@")},
			wantToken:  nil,
			wantLength: 5,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid @, continued",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("#doo@a")},
			wantToken:  nil,
			wantLength: 5,
			wantErr:    assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotToken, gotLength, gotErr := tt.sut.FindClassOrIDToken(tt.args.runes)
			tt.wantErr(t, gotErr)

			assert.Equalf(t, tt.wantToken, gotToken, "FindClassOrIDToken(%v)", tt.args.runes)
			assert.Equalf(t, tt.wantLength, gotLength, "FindClassOrIDToken(%v)", tt.args.runes)
		})
	}
}

func TestLexer_FindAttribute(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name             string
		sut              *Lexer
		args             args
		wantName         string
		wantValue        string
		wantHasEqualSign bool
		wantLength       int
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:             "empty",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune{}},
			wantName:         "",
			wantValue:        "",
			wantHasEqualSign: false,
			wantLength:       0,
			wantErr:          assert.Error,
		},
		{
			name:             "name only",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune("foobar")},
			wantName:         "foobar",
			wantValue:        "",
			wantHasEqualSign: false,
			wantLength:       6,
			wantErr:          assert.NoError,
		},
		{
			name:             "weird name only",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune("foo12-Bar")},
			wantName:         "foo12-Bar",
			wantValue:        "",
			wantHasEqualSign: false,
			wantLength:       9,
			wantErr:          assert.NoError,
		},
		{
			name:             "weird name only, followed by a space",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune("foo12-Bar ")},
			wantName:         "foo12-Bar",
			wantValue:        "",
			wantHasEqualSign: false,
			wantLength:       9,
			wantErr:          assert.NoError,
		},
		{
			name:             "name and sign",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune("foobar=")},
			wantName:         "foobar",
			wantValue:        "",
			wantHasEqualSign: true,
			wantLength:       7,
			wantErr:          assert.NoError,
		},
		{
			name:             "name, sign, and a space",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune("foobar= ")},
			wantName:         "foobar",
			wantValue:        "",
			wantHasEqualSign: true,
			wantLength:       7,
			wantErr:          assert.NoError,
		},
		{
			name:             "complete with unquoted value",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune("foobar=barfoo!")},
			wantName:         "foobar",
			wantValue:        "barfoo!",
			wantHasEqualSign: true,
			wantLength:       14,
			wantErr:          assert.NoError,
		},
		{
			name:             "complete with quoted value",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune(`foobar="bar! foo"`)},
			wantName:         "foobar",
			wantValue:        "bar! foo",
			wantHasEqualSign: true,
			wantLength:       17,
			wantErr:          assert.NoError,
		},
		{
			name:             "! instead of =",
			sut:              NewLexer(ModeHTML),
			args:             args{runes: []rune("foobar!")},
			wantName:         "",
			wantValue:        "",
			wantHasEqualSign: false,
			wantLength:       6,
			wantErr:          assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotName, gotValue, gotHasEqualSign, gotLength, err := tt.sut.FindAttribute(tt.args.runes)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantName, gotName)
			assert.Equal(t, tt.wantValue, gotValue)
			assert.Equal(t, tt.wantHasEqualSign, gotHasEqualSign)
			assert.Equal(t, tt.wantLength, gotLength)
		})
	}
}

func TestLexer_NextAttributeToken(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantLength int
		wantTokens AttrList
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "broken",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[`),
			},
			wantTokens: nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name: "broken 2",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[foobar`),
			},
			wantTokens: nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name: "empty",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[]`),
			},
			wantTokens: nil,
			wantLength: 2,
			wantErr:    assert.NoError,
		},
		{
			name: "simple",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[foo=3]`),
			},
			wantTokens: []*Attr{
				NewAttr("foo", "3"),
			},
			wantLength: 7,
			wantErr:    assert.NoError,
		},
		{
			name: "simple, continued",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[title="Hello world!"].class`),
			},
			wantTokens: []*Attr{
				NewAttr("title", "Hello world!"),
			},
			wantLength: 22,
			wantErr:    assert.NoError,
		},
		{
			name: "simple, no equal sign, continued",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[title].class`),
			},
			wantTokens: []*Attr{
				NewAttr("title", "").HasNoEqualSign(),
			},
			wantLength: 7,
			wantErr:    assert.NoError,
		},
		{
			name: "multiple with attributes with no equal sign",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[title colspan= foo]`),
			},
			wantTokens: []*Attr{
				NewAttr("title", "").HasNoEqualSign(),
				NewAttr("colspan", ""),
				NewAttr("foo", "").HasNoEqualSign(),
			},
			wantLength: 20,
			wantErr:    assert.NoError,
		},
		{
			name: "complex",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[title="Hello world!" colspan=3]`),
			},
			wantTokens: []*Attr{
				NewAttr("title", "Hello world!"),
				NewAttr("colspan", "3"),
			},
			wantLength: 32,
			wantErr:    assert.NoError,
		},
		{
			name: "complex, continued",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[title="Hello world!" colspan=3].class`),
			},
			wantTokens: []*Attr{
				NewAttr("title", "Hello world!"),
				NewAttr("colspan", "3"),
			},
			wantLength: 32,
			wantErr:    assert.NoError,
		},
		{
			name: "invalid first character",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`abc`),
			},
			wantTokens: nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name: "invalid second character",
			sut:  NewLexer(ModeHTML),
			args: args{
				runes: []rune(`[!]`),
			},
			wantTokens: nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotTokens, gotLength, err := tt.sut.FindAttributeTokens(tt.args.runes)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantLength, gotLength)
			assert.Equal(t, tt.wantTokens, gotTokens)
		})
	}
}

func TestLexer_FindClassToken(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantLength int
		wantToken  *AttrValue
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "simple",
			sut:  &Lexer{},
			args: args{
				runes: []rune(`.foo`),
			},
			wantLength: 4,
			wantToken:  NewClass("foo"),
			wantErr:    assert.NoError,
		},
		{
			name: "complex",
			sut:  &Lexer{},
			args: args{
				runes: []rune(`.foo$$$@-3`),
			},
			wantLength: 10,
			wantToken:  NewClass("foo").SetNumbering("$$$").SetReverse().SetStart(3),
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotToken, gotLength, err := tt.sut.FindClassToken(tt.args.runes)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantLength, gotLength)
			assert.Equal(t, tt.wantToken, gotToken)
		})
	}
}

func TestLexer_FindIDToken(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantLength int
		wantToken  *AttrValue
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "simple",
			sut:  &Lexer{},
			args: args{
				runes: []rune(`#foo`),
			},
			wantLength: 4,
			wantToken:  NewID("foo"),
			wantErr:    assert.NoError,
		},
		{
			name: "complex",
			sut:  &Lexer{},
			args: args{
				runes: []rune(`#foo$$$@-3`),
			},
			wantLength: 10,
			wantToken:  NewID("foo").SetNumbering("$$$").SetReverse().SetStart(3),
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotToken, gotLength, err := tt.sut.FindIDToken(tt.args.runes)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantLength, gotLength)
			assert.Equal(t, tt.wantToken, gotToken)
		})
	}
}

func TestLexer_FindAllAttributeTokens(t *testing.T) {
	t.Parallel()

	type args struct {
		token *TagToken
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantLength int
		wantToken  *TagToken
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "empty is valid",
			sut:  NewLexer(ModeHTML),
			args: args{
				token: &TagToken{},
				runes: []rune{},
			},
			wantLength: 0,
			wantToken:  &TagToken{},
			wantErr:    assert.NoError,
		},
		{
			name: "simple id",
			sut:  NewLexer(ModeHTML),
			args: args{
				token: &TagToken{},
				runes: []rune("#foo"),
			},
			wantLength: 4,
			wantToken: &TagToken{
				ID: NewID("foo"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "only the last id remains",
			sut:  NewLexer(ModeHTML),
			args: args{
				token: &TagToken{},
				runes: []rune("#foo#bar"),
			},
			wantLength: 4,
			wantToken: &TagToken{
				ID: NewID("foo"),
			},
			wantErr: assert.Error,
		},
		{
			name: "complex id",
			sut:  NewLexer(ModeHTML),
			args: args{
				token: &TagToken{},
				runes: []rune("#foo$$$@-3"),
			},
			wantLength: 10,
			wantToken: &TagToken{
				ID: NewID("foo").SetNumbering("$$$").SetReverse().SetStart(3),
			},
			wantErr: assert.NoError,
		},
		{
			name: "simple class",
			sut:  NewLexer(ModeHTML),
			args: args{
				token: &TagToken{},
				runes: []rune(".foo"),
			},
			wantLength: 4,
			wantToken: &TagToken{
				Classes: AttrValues{NewClass("foo")},
			},
			wantErr: assert.NoError,
		},
		{
			name: "attributes",
			sut:  NewLexer(ModeHTML),
			args: args{
				token: &TagToken{},
				runes: []rune(`[title="Hello world!" colspan=3]`),
			},
			wantLength: 32,
			wantToken: &TagToken{
				Attributes: AttrList{
					NewAttr("title", "Hello world!"),
					NewAttr("colspan", "3"),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.sut.FindAllAttributeTokens(tt.args.token, tt.args.runes)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantLength, got)
			assert.Equal(t, tt.wantToken, tt.args.token)
		})
	}
}

func TestLexer_FindRepeat(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantLength int
		wantRepeat int
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "empty is valid",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune{}},
			wantLength: 0,
			wantRepeat: 1,
			wantErr:    assert.NoError,
		},
		{
			name:       "single *",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("*")},
			wantLength: 0,
			wantRepeat: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "valid",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("*123")},
			wantLength: 4,
			wantRepeat: 123,
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotLength, err := tt.sut.FindRepeat(tt.args.runes)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantLength, gotLength)
			assert.Equal(t, tt.wantRepeat, got)
		})
	}
}

func TestLexer_NextTextToken(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantToken  *Text
		wantLength int
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "empty",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.NoError,
		},
		{
			name:       "skip",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("abc")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.NoError,
		},
		{
			name:       "invalid non-empty",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("{foo=bar")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "valid empty",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("{}")},
			wantToken:  nil,
			wantLength: 2,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid empty",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("{}")},
			wantToken:  nil,
			wantLength: 2,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid empty, continued",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("{}next")},
			wantToken:  nil,
			wantLength: 2,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid non-empty, simple",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("{foo bar}")},
			wantToken:  NewText("foo bar"),
			wantLength: 9,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid non-empty, simple, continued",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("{foo bar}, next")},
			wantToken:  NewText("foo bar"),
			wantLength: 9,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid non-empty, complex",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("{1928 !>&/ődisü BAR__?+ő}")},
			wantToken:  NewText("1928 !>&/ődisü BAR__?+ő"),
			wantLength: 25,
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotToken, gotLength, gotErr := tt.sut.NextTextToken(tt.args.runes)

			tt.wantErr(t, gotErr)
			assert.Equal(t, tt.wantToken, gotToken)
			assert.Equal(t, tt.wantLength, gotLength)
		})
	}
}

func TestLexer_NextTagToken(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantToken  *TagToken
		wantLength int
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "simple",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div")},
			wantToken:  NewTagToken("div", 1),
			wantLength: 3,
			wantErr:    assert.NoError,
		},
		{
			name:       "longer",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div+p")},
			wantToken:  NewTagToken("div", 1),
			wantLength: 3,
			wantErr:    assert.NoError,
		},
		{
			name:       "dash not supported",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div-p")},
			wantToken:  NewTagToken("div", 1),
			wantLength: 3,
			wantErr:    assert.NoError,
		},
		{
			name:       "underscore not supported",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div_p")},
			wantToken:  NewTagToken("div", 1),
			wantLength: 3,
			wantErr:    assert.NoError,
		},
		{
			name:       "strange character",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("divőp")},
			wantToken:  NewTagToken("div", 1),
			wantLength: 3,
			wantErr:    assert.NoError,
		},
		{
			name:       "empty",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid character",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("!p")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
		{
			name:       "valid multiplier",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div*25")},
			wantToken:  NewTagToken("div", 25),
			wantLength: 6,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid multiplier, continued",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div*25>p")},
			wantToken:  NewTagToken("div", 25),
			wantLength: 6,
			wantErr:    assert.NoError,
		},
		{
			name:       "invalid * if not continued",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div*")},
			wantToken:  nil,
			wantLength: 3,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid * without valid number",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div*p")},
			wantToken:  nil,
			wantLength: 3,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid empty class",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("div.")},
			wantToken:  nil,
			wantLength: 3,
			wantErr:    assert.Error,
		},
		{
			name: "valid with simple id",
			sut:  NewLexer(ModeHTML),
			args: args{runes: []rune("div#foo")},
			wantToken: NewTagToken("div", 1).
				SetID(NewID("foo")),
			wantLength: 7,
			wantErr:    assert.NoError,
		},
		{
			name: "valid with complex id",
			sut:  NewLexer(ModeHTML),
			args: args{runes: []rune("div#foo@-25")},
			wantToken: NewTagToken("div", 1).
				SetID(NewID("foo").SetReverse().SetStart(25)),
			wantLength: 11,
			wantErr:    assert.NoError,
		},
		{
			name: "valid with complex id and classes",
			sut:  NewLexer(ModeHTML),
			args: args{runes: []rune("div.bar#foo@-25.baz$$$@3")},
			wantToken: NewTagToken("div", 1).
				SetID(NewID("foo").SetReverse().SetStart(25)).
				AddClass(NewClass("bar")).
				AddClass(NewClass("baz").SetNumbering("$$$").SetStart(3)),
			wantLength: 24,
			wantErr:    assert.NoError,
		},
		{
			name: "valid, super complex attribute list",
			sut:  NewLexer(ModeHTML),
			args: args{runes: []rune(`div#footer.foo@2.bar2$@22.baz$$$@-3[qux="byE bye" quix quuix=]*4`)},
			wantToken: NewTagToken("div", 4).
				SetID(NewID("footer")).
				AddClass(NewClass("foo").SetStart(2)).
				AddClass(NewClass("bar2").SetNumbering("$").SetStart(22)).
				AddClass(NewClass("baz").SetNumbering("$$$").SetReverse().SetStart(3)).
				AddAttribute(NewAttr("qux", "byE bye")).
				AddAttribute(NewAttr("quix", "").HasNoEqualSign()).
				AddAttribute(NewAttr("quuix", "")),
			wantLength: 64,
			wantErr:    assert.NoError,
		},
		{
			name:       "xml tag is recognized in html mode",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(`X:Foo`)},
			wantToken:  NewTagToken("X:Foo", 1),
			wantLength: 5,
			wantErr:    assert.NoError,
		},
		{
			name:       "xml tag is recognized in xml mode",
			sut:        NewLexer(ModeXML),
			args:       args{runes: []rune(`X:Foo`)},
			wantToken:  NewTagToken("X:Foo", 1),
			wantLength: 5,
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotToken, gotLength, gotErr := tt.sut.NextTagToken(tt.args.runes)

			tt.wantErr(t, gotErr)
			assert.Equal(t, tt.wantLength, gotLength)
			assert.Equal(t, tt.wantToken, gotToken)
		})
	}
}

func TestLexer_NextDirectiveToken(t *testing.T) {
	t.Parallel()

	type args struct {
		runes []rune
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantToken  *DirectiveToken
		wantLength int
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "valid +",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("+")},
			wantToken:  NewDirectiveToken(Add, 1),
			wantLength: 1,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid +, continued",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("+next")},
			wantToken:  NewDirectiveToken(Add, 1),
			wantLength: 1,
			wantErr:    assert.NoError,
		},
		{
			name:       "invalid ++",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("++")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid ++++++",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("++")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "valid >",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(">")},
			wantToken:  NewDirectiveToken(Dive, 1),
			wantLength: 1,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid >, continued",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(">next")},
			wantToken:  NewDirectiveToken(Dive, 1),
			wantLength: 1,
			wantErr:    assert.NoError,
		},
		{
			name:       "invalid >>",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(">>")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid >>>>>>",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(">>")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "valid, single ^",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("^")},
			wantToken:  NewDirectiveToken(Ascend, 1),
			wantLength: 1,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid, multiple ^",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("^^^")},
			wantToken:  NewDirectiveToken(Ascend, 3),
			wantLength: 3,
			wantErr:    assert.NoError,
		},
		{
			name:       "valid, multiple ^, continued",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("^^^next")},
			wantToken:  NewDirectiveToken(Ascend, 3),
			wantLength: 3,
			wantErr:    assert.NoError,
		},
		{
			name:       "invalid, ^>",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("^>")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "invalid, >^",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune(">^")},
			wantToken:  nil,
			wantLength: 1,
			wantErr:    assert.Error,
		},
		{
			name:       "empty",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune{}},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.NoError,
		},
		{
			name:       "invalid",
			sut:        NewLexer(ModeHTML),
			args:       args{runes: []rune("%")},
			wantToken:  nil,
			wantLength: 0,
			wantErr:    assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotTokens, gotLength, gotErr := tt.sut.NextDirectiveToken(tt.args.runes)

			tt.wantErr(t, gotErr)
			assert.Equalf(t, tt.wantToken, gotTokens, "NextDirectiveToken(%v)", tt.args.runes)
			assert.Equalf(t, tt.wantLength, gotLength, "NextDirectiveToken(%v)", tt.args.runes)
		})
	}
}

func TestLexer_Tokenize(t *testing.T) {
	t.Parallel()

	type args struct {
		s       string
		inGroup bool
	}
	tests := []struct {
		name       string
		sut        *Lexer
		args       args
		wantTokens []Token
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "ksd",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ksd", inGroup: false},
			wantTokens: []Token{
				NewTagToken("ksd", 1),
			},
			wantErr: assert.NoError,
		},
		{
			name: "ksd>idu",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ksd>idu", inGroup: false},
			wantTokens: []Token{
				NewTagToken("ksd", 1).
					AddChildren(NewTagToken("idu", 1)),
			},
			wantErr: assert.NoError,
		},
		{
			name: "ksd>idu+blabla",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ksd>idu+blabla"},
			wantTokens: []Token{
				NewTagToken("ksd", 1).
					AddChildren(
						NewTagToken("idu", 1),
						NewTagToken("blabla", 1),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "ksd.foo>idu.bar+blabla.baz.quix",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ksd.foo>idu.bar+blabla.baz.quix"},
			wantTokens: []Token{
				NewTagToken("ksd", 1).
					AddClass(NewClass("foo")).
					AddChildren(
						NewTagToken("idu", 1).
							AddClass(NewClass("bar")),
						NewTagToken("blabla", 1).
							AddClass(NewClass("baz")).
							AddClass(NewClass("quix")),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 1 - children",
			sut:  NewLexer(ModeHTML),
			args: args{s: "div>ul>li"},
			wantTokens: []Token{
				NewTagToken("div", 1).
					AddChildren(
						NewTagToken("ul", 1).
							AddChildren(
								NewTagToken("li", 1),
							),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 2 - siblings",
			sut:  NewLexer(ModeHTML),
			args: args{s: "div+p+bq"},
			wantTokens: []Token{
				NewTagToken("div", 1),
				NewTagToken("p", 1),
				NewTagToken("bq", 1),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 3 - mixed",
			sut:  NewLexer(ModeHTML),
			args: args{s: "div+div>p>span+em"},
			wantTokens: []Token{
				NewTagToken("div", 1),
				NewTagToken("div", 1).
					AddChildren(
						NewTagToken("p", 1).
							AddChildren(
								NewTagToken("span", 1),
								NewTagToken("em", 1),
							),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 4 - climb up",
			sut:  NewLexer(ModeHTML),
			args: args{s: "div+div>p>span+em^bq"},
			wantTokens: []Token{
				NewTagToken("div", 1),
				NewTagToken("div", 1).
					AddChildren(
						NewTagToken("p", 1).
							AddChildren(
								NewTagToken("span", 1),
								NewTagToken("em", 1),
							),
						NewTagToken("bq", 1),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 5 - climb up repeated",
			sut:  NewLexer(ModeHTML),
			args: args{s: "div+div>p>span+em^^bq"},
			wantTokens: []Token{
				NewTagToken("div", 1),
				NewTagToken("div", 1).
					AddChildren(
						NewTagToken("p", 1).
							AddChildren(
								NewTagToken("span", 1),
								NewTagToken("em", 1),
							),
					),
				NewTagToken("bq", 1),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 6 - ul>li*5",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ul>li*5"},
			wantTokens: []Token{
				NewTagToken("ul", 1).
					AddChildren(
						NewTagToken("li", 5),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 7 - groups",
			sut:  NewLexer(ModeHTML),
			args: args{s: "div>(header>ul>li*2>a)+footer>p"},
			wantTokens: []Token{
				NewTagToken("div", 1).AddChildren(
					NewGroupToken(1).AddChildren(
						NewTagToken("header", 1).AddChildren(
							NewTagToken("ul", 1).AddChildren(
								NewTagToken("li", 2).AddChildren(
									NewTagToken("a", 1),
								),
							),
						),
					),
					NewTagToken("footer", 1).
						AddChildren(NewTagToken("p", 1)),
				),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 8 - ids and classes",
			sut:  NewLexer(ModeHTML),
			args: args{s: "div#header+div.page+div#footer.class1.class2.class3"},
			wantTokens: []Token{
				NewTagToken("div", 1).
					SetID(NewID("header")),
				NewTagToken("div", 1).
					AddClass(NewClass("page")),
				NewTagToken("div", 1).
					SetID(NewID("footer")).
					AddClass(NewClass("class1")).
					AddClass(NewClass("class2")).
					AddClass(NewClass("class3")),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 9 - attributes",
			sut:  NewLexer(ModeHTML),
			args: args{s: "td[title=\"Hello world!\" colspan=3]"},
			wantTokens: []Token{
				NewTagToken("td", 1).
					AddAttribute(NewAttr("title", "Hello world!")).
					AddAttribute(NewAttr("colspan", "3")),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 10 - item numbering #1",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ul>li.item$*5"},
			wantTokens: []Token{
				NewTagToken("ul", 1).
					AddChildren(
						NewTagToken("li", 5).
							AddClass(NewClass("item").SetNumbering("$")),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 11 - item numbering #2",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ul>li.item$@-*5"},
			wantTokens: []Token{
				NewTagToken("ul", 1).
					AddChildren(
						NewTagToken("li", 5).
							AddClass(NewClass("item").SetNumbering("$").SetReverse()),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 12 - item numbering #3",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ul>li.item$@3*5"},
			wantTokens: []Token{
				NewTagToken("ul", 1).
					AddChildren(
						NewTagToken("li", 5).
							AddClass(NewClass("item").SetNumbering("$").SetStart(3)),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 13 - item numbering #4",
			sut:  NewLexer(ModeHTML),
			args: args{s: "ul>li.item$@-3*5"},
			wantTokens: []Token{
				NewTagToken("ul", 1).
					AddChildren(
						NewTagToken("li", 5).
							AddClass(NewClass("item").SetNumbering("$").SetReverse().SetStart(3)),
					),
			},
			wantErr: assert.NoError,
		},
		{
			name: "official 14 - text",
			sut:  NewLexer(ModeHTML),
			args: args{s: "a{Click me}"},
			wantTokens: []Token{
				NewTagToken("a", 1).
					SetText(NewText("Click me")),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wantLength := len([]rune(tt.args.s))

			gotTokens, gotLength, err := tt.sut.Tokenize([]rune(tt.args.s), tt.args.inGroup)

			tt.wantErr(t, err)
			assert.Equal(t, wantLength, gotLength)
			assert.Equal(t, tt.wantTokens, gotTokens)
		})
	}
}
