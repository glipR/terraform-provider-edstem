package md2ed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"terraform-provider-edstem/internal/client"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func renderEmphasis(w io.Writer, p *ast.Emph, entering bool) {
	if entering {
		io.WriteString(w, "<italic>")
	} else {
		io.WriteString(w, "</italic>")
	}
}

func renderStrong(w io.Writer, p *ast.Strong, entering bool) {
	if entering {
		io.WriteString(w, "<bold>")
	} else {
		io.WriteString(w, "</bold>")
	}
}

func renderParagraph(w io.Writer, p *ast.Paragraph, entering bool, attribute *ast.Attribute) {
	mode := "paragraph"
	callout_type := "info"
	if attribute != nil {
		for _, class := range attribute.Classes {
			if string(class) == "callout" {
				mode = "callout"
				callout_type = string(attribute.Attrs["type"])
			}
		}
	}
	if entering {
		if mode == "paragraph" {
			io.WriteString(w, "<paragraph>")
		} else {
			io.WriteString(w, fmt.Sprintf("<callout type=\"%s\">", callout_type))
		}
	} else {
		if mode == "paragraph" {
			io.WriteString(w, "</paragraph>")
		} else {
			io.WriteString(w, "</callout>")
		}
	}
}

func renderHeading(w io.Writer, p *ast.Heading, entering bool) {
	if entering {
		io.WriteString(w, fmt.Sprintf("<heading level=\"%d\">", p.Level))
	} else {
		io.WriteString(w, "</heading>")
	}
}

func renderLink(w io.Writer, p *ast.Link, entering bool) {
	if entering {
		io.WriteString(w, fmt.Sprintf("<link href=\"%s\">", strings.ReplaceAll(string(p.Destination), "&", "&amp;")))
	} else {
		io.WriteString(w, "</link>")
	}
}

func renderList(w io.Writer, p *ast.List, entering bool) {
	if entering {
		style := "bullet"
		if p.ListFlags == ast.ListTypeOrdered {
			style = "number"
		}
		io.WriteString(w, fmt.Sprintf("<list style=\"%s\">", style))
	} else {
		io.WriteString(w, "</list>")
	}
}

func renderListItem(w io.Writer, p *ast.ListItem, entering bool) {
	if entering {
		io.WriteString(w, "<list-item>")
	} else {
		io.WriteString(w, "</list-item>")
	}
}

func renderCodeBlock(w io.Writer, p *ast.CodeBlock, entering bool) {

	if len(p.Info) == 0 {
		io.WriteString(w, "<pre>")
		io.WriteString(w, strings.ReplaceAll(string(bytes.Trim(p.Literal, "\n")), "\n", "\\n"))
		io.WriteString(w, "</pre>")
	} else {
		fields := strings.FieldsFunc(string(p.Info), func(s rune) bool { return s == '.' })
		lang := fields[0]
		extra := ""

		if len(fields) > 1 {
			extra = strings.Join(fields[1:], " ")
		}

		io.WriteString(w, fmt.Sprintf("<snippet language=\"%s\" %s><snippet-file id=\"code\">", lang, extra))
		io.WriteString(w, strings.ReplaceAll(string(bytes.Trim(p.Literal, "\n")), "\n", "\\n"))
		io.WriteString(w, "</snippet-file></snippet>")
	}

}

func renderMath(w io.Writer, p *ast.Math, entering bool) {
	io.WriteString(w, "$")
	io.WriteString(w, string(p.Literal))
	io.WriteString(w, "$")
}

func renderMathBlock(w io.Writer, p *ast.MathBlock, entering bool) {
	io.WriteString(w, "<paragraph>$$</paragraph>")
	if entering {
		io.WriteString(w, "<paragraph>"+string(p.Literal)+"</paragraph>")
	}
}

func renderImgBlock(w io.Writer, p *ast.Image, id string, alt string, width *string, height *string, entering bool) {
	io.WriteString(w, fmt.Sprintf("<figure><image src=\"https://static.au.edusercontent.com/files/%s\"", id))
	if alt != "" {
		io.WriteString(w, fmt.Sprintf(" alt=\"%s\"", alt))
	}
	if width != nil {
		io.WriteString(w, fmt.Sprintf(" width=\"%s\"", *width))
	}
	if height != nil {
		io.WriteString(w, fmt.Sprintf(" height=\"%s\"", *height))
	}
	io.WriteString(w, "/></figure>")
}

func uploadImg(w io.Writer, img *ast.Image, para *ast.Paragraph) error {
	path := string(img.Destination)
	alt_text := string(img.Children[0].AsLeaf().Literal)

	dat, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Request an image link
	// The course id doesn't matter.
	var course_id = ""
	var token = os.Getenv("EDSTEM_TOKEN")
	var c, _ = client.NewClient(&course_id, &token)

	boundary := "----WebKitFormBoundaryplBATvmbbo4b7Pet"
	req_text := fmt.Sprintf("--%s\nContent-Disposition: form-data; name=\"attachment\"; filename=\"%s\"\nContent-Type: image/png\n\n%s\n--%s--\n", boundary, path, dat, boundary)
	actual_req := bytes.Buffer{}
	actual_req.Write([]byte(req_text))

	body, err := c.HTTPRequest("files", "POST", actual_req, &boundary)
	if err != nil {
		return err
	}

	resp_file := &ImgPostResponse{}
	err = json.NewDecoder(body).Decode(resp_file)
	if err != nil {
		return err
	}

	var width *string
	map_width, ok := para.Attrs["width"]
	string_width := string(map_width)
	if ok {
		width = &string_width
	}
	var height *string
	map_height, ok := para.Attrs["height"]
	string_height := string(map_height)
	if ok {
		height = &string_height
	}

	renderImgBlock(w, img, resp_file.File.ID, alt_text, width, height, true)
	return nil
}

type ImgPostResponse struct {
	File ImgPostResponseFileData `json:"file"`
}

type ImgPostResponseFileData struct {
	ID string `json:"id"`
}

func customHTMLRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if emph, ok := node.(*ast.Emph); ok {
		renderEmphasis(w, emph, entering)
		return ast.GoToNext, true
	}
	if strong, ok := node.(*ast.Strong); ok {
		renderStrong(w, strong, entering)
		return ast.GoToNext, true
	}
	if para, ok := node.(*ast.Paragraph); ok {
		for _, child := range node.AsContainer().Children {
			if img, ok := child.(*ast.Image); ok {
				// Images in paragraphs don't render.
				if entering {
					err := uploadImg(w, img, para)
					if err != nil {
						fmt.Println("ERROR", err)
						return ast.Terminate, false
					}
				}
				return ast.SkipChildren, true
			}
		}
		renderParagraph(w, para, entering, para.AsContainer().Attribute)
		return ast.GoToNext, true
	}
	if heading, ok := node.(*ast.Heading); ok {
		renderHeading(w, heading, entering)
		return ast.GoToNext, true
	}
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCodeBlock(w, code, entering)
		return ast.GoToNext, true
	}
	if math, ok := node.(*ast.Math); ok {
		renderMath(w, math, entering)
		return ast.GoToNext, true
	}
	if math, ok := node.(*ast.MathBlock); ok {
		renderMathBlock(w, math, entering)
		return ast.GoToNext, true
	}
	if link, ok := node.(*ast.Link); ok {
		renderLink(w, link, entering)
		return ast.GoToNext, true
	}
	if list, ok := node.(*ast.List); ok {
		renderList(w, list, entering)
		return ast.GoToNext, true
	}
	if listitem, ok := node.(*ast.ListItem); ok {
		renderListItem(w, listitem, entering)
		return ast.GoToNext, true
	}
	if _, ok := node.(*ast.Image); ok {
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func RenderMDToEd(content string) string {

	opts := html.RendererOptions{
		Flags:          html.FlagsNone,
		RenderNodeHook: customHTMLRenderHook,
	}
	renderer := html.NewRenderer(opts)
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions | parser.Attributes)
	html := markdown.ToHTML([]byte(content), p, renderer)

	return "<document version=\"2.0\">" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(string(html), "\r", ""), "\n", ""), "\\n", "\n") + "</document>"
}
