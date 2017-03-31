package html2html

import (
	"bytes"
	"errors"
	"strings"
)

type TokenType int

const (
	TypeTagToken TokenType = 1 + iota
	TypeDoctypeToken
	TypeTextToken
	TypeCommentToken
)

var ErrInvalidTokenType = errors.New("invalid token type")

type Token interface {
	setParent(parent Tag)

	Type() TokenType
	Parent() Tag
	BuildHTML(buf *bytes.Buffer)

	Tag() Tag
	TextToken() TextToken
}

type Tag interface {
	Token

	IsDocumentRoot() bool

	Name() string
	Tokens() []Token
	SetTokens(tokens []Token)

	AddChildTokens(childlen ...Token)
	UnshiftChileToken(token Token)
	ReplateChildToken(from Token, to Token)
	GetElementsByTagName(tagName string) []Tag
	FindAncestor(tagName string) Tag

	Attrs() []*Attr
	SetAttrs(attrs []*Attr)
	AddAttr(key, value string)
	GetAttr(attrKey string) *Attr
	RemoveAttr(attrKey string)
	HasAttr(attrKey string) bool
	HasAttrValue(attrKey string, attrValue string) bool
	HasAttrValueCaseInsensitive(attrKey string, attrValue string) bool
	AddText(text string)
	AddComment(text string)
}

type TextToken interface {
	Token

	Text() string
}

func CreateDocumentRoot() Tag {
	return &tagImpl{documentRoot: true}
}

func CreateElement(tagName string) Tag {
	return &tagImpl{name: tagName}
}

func CreateElementSelfClosing(tagName string) Tag {
	return &tagImpl{name: tagName, selfClosing: true}
}

func CreateDoctypeToken(text string) Token {
	return &textTokenImpl{tokenType: TypeDoctypeToken, text: text}
}

func CreateTextToken(text string) Token {
	return &textTokenImpl{tokenType: TypeTextToken, text: text}
}

func CreateCommentToken(text string) Token {
	return &textTokenImpl{tokenType: TypeCommentToken, text: text}
}

var _ Tag = &tagImpl{}

type tagImpl struct {
	documentRoot bool
	parent       Tag
	name         string
	selfClosing  bool
	tokens       []Token
	attrs        []*Attr
}

func (tag *tagImpl) Type() TokenType {
	return TypeTagToken
}

func (tag *tagImpl) Parent() Tag {
	return tag.parent
}

func (tag *tagImpl) BuildHTML(buf *bytes.Buffer) {
	// root node doesn't have name & attrs.
	if tag.name != "" {
		buf.WriteString("<")
		buf.WriteString(tag.name)

		for _, attr := range tag.attrs {
			buf.WriteString(" ")
			buf.WriteString(attr.Key)
			if attr.Value != "" {
				buf.WriteString("=\"")
				buf.WriteString(attr.Value)
				buf.WriteString("\"")
			}
		}
		if tag.selfClosing {
			buf.WriteString("/>")
		} else {
			buf.WriteString(">")
		}
	}

	if tag.selfClosing {
		return
	}

	if len(tag.tokens) == 0 {
		for _, voidElement := range VoidElements {
			if voidElement == tag.name {
				return
			}
		}
	}

	for _, token := range tag.tokens {
		token.BuildHTML(buf)
	}

	if tag.name != "" {
		buf.WriteString("</")
		buf.WriteString(tag.name)
		buf.WriteString(">")
	}
}

func (tag *tagImpl) setParent(parent Tag) {
	current := tag.Parent()
	for {
		if current == nil {
			break
		}
		if current == tag {
			panic("recursive dom dependencies")
		}
		current = current.Parent()
	}

	tag.parent = parent
}

func (tag *tagImpl) Tag() Tag {
	return tag
}

func (tag *tagImpl) TextToken() TextToken {
	panic(ErrInvalidTokenType)
}

func (tag *tagImpl) IsDocumentRoot() bool {
	return tag.documentRoot
}

func (tag *tagImpl) Name() string {
	return tag.name
}

func (tag *tagImpl) Tokens() []Token {
	return tag.tokens
}

func (tag *tagImpl) SetTokens(tokens []Token) {
	for _, token := range tag.tokens {
		token.setParent(nil)
	}

	tag.tokens = tokens

	for _, token := range tag.tokens {
		token.setParent(tag)
	}
}

func (tag *tagImpl) AddChildTokens(childlen ...Token) {
	for _, token := range childlen {
		token.setParent(tag)
	}

	tag.tokens = append(tag.tokens, childlen...)
}

func (tag *tagImpl) UnshiftChileToken(token Token) {
	token.setParent(tag)

	newTokens := make([]Token, 0, len(tag.tokens)+1)
	newTokens = append(newTokens, token)
	newTokens = append(newTokens, tag.tokens...)

	tag.tokens = newTokens
}

func (tag *tagImpl) ReplateChildToken(from Token, to Token) {
	from.setParent(nil)
	to.setParent(tag)

	for idx, token := range tag.tokens {
		if token == from {
			tag.tokens[idx] = to
			break
		}
	}
}

func (tag *tagImpl) Attrs() []*Attr {
	return tag.attrs
}

func (tag *tagImpl) SetAttrs(attrs []*Attr) {
	tag.attrs = attrs
}

func (tag *tagImpl) AddAttr(key, value string) {
	tag.attrs = append(tag.attrs, &Attr{Key: key, Value: value})
}

func (tag *tagImpl) GetAttr(attrKey string) *Attr {
	for _, attr := range tag.attrs {
		if attr.Key == attrKey {
			return attr
		}
	}

	return nil
}

func (tag *tagImpl) RemoveAttr(attrKey string) {
	newAttrs := make([]*Attr, 0, len(tag.attrs))
	for _, attr := range tag.attrs {
		if attr.Key == attrKey {
			continue
		}
		newAttrs = append(newAttrs, attr)
	}
	tag.attrs = newAttrs
}

func (tag *tagImpl) HasAttr(attrKey string) bool {
	return tag.GetAttr(attrKey) != nil
}

func (tag *tagImpl) HasAttrValue(attrKey string, attrValue string) bool {
	attr := tag.GetAttr(attrKey)
	if attr == nil {
		return false
	}
	return attr.Value == attrValue
}

func (tag *tagImpl) HasAttrValueCaseInsensitive(attrKey string, attrValue string) bool {
	attr := tag.GetAttr(attrKey)
	if attr == nil {
		return false
	}
	return strings.ToLower(attr.Value) == strings.ToLower(attrValue)
}

func (tag *tagImpl) AddText(text string) {
	textToken := CreateTextToken(text)
	tag.AddChildTokens(textToken)
}

func (tag *tagImpl) AddComment(text string) {
	comment := CreateCommentToken(text)
	tag.AddChildTokens(comment)
}

func (tag *tagImpl) GetElementsByTagName(tagName string) []Tag {
	var result []Tag

	tagName = strings.ToLower(tagName)

	for _, token := range tag.tokens {
		if token.Type() != TypeTagToken {
			continue
		}

		child := token.Tag()
		if strings.ToLower(child.Name()) == tagName {
			result = append(result, child)
		}
		result = append(result, child.GetElementsByTagName(tagName)...)
	}

	return result
}

func (tag *tagImpl) FindAncestor(tagName string) Tag {
	var target Tag = tag
	for {
		parent := target.Parent()
		if parent == nil {
			break
		}

		if parent.Name() == tagName {
			return parent
		}

		target = parent
	}

	return nil
}

var _ TextToken = &textTokenImpl{}

type textTokenImpl struct {
	parent    Tag
	tokenType TokenType
	text      string
}

func (textToken *textTokenImpl) setParent(parent Tag) {
	textToken.parent = parent
}

func (textToken *textTokenImpl) Type() TokenType {
	return textToken.tokenType
}

func (textToken *textTokenImpl) Parent() Tag {
	return textToken.parent
}

func (textToken *textTokenImpl) BuildHTML(buf *bytes.Buffer) {
	switch textToken.tokenType {
	case TypeDoctypeToken:
		buf.WriteString("<!DOCTYPE ")
		buf.WriteString(textToken.text)
		buf.WriteString(">")
	case TypeTextToken:
		buf.WriteString(textToken.text)
	case TypeCommentToken:
		buf.WriteString("<!--")
		buf.WriteString(textToken.text)
		buf.WriteString("-->")
	}
}

func (textToken *textTokenImpl) Tag() Tag {
	panic(ErrInvalidTokenType)
}

func (textToken *textTokenImpl) TextToken() TextToken {
	return textToken
}

func (textToken *textTokenImpl) Text() string {
	return textToken.text
}

type Attr struct {
	Key   string
	Value string
}
