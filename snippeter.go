package main

type Snippeter struct {
	mode Mode
}

func NewSnippeter(mode Mode) *Snippeter {
	return &Snippeter{
		mode: mode,
	}
}

// nolint:gochecknoglobals
var htmlTagAbbreviations = map[string]string{
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
func (s *Snippeter) Walk(tokens ...Token) []Token {
	for _, token := range tokens {
		if tagToken, ok := token.(*TagToken); ok {
			s.ApplySnippets(tagToken)
		}

		for _, child := range token.GetChildren() {
			s.Walk(child)
		}

		tagToken, ok := token.(*TagToken)
		if !ok {
			continue
		}

		s.ApplySnippets(tagToken)
	}

	return tokens
}

// nolint: unparam, ireturn
func (s *Snippeter) ApplySnippets(token *TagToken) Token {
	// nolint: exhaustive
	switch s.mode {
	case ModeHTML, ModeHTMX:
		if mappedName, ok := htmlTagAbbreviations[token.Name]; ok {
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
	case "a:get", "a:post", "a:put", "a:patch", "a:delete":
		method := token.Name[2:]
		token.
			SetName("a").
			FallbackAttribute(NewDefaultAttr("href", "https://")).
			FallbackAttribute(NewDefaultAttr("hx-"+method, "https://")).
			FallbackAttribute(NewAttr("hx-trigger", "click")).
			FallbackAttribute(NewDefaultAttr("hx-target", "")).
			FallbackAttribute(NewAttr("hx-swap", "innerHTML"))

	case "button:get", "button:post", "button:put", "button:patch", "button:delete":
		method := token.Name[7:]
		token.
			SetName("button").
			FallbackAttribute(NewDefaultAttr("hx-"+method, "https://")).
			FallbackAttribute(NewAttr("hx-trigger", "click")).
			FallbackAttribute(NewDefaultAttr("hx-target", "")).
			FallbackAttribute(NewAttr("hx-swap", "innerHTML"))

	case "input:q", "input:search":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("name", "q")).
			FallbackAttribute(NewAttr("type", "search")).
			FallbackAttribute(NewDefaultAttr("hx-get", "")).
			FallbackAttribute(NewDefaultAttr("hx-trigger", "keyup changed delay:500ms")).
			FallbackAttribute(NewDefaultAttr("hx-target", "")).
			FallbackAttribute(NewAttr("hx-swap", "innerHTML")).
			FallbackAttribute(NewDefaultAttr("placeholder", ""))

	case "script:htmx":
		token.
			SetName("script").
			FallbackAttribute(NewDefaultAttr("src", "https://unpkg.com/htmx.org@1.9.10"))
	}
}

// nolint: funlen, gocyclo, cyclop, maintidx
func (s *Snippeter) ApplyHTMLSnippets(token *TagToken) {
	switch token.Name {
	case "a":
		token.FallbackAttribute(NewDefaultAttr("href", "#"))

	case "a:blank":
		token.
			SetName("a").
			FallbackAttribute(NewDefaultAttr("href", "https://")).
			FallbackAttribute(NewAttr("target", "_blank")).
			FallbackAttribute(NewAttr("rel", "noopener noreferrer"))

	case "a:link":
		token.
			SetName("a").
			FallbackAttribute(NewDefaultAttr("href", "https://"))

	case "a:mail":
		token.
			SetName("a").
			FallbackAttribute(NewDefaultAttr("href", "mailto:"))

	case "a:tel":
		token.
			SetName("a").
			FallbackAttribute(NewDefaultAttr("href", "tel:+"))

	case "abbr", "acr", "acronym":
		token.
			SetName("abbr").
			FallbackAttribute(NewAttr("title", ""))

	case "bdo":
		token.
			SetName("bdo").
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
			FallbackAttribute(NewAttr("href", "style.css"))

	case "link:print":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "stylesheet")).
			FallbackAttribute(NewAttr("href", "style.css")).
			FallbackAttribute(NewAttr("media", "print"))

	case "link:favicon":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "shortcut icon")).
			FallbackAttribute(NewAttr("type", "image/x-icon")).
			FallbackAttribute(NewAttr("href", "favicon.ico"))

	case "link:mf", "link:manifest":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "manifest")).
			FallbackAttribute(NewAttr("href", "manifest.json"))

	case "link:touch":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "apple-touch-icon")).
			FallbackAttribute(NewAttr("href", "favicon.png"))

	case "link:rss":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "alternate")).
			FallbackAttribute(NewAttr("type", "application/rss+xml")).
			FallbackAttribute(NewAttr("title", "RSS")).
			FallbackAttribute(NewAttr("href", "rss.xml"))

	case "link:atom":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "alternate")).
			FallbackAttribute(NewAttr("type", "application/atom+xml")).
			FallbackAttribute(NewAttr("title", "Atom")).
			FallbackAttribute(NewAttr("href", "atom.xml"))

	case "link:im", "link:import":
		token.
			SetName("link").
			FallbackAttribute(NewAttr("rel", "import")).
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
			FallbackAttribute(NewAttr("src", ""))

	case "img":
		token.
			FallbackAttribute(NewAttr("src", "")).
			FallbackAttribute(NewAttr("alt", ""))

	case "img:s", "img:srcset", "ri:d", "ri:dpr":
		token.
			SetName("img").
			FallbackAttribute(NewAttr("srcset", "")).
			FallbackAttribute(NewAttr("src", "")).
			FallbackAttribute(NewAttr("alt", ""))

	case "img:z", "img:sizes", "ri:v", "ri:viewport":
		token.
			SetName("img").
			FallbackAttribute(NewAttr("sizes", "")).
			FallbackAttribute(NewAttr("srcset", "")).
			FallbackAttribute(NewAttr("src", "")).
			FallbackAttribute(NewAttr("alt", ""))

	case "src":
		token.
			SetName("source")

	case "src:sc", "source:src":
		token.
			SetName("source").
			FallbackAttribute(NewAttr("src", "")).
			FallbackAttribute(NewAttr("type", ""))

	case "src:s", "source:srcset":
		token.
			SetName("source").
			FallbackAttribute(NewAttr("srcset", ""))

	case "src:t", "source:type":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("type", "image/"))

	case "src:z", "source:sizes":
		token.
			SetName("source").
			FallbackAttribute(NewAttr("sizes", "")).
			FallbackAttribute(NewAttr("srcset", ""))

	case "src:m", "source:media":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("media", "(min-width: )")).
			FallbackAttribute(NewAttr("srcset", ""))

	case "src:mt", "source:media:type":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("media", "(min-width: )")).
			FallbackAttribute(NewAttr("srcset", "")).
			FallbackAttribute(NewDefaultAttr("type", "image/"))

	case "src:mz", "source:media:sizes":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("media", "(min-width: )")).
			FallbackAttribute(NewAttr("sizes", "")).
			FallbackAttribute(NewAttr("srcset", ""))

	case "src:zt", "source:sizes:type":
		token.
			SetName("source").
			FallbackAttribute(NewAttr("sizes", "")).
			FallbackAttribute(NewAttr("srcset", "")).
			FallbackAttribute(NewDefaultAttr("type", "image/"))

	case "iframe":
		token.
			FallbackAttribute(NewAttr("src", "")).
			FallbackAttribute(NewDefaultAttr("frameborder", "0"))

	case "embed":
		token.
			FallbackAttribute(NewAttr("src", "")).
			FallbackAttribute(NewAttr("type", ""))

	case "object":
		token.
			FallbackAttribute(NewAttr("data", "")).
			FallbackAttribute(NewAttr("type", ""))

	case "map":
		token.
			FallbackAttribute(NewAttr("name", ""))

	case "area", "area:d", "area:c", "area:r", "area:p":
		token.
			FallbackAttribute(NewAttr("coords", "")).
			FallbackAttribute(NewAttr("href", "")).
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
			token.FallbackAttribute(NewAttr("shape", ""))
		}

		token.SetName("area")

	case "form":
		token.
			FallbackAttribute(NewDefaultAttr("action", "post"))

	case "form:get":
		token.
			SetName("form").
			FallbackAttribute(NewAttr("action", "")).
			FallbackAttribute(NewAttr("method", "get"))

	case "form:post":
		token.
			SetName("form").
			FallbackAttribute(NewAttr("action", "")).
			FallbackAttribute(NewAttr("method", "post"))

	case "label":
		token.
			FallbackAttribute(NewAttr("for", ""))

	case "input":
		token.
			FallbackAttribute(NewAttr("type", "text")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:h", "input:hidden":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "hidden")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:t", "input:text":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "text")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:search":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "search")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:email":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "email")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:url":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "url")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:p", "input:password":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "password")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:datetime":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "datetime")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:date":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "date")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:datetime-local":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "datetime-local")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:month":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "month")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:week":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "week")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:time":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "time")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:tel":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "tel")).
			FallbackAttribute(NewAttr("name", "")).
			FallbackAttribute(NewAttr("pattern", "[0-9]{3}-[0-9]{2}-[0-9]{3}"))

	case "input:number":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "number")).
			FallbackAttribute(NewAttr("name", "")).
			FallbackAttribute(NewAttr("min", "")).
			FallbackAttribute(NewDefaultAttr("max", ""))

	case "input:color":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "password")).
			FallbackAttribute(NewAttr("Value", "")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:c", "input:checkbox":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "checkbox")).
			FallbackAttribute(NewAttr("Value", "")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:r", "input:radio":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "radio")).
			FallbackAttribute(NewAttr("Value", "")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:range":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "range")).
			FallbackAttribute(NewAttr("name", "")).
			FallbackAttribute(NewAttr("min", "")).
			FallbackAttribute(NewAttr("max", ""))

	case "input:f", "input:file":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "file")).
			FallbackAttribute(NewAttr("name", ""))

	case "input:s", "input:submit":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "submit")).
			FallbackAttribute(NewAttr("Value", ""))

	case "input:i", "input:image":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "image")).
			FallbackAttribute(NewAttr("alt", "")).
			FallbackAttribute(NewAttr("src", ""))

	case "input:b", "input:btn", "input:button":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "button")).
			FallbackAttribute(NewAttr("Value", ""))

	case "input:reset":
		token.
			SetName("input").
			FallbackAttribute(NewAttr("type", "reset"))

	case "select":
		token.FallbackAttribute(NewAttr("name", ""))

	case "select:d", "select:disabled":
		token.
			SetName("select").
			FallbackAttribute(NewAttr("name", "")).
			FallbackAttribute(NewAttr("disabled", "").HasNoEqualSign())

	case "opt", "option":
		token.
			SetName("select").
			FallbackAttribute(NewAttr("Value", ""))

	case "textarea":
		token.
			FallbackAttribute(NewAttr("name", "")).
			FallbackAttribute(NewAttr("cols", "")).
			FallbackAttribute(NewAttr("rows", ""))

	case "marquee":
		token.
			FallbackAttribute(NewAttr("behavior", "")).
			FallbackAttribute(NewAttr("direction", ""))

	case "menu:c":
		token.
			SetName("menu").
			FallbackAttribute(NewAttr("type", "context"))

	case "menu:t":
		token.
			SetName("menu").
			FallbackAttribute(NewAttr("type", "toolbar"))

	case "video":
		token.
			FallbackAttribute(NewAttr("src", ""))

	case "audio":
		token.
			FallbackAttribute(NewAttr("src", ""))

	case "btn:s", "button:s", "button:submit":
		token.
			SetName("button").
			FallbackAttribute(NewAttr("type", "submit"))

	case "btn:r", "button:l", "button:reset":
		token.
			SetName("button").
			FallbackAttribute(NewAttr("type", "reset"))

	case "btn:b", "button:b", "button:button":
		token.
			SetName("button").
			FallbackAttribute(NewAttr("type", "button"))

	case "btn:d", "button:d", "button:disabled":
		token.
			SetName("button").
			FallbackAttribute(NewAttr("disabled", "").HasNoEqualSign())

	case "fst:d", "fset:d", "fieldset:d", "fieldset:disabled":
		token.
			SetName("fieldset").
			FallbackAttribute(NewAttr("disabled", "").HasNoEqualSign())

	case "data":
		token.
			FallbackAttribute(NewAttr("Value", ""))

	case "meter":
		token.
			FallbackAttribute(NewAttr("Value", ""))

	case "time":
		token.
			FallbackAttribute(NewAttr("datetime", ""))
	}
}
