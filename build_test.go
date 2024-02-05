package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	t.Parallel()

	type args struct {
		tokens       []Token
		mode         Mode
		indentation  string
		depth        int
		multiline    bool
		num          int
		siblingCount int
	}
	tests := []struct {
		name     string
		args     args
		wantHTML string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "p.foo",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).AddClass(NewClass("foo")),
				},
				mode:        ModeHTML,
				indentation: "",
				depth:       0,
				multiline:   false,
				num:         1,
			},
			wantHTML: `<p class="foo"></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "p.foo - xml",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).AddClass(NewClass("foo")),
				},
				mode:         ModeXML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p class="foo" />`,
			wantErr:  assert.NoError,
		},
		{
			name: "p.foo>span.bar",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).
						AddClass(NewClass("foo")).
						AddChildren(
							NewTagToken("span", 1).
								AddClass(NewClass("bar")),
						),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p class="foo"><span class="bar"></span></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "p#foo+p#bar",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).SetID(NewID("foo")),
					NewTagToken("p", 1).SetID(NewID("bar")),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p id="foo"></p><p id="bar"></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "p#foo+p#bar>span",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).SetID(NewID("foo")),
					NewTagToken("p", 1).SetID(NewID("bar")).AddChildren(
						NewTagToken("span", 1),
					),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p id="foo"></p><p id="bar"><span></span></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "p#foo+p#bar>span - xml / multiline",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).SetID(NewID("foo")),
					NewTagToken("p", 1).SetID(NewID("bar")).AddChildren(
						NewTagToken("span", 1),
					),
				},
				mode:         ModeXML,
				indentation:  "    ",
				depth:        0,
				multiline:    true,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p id="foo" />
<p id="bar">
    <span />
</p>
`,
			wantErr: assert.NoError,
		},
		{
			name: "p#foo>p#bar[baz=qux]",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).SetID(NewID("foo")).AddChildren(
						NewTagToken("p", 1).
							SetID(NewID("bar")).
							AddAttribute(NewAttr("baz", "qux")),
					),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p id="foo"><p id="bar" baz="qux"></p></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "div#foo>p#bar{This is a text}+div#baz>ul#qux",
			args: args{
				tokens: []Token{
					NewTagToken("div", 1).SetID(NewID("foo")).AddChildren(
						NewTagToken("p", 1).
							SetID(NewID("bar")).
							SetText(NewText("This is a text")),
						NewTagToken("div", 1).SetID(NewID("baz")).AddChildren(
							NewTagToken("ul", 1).
								SetID(NewID("qux")),
						),
					),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<div id="foo"><p id="bar">This is a text</p><div id="baz"><ul id="qux"></ul></div></div>`,
			wantErr:  assert.NoError,
		},
		{
			name: "div#foo>p#bar{This is a text}+div#baz>ul#qux^^p.quix",
			args: args{
				tokens: []Token{
					NewTagToken("div", 1).SetID(NewID("foo")).AddChildren(
						NewTagToken("p", 1).SetID(NewID("bar")).SetText(NewText("This is a text")),
						NewTagToken("div", 1).SetID(NewID("baz")).AddChildren(
							NewTagToken("ul", 1).SetID(NewID("qux")),
						),
					),
					NewTagToken("p", 1).AddClass(NewClass("quix")),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<div id="foo"><p id="bar">This is a text</p><div id="baz"><ul id="qux"></ul></div></div><p class="quix"></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "p#foo+p.bar*3",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).SetID(NewID("foo")),
					NewTagToken("p", 3).AddClass(NewID("bar")),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p id="foo"></p><p class="bar"></p><p class="bar"></p><p class="bar"></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "p#foo+p.bar$*3",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).SetID(NewID("foo")),
					NewTagToken("p", 3).AddClass(NewID("bar").SetNumbering("$")),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p id="foo"></p><p class="bar1"></p><p class="bar2"></p><p class="bar3"></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "p#foo+p.bar$$$@-12*3",
			args: args{
				tokens: []Token{
					NewTagToken("p", 1).SetID(NewID("foo")),
					NewTagToken("p", 3).AddClass(NewID("bar").SetNumbering("$$$").SetStart(12).SetReverse()),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<p id="foo"></p><p class="bar014"></p><p class="bar013"></p><p class="bar012"></p>`,
			wantErr:  assert.NoError,
		},
		{
			name: "ul>(li>a)*3",
			args: args{
				tokens: []Token{
					NewTagToken("ul", 1).AddChildren(
						NewGroupToken(3).AddChildren(
							NewTagToken("li", 1).AddChildren(
								NewTagToken("a", 1),
							),
						),
					),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<ul><li><a></a></li><li><a></a></li><li><a></a></li></ul>`,
			wantErr:  assert.NoError,
		},
		{
			name: "tbody>tr*3>td*2",
			args: args{
				tokens: []Token{
					NewTagToken("tbody", 1).AddChildren(
						NewTagToken("tr", 3).AddChildren(
							NewTagToken("td", 2),
						),
					),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: strings.Join([]string{
				`<tbody>`,
				`<tr><td></td><td></td></tr>`,
				`<tr><td></td><td></td></tr>`,
				`<tr><td></td><td></td></tr>`,
				`</tbody>`,
			}, ""),
			wantErr: assert.NoError,
		},
		{
			name: "body>((table>tbody>tr*3>td*4)+hr)*2",
			args: args{
				tokens: []Token{
					NewTagToken("body", 1).AddChildren(
						NewGroupToken(2).AddChildren(
							NewTagToken("table", 1).AddChildren(
								NewTagToken("tbody", 1).AddChildren(
									NewTagToken("tr", 3).AddChildren(
										NewTagToken("td", 4),
									),
								),
							),
							NewTagToken("hr", 1),
						),
					),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: strings.Join([]string{
				`<body>`,
				`<table><tbody>`,
				`<tr><td></td><td></td><td></td><td></td></tr>`,
				`<tr><td></td><td></td><td></td><td></td></tr>`,
				`<tr><td></td><td></td><td></td><td></td></tr>`,
				`</tbody></table>`,
				`<hr>`,
				`<table><tbody>`,
				`<tr><td></td><td></td><td></td><td></td></tr>`,
				`<tr><td></td><td></td><td></td><td></td></tr>`,
				`<tr><td></td><td></td><td></td><td></td></tr>`,
				`</tbody></table>`,
				`<hr>`,
				`</body>`,
			}, ""),
			wantErr: assert.NoError,
		},
		{
			name: "ul>li.item$$*3>a",
			args: args{
				tokens: []Token{
					NewTagToken("ul", 1).AddChildren(
						NewTagToken("li", 3).
							AddClass(NewClass("item").SetNumbering("$$")).
							AddChildren(NewTagToken("a", 1)),
					),
				},
				mode:         ModeHTML,
				indentation:  "  ",
				depth:        0,
				multiline:    true,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<ul>
  <li class="item01">
    <a></a>
  </li>
  <li class="item02">
    <a></a>
  </li>
  <li class="item03">
    <a></a>
  </li>
</ul>
`,
			wantErr: assert.NoError,
		},
		{
			name: "ul>(li.item$$>a)*3",
			args: args{
				tokens: []Token{
					NewTagToken("ul", 1).AddChildren(
						NewGroupToken(3).AddChildren(
							NewTagToken("li", 1).
								AddClass(NewClass("item").SetNumbering("$$")).
								AddChildren(NewTagToken("a", 1)),
						),
					),
				},
				mode:         ModeHTML,
				indentation:  "  ",
				depth:        0,
				multiline:    true,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<ul>
  <li class="item01">
    <a></a>
  </li>
  <li class="item02">
    <a></a>
  </li>
  <li class="item03">
    <a></a>
  </li>
</ul>
`,
			wantErr: assert.NoError,
		},
		{
			name: "div.div$$*3>p.item$$>span",
			args: args{
				tokens: []Token{
					NewTagToken("div", 3).
						AddClass(NewClass("div").SetNumbering("$$").SetReverse()).
						AddChildren(
							NewTagToken("p", 1).
								AddClass(NewClass("item").SetNumbering("$$")).
								AddChildren(NewTagToken("span", 1)),
						),
				},
				mode:         ModeHTML,
				indentation:  "  ",
				depth:        0,
				multiline:    true,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<div class="div03">
  <p class="item01">
    <span></span>
  </p>
</div>
<div class="div02">
  <p class="item02">
    <span></span>
  </p>
</div>
<div class="div01">
  <p class="item03">
    <span></span>
  </p>
</div>
`,
			wantErr: assert.NoError,
		},
		{
			name: `span.span$$@-8`,
			args: args{
				tokens: []Token{
					NewTagToken("span", 1).
						AddClass(NewClass("span").SetNumbering("$$").SetStart(8).SetReverse()),
				},
				mode:         ModeHTML,
				indentation:  "",
				depth:        0,
				multiline:    false,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<span class="span08"></span>`,
		},
		{
			name: `td*2>span.span$$@-8`,
			args: args{
				tokens: []Token{
					NewTagToken("td", 2).AddChildren(
						NewTagToken("span", 1).
							AddClass(NewClass("span").SetNumbering("$$").SetStart(8).SetReverse()),
					),
				},
				mode:         ModeHTML,
				indentation:  "  ",
				depth:        1,
				multiline:    true,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `  <td>
    <span class="span09"></span>
  </td>
  <td>
    <span class="span08"></span>
  </td>
`,
		},
		{
			name: `tr*3>td*2>span.span$$@-8`,
			args: args{
				tokens: []Token{
					NewTagToken("tr", 3).AddChildren(
						NewTagToken("td", 2).AddChildren(
							NewTagToken("span", 1).
								AddClass(NewClass("span").SetNumbering("$$").SetStart(8).SetReverse()),
						),
					),
				},
				mode:         ModeHTML,
				indentation:  "  ",
				depth:        0,
				multiline:    true,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<tr>
  <td>
    <span class="span09"></span>
  </td>
  <td>
    <span class="span08"></span>
  </td>
</tr>
<tr>
  <td>
    <span class="span09"></span>
  </td>
  <td>
    <span class="span08"></span>
  </td>
</tr>
<tr>
  <td>
    <span class="span09"></span>
  </td>
  <td>
    <span class="span08"></span>
  </td>
</tr>
`,
		},
		{
			name: `html>body>(table.table$$>tbody>tr*3>td.hi$*3+td.hello$$@-2)*2`,
			args: args{
				tokens: []Token{
					NewTagToken("html", 1).AddChildren(
						NewTagToken("body", 1).AddChildren(
							NewGroupToken(2).AddChildren(
								NewTagToken("table", 1).
									AddClass(NewClass("table").SetNumbering("$$")).
									AddChildren(
										NewTagToken("tbody", 1).AddChildren(
											NewTagToken("tr", 3).AddChildren(
												NewTagToken("td", 4).AddClass(NewClass("hi").SetNumbering("$")),
												NewTagToken("td", 1).AddClass(NewClass("hello").SetNumbering("$$").SetStart(2).SetReverse()),
											),
										),
									),
							),
						),
					),
				},
				mode:         ModeHTML,
				indentation:  "  ",
				depth:        0,
				multiline:    true,
				num:          1,
				siblingCount: 1,
			},
			wantHTML: `<html>
  <body>
    <table class="table01">
      <tbody>
        <tr>
          <td class="hi1"></td>
          <td class="hi2"></td>
          <td class="hi3"></td>
          <td class="hi4"></td>
          <td class="hello04"></td>
        </tr>
        <tr>
          <td class="hi1"></td>
          <td class="hi2"></td>
          <td class="hi3"></td>
          <td class="hi4"></td>
          <td class="hello03"></td>
        </tr>
        <tr>
          <td class="hi1"></td>
          <td class="hi2"></td>
          <td class="hi3"></td>
          <td class="hi4"></td>
          <td class="hello02"></td>
        </tr>
      </tbody>
    </table>
    <table class="table02">
      <tbody>
        <tr>
          <td class="hi1"></td>
          <td class="hi2"></td>
          <td class="hi3"></td>
          <td class="hi4"></td>
          <td class="hello04"></td>
        </tr>
        <tr>
          <td class="hi1"></td>
          <td class="hi2"></td>
          <td class="hi3"></td>
          <td class="hi4"></td>
          <td class="hello03"></td>
        </tr>
        <tr>
          <td class="hi1"></td>
          <td class="hi2"></td>
          <td class="hi3"></td>
          <td class="hi4"></td>
          <td class="hello02"></td>
        </tr>
      </tbody>
    </table>
  </body>
</html>
`,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			elements := Build(tt.args.tokens, tt.args.num, tt.args.siblingCount)

			gotHTML := elements.HTML(tt.args.mode, tt.args.indentation, tt.args.depth, tt.args.multiline)

			assert.Equal(t, tt.wantHTML, gotHTML)
		})
	}
}

func TestCreateElem(t *testing.T) {
	t.Parallel()

	type args struct {
		token        *TagToken
		num          int
		siblingCount int
	}
	tests := []struct {
		name string
		args args
		want ElemList
	}{
		{
			name: "default",
			args: args{
				token: NewTagToken("div", 1).
					SetID(NewID("foo").SetStart(3).SetNumbering("$$")).
					AddClass(NewClass("bar").SetReverse().SetStart(2)).
					AddAttribute(NewAttr("baz", "qux")).
					SetText(NewText("This is a text")),
				num:          1,
				siblingCount: 1,
			},
			want: ElemList{
				&Elem{
					Name:         "div",
					ID:           NewID("foo").SetStart(3).SetNumbering("$$"),
					Classes:      AttrValues{NewClass("bar").SetReverse().SetStart(2)},
					Attributes:   AttrList{NewAttr("baz", "qux")},
					Text:         NewText("This is a text"),
					Num:          1,
					SiblingCount: 1,
				},
			},
		},
		{
			name: "default 3x",
			args: args{
				token: NewTagToken("div", 3).
					SetID(NewID("foo").SetStart(3).SetNumbering("$$")).
					AddClass(NewClass("bar").SetReverse().SetStart(2)).
					AddAttribute(NewAttr("baz", "qux")).
					SetText(NewText("This is a text")),
				num:          1,
				siblingCount: 1,
			},
			want: ElemList{
				&Elem{
					Name:         "div",
					ID:           NewID("foo").SetStart(3).SetNumbering("$$"),
					Classes:      AttrValues{NewClass("bar").SetReverse().SetStart(2)},
					Attributes:   AttrList{NewAttr("baz", "qux")},
					Text:         NewText("This is a text"),
					Num:          1,
					SiblingCount: 3,
				},
				&Elem{
					Name:         "div",
					ID:           NewID("foo").SetStart(3).SetNumbering("$$"),
					Classes:      AttrValues{NewClass("bar").SetReverse().SetStart(2)},
					Attributes:   AttrList{NewAttr("baz", "qux")},
					Text:         NewText("This is a text"),
					Num:          2,
					SiblingCount: 3,
				},
				&Elem{
					Name:         "div",
					ID:           NewID("foo").SetStart(3).SetNumbering("$$"),
					Classes:      AttrValues{NewClass("bar").SetReverse().SetStart(2)},
					Attributes:   AttrList{NewAttr("baz", "qux")},
					Text:         NewText("This is a text"),
					Num:          3,
					SiblingCount: 3,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := BuildFromTag(tt.args.token, tt.args.num, tt.args.siblingCount)

			assert.Equal(t, tt.want, got)
		})
	}
}
