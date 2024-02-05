package main

type Snippeter struct {
	mode Mode
}

func NewSnippeter(mode Mode) *Snippeter {
	return &Snippeter{
		mode: mode,
	}
}

var m = map[string]string{
	"bq":    "blockquote",
	"fig":   "figure",
	"figc":  "figcaption",
	"pic":   "picture",
	"ifr":   "iframe",
	"emb":   "embed",
	"obj":   "object",
	"cap":   "caption",
	"colg":  "colgroup",
	"fst":   "fieldset",
	"btn":   "button",
	"optg":  "optgroup",
	"tarea": "textarea",
	"leg":   "legend",
	"sect":  "section",
	"art":   "article",
	"hdr":   "header",
	"ftr":   "footer",
	"adr":   "address",
	"dlg":   "dialog",
	"str":   "strong",
	"prog":  "progress",
	"mn":    "main",
	"tem":   "template",
	"fset":  "fieldset",
	"datal": "datalist",
	"kg":    "keygen",
	"out":   "output",
	"det":   "details",
	"sum":   "summary",
	"cmd":   "command",
}

// based on https://github.com/emmetio/emmet/blob/master/src/snippets/html.json
func (s *Snippeter) Walk(tokens []Token) []Token {
	for _, token := range tokens {
		if tagToken, ok := token.(*TagToken); ok {
			s.ApplySnippets(tagToken)
		}

		for _, child := range token.GetChildren() {
			tagToken, ok := child.(*TagToken)
			if !ok {
				continue
			}

			s.ApplySnippets(tagToken)
		}
	}

	return tokens
}

func (s *Snippeter) ApplySnippets(token *TagToken) Token {
	switch s.mode {
	case ModeHTML, ModeHTMX:
		if mappedName, ok := m[token.Name]; ok {
			token.SetName(mappedName)
		}

		if s.mode == ModeHTMX {
			s.ApplyHTMXSnippets(token)
		}

		s.ApplyHTMLSnippets(token)
	}

	return token
}

func (s *Snippeter) ApplyHTMXSnippets(token *TagToken) {
	switch token.Name {
	case "button:get", "button:post", "button:put", "button:patch", "button:delete":
		token.
			SetName("button").
			FallbackAttribute(NewAttr("hx-trigger", "click")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("hx-"+token.Name[7:], "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("hx-target", "")).
			FallbackAttribute(NewAttr("hx-swap", "innerHTML"))

	case "input:q", "input:search":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("name", "q")).
			FallbackAttribute(NewAttr("type", "search")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("hx-get", "")).
			FallbackAttribute(NewAttr("hx-trigger", "keyup changed delay:500ms")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("hx-target", "")).
			FallbackAttribute(NewAttr("hx-swap", "innerHTML")).
			FallbackAttribute(NewAttr("placeholder", ""))

	case "script:htmx":
		token.
			SetName("script").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", "https://unpkg.com/htmx.org@1.9.10"))

	}
}

func (s *Snippeter) ApplyHTMLSnippets(token *TagToken) {
	switch token.Name {
	case "a":
		token.FallbackAttribute(NewAttr("href", "#"))

	case "a:blank":
		token.
			SetName("a").
			FallbackAttribute(NewAttr("href", "https://")).
			FallbackAttribute(NewAttr("target", "_blank")).
			FallbackAttribute(NewAttr("rel", "noopener noreferrer"))

	case "a:link":
		token.
			SetName("a").
			FallbackAttribute(NewAttr("href", "https://"))

	case "a:mail":
		token.
			SetName("a").
			FallbackAttribute(NewAttr("href", "mailto:"))

	case "a:tel":
		token.
			SetName("a").
			FallbackAttribute(NewAttr("href", "tel:+"))

	case "abbr", "acr", "acronym":
		token.
			SetName("abbr").
			FallbackAttribute(NewAttr("title", ""))

	case "bdo":
		token.
			SetName("bdo").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("dir", ""))

	case "bdo:r":
		token.
			SetName("bdo").
			FallbackAttribute(NewAttr("dir", "rtl"))

	case "bdo:l":
		token.
			SetName("bdo").
			FallbackAttribute(NewAttr("dir", "ltr"))

	case "link":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "stylesheet")).
			FallbackAttribute(NewAttr("href", ""))

	case "link:css":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "stylesheet")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "style.css"))

	case "link:print":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "stylesheet")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "style.css")).
			FallbackAttribute(NewAttr("media", "print"))

	case "link:favicon":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "shortcut icon")).
			FallbackAttribute(NewAttr("type", "image/x-icon")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "favicon.ico"))

	case "link:mf", "link:manifest":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "manifest")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "manifest.json"))

	case "link:touch":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "apple-touch-icon")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "favicon.png"))

	case "link:rss":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "alternate")).
			FallbackAttribute(NewAttr("type", "application/rss+xml")).
			FallbackAttribute(NewAttr("title", "RSS")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "rss.xml"))

	case "link:atom":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "alternate")).
			FallbackAttribute(NewAttr("type", "application/atom+xml")).
			FallbackAttribute(NewAttr("title", "Atom")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "atom.xml"))

	case "link:im", "link:import":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "import")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "component.html"))

	case "meta:utf":
		token.
			SetName("meta").
			FallbackAttribute(NewAttr("http-equiv", "Content-Type")).
			FallbackAttribute(NewAttr("content", "text/html;charset=UTF-8"))

	case "meta:vp":
		token.
			SetName("meta").
			FallbackAttribute(NewAttr("name", "viewport")).
			FallbackAttribute(NewAttr("content", "width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"))

	case "meta:compat":
		token.
			SetName("meta").
			FallbackAttribute(NewAttr("http-equiv", "X-UA-Compatible")).
			FallbackAttribute(NewAttr("content", "IE=7"))

	case "script:src":
		token.
			SetName("script").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", ""))

	case "img":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("alt", ""))

	case "img:s", "img:srcset", "ri:d", "ri:dpr":
		token.
			SetName("img").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("alt", ""))

	case "img:z", "img:sizes", "ri:v", "ri:viewport":
		token.
			SetName("img").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("sizes", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("alt", ""))

	case "src":
		token.
			SetName("source")

	case "src:sc", "source:src":
		token.
			SetName("source").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", ""))

	case "src:s", "source:srcset":
		token.
			SetName("source").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", ""))

	case "src:t", "source:type":
		token.
			SetName("source").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", "image/"))

	case "src:z", "source:sizes":
		token.
			SetName("source").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("sizes", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", ""))

	case "src:m", "source:media":
		token.
			SetName("source").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("media", "(min-width: )")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", ""))

	case "src:mt", "source:media:type":
		token.
			SetName("source").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("media", "(min-width: )")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", "image/"))

	case "src:mz", "source:media:sizes":
		token.
			SetName("source").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("media", "(min-width: )")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("sizes", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", ""))

	case "src:zt", "source:sizes:type":
		token.
			SetName("source").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("sizes", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("srcset", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", "image/"))

	case "iframe":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("frameborder", "0"))

	case "embed":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", ""))

	case "object":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("data", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", ""))

	case "map":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "area", "area:d", "area:c", "area:r", "area:p":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("coords", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("href", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("alt", ""))

		switch token.Name {
		case "area:d":
			token.FallbackAttribute(NewAttr("shape", "default"))

		case "area:c":
			token.FallbackAttribute(NewAttr("shape", "circle"))

		case "area:r":
			token.FallbackAttribute(NewAttr("shape", "rect"))

		case "area:p":
			token.FallbackAttribute(NewAttr("shape", "poly"))

		default:
			// TODO: Enable tab selection
			token.FallbackAttribute(NewAttr("shape", ""))
		}

		token.SetName("area")

	case "form":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("action", ""))

	case "form:get":
		token.
			SetName("form").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("action", "")).
			FallbackAttribute(NewAttr("method", "get"))

	case "form:post":
		token.
			SetName("form").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("action", "")).
			FallbackAttribute(NewAttr("method", "post"))

	case "label":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("for", ""))

	case "input":
		token.
			FallbackAttribute(NewAttr("type", "text")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:h", "input:hidden":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "hidden")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:t", "input:text":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "text")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:search":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "search")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:email":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "email")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:url":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "url")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:p", "input:password":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "password")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:datetime":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "datetime")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:date":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "date")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:datetime-local":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "datetime-local")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:month":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "month")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:week":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "week")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:time":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "time")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:tel":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "tel")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", "")).
			FallbackAttribute(NewAttr("pattern", "[0-9]{3}-[0-9]{2}-[0-9]{3}"))

	case "input:number":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "number")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("min", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("max", ""))

	case "input:color":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "password")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("value", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:c", "input:checkbox":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "checkbox")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("value", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:r", "input:radio":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "radio")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("value", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:range":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "range")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("min", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("max", ""))

	case "input:f", "input:file":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "file")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", ""))

	case "input:s", "input:submit":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "submit")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("value", ""))

	case "input:i", "input:image":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "image")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("alt", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", ""))

	case "input:b", "input:btn", "input:button":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "button")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("value", ""))

	case "input:reset":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "reset"))

	case "select":
		// TODO: Enable tab selection
		token.FallbackAttribute(NewAttr("name", ""))

	case "select:d", "select:disabled":
		token.
			SetName("select").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", "")).
			FallbackAttribute(NewAttr("disabled", "").HasNoEqualSign())

	case "opt", "option":
		token.
			SetName("select").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("value", ""))

	case "textarea":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("name", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("cols", "30")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("rows", "10"))

	case "marquee":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("behavior", "")).
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("direction", ""))

	case "menu:c":
		token.
			SetName("menu").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", "context"))

	case "menu:t":
		token.
			SetName("menu").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", "toolbar"))

	case "video":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", ""))

	case "audio":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("src", ""))

	case "btn:s", "button:s", "button:submit":
		token.
			SetName("button").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", "submit"))

	case "btn:r", "button:l", "button:reset":
		token.
			SetName("button").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", "reset"))

	case "btn:b", "button:b", "button:button":
		token.
			SetName("button").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("type", "button"))

	case "btn:d", "button:d", "button:disabled":
		token.
			SetName("button").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("disabled", "").HasNoEqualSign())

	case "fst:d", "fset:d", "fieldset:d", "fieldset:disabled":
		token.
			SetName("fieldset").
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("disabled", "").HasNoEqualSign())

	case "data":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("value", ""))

	case "meter":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("value", ""))

	case "time":
		token.
			// TODO: Enable tab selection
			FallbackAttribute(NewAttr("datetime", ""))
	}
}
