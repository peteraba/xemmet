package main

import (
	"fmt"
	"strings"

	"github.com/peteraba/xemmet/counter"
)

// nolint: gochecknoglobals
var shortHTMLTagNames = map[string]struct{}{
	"br":      {},
	"hr":      {},
	"img":     {},
	"input":   {},
	"link":    {},
	"meta":    {},
	"area":    {},
	"base":    {},
	"col":     {},
	"command": {},
	"embed":   {},
	"keygen":  {},
	"param":   {},
	"source":  {},
	"video":   {},
	"audio":   {},
	"track":   {},
	"wbr":     {},
}

type Elem struct {
	Name         string
	Classes      AttrValues
	ID           *AttrValue
	Attributes   AttrList
	Text         *Text
	Num          int
	SiblingCount int
	Children     ElemList
}

func (e Elem) isEmptyTag() bool {
	return len(e.Children) == 0 && e.Text.IsEmpty()
}

func (e Elem) isShortTagHTML(mode Mode) bool {
	if mode != ModeHTML {
		return false
	}

	if len(e.Children) > 0 || !e.Text.IsEmpty() {
		return false
	}

	_, ok := shortHTMLTagNames[e.Name]

	return ok
}

func (e Elem) isShortTagXML(mode Mode) bool {
	if mode != ModeXML {
		return false
	}

	if len(e.Children) > 0 || !e.Text.IsEmpty() {
		return false
	}

	return true
}

func (e Elem) HTML(builder *strings.Builder, mode Mode, indentation string, depth int, multiline bool, tabStopWrapper string) {
	xmlShortTag := e.isShortTagXML(mode) && tabStopWrapper == ""
	htmlShortTag := e.isShortTagHTML(mode) && tabStopWrapper == ""
	shortTag := xmlShortTag || htmlShortTag
	emptyTag := e.isEmptyTag()

	currentIndentation := ""
	if indentation != "" {
		currentIndentation = strings.Repeat(indentation, depth)
	}

	if e.Name == "" {
		e.TextOnly(builder, currentIndentation, "", multiline)

		return
	}

	e.OpeningTag(builder, currentIndentation, xmlShortTag, multiline, tabStopWrapper)

	if multiline && (!emptyTag || shortTag) {
		builder.WriteString("\n")
	}

	if !shortTag {
		e.TextOnly(builder, currentIndentation, indentation, multiline)

		e.TabStop(builder, tabStopWrapper)

		e.RenderChildren(builder, mode, indentation, depth, multiline, tabStopWrapper)

		e.ClosingTag(builder, currentIndentation, multiline, emptyTag)
	}
}

func (e Elem) GetText() string {
	if e.Text == nil {
		return ""
	}

	return e.Text.GetValue()
}

func (e Elem) TextOnly(builder *strings.Builder, currentIndentation, indentationExtra string, multiline bool) {
	if e.Text.IsEmpty() || !multiline {
		builder.WriteString(e.GetText())

		return
	}

	builder.WriteString(currentIndentation)
	builder.WriteString(indentationExtra)
	builder.WriteString(e.GetText())

	if multiline {
		builder.WriteString("\n")
	}
}

func (e Elem) OpeningTag(builder *strings.Builder, currentIndentation string, xmlShortTag, multiline bool, tabStopWrapper string) {
	if multiline {
		builder.WriteString(currentIndentation)
	}

	builder.WriteString("<")
	builder.WriteString(e.Name)

	id := e.GetID()

	if id != "" {
		builder.WriteString(" id=\"")
		builder.WriteString(id)
		builder.WriteString("\"")
	}

	if len(e.Attributes) > 0 {
		builder.WriteString(" ")
		builder.WriteString(e.GetAttrs(tabStopWrapper))
	}

	if len(e.Classes) > 0 {
		builder.WriteString(" class=\"")
		builder.WriteString(e.GetClass())
		builder.WriteString("\"")
	}

	if xmlShortTag {
		builder.WriteString(" /")
	}

	builder.WriteString(">")
}

func (e Elem) TabStop(builder *strings.Builder, tabStopWrapper string) {
	if len(e.Children) != 0 {
		return
	}

	builder.WriteString(newTabStop(tabStopWrapper))
}

func newTabStop(tabStopWrapper string) string {
	if tabStopWrapper == "" {
		return ""
	}

	return fmt.Sprintf("%sSTOP%d%s", tabStopWrapper, counter.GetGlobalTabStopCounter(), tabStopWrapper)
}

func (e Elem) RenderChildren(builder *strings.Builder, mode Mode, indentation string, depth int, multiline bool, tabStopWrapper string) {
	if len(e.Children) == 0 {
		return
	}

	nextDepth := depth + 1

	for _, child := range e.Children {
		child.HTML(builder, mode, indentation, nextDepth, multiline, tabStopWrapper)
	}
}

func (e Elem) ClosingTag(builder *strings.Builder, currentIndentation string, multiline, emptyTag bool) {
	if multiline && !emptyTag {
		builder.WriteString(currentIndentation)
	}

	builder.WriteString("</")
	builder.WriteString(e.Name)
	builder.WriteString(">")

	if multiline {
		builder.WriteString("\n")
	}
}

func (e Elem) Clone(num, siblingCount int) Elem {
	return Elem{
		Name:         e.Name,
		ID:           e.ID.Clone(),
		Classes:      e.Classes.Clone(),
		Attributes:   e.Attributes.Clone(),
		Text:         e.Text.Clone(),
		Num:          num,
		SiblingCount: siblingCount,
		Children:     e.Children.Clone(num, siblingCount),
	}
}

func (e Elem) GetNum() int {
	return e.Num
}

func (e Elem) GetID() string {
	if e.ID == nil {
		return ""
	}

	return e.ID.getValue(e.Num, e.SiblingCount)
}

func (e Elem) GetClass() string {
	if len(e.Classes) == 0 {
		return ""
	}

	classes := make([]string, 0, len(e.Classes))
	for _, class := range e.Classes {
		classes = append(classes, class.getValue(e.Num, e.SiblingCount))
	}

	return strings.Join(classes, " ")
}

func (e Elem) GetAttrs(tabStopWrapper string) string {
	attrs := []string{}
	for _, attr := range e.Attributes {
		// TODO: escape attribute values
		attrs = append(attrs, fmt.Sprintf(`%s="%s"`, attr.Name, attr.GetValue(tabStopWrapper)))
	}

	return strings.Join(attrs, " ")
}

type ElemList []*Elem

func (el ElemList) Clone(num, siblingCount int) ElemList {
	newEl := make(ElemList, 0, len(el))

	for _, e := range el {
		clone := e.Clone(num, siblingCount)
		newEl = append(newEl, &clone)
	}

	return newEl
}

func (el ElemList) HTML(builder *strings.Builder, mode Mode, indentation string, depth int, multiline bool, tabStopWrapper string) {
	for _, e := range el {
		e.HTML(builder, mode, indentation, depth, multiline, tabStopWrapper)
	}
}
