package md2ed

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var img_blocks = 0

func resolveNodes(n *html.Node, content_folder string) string {
	preblocks := make([]string, 0)
	blocks := make([]string, 0)
	endblocks := make([]string, 0)
	combinator := "\n\n"
	/*if n.Parent != nil {
		fmt.Println(n.Parent.Data, "->", n.Data)
	}*/
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		blocks = append(blocks, resolveNodes(c, content_folder))
	}
	if n.Type == html.ElementNode {
		if n.Data == "html" {
			// Nothing other than head and body
			combinator = ""
		} else if n.Data == "paragraph" {
			combinator = ""
		} else if n.Data == "callout" {
			callout_type := ""
			for _, attr := range n.Attr {
				if attr.Key == "type" {
					callout_type = attr.Val
				}
			}
			preblocks = append(preblocks, fmt.Sprintf("{.callout type=\"%s\"}\n", callout_type))
			combinator = ""
		} else if n.Data == "a" {
			link_value := ""
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link_value = attr.Val
				}
			}
			return fmt.Sprintf("[%s](%s)", n.FirstChild.Data, link_value)
		} else if n.Data == "list" {
			combinator = "\n"
		} else if n.Data == "list-item" {
			preblocks = append(preblocks, "* ")
			combinator = ""
		} else if n.Data == "break" {
			return "\n"
		} else if n.Data == "bold" {
			combinator = ""
			preblocks = append(preblocks, "**")
			endblocks = append(endblocks, "**")
		} else if n.Data == "italic" {
			combinator = ""
			preblocks = append(preblocks, "*")
			endblocks = append(endblocks, "*")
		} else if n.Data == "underline" {
			combinator = ""
			preblocks = append(preblocks, "<underline>")
			endblocks = append(endblocks, "</underline>")
		} else if n.Data == "code" {
			combinator = ""
			preblocks = append(preblocks, "`")
			endblocks = append(endblocks, "`")
		} else if n.Data == "img" {
			src := ""
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					src = attr.Val
				}
			}
			// Just guess its a png. TODO: Figure this out
			image_file := fmt.Sprintf("image%d.png", img_blocks)
			img_blocks += 1
			out, err := os.Create(path.Join(content_folder, image_file))
			if err != nil {
				return ""
			}
			defer out.Close()
			resp, err := http.Get(src)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()

			io.Copy(out, resp.Body)
			preblocks = append(preblocks, "![](")
			blocks = append(blocks, image_file)
			endblocks = append(endblocks, ")")
			combinator = ""
		} else if n.Data == "heading" {
			level := 1
			for _, attr := range n.Attr {
				if attr.Key == "level" {
					level, _ = strconv.Atoi(attr.Val)
				}
			}
			combinator = ""
			for i := 0; i < level; i++ {
				preblocks = append(preblocks, "#")
			}
			preblocks = append(preblocks, " ")
		} else if n.Data == "pre" {
			combinator = ""
			preblocks = append(preblocks, "```\n")
			endblocks = append(endblocks, "```")
		} else if n.Data == "snippet" {
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
			preblocks = append(preblocks, fmt.Sprintf("```%s%s", language, strings.Join(extras, "")))
			endblocks = append(endblocks, "```")
		} else if n.Data == "head" || n.Data == "body" || n.Data == "document" || n.Data == "snippet-file" || n.Data == "figure" {
			// Nothing
		} else if n.Data == "table" {
			// This should contain a `thead` child which tells us how many columns.
			combinator = "\n"
			preblocks = append(preblocks, blocks[0])
			blocks = blocks[1:]
			bar_count := strings.Count(preblocks[len(preblocks)-1], "|")
			bars := make([]string, 0)
			for i := 1; i < bar_count; i++ {
				bars = append(bars, "---")
			}
			preblocks = append(preblocks, "|"+strings.Join(bars, "|")+"|")
		} else if n.Data == "thead" || n.Data == "tbody" {
			combinator = "\n"
		} else if n.Data == "tr" {
			preblocks = append(preblocks, "")
			endblocks = append(endblocks, "")
			combinator = "|"
		} else if n.Data == "td" || n.Data == "th" {
			combinator = ""
		} else {
			fmt.Println("Unhandled node element", n.Data, " has parent ", n.Parent.Data)
		}
	}
	if n.Type == html.TextNode {
		return n.Data
	}
	blocks = append(preblocks, blocks...)
	blocks = append(blocks, endblocks...)
	return strings.Join(blocks, combinator)
}

func RenderEdToMD(content string, content_folder string) string {
	content = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(content,
		"</link", "</a"),
		"<link", "<a"),
		"<break/>", "<break></break>")
	node, _ := html.Parse(strings.NewReader(content))
	return resolveNodes(node, content_folder)
}
