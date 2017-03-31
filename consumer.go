package html2html

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

var _ TokenConsumer = &DefaultConsumer{}
var _ TokenConsumer = &vacuumConsumer{}

type DefaultConsumer struct {
	conv          Converter
	attrsConsumer TagAttrsConsumer
}

func (consumer *DefaultConsumer) ConsumeToken(parent Tag, tokenizer *html.Tokenizer, token html.Token) (html.Token, error) {
	nextConsumer := consumer.conv.ConsumerByTokenType(token.Type)
	if nextConsumer == nil {
		// skip
	} else if c, ok := nextConsumer.(*DefaultConsumer); ok {
		return c.ConsumeTokenImpl(parent, tokenizer, token)
	} else {
		return nextConsumer.ConsumeToken(parent, tokenizer, token)
	}

	switch token.Type {
	case html.StartTagToken, html.SelfClosingTagToken:
		nextConsumer = consumer.conv.ConsumerByTagName(token.Data)
		if nextConsumer == nil {
			// skip
		} else if c, ok := nextConsumer.(*DefaultConsumer); ok {
			return c.ConsumeTokenImpl(parent, tokenizer, token)
		} else {
			return nextConsumer.ConsumeToken(parent, tokenizer, token)
		}
	}

	return consumer.ConsumeTokenImpl(parent, tokenizer, token)
}

func (consumer *DefaultConsumer) ConsumeAttrs(tag Tag, token html.Token) error {
	if consumer.attrsConsumer != nil {
		err := consumer.attrsConsumer.ConsumeAttrs(tag, token)
		if err != nil {
			return err
		}
	} else {
		for _, attr := range token.Attr {
			tag.AddAttr(attr.Key, attr.Val)
		}
	}

	return nil
}

func (consumer *DefaultConsumer) ConsumeTokenImpl(parent Tag, tokenizer *html.Tokenizer, token html.Token) (html.Token, error) {
	var tagNameBytes []byte

	switch token.Type {
	case html.ErrorToken:
		if err := tokenizer.Err(); err != nil {
			return token, err
		}

		return token, fmt.Errorf("unknown state")

	case html.DoctypeToken:
		parent.AddChildTokens(CreateDoctypeToken(token.Data))

		tokenizer.Next()
		return tokenizer.Token(), nil

	case html.TextToken:
		parent.AddChildTokens(CreateTextToken(token.Data))

		tokenizer.Next()
		return tokenizer.Token(), nil

	case html.CommentToken:
		parent.AddChildTokens(CreateCommentToken(token.Data))

		tokenizer.Next()
		return tokenizer.Token(), nil

	case html.SelfClosingTagToken:
		child := CreateElementSelfClosing(token.Data)
		err := consumer.ConsumeAttrs(child, token)
		if err != nil {
			return token, err
		}
		parent.AddChildTokens(child)

		tokenizer.Next()
		return tokenizer.Token(), nil

	case html.EndTagToken:
		if consumer.conv.RaiseErrorOnInvalidEndTag() {
			return token, fmt.Errorf("unexpected end tag: %s", string(tagNameBytes))
		}

		// ignore
		tokenizer.Next()
		return tokenizer.Token(), nil

	case html.StartTagToken:
		child := CreateElement(token.Data)
		err := consumer.ConsumeAttrs(child, token)
		if err != nil {
			return token, err
		}
		parent.AddChildTokens(child)

		return consumer.ConsumeElementBody(child, tokenizer, token, token.Data)

	default:
		return token, fmt.Errorf("unknown tokenType: %s", token.Type.String())
	}
}

func (consumer *DefaultConsumer) ConsumeElementBody(tag Tag, tokenizer *html.Tokenizer, token html.Token, tagName string) (html.Token, error) {
	startTagName := token.Data
	if tagName == "" {
		tagName = startTagName
	}
	for _, voidElement := range VoidElements {
		if tagName == voidElement {
			// VoidElementにはbodyがないのでさくっと終わらせて次を読ませる
			tokenizer.Next()
			return tokenizer.Token(), nil
		}
	}
	var err error
	for {
		tokenizer.Next()
		token = tokenizer.Token()
		for {
			if token.Type == html.EndTagToken {
				break
			}

			token, err = consumer.conv.DefaultConsumer().ConsumeToken(tag, tokenizer, token)
			if err == io.EOF {
				if consumer.conv.RaiseErrorOnInvalidEndTag() {
					return token, fmt.Errorf("end tag `%s` is not comming", startTagName)
				}

				// missing tags are automatically interpolated by builder
				return token, err

			} else if err != nil {
				return token, err
			}
		}

		endTagName := token.Data

		if startTagName == endTagName {
			break
		} else if consumer.conv.RaiseErrorOnInvalidEndTag() {
			return token, fmt.Errorf("unexpected end tag: %s, expected: %s", endTagName, startTagName)
		}
	}
	tokenizer.Next()
	return tokenizer.Token(), nil
}

func (consumer *DefaultConsumer) SetTagAttrsConsumer(attrsConsumer TagAttrsConsumer) {
	consumer.attrsConsumer = attrsConsumer
}

// NewVacuumConsumer は渡されたTokenから始まる要素を全て捨てるTokenConsumerを返す
func NewVacuumConsumer(conv Converter) TokenConsumer {
	return &vacuumConsumer{conv}
}

type vacuumConsumer struct {
	conv Converter
}

func (consumer *vacuumConsumer) ConsumeToken(parent Tag, tokenizer *html.Tokenizer, token html.Token) (html.Token, error) {
	tag := CreateElement(token.Data)
	consumer.conv.SetTagNameConsumer(tag.Name(), nil)
	defer func() {
		consumer.conv.SetTagNameConsumer(tag.Name(), consumer)
	}()
	return consumer.conv.DefaultConsumer().ConsumeToken(tag, tokenizer, token)
}
