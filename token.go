package main

type TokenType string

const (
	ID    TokenType = "ID"
	Class TokenType = "Class"
)

type Token interface {
	GetRepeat() int
	GetParent() Token
	SetParent(t Token) Token
	AddChildren(children ...Token) Token
	GetChildren() []Token
}

type TagToken struct {
	Name       string
	Repeat     int
	Classes    AttrValues
	ID         *AttrValue
	Attributes AttrList
	Text       *Text
	Parent     Token
	Children   []Token
}

func NewTagToken(name string, repeat int) *TagToken {
	return &TagToken{
		Name:   name,
		Repeat: repeat,
	}
}

func (t *TagToken) SetName(name string) *TagToken {
	t.Name = name

	return t
}

func (t *TagToken) SetID(id *AttrValue) *TagToken {
	t.ID = id

	return t
}

func (t *TagToken) AddClass(class *AttrValue) *TagToken {
	t.Classes = append(t.Classes, class)

	return t
}

func (t *TagToken) AddAttribute(attribute *Attr) *TagToken {
	t.Attributes = append(t.Attributes, attribute)

	return t
}

func (t *TagToken) FallbackAttribute(attribute *Attr) *TagToken {
	for _, a := range t.Attributes {
		if a.Name == attribute.Name {
			return t
		}
	}

	t.AddAttribute(attribute)

	return t
}

func (t *TagToken) SetAttribute(attribute *Attr) *TagToken {
	for i, a := range t.Attributes {
		if a.Name == attribute.Name {
			t.Attributes[i] = attribute

			return t
		}
	}

	t.AddAttribute(attribute)

	return t
}

func (t *TagToken) SetText(text *Text) *TagToken {
	t.Text = text

	return t
}

// nolint: ireturn
func (t *TagToken) GetParent() Token {
	return t.Parent
}

// nolint: ireturn
func (t *TagToken) SetParent(parent Token) Token {
	t.Parent = parent

	return t
}

// nolint: ireturn
func (t *TagToken) AddChildren(children ...Token) Token {
	for _, child := range children {
		child.SetParent(t)
	}

	t.Children = append(t.Children, children...)

	return t
}

func (t *TagToken) GetChildren() []Token {
	return t.Children
}

func (t *TagToken) GetRepeat() int {
	return t.Repeat
}

type GroupToken struct {
	Children []Token
	Parent   Token
	Repeat   int
}

func NewGroupToken(repeat int, tokens ...Token) *GroupToken {
	token := &GroupToken{
		Children: nil,
		Repeat:   repeat,
	}

	token.AddChildren(tokens...)

	return token
}

func (g *GroupToken) GetRepeat() int {
	return g.Repeat
}

// nolint: ireturn
func (g *GroupToken) GetParent() Token {
	return g.Parent
}

// nolint: ireturn
func (g *GroupToken) SetParent(parent Token) Token {
	g.Parent = parent

	return g
}

// nolint: ireturn
func (g *GroupToken) AddChildren(children ...Token) Token {
	for _, child := range children {
		child.SetParent(g)
	}

	g.Children = append(g.Children, children...)

	return g
}

func (g *GroupToken) GetChildren() []Token {
	return g.Children
}
