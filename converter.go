package html2html

import (
	"bytes"
	"io"

	"golang.org/x/net/context"
	"golang.org/x/net/html"
)

var _ Converter = &defaultConverter{}

type defaultConverter struct {
	defaultConsumer           TokenConsumer
	raiseErrorOnInvalidEndTag bool

	c                 context.Context
	tokenTypeConsumer map[html.TokenType]TokenConsumer
	tagConsumer       map[string]TokenConsumer
}

func (conv *defaultConverter) DefaultConsumer() TokenConsumer {
	return conv.defaultConsumer
}

func (conv *defaultConverter) RaiseErrorOnInvalidEndTag() bool {
	return conv.raiseErrorOnInvalidEndTag
}

func (conv *defaultConverter) ConsumerByTokenType(tokenType html.TokenType) TokenConsumer {
	return conv.tokenTypeConsumer[tokenType]
}

func (conv *defaultConverter) ConsumerByTagName(tagName string) TokenConsumer {
	return conv.tagConsumer[tagName]
}

func (conv *defaultConverter) SetRaiseErrorOnInvalidEndTag(newVal bool) {
	conv.raiseErrorOnInvalidEndTag = newVal
}

func (conv *defaultConverter) SetDefaultConsumer(consumer TokenConsumer) {
	conv.defaultConsumer = consumer
}

func (conv *defaultConverter) SetTokenTypeConsumer(tokenType html.TokenType, consumer TokenConsumer) {
	if conv.tokenTypeConsumer == nil {
		conv.tokenTypeConsumer = make(map[html.TokenType]TokenConsumer)
	}
	conv.tokenTypeConsumer[tokenType] = consumer
}

func (conv *defaultConverter) SetTagNameConsumer(tagName string, consumer TokenConsumer) {
	if conv.tagConsumer == nil {
		conv.tagConsumer = make(map[string]TokenConsumer)
	}
	conv.tagConsumer[tagName] = consumer
}

func (conv *defaultConverter) Parse(r io.Reader) (Tag, error) {
	tokenizer := html.NewTokenizer(r)

	tokenizer.Next()
	token := tokenizer.Token()

	var err error
	root := CreateDocumentRoot()
	for {
		// break the loop by io.EOF
		token, err = conv.DefaultConsumer().ConsumeToken(root, tokenizer, token)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	if err != io.EOF && err != nil {
		return nil, err
	}

	return root, nil
}

func (conv *defaultConverter) Convert(r io.Reader) (string, error) {
	tag, err := conv.Parse(r)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	tag.BuildHTML(buf)
	return buf.String(), nil
}
