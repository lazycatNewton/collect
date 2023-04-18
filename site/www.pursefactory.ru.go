package site

func init() {
	register("luxury-product.com", &Site{
		Domain:          "luxury-product.com",
		ArticleSelector: "#main > ul > li",
		TitleSelector:   "#main > div:nth-child(2) > div.summary.entry-summary > h1",
		ContentSelector: "#tab-description",
		ImgSelector:     "#main > div:nth-child(2) > div.images > div.thumbnails.slider > ul > li",
	})
}
