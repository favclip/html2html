package html2html

import (
	"io"

	"golang.org/x/net/html"
)

// VoidElements is start with html.StartTagToken, but it is not with html.EndTagToken.
// see https://www.w3.org/TR/html5/syntax.html#void-elements
var VoidElements = []string{
	"area",
	"base",
	"br",
	"col",
	"embed",
	"hr",
	"img",
	"input",
	"keygen",
	"link",
	"meta",
	"param",
	"source",
	"track",
	"wbr",
}

var _ TokenConsumer = &DefaultConsumer{}
var _ TagAttrsConsumer = &DefaultConsumer{}

func IsVoidElement(tag Tag) bool {
	for _, ve := range VoidElements {
		if tag.Name() == ve {
			return true
		}
	}

	return false
}

// TokenConsumer はHTMLのTokenを食べて何らかの文字列を組み立てる
type TokenConsumer interface {
	// ConsumeToken は渡されたTokenをキリのいいところまで処理し、次に処理するべきTokenを返す
	ConsumeToken(parent Tag, tokenizer *html.Tokenizer, token html.Token) (html.Token, error)
}

type TagAttrsConsumer interface {
	ConsumeAttrs(tag Tag, token html.Token) error
}

type Converter interface {
	DefaultConsumer() TokenConsumer
	RaiseErrorOnInvalidEndTag() bool

	ConsumerByTokenType(tokenType html.TokenType) TokenConsumer
	ConsumerByTagName(tagName string) TokenConsumer

	SetRaiseErrorOnInvalidEndTag(newVal bool)
	SetDefaultConsumer(consumer TokenConsumer)
	SetTokenTypeConsumer(tokenType html.TokenType, consumer TokenConsumer)
	SetTagNameConsumer(tagName string, consumer TokenConsumer)

	Parse(r io.Reader) (Tag, error)
	Convert(r io.Reader) (string, error)
}

func NewConverter() Converter {
	var conv Converter = &defaultConverter{}
	conv.SetDefaultConsumer(&DefaultConsumer{
		conv: conv,
	})
	conv.SetRaiseErrorOnInvalidEndTag(true)

	return conv
}
