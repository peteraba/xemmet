package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElem_HTML(t *testing.T) {
	t.Parallel()

	type args struct {
		mode           Mode
		indentation    string
		depth          int
		multiline      bool
		tabStopWrapper string
	}
	tests := []struct {
		name string
		sut  Elem
		args args
		want string
	}{
		{
			name: "empty",
			sut:  Elem{},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: "",
		},
		{
			name: "text",
			sut: Elem{
				Text: NewText("foo"),
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: "foo",
		},
		{
			name: "simple div",
			sut: Elem{
				Name: "div",
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: "<div></div>",
		},
		{
			name: "simple div in XML mode",
			sut: Elem{
				Name: "div",
			},
			args: args{
				mode:           ModeXML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: "<div />",
		},
		{
			name: "simple br in HTML mode",
			sut: Elem{
				Name: "br",
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: "<br>",
		},
		{
			name: "simple br in XML mode",
			sut: Elem{
				Name: "br",
			},
			args: args{
				mode:           ModeXML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: "<br />",
		},
		{
			name: "div with attributes",
			sut: Elem{
				Name:    "div",
				ID:      NewID("foo"),
				Classes: AttrValues{NewClass("bar"), NewClass("baz")},
				Attributes: AttrList{
					NewAttr("style", "background-color: red;"),
					NewAttr("hello", "bye"),
				},
				Text:         NewText("Hello, World!"),
				Num:          1,
				SiblingCount: 1,
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: `<div id="foo" class="bar baz" style="background-color: red;" hello="bye">Hello, World!</div>`,
		},
		{
			name: "div with children and attributes, single line",
			sut: Elem{
				Name:    "div",
				ID:      NewID("foo"),
				Classes: AttrValues{NewClass("bar"), NewClass("baz")},
				Attributes: AttrList{
					NewAttr("style", "background-color: red;"),
					NewAttr("hello", "bye"),
				},
				Children: []*Elem{
					{
						Name: "p",
						Text: NewText("Hello, "),
						Children: []*Elem{
							{
								Name: "span",
								Text: NewText("World"),
							},
							{
								Text: NewText("!"),
							},
						},
					},
					{
						Name: "p",
						Text: NewText("Aloha!"),
					},
				},
				Num:          1,
				SiblingCount: 1,
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: `<div id="foo" class="bar baz" style="background-color: red;" hello="bye"><p>Hello, <span>World</span>!</p><p>Aloha!</p></div>`,
		},
		{
			name: "div with children and attributes, multiline",
			sut: Elem{
				Name: "div",
				ID:   NewID("foo"),
				Classes: AttrValues{
					NewClass("bar"),
					NewClass("baz"),
				},
				Attributes: AttrList{
					NewAttr("style", "background-color: red;"),
					NewAttr("hello", "bye"),
				},
				Children: []*Elem{
					{
						Name: "p",
						Text: NewText("Hello, "),
						Children: []*Elem{
							{
								Name: "span",
								Text: NewText("World"),
							},
							{
								Text: NewText("!"),
							},
						},
					},
					{
						Name: "p",
						Text: NewText("Aloha!"),
					},
				},
				Num:          1,
				SiblingCount: 1,
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "  ",
				depth:          0,
				multiline:      true,
				tabStopWrapper: "",
			},
			want: `<div id="foo" class="bar baz" style="background-color: red;" hello="bye">
  <p>
    Hello, 
    <span>
      World
    </span>
    !
  </p>
  <p>
    Aloha!
  </p>
</div>
`,
		},
		{
			name: "attributes #1",
			sut: Elem{
				Name: "div",
				ID:   NewID("foo").SetStart(5).SetReverse(),
				Classes: AttrValues{
					NewClass("bar").SetStart(5),
					NewClass("baz").SetStart(5).SetNumbering("$$"),
				},
				Text:         NewText("Hello, World!"),
				Num:          1,
				SiblingCount: 1,
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: `<div id="foo5" class="bar5 baz05">Hello, World!</div>`,
		},
		{
			name: "attributes #2",
			sut: Elem{
				Name: "div",
				ID:   NewID("foo").SetStart(5).SetReverse(),
				Classes: AttrValues{
					NewClass("bar").SetStart(5),
					NewClass("baz").SetStart(5).SetNumbering("$$"),
				},
				Text:         NewText("Hello, World!"),
				Num:          1,
				SiblingCount: 2,
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: `<div id="foo6" class="bar5 baz05">Hello, World!</div>`,
		},
		{
			name: "attributes #3",
			sut: Elem{
				Name: "div",
				ID:   NewID("foo").SetStart(5).SetReverse(),
				Classes: AttrValues{
					NewClass("bar").SetStart(5),
					NewClass("baz").SetStart(5).SetNumbering("$$"),
				},
				Text:         NewText("Hello, World!"),
				Num:          2,
				SiblingCount: 2,
			},
			args: args{
				mode:           ModeHTML,
				indentation:    "",
				depth:          0,
				multiline:      false,
				tabStopWrapper: "",
			},
			want: `<div id="foo5" class="bar6 baz06">Hello, World!</div>`,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.sut.HTML(tt.args.mode, tt.args.indentation, tt.args.depth, tt.args.multiline, tt.args.tabStopWrapper)

			assert.Equal(t, tt.want, got)
		})
	}
}
