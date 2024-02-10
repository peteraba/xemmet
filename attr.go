package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

func Number(idx, repeat, start int, reverse bool, numbering string) string {
	if idx == 0 {
		panic("should not happen")
	}

	if !reverse && start == 1 && numbering == "" {
		return ""
	}

	count := idx + start - 1
	if reverse {
		count = repeat - idx + start
	}

	if numbering == "" {
		numbering = "$"
	}

	str := fmt.Sprintf("%0*d", len(numbering), count)

	return str
}

type AttrValue struct {
	Type      TokenType
	Value     string
	Numbering string
	Start     int
	Reverse   bool
}

func NewClass(value string) *AttrValue {
	return &AttrValue{
		Type:  Class,
		Value: value,
		Start: 1,
	}
}

func NewID(value string) *AttrValue {
	return &AttrValue{
		Type:  ID,
		Value: value,
		Start: 1,
	}
}

func (v *AttrValue) SetStart(start int) *AttrValue {
	v.Start = start

	return v
}

func (v *AttrValue) SetNumbering(numbering string) *AttrValue {
	v.Numbering = numbering

	return v
}

func (v *AttrValue) SetReverse(values ...bool) *AttrValue {
	if len(values) > 0 {
		v.Reverse = values[0]
	} else {
		v.Reverse = true
	}

	return v
}

func (v *AttrValue) getValue(idx, repeat int) string {
	if v == nil {
		return ""
	}

	return v.Value + Number(idx, repeat, v.Start, v.Reverse, v.Numbering)
}

func (v *AttrValue) Clone() *AttrValue {
	if v == nil {
		return nil
	}

	return &AttrValue{
		Type:      v.Type,
		Value:     v.Value,
		Numbering: v.Numbering,
		Start:     v.Start,
		Reverse:   v.Reverse,
	}
}

type AttrValues []*AttrValue

func (av AttrValues) Clone() AttrValues {
	newAv := make(AttrValues, 0, len(av))

	for _, a := range av {
		newAv = append(newAv, a.Clone())
	}

	return newAv
}

type Attr struct {
	Type         TokenType
	Name         string
	Value        string
	DefaultValue string
	HasEqualSign bool
}

func NewDefaultAttr(name, defaultValue string) *Attr {
	return &Attr{
		Name:         name,
		DefaultValue: defaultValue,
		HasEqualSign: true,
	}
}

func NewAttr(name, value string) *Attr {
	return &Attr{
		Name:         name,
		Value:        value,
		HasEqualSign: true,
	}
}

func (a *Attr) HasNoEqualSign() *Attr {
	if a.Value != "" {
		panic(fmt.Sprintf("attribute '%s' has no equal sign, but has a Value of '%s'", a.Name, a.Value))
	}

	a.HasEqualSign = false

	return a
}

func (a *Attr) Clone() *Attr {
	return &Attr{
		Name:         a.Name,
		Value:        a.Value,
		HasEqualSign: a.HasEqualSign,
	}
}

func (a *Attr) GetValue(counter *Counter, tabStopWrapper string) string {
	if a == nil || a.Value == "" {
		return a.DefaultValue + newTabStop(tabStopWrapper, counter.Get())
	}

	if len(a.Value) < 5 || a.Value[:5] != "lorem" {
		return a.Value
	}

	return lorem(a.Value)
}

type AttrList []*Attr

func (al AttrList) Clone() AttrList {
	newAl := make(AttrList, 0, len(al))

	for _, a := range al {
		newAl = append(newAl, a.Clone())
	}

	return newAl
}

func (al AttrList) Has(name string) bool {
	for _, a := range al {
		if a.Name == name {
			return true
		}
	}

	return false
}

type Text struct {
	value string
}

func (t *Text) IsEmpty() bool {
	return t == nil || t.value == ""
}

func (t *Text) Clone() *Text {
	if t == nil {
		return nil
	}

	return &Text{
		value: t.value,
	}
}

const loremKeyword = "lorem"

func (t *Text) GetValue() string {
	if t == nil {
		return ""
	}

	if !strings.HasPrefix(t.value, loremKeyword) {
		return t.value
	}

	return lorem(t.value)
}

func NewText(value string) *Text {
	return &Text{
		value: value,
	}
}

const defaultWordCount = 5

func lorem(expression string) string {
	var (
		words = defaultWordCount
		err   error
	)

	if len(expression) > len(loremKeyword) {
		words, err = strconv.Atoi(expression[len(loremKeyword):])
		if err != nil {
			words = defaultWordCount
		}
	}

	return gofakeit.LoremIpsumSentence(words)
}
