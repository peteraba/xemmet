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

func (e Elem) HTML(mode Mode, indentation string, depth int, multiline bool, tabStopWrapper string) string {
	var builder strings.Builder

	xmlShortTag := e.isShortTagXML(mode) && tabStopWrapper == ""
	htmlShortTag := e.isShortTagHTML(mode) && tabStopWrapper == ""
	shortTag := xmlShortTag || htmlShortTag
	emptyTag := e.isEmptyTag()

	currentIndentation := ""
	if indentation != "" {
		currentIndentation = strings.Repeat(indentation, depth)
	}

	if e.Name == "" {
		return e.TextOnly(currentIndentation, "", multiline)
	}

	builder.WriteString(e.OpeningTag(currentIndentation, xmlShortTag, multiline, tabStopWrapper))

	if multiline && (!emptyTag || shortTag) {
		builder.WriteString("\n")
	}

	if !shortTag {
		builder.WriteString(e.TextOnly(currentIndentation, indentation, multiline))

		builder.WriteString(e.TabStop(tabStopWrapper))

		builder.WriteString(e.RenderChildren(mode, indentation, depth, multiline, tabStopWrapper))

		builder.WriteString(e.ClosingTag(currentIndentation, multiline, emptyTag))
	}

	return builder.String()
}

func (e Elem) GetText() string {
	if e.Text == nil {
		return ""
	}

	return e.Text.GetValue()
}

func (e Elem) TextOnly(currentIndentation, indentationExtra string, multiline bool) string {
	if e.Text.IsEmpty() || !multiline {
		return e.GetText()
	}

	if multiline {
		return currentIndentation + indentationExtra + e.GetText() + "\n"
	}

	return currentIndentation + indentationExtra + e.GetText()
}

func (e Elem) OpeningTag(currentIndentation string, xmlShortTag, multiline bool, tabStopWrapper string) string {
	var builder strings.Builder

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

	return builder.String()
}

func (e Elem) TabStop(tabStopWrapper string) string {
	if len(e.Children) != 0 {
		return ""
	}

	return newTabStop(tabStopWrapper)
}

func newTabStop(tabStopWrapper string) string {
	if tabStopWrapper == "" {
		return ""
	}

	return fmt.Sprintf("%sSTOP%d%s", tabStopWrapper, counter.GetGlobalTabStopCounter(), tabStopWrapper)
}

func (e Elem) RenderChildren(mode Mode, indentation string, depth int, multiline bool, tabStopWrapper string) string {
	if len(e.Children) == 0 {
		return ""
	}

	var builder strings.Builder

	nextDepth := depth + 1

	for _, child := range e.Children {
		builder.WriteString(child.HTML(mode, indentation, nextDepth, multiline, tabStopWrapper))
	}

	return builder.String()
}

func (e Elem) ClosingTag(currentIndentation string, multiline, emptyTag bool) string {
	var builder strings.Builder

	if multiline && !emptyTag {
		builder.WriteString(currentIndentation)
	}

	builder.WriteString("</")
	builder.WriteString(e.Name)
	builder.WriteString(">")

	if multiline {
		builder.WriteString("\n")
	}

	return builder.String()
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

func (el ElemList) HTML(mode Mode, indentation string, depth int, multiline bool, tabStopWrapper string) string {
	var builder strings.Builder

	for _, e := range el {
		builder.WriteString(e.HTML(mode, indentation, depth, multiline, tabStopWrapper))
	}

	return builder.String()
}
