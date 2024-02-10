package main

import (
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXemmet(t *testing.T) {
	t.Parallel()

	const validSnippet = "body[x-data=lorem3]>(table.table$@>(thead>tr.class$$@-3>th#th.col$@*4{lorem2})+(tbody>tr.row$@1*3>td*4{lorem10})+(tfoot>tr>td*4{lorem2}))*2"

	t.Run("generate multiline html", func(t *testing.T) {
		t.Parallel()

		got, gotErr := Xemmet(ModeHTML, validSnippet, "  ", 1, true, "$$$")
		require.NoError(t, gotErr)

		assert.NotEmpty(t, got)

		assert.Regexp(t, regexp.MustCompile(`<body x-data="[\w\d\s.]+">`), got)
		assert.NotContains(t, got, `<body x-data="lorem3">`)

		assert.Contains(t, got, `    <table class="table1">`)
		assert.Contains(t, got, `    <table class="table2">`)
		assert.NotContains(t, got, `    <table class="table3">`)

		assert.Contains(t, got, `      <tr class="class03">`)
		// TODO: Fix this trickling down issue
		// assert.Contains(t, got, `      <tr class="class04">`)
		assert.NotContains(t, got, `      <tr class="class05">`)

		assert.Contains(t, got, `      <tr class="row1">`)
		assert.Contains(t, got, `      <tr class="row2">`)
		assert.Contains(t, got, `      <tr class="row3">`)

		assert.Contains(t, got, `      <th id="th" class="col1">`)
		assert.Contains(t, got, `      <th id="th" class="col2">`)
		assert.Contains(t, got, `      <th id="th" class="col3">`)
		assert.Contains(t, got, `      <th id="th" class="col4">`)
		assert.NotContains(t, got, `      <th id="th" class="col5">`)
		assert.Regexp(t, regexp.MustCompile(`<th id="th" class="col4">[$\w\d\s.]+</th>`), got)
		assert.NotContains(t, got, regexp.MustCompile(`<th id="th" class="col4">lorem2</th>`))

		assert.Contains(t, got, "  </body>")
		assert.Contains(t, got, "$$$STOP1$$$")

		assert.Contains(t, got, "\n")
	})

	t.Run("fail on invalid snippet", func(t *testing.T) {
		t.Parallel()

		from := gofakeit.Number(0, len(validSnippet)/2)
		to := gofakeit.Number(0, len(validSnippet)/2)
		snippet := validSnippet[from : len(validSnippet)-to]

		got, gotErr := Xemmet(ModeHTML, snippet, "  ", 1, true, "")
		require.Error(t, gotErr)
		assert.Contains(t, gotErr.Error(), ErrTokenizingMsg)
		assert.Empty(t, got)
	})

	type args struct {
		mode           Mode
		snippet        string
		indentation    string
		depth          int
		multiline      bool
		tabStopWrapper string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "very simple xml - tab stops, depth = 2",
			args: args{
				mode:           ModeXML,
				snippet:        `collection[foo=bar]>item*3+item[bar=]`,
				indentation:    "  ",
				depth:          2,
				multiline:      true,
				tabStopWrapper: "$",
			},
			want: `<collection foo="bar">
      <item>$STOP1$</item>
      <item>$STOP2$</item>
      <item>$STOP3$</item>
      <item bar="$STOP4$">$STOP5$</item>
    </collection>`,
			wantErr: RequireNoError,
		},
		{
			name: "very simple html - no tab stops, depth = 2, anchor used",
			args: args{
				mode:           ModeHTML,
				snippet:        `div.container>h1.h1+ul.list>li.item#item$$*3^a:blank.button+br`,
				indentation:    "  ",
				depth:          2,
				multiline:      true,
				tabStopWrapper: "",
			},
			want: `<div class="container">
      <h1 class="h1"></h1>
      <ul class="list">
        <li id="item01" class="item"></li>
        <li id="item02" class="item"></li>
        <li id="item03" class="item"></li>
      </ul>
      <a href="https://" target="_blank" rel="noopener noreferrer" class="button"></a>
      <br>
    </div>`,
			wantErr: RequireNoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := Xemmet(tt.args.mode, tt.args.snippet, tt.args.indentation, tt.args.depth, true, tt.args.tabStopWrapper)
			tt.wantErr(t, gotErr)

			assert.Equal(t, tt.want, got)
		})
	}
}
