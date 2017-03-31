package html2html

import (
	"strings"
	"testing"
)

func TestConverterConvert(t *testing.T) {
	html := `Hi!<br/>
Welcome to <a href="https://github.com/favclip/">favclip!</a><br/>
<img src="https://avatars1.githubusercontent.com/u/15679?v=3&s=460" alt="a2c"/>
`

	ampHTML, err := NewConverter().Convert(strings.NewReader(html))
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `Hi!<br/>
Welcome to <a href="https://github.com/favclip/">favclip!</a><br/>
<img src="https://avatars1.githubusercontent.com/u/15679?v=3&s=460" alt="a2c"/>
`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestConverterConvert_fixBrokenEndTag1(t *testing.T) {
	html := `<i>Hi!`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(false)

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<i>Hi!</i>`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestConverterConvert_fixBrokenEndTag2(t *testing.T) {
	html := `<i>Hi!</b>`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(false)

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<i>Hi!</i>`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestConverterConvert_fixBrokenEndTag3(t *testing.T) {
	html := `<i><b><strike>Hi!</b></i></strike>`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(false)

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<i><b><strike>Hi!</strike></b></i>`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestConverterConvert_fixBrokenEndTag4(t *testing.T) {
	html := `<i>Hi!<b> world</i>`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(false)

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<i>Hi!<b> world</b></i>`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestConverterConvert_fixBrokenEndTag5(t *testing.T) {
	html := `<i>Hi!</b></i>`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(false)

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<i>Hi!</i>`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestConverterConvert_fixUnexpectedAttr(t *testing.T) {
	html := `<a href=http://example.com>Hi!</a>`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(false)

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<a href="http://example.com">Hi!</a>`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestConverterConvert_voidElements1(t *testing.T) {
	html := `<img src="https://avatars1.githubusercontent.com/u/15679?v=3&s=460">`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(true)

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<img src="https://avatars1.githubusercontent.com/u/15679?v=3&s=460">`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestConverterConvert_voidElements2(t *testing.T) {
	html := `<a><img src="https://avatars1.githubusercontent.com/u/15679?v=3&s=460"></a>`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(true)

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<a><img src="https://avatars1.githubusercontent.com/u/15679?v=3&s=460"></a>`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}

func TestVacuumConsumer(t *testing.T) {
	html := `<script src="foo.js"></script><h1>Hi!</h1>`

	conv := NewConverter()
	conv.SetRaiseErrorOnInvalidEndTag(false)
	conv.SetTagNameConsumer("script", NewVacuumConsumer(conv))

	r := strings.NewReader(html)
	ampHTML, err := conv.Convert(r)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	expected := `<h1>Hi!</h1>`
	if ampHTML != expected {
		t.Log("expected:\n", expected, "actual:\n", ampHTML)
		t.Fail()
	}
}
