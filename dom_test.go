package html2html

import (
	"bytes"
	"testing"
)

func TestNewTagBuilder(t *testing.T) {
	root := CreateDocumentRoot()
	root.AddChildTokens(CreateDoctypeToken("html"))
	html := CreateElement("html")
	root.AddChildTokens(html)
	html.AddAttr("amp", "")
	html.AddAttr("lang", "en")
	{
		head := CreateElement("head")
		html.AddChildTokens(head)
		{
			meta := CreateElement("meta")
			head.AddChildTokens(meta)
			meta.AddAttr("charset", "utf-8")
		}
		{
			script := CreateElement("script")
			head.AddChildTokens(script)
			script.AddAttr("async", "")
			script.AddAttr("src", "https://cdn.ampproject.org/v0.js")
		}
		{
			title := CreateElement("title")
			head.AddChildTokens(title)
			title.AddText("Hello, AMPs")
		}
		{
			link := CreateElement("link")
			head.AddChildTokens(link)
			link.AddAttr("rel", "canonical")
			link.AddAttr("href", "http://example.ampproject.org/article-metadata.html")
		}
		{
			meta := CreateElement("meta")
			head.AddChildTokens(meta)
			meta.AddAttr("name", "viewport")
			meta.AddAttr("content", "width=device-width,minimum-scale=1,initial-scale=1")
		}
		{
			script := CreateElement("script")
			head.AddChildTokens(script)
			script.AddAttr("type", "application/ld+json")
			script.AddText(`{"@context": "http://schema.org","@type": "NewsArticle","headline": "Open-source framework for publishing content","datePublished": "2015-10-07T12:02:41Z","image": ["logo.jpg"]}`)
		}
		{
			style := CreateElement("style")
			head.AddChildTokens(style)
			style.AddAttr("amp-boilerplate", "")
			style.AddText("body{-webkit-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-moz-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-ms-animation:-amp-start 8s steps(1,end) 0s 1 normal both;animation:-amp-start 8s steps(1,end) 0s 1 normal both}@-webkit-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-moz-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-ms-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-o-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}")
		}
		{
			noscript := CreateElement("noscript")
			head.AddChildTokens(noscript)
			{
				style := CreateElement("style")
				noscript.AddChildTokens(style)
				style.AddAttr("amp-boilerplate", "")
				style.AddText("body{-webkit-animation:none;-moz-animation:none;-ms-animation:none;animation:none}")
			}
		}
	}
	{
		body := CreateElement("body")
		html.AddChildTokens(body)
		{
			h1 := CreateElement("h1")
			body.AddChildTokens(h1)
			h1.AddText("Welcome to the mobile web")
		}
	}

	buf := bytes.NewBufferString("")
	root.BuildHTML(buf)

	htmlStr := buf.String()

	if htmlStr != `<!DOCTYPE html><html amp lang="en"><head><meta charset="utf-8"><script async src="https://cdn.ampproject.org/v0.js"></script><title>Hello, AMPs</title><link rel="canonical" href="http://example.ampproject.org/article-metadata.html"><meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1"><script type="application/ld+json">{"@context": "http://schema.org","@type": "NewsArticle","headline": "Open-source framework for publishing content","datePublished": "2015-10-07T12:02:41Z","image": ["logo.jpg"]}</script><style amp-boilerplate>body{-webkit-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-moz-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-ms-animation:-amp-start 8s steps(1,end) 0s 1 normal both;animation:-amp-start 8s steps(1,end) 0s 1 normal both}@-webkit-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-moz-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-ms-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-o-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}</style><noscript><style amp-boilerplate>body{-webkit-animation:none;-moz-animation:none;-ms-animation:none;animation:none}</style></noscript></head><body><h1>Welcome to the mobile web</h1></body></html>` {
		t.Error("unexpected", htmlStr)
	}
}
