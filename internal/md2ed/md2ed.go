package md2ed

import (
	"bytes"
	"fmt"
	"io"
	"strings"

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

func renderCodeBlock(w io.Writer, p *ast.CodeBlock, entering bool) {

	if len(p.Info) == 0 {
		io.WriteString(w, "<pre>")
		io.WriteString(w, string(p.Literal))
		io.WriteString(w, "</pre>")
	} else {
		fields := strings.FieldsFunc(string(p.Info), func(s rune) bool { return s == '.' })
		lang := fields[0]
		extra := ""

		if len(fields) > 1 {
			extra = strings.Join(fields[1:], " ")
		}

		io.WriteString(w, fmt.Sprintf("<snippet language=\"%s\" %s><snippet-file id=\"code\">", lang, extra))
		io.WriteString(w, string(bytes.Trim(p.Literal, "\n")))
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
	return ast.GoToNext, false
}

func RenderMDToEd(content string) string {

	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: customHTMLRenderHook,
	}
	renderer := html.NewRenderer(opts)
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions | parser.Attributes)
	html := markdown.ToHTML([]byte(content), p, renderer)

	return "<document version=\"2.0\">" + strings.Replace(string(html), "\r", "", -1) + "</document>"
}
