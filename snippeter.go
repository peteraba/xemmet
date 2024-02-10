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
			FallbackAttribute(NewDefaultAttr("hx-trigger", "click")).
			FallbackAttribute(NewDefaultAttr("hx-target", "")).
			FallbackAttribute(NewDefaultAttr("hx-swap", "innerHTML"))

	case "button:get", "button:post", "button:put", "button:patch", "button:delete":
		method := token.Name[7:]
		token.
			SetName("button").
			FallbackAttribute(NewDefaultAttr("hx-"+method, "https://")).
			FallbackAttribute(NewDefaultAttr("hx-trigger", "click")).
			FallbackAttribute(NewDefaultAttr("hx-target", "")).
			FallbackAttribute(NewDefaultAttr("hx-swap", "innerHTML"))

	case "input:q", "input:search":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("name", "q")).
			FallbackAttribute(NewDefaultAttr("type", "search")).
			FallbackAttribute(NewDefaultAttr("hx-get", "")).
			FallbackAttribute(NewDefaultAttr("hx-trigger", "keyup changed delay:500ms")).
			FallbackAttribute(NewDefaultAttr("hx-target", "")).
			FallbackAttribute(NewDefaultAttr("hx-swap", "innerHTML")).
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
			FallbackAttribute(NewDefaultAttr("target", "_blank")).
			FallbackAttribute(NewDefaultAttr("rel", "noopener noreferrer"))

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
			FallbackAttribute(NewDefaultAttr("title", ""))

	case "bdo":
		token.
			SetName("bdo").
			FallbackAttribute(NewDefaultAttr("dir", ""))

	case "bdo:r":
		token.
			SetName("bdo").
			FallbackAttribute(NewDefaultAttr("dir", "rtl"))

	case "bdo:l":
		token.
			SetName("bdo").
			FallbackAttribute(NewDefaultAttr("dir", "ltr"))

	case "link":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "stylesheet")).
			FallbackAttribute(NewDefaultAttr("href", ""))

	case "link:css":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "stylesheet")).
			FallbackAttribute(NewDefaultAttr("href", "style.css"))

	case "link:print":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "stylesheet")).
			FallbackAttribute(NewDefaultAttr("href", "style.css")).
			FallbackAttribute(NewDefaultAttr("media", "print"))

	case "link:favicon":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "shortcut icon")).
			FallbackAttribute(NewDefaultAttr("type", "image/x-icon")).
			FallbackAttribute(NewDefaultAttr("href", "favicon.ico"))

	case "link:mf", "link:manifest":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "manifest")).
			FallbackAttribute(NewDefaultAttr("href", "manifest.json"))

	case "link:touch":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "apple-touch-icon")).
			FallbackAttribute(NewDefaultAttr("href", "favicon.png"))

	case "link:rss":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "alternate")).
			FallbackAttribute(NewDefaultAttr("type", "application/rss+xml")).
			FallbackAttribute(NewDefaultAttr("title", "RSS")).
			FallbackAttribute(NewDefaultAttr("href", "rss.xml"))

	case "link:atom":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "alternate")).
			FallbackAttribute(NewDefaultAttr("type", "application/atom+xml")).
			FallbackAttribute(NewDefaultAttr("title", "Atom")).
			FallbackAttribute(NewDefaultAttr("href", "atom.xml"))

	case "link:im", "link:import":
		token.
			SetName("link").
			FallbackAttribute(NewDefaultAttr("rel", "import")).
			FallbackAttribute(NewDefaultAttr("href", "component.html"))

	case "meta:utf":
		token.
			SetName("meta").
			FallbackAttribute(NewDefaultAttr("http-equiv", "Content-Type")).
			FallbackAttribute(NewDefaultAttr("content", "text/html;charset=UTF-8"))

	case "meta:vp":
		token.
			SetName("meta").
			FallbackAttribute(NewDefaultAttr("name", "viewport")).
			FallbackAttribute(NewDefaultAttr("content", "width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"))

	case "meta:compat":
		token.
			SetName("meta").
			FallbackAttribute(NewDefaultAttr("http-equiv", "X-UA-Compatible")).
			FallbackAttribute(NewDefaultAttr("content", "IE=7"))

	case "script:src":
		token.
			SetName("script").
			FallbackAttribute(NewDefaultAttr("src", ""))

	case "img":
		token.
			FallbackAttribute(NewDefaultAttr("src", "")).
			FallbackAttribute(NewDefaultAttr("alt", ""))

	case "img:s", "img:srcset", "ri:d", "ri:dpr":
		token.
			SetName("img").
			FallbackAttribute(NewDefaultAttr("srcset", "")).
			FallbackAttribute(NewDefaultAttr("src", "")).
			FallbackAttribute(NewDefaultAttr("alt", ""))

	case "img:z", "img:sizes", "ri:v", "ri:viewport":
		token.
			SetName("img").
			FallbackAttribute(NewDefaultAttr("sizes", "")).
			FallbackAttribute(NewDefaultAttr("srcset", "")).
			FallbackAttribute(NewDefaultAttr("src", "")).
			FallbackAttribute(NewDefaultAttr("alt", ""))

	case "src":
		token.
			SetName("source")

	case "src:sc", "source:src":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("src", "")).
			FallbackAttribute(NewDefaultAttr("type", ""))

	case "src:s", "source:srcset":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("srcset", ""))

	case "src:t", "source:type":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("type", "image/"))

	case "src:z", "source:sizes":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("sizes", "")).
			FallbackAttribute(NewDefaultAttr("srcset", ""))

	case "src:m", "source:media":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("media", "(min-width: )")).
			FallbackAttribute(NewDefaultAttr("srcset", ""))

	case "src:mt", "source:media:type":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("media", "(min-width: )")).
			FallbackAttribute(NewDefaultAttr("srcset", "")).
			FallbackAttribute(NewDefaultAttr("type", "image/"))

	case "src:mz", "source:media:sizes":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("media", "(min-width: )")).
			FallbackAttribute(NewDefaultAttr("sizes", "")).
			FallbackAttribute(NewDefaultAttr("srcset", ""))

	case "src:zt", "source:sizes:type":
		token.
			SetName("source").
			FallbackAttribute(NewDefaultAttr("sizes", "")).
			FallbackAttribute(NewDefaultAttr("srcset", "")).
			FallbackAttribute(NewDefaultAttr("type", "image/"))

	case "iframe":
		token.
			FallbackAttribute(NewDefaultAttr("src", "")).
			FallbackAttribute(NewDefaultAttr("frameborder", "0"))

	case "embed":
		token.
			FallbackAttribute(NewDefaultAttr("src", "")).
			FallbackAttribute(NewDefaultAttr("type", ""))

	case "object":
		token.
			FallbackAttribute(NewDefaultAttr("data", "")).
			FallbackAttribute(NewDefaultAttr("type", ""))

	case "map":
		token.
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "area", "area:d", "area:c", "area:r", "area:p":
		token.
			FallbackAttribute(NewDefaultAttr("coords", "")).
			FallbackAttribute(NewDefaultAttr("href", "")).
			FallbackAttribute(NewDefaultAttr("alt", ""))

		switch token.Name {
		case "area:d":
			token.FallbackAttribute(NewDefaultAttr("shape", "default"))

		case "area:c":
			token.FallbackAttribute(NewDefaultAttr("shape", "circle"))

		case "area:r":
			token.FallbackAttribute(NewDefaultAttr("shape", "rect"))

		case "area:p":
			token.FallbackAttribute(NewDefaultAttr("shape", "poly"))

		default:
			token.FallbackAttribute(NewDefaultAttr("shape", ""))
		}

		token.SetName("area")

	case "form":
		token.
			FallbackAttribute(NewDefaultAttr("action", ""))

	case "form:get":
		token.
			SetName("form").
			FallbackAttribute(NewDefaultAttr("action", "")).
			FallbackAttribute(NewDefaultAttr("method", "get"))

	case "form:post":
		token.
			SetName("form").
			FallbackAttribute(NewDefaultAttr("action", "")).
			FallbackAttribute(NewDefaultAttr("method", "post"))

	case "label":
		token.
			FallbackAttribute(NewDefaultAttr("for", ""))

	case "input":
		token.
			FallbackAttribute(NewDefaultAttr("type", "text")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:h", "input:hidden":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "hidden")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:t", "input:text":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "text")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:search":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "search")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:email":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "email")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:url":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "url")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:p", "input:password":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "password")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:datetime":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "datetime")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:date":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "date")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:datetime-local":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "datetime-local")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:month":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "month")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:week":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "week")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:time":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "time")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:tel":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "tel")).
			FallbackAttribute(NewDefaultAttr("name", "")).
			FallbackAttribute(NewDefaultAttr("pattern", "[0-9]{3}-[0-9]{2}-[0-9]{3}"))

	case "input:number":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "number")).
			FallbackAttribute(NewDefaultAttr("name", "")).
			FallbackAttribute(NewDefaultAttr("min", "")).
			FallbackAttribute(NewDefaultAttr("max", ""))

	case "input:color":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "password")).
			FallbackAttribute(NewDefaultAttr("Value", "")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:c", "input:checkbox":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "checkbox")).
			FallbackAttribute(NewDefaultAttr("Value", "")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:r", "input:radio":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "radio")).
			FallbackAttribute(NewDefaultAttr("Value", "")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:range":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "range")).
			FallbackAttribute(NewDefaultAttr("name", "")).
			FallbackAttribute(NewDefaultAttr("min", "")).
			FallbackAttribute(NewDefaultAttr("max", ""))

	case "input:f", "input:file":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "file")).
			FallbackAttribute(NewDefaultAttr("name", ""))

	case "input:s", "input:submit":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "submit")).
			FallbackAttribute(NewDefaultAttr("Value", ""))

	case "input:i", "input:image":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "image")).
			FallbackAttribute(NewDefaultAttr("alt", "")).
			FallbackAttribute(NewDefaultAttr("src", ""))

	case "input:b", "input:btn", "input:button":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "button")).
			FallbackAttribute(NewDefaultAttr("Value", ""))

	case "input:reset":
		token.
			SetName("input").
			FallbackAttribute(NewDefaultAttr("type", "reset"))

	case "select":
		token.FallbackAttribute(NewDefaultAttr("name", ""))

	case "select:d", "select:disabled":
		token.
			SetName("select").
			FallbackAttribute(NewDefaultAttr("name", "")).
			FallbackAttribute(NewDefaultAttr("disabled", "").HasNoEqualSign())

	case "opt", "option":
		token.
			SetName("select").
			FallbackAttribute(NewDefaultAttr("Value", ""))

	case "textarea":
		token.
			FallbackAttribute(NewDefaultAttr("name", "")).
			FallbackAttribute(NewDefaultAttr("cols", "30")).
			FallbackAttribute(NewDefaultAttr("rows", "10"))

	case "marquee":
		token.
			FallbackAttribute(NewDefaultAttr("behavior", "")).
			FallbackAttribute(NewDefaultAttr("direction", ""))

	case "menu:c":
		token.
			SetName("menu").
			FallbackAttribute(NewDefaultAttr("type", "context"))

	case "menu:t":
		token.
			SetName("menu").
			FallbackAttribute(NewDefaultAttr("type", "toolbar"))

	case "video":
		token.
			FallbackAttribute(NewDefaultAttr("src", ""))

	case "audio":
		token.
			FallbackAttribute(NewDefaultAttr("src", ""))

	case "btn:s", "button:s", "button:submit":
		token.
			SetName("button").
			FallbackAttribute(NewDefaultAttr("type", "submit"))

	case "btn:r", "button:l", "button:reset":
		token.
			SetName("button").
			FallbackAttribute(NewDefaultAttr("type", "reset"))

	case "btn:b", "button:b", "button:button":
		token.
			SetName("button").
			FallbackAttribute(NewDefaultAttr("type", "button"))

	case "btn:d", "button:d", "button:disabled":
		token.
			SetName("button").
			FallbackAttribute(NewDefaultAttr("disabled", "").HasNoEqualSign())

	case "fst:d", "fset:d", "fieldset:d", "fieldset:disabled":
		token.
			SetName("fieldset").
			FallbackAttribute(NewDefaultAttr("disabled", "").HasNoEqualSign())

	case "data":
		token.
			FallbackAttribute(NewDefaultAttr("Value", ""))

	case "meter":
		token.
			FallbackAttribute(NewDefaultAttr("Value", ""))

	case "time":
		token.
			FallbackAttribute(NewDefaultAttr("datetime", ""))
	}
}
