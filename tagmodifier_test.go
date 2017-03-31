package html2html

import (
	"bytes"
	"strings"
	"testing"
)

func TestTagReplacer(t *testing.T) {
	html := "<b><strike>Foobar<strike>FizzBuzz"
	r := strings.NewReader(html)
	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(false)
	tag, err := conv.Parse(r)
	if err != nil {
		t.Fatal(err)
	}

	var modifier func(token Tag) (Token, error)
	modifier = func(token Tag) (Token, error) {
		if token.Type() != TypeTagToken {
			return nil, nil
		}

		tag := token.Tag()
		if tag.Name() != "strike" {
			return nil, nil
		}
		altTag := CreateElement("span")
		altTag.AddAttr("class", "strike")
		altTag.AddChildTokens(tag.Tokens()...)

		for _, token := range altTag.Tokens() {
			if token.Type() != TypeTagToken {
				continue
			}

			childTag := token.Tag()
			childAltTag, err := TagReplacer(childTag, modifier)
			if err != nil {
				return nil, err
			} else if childAltTag != nil {
				altTag.ReplateChildToken(childTag, childAltTag)
			}
		}

		return altTag, nil
	}
	_, err = TagReplacer(tag, modifier)
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBufferString("")
	tag.BuildHTML(buf)
	if v := buf.String(); v != `<b><span class="strike">Foobar<span class="strike">FizzBuzz</span></span></b>` {
		t.Error("unexpected", v)
	}
}
