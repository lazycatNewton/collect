package site

import (
	"fmt"

	"collect/article"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

var data = make(map[string]*Site)

type Site struct {
	Domain          string
	ArticleSelector string
	ArticleURL      string
	TitleSelector   string
	ContentSelector string
	ImgSelector     string
}

func GetArticleList(domain string) ([]string, error) {
	list := make([]string, 0)
	r, err := getSite(domain)
	if err != nil {
		return nil, err
	}
	doc := &goquery.Document{}
	if doc, err = goquery.NewDocumentFromReader(r); err != nil {
		return nil, err
	}

	doc.Find("#main > ul > li").Each(func(i int, selection *goquery.Selection) {
		url, ok := selection.Find("a").Attr("href")
		if ok && url != "" {
			list = append(list, url)
		}
	})

	return list, nil
}

func Collect(domain string) (*article.Article, error) {
	var (
		err error
	)

	s := data["luxury-product.com"]
	if s == nil {
		return nil, fmt.Errorf("%s is not support", domain)
	}

	r, e := getSite(fmt.Sprintf(domain))
	if e != nil {
		return nil, e
	}

	var (
		art  article.Article
		doc  *goquery.Document
		imgs string
	)

	if doc, err = goquery.NewDocumentFromReader(r); err != nil {
		return nil, err
	}

	doc.Find(s.TitleSelector).Each(func(i int, selection *goquery.Selection) {
		if title := selection.Text(); title != "" {
			art.Title = title
		}
	})

	doc.Find(s.ImgSelector).Each(func(i int, selection *goquery.Selection) {
		img, ok := selection.Find("a").Attr("href")
		if ok {
			img = fmt.Sprintf("![](%v)", img)
			imgs = imgs + "\n" + img
		}
	})

	doc.Find(s.ContentSelector).Each(func(i int, selection *goquery.Selection) {
		converter := md.NewConverter(domain, true, nil)
		art.Content = converter.Convert(selection)

		art.Content = art.Content + imgs
	})

	return &art, nil
}

func register(domain string, site *Site) {
	data[domain] = site
}
