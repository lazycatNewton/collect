package article

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

type Article struct {
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
	Desc    string   `json:"desc"`
	Avatar  string   `json:"avatar"`
	Content string   `json:"content"`

	Images []string `json:"-"`
}

func (art *Article) Purge() {
	art.parse()
	art.replaceImage()
}

func (art *Article) parse() {
	node := markdown.Parse([]byte(art.Content), parser.New())
	art.Desc = art.getDesc(node)
	art.Images = art.getImage(node)

	if len(art.Images) != 0 {
		art.Avatar = art.Images[0]
	}
}

func (art *Article) getDesc(node ast.Node) string {
	for _, v := range node.GetChildren() {
		if text, ok := v.(*ast.Text); ok && text != nil && len(text.Leaf.Literal) != 0 {
			return string(text.Leaf.Literal)
		}
		if text := art.getDesc(v); text != "" {
			return text
		}
	}
	return ""
}

func (art *Article) getImage(node ast.Node) []string {
	var arr []string

	for _, v := range node.GetChildren() {
		if image, ok := v.(*ast.Image); ok && image != nil && len(image.Destination) != 0 {
			arr = append(arr, string(image.Destination))
			continue
		}
		if list := art.getImage(v); len(list) != 0 {
			arr = append(arr, list...)
		}
	}
	return arr
}
