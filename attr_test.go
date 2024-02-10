package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber(t *testing.T) {
	t.Parallel()

	type args struct {
		repeat    int
		start     int
		reverse   bool
		numbering string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				repeat:    3,
				start:     1,
				reverse:   false,
				numbering: "",
			},
			want: []string{"", "", ""},
		},
		{
			name: "default",
			args: args{
				repeat:    3,
				start:     1,
				reverse:   false,
				numbering: "$",
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "reverse",
			args: args{
				repeat:    3,
				start:     1,
				reverse:   true,
				numbering: "$",
			},
			want: []string{"3", "2", "1"},
		},
		{
			name: "reverse, start at 3",
			args: args{
				repeat:    3,
				start:     3,
				reverse:   true,
				numbering: "$",
			},
			want: []string{"5", "4", "3"},
		},
		{
			name: "reverse, start at 11",
			args: args{
				repeat:    3,
				start:     11,
				reverse:   true,
				numbering: "$",
			},
			want: []string{"13", "12", "11"},
		},
		{
			name: "start at 8, 2 digits minimum",
			args: args{
				repeat:    3,
				start:     8,
				reverse:   false,
				numbering: "$$",
			},
			want: []string{"08", "09", "10"},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := make([]string, 0, tt.args.repeat)

			for i := 1; i <= tt.args.repeat; i++ {
				got = append(got, Number(i, tt.args.repeat, tt.args.start, tt.args.reverse, tt.args.numbering))
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAttr_GetValue(t *testing.T) {
	t.Parallel()

	t.Run("lorem ipsum", func(t *testing.T) {
		t.Parallel()

		got := NewAttr("foo", "lorem").GetValue("")

		assert.NotEmpty(t, got)
		assert.Equal(t, 4, strings.Count(got, " "))
	})

	t.Run("lorem ipsum 25", func(t *testing.T) {
		t.Parallel()

		got := NewAttr("foo", "lorem25").GetValue("")

		assert.NotEmpty(t, got)
		assert.Equal(t, 24, strings.Count(got, " "))
	})

	type args struct {
		tabStopWrapper string
	}
	tests := []struct {
		name string
		sut  *Attr
		args args
		want string
	}{
		{
			name: "short",
			sut:  NewAttr("foo", "bar"),
			args: args{
				tabStopWrapper: "",
			},
			want: "bar",
		},
		{
			name: "not short",
			sut:  NewAttr("foo", "this is long enough"),
			args: args{
				tabStopWrapper: "",
			},
			want: "this is long enough",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.sut.GetValue(tt.args.tabStopWrapper)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestText_GetValue(t *testing.T) {
	t.Parallel()

	t.Run("lorem ipsum", func(t *testing.T) {
		t.Parallel()

		got := NewText("lorem").GetValue()

		assert.NotEmpty(t, got)
		assert.Equal(t, 4, strings.Count(got, " "))
	})

	t.Run("lorem ipsum 25", func(t *testing.T) {
		t.Parallel()

		got := NewText("lorem25").GetValue()

		assert.NotEmpty(t, got)
		assert.Equal(t, 24, strings.Count(got, " "))
	})

	tests := []struct {
		name string
		sut  *Text
		want string
	}{
		{
			name: "short",
			sut:  NewText("bar"),
			want: "bar",
		},
		{
			name: "not short",
			sut:  NewText("this is long enough"),
			want: "this is long enough",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.sut.GetValue()

			assert.Equal(t, tt.want, got)
		})
	}
}
