package md2ed

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func resolveNodes(n *html.Node) string {
	blocks := make([]string, 0)
	endblocks := make([]string, 0)
	combinator := "\n\n"
	/*if n.Parent != nil {
		fmt.Println(n.Parent.Data, "->", n.Data)
	}*/
	if n.Type == html.ElementNode {
		if n.Data == "html" {
			// Nothing other than head and body
			combinator = ""
		}
		if n.Data == "paragraph" {
			combinator = ""
		}
		if n.Data == "callout" {
			callout_type := ""
			for _, attr := range n.Attr {
				if attr.Key == "type" {
					callout_type = attr.Val
				}
			}
			blocks = append(blocks, fmt.Sprintf("{.callout type=\"%s\"}\n", callout_type))
			combinator = ""
		}
		if n.Data == "a" {
			link_value := ""
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link_value = attr.Val
				}
			}
			return fmt.Sprintf("[%s](%s)", n.FirstChild.Data, link_value)
		}
		if n.Data == "list" {
			combinator = "\n"
		}
		if n.Data == "list-item" {
			blocks = append(blocks, "* ")
			combinator = ""
		}
		if n.Data == "break" {
			return "\n"
		}

		if n.Data == "bold" {
			combinator = ""
			blocks = append(blocks, "**")
			endblocks = append(endblocks, "**")
		}
		if n.Data == "italic" {
			combinator = ""
			blocks = append(blocks, "*")
			endblocks = append(endblocks, "*")
		}
		if n.Data == "underline" {
			combinator = ""
			blocks = append(blocks, "<underline>")
			endblocks = append(endblocks, "</underline>")
		}
		if n.Data == "code" {
			combinator = ""
			blocks = append(blocks, "`")
			endblocks = append(endblocks, "`")
		}
		if n.Data == "heading" {
			level := 1
			for _, attr := range n.Attr {
				if attr.Key == "level" {
					level, _ = strconv.Atoi(attr.Val)
				}
			}
			combinator = ""
			for i := 0; i < level; i++ {
				blocks = append(blocks, "#")
			}
			blocks = append(blocks, " ")
		}
		if n.Data == "pre" {
			combinator = ""
			blocks = append(blocks, "```\n")
			endblocks = append(endblocks, "```")
		}
		if n.Data == "snippet" {
			extras := make([]string, 0)
			language := ""
			for _, attr := range n.Attr {
				if attr.Key == "language" {
					language = attr.Val
				} else {
					extras = append(extras, fmt.Sprintf(".%s=\"%s\"", attr.Key, attr.Val))
				}
			}
			combinator = "\n"
			blocks = append(blocks, fmt.Sprintf("```%s%s", language, strings.Join(extras, "")))
			endblocks = append(endblocks, "```")
		}
	}
	if n.Type == html.TextNode {
		return n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		blocks = append(blocks, resolveNodes(c))
	}
	blocks = append(blocks, endblocks...)
	return strings.Join(blocks, combinator)
}

func RenderEdToMD(content string) string {
	content = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(content,
		"</link", "</a"),
		"<link", "<a"),
		"<break/>", "<break></break>")
	node, _ := html.Parse(strings.NewReader(content))
	return resolveNodes(node)
}
