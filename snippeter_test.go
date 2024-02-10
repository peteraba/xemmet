package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnippeter_Walk(t *testing.T) {
	t.Parallel()

	type fields struct {
		mode Mode
	}
	type args struct {
		tokens []Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *TagToken
	}{
		{
			name: "a",
			fields: fields{
				mode: ModeHTML,
			},
			args: args{
				tokens: []Token{
					NewTagToken("a", 1),
				},
			},
			want: NewTagToken("a", 1).
				AddAttribute(NewAttr("href", "#")),
		},
		{
			name: "a - overwrite",
			fields: fields{
				mode: ModeHTML,
			},
			args: args{
				tokens: []Token{
					NewTagToken("a", 1).
						AddAttribute(NewAttr("href", "foo")),
				},
			},
			want: NewTagToken("a", 1).
				AddAttribute(NewAttr("href", "foo")),
		},
		{
			name: "area",
			fields: fields{
				mode: ModeHTML,
			},
			args: args{
				tokens: []Token{
					NewTagToken("area", 1),
				},
			},
			want: NewTagToken("area", 1).
				AddAttribute(NewAttr("coords", "")).
				AddAttribute(NewAttr("href", "")).
				AddAttribute(NewAttr("alt", "")).
				AddAttribute(NewAttr("shape", "")),
		},
		{
			name: "area:d",
			fields: fields{
				mode: ModeHTML,
			},
			args: args{
				tokens: []Token{
					NewTagToken("area:d", 1),
				},
			},
			want: NewTagToken("area", 1).
				AddAttribute(NewAttr("coords", "")).
				AddAttribute(NewAttr("href", "")).
				AddAttribute(NewAttr("alt", "")).
				AddAttribute(NewAttr("shape", "default")),
		},
		{
			name: "bq",
			fields: fields{
				mode: ModeHTML,
			},
			args: args{
				tokens: []Token{
					NewTagToken("bq", 1),
				},
			},
			want: NewTagToken("blockquote", 1),
		},
		{
			name: "ifr",
			fields: fields{
				mode: ModeHTML,
			},
			args: args{
				tokens: []Token{
					NewTagToken("ifr", 1),
				},
			},
			want: NewTagToken("iframe", 1).
				AddAttribute(NewAttr("src", "")).
				AddAttribute(NewAttr("frameborder", "0")),
		},
		{
			name: "input:q",
			fields: fields{
				mode: ModeHTMX,
			},
			args: args{
				tokens: []Token{
					NewTagToken("input:q", 1),
				},
			},
			want: NewTagToken("input", 1).
				AddAttribute(NewAttr("name", "q")).
				AddAttribute(NewAttr("type", "search")).
				AddAttribute(NewAttr("hx-get", "")).
				AddAttribute(NewAttr("hx-trigger", "keyup changed delay:500ms")).
				AddAttribute(NewAttr("hx-target", "")).
				AddAttribute(NewAttr("hx-swap", "innerHTML")).
				AddAttribute(NewAttr("placeholder", "")),
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sut := NewSnippeter(tt.fields.mode)

			got := sut.Walk(tt.args.tokens)

			assert.NotNil(t, got)
			assert.Len(t, got, 1)
			assert.IsType(t, &TagToken{}, got[0])

			g, ok := got[0].(*TagToken)
			assert.True(t, ok)
			assert.Equal(t, tt.want.Name, g.Name)
			assert.Equal(t, tt.want.ID, g.ID)
			assert.Equal(t, tt.want.Classes, g.Classes)
			assert.Equal(t, tt.want.Attributes, g.Attributes)
		})
	}
}
