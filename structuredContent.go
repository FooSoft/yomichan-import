package yomichan

type contentAttr struct {
	lang               string
	fontStyle          string   // normal, italic
	fontWeight         string   // normal, bold
	fontSize           string   // small, medium, large, smaller, 80%, 125%, etc.
	textDecorationLine []string // underline, overline, line-through
	verticalAlign      string   // baseline, sub, super, text-top, text-bottom, middle, top, bottom
	textAlign          string   // start, end, left, right, center, justify, justify-all, match-parent
	marginTop          int
	marginLeft         int
	marginRight        int
	marginBottom       int
	listStyleType      string
	data               map[string]string
}

// if the array contains adjacent strings, concatenate them.
// ex: ["one", "two", content_structure, "four"] -> ["onetwo", content_structure, "four"]
// if the array only contains strings, return a concatenated string.
// ex: ["one", "two"] -> "onetwo"
func contentReduce(contents []any) any {
	if len(contents) == 1 {
		return contents[0]
	}
	newContents := []any{}
	var accumulator string
	for _, content := range contents {
		switch v := content.(type) {
		case string:
			accumulator = accumulator + v
		default:
			if accumulator != "" {
				newContents = append(newContents, accumulator)
				accumulator = ""
			}
			newContents = append(newContents, content)
		}
	}
	if accumulator != "" {
		newContents = append(newContents, accumulator)
	}
	if len(newContents) == 1 {
		return newContents[0]
	} else {
		return newContents
	}
}

func contentStructure(contents ...any) map[string]any {
	return map[string]any{
		"type":    "structured-content",
		"content": contentReduce(contents),
	}
}

func contentRuby(attr contentAttr, ruby string, contents ...any) map[string]any {
	rubyContent := map[string]any{
		"tag": "ruby",
		"content": []any{
			contentReduce(contents),
			map[string]string{"tag": "rp", "content": "("},
			map[string]string{"tag": "rt", "content": ruby},
			map[string]string{"tag": "rp", "content": ")"},
		},
	}
	if attr.lang != "" {
		rubyContent["lang"] = attr.lang
	}
	if len(attr.data) != 0 {
		rubyContent["data"] = attr.data
	}
	return rubyContent
}

func contentInternalLink(attr contentAttr, query string, contents ...any) map[string]any {
	linkContent := map[string]any{
		"tag":  "a",
		"href": "?query=" + query + "&wildcards=off",
	}
	if len(contents) == 0 {
		linkContent["content"] = query
	} else {
		linkContent["content"] = contentReduce(contents)
	}
	if attr.lang != "" {
		linkContent["lang"] = attr.lang
	}
	if len(attr.data) != 0 {
		linkContent["data"] = attr.data
	}
	return linkContent
}

func contentSpan(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "span", contents...)
}

func contentDiv(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "div", contents...)
}

func contentListItem(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "li", contents...)
}

func contentOrderedList(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "ol", contents...)
}

func contentUnorderedList(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "ul", contents...)
}

func contentTable(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "table", contents...)
}

func contentTableHead(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "thead", contents...)
}

func contentTableBody(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "tbody", contents...)
}

func contentTableRow(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "tr", contents...)
}

func contentTableHeadCell(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "th", contents...)
}

func contentTableCell(attr contentAttr, contents ...any) map[string]any {
	return contentStyledContainer(attr, "td", contents...)
}

func contentStyledContainer(attr contentAttr, tag string, contents ...any) map[string]any {
	container := map[string]any{"tag": tag}
	container["content"] = contentReduce(contents)
	if attr.lang != "" {
		container["lang"] = attr.lang
	}
	if len(attr.data) != 0 {
		container["data"] = attr.data
	}
	style := contentStyle(attr)
	if len(style) != 0 {
		container["style"] = style
	}
	return container
}

func contentStyle(attr contentAttr) map[string]any {
	style := make(map[string]any)
	if attr.fontStyle != "" {
		style["fontStyle"] = attr.fontStyle
	}
	if attr.fontWeight != "" {
		style["fontWeight"] = attr.fontWeight
	}
	if attr.fontSize != "" {
		style["fontSize"] = attr.fontSize
	}
	if len(attr.textDecorationLine) != 0 {
		style["textDecorationLine"] = attr.textDecorationLine
	}
	if attr.verticalAlign != "" {
		style["verticalAlign"] = attr.verticalAlign
	}
	if attr.textAlign != "" {
		style["textAlign"] = attr.textAlign
	}
	if attr.marginTop != 0 {
		style["marginTop"] = attr.marginTop
	}
	if attr.marginLeft != 0 {
		style["marginLeft"] = attr.marginLeft
	}
	if attr.marginRight != 0 {
		style["marginRight"] = attr.marginRight
	}
	if attr.marginBottom != 0 {
		style["marginBottom"] = attr.marginBottom
	}
	if attr.listStyleType != "" {
		style["listStyleType"] = attr.listStyleType
	}
	return style
}
