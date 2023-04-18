package main

import (
	"fmt"
	"strconv"
	"sync"

	"collect/site"
)

var (
	wg sync.WaitGroup
	mu sync.Mutex

	sites = [][]string{
		{"https://luxury-product.com/product-category/bags/womens-bags/womens-bags-chanel/page/%v/", "Chanel", "188"},
		{"https://luxury-product.com/product-category/bags/womens-bags/womens-bags-dior/page/%v/", "Dior", "22"},
		{"https://luxury-product.com/product-category/bags/womens-bags/womens-bags-fendi/page/%v/", "Fendi", "24"},
		{"https://luxury-product.com/product-category/bags/womens-bags/womens-bags-gucci/page/%v/", "Gucci", "32"},
		{"https://luxury-product.com/product-category/bags/womens-bags/womens-bags-louis-vuitton/page/%v/", "Louis Vuitton", "49"},
		{"https://luxury-product.com/product-category/bags/womens-bags/womens-bags-saint-laurent/page/%v/", "Saint Laurent", "15"},
	}
)

type CollectRange struct {
	Domain string `json:"domain"`
	Tag    string `json:"tag"`
	Page   int    `json:"page"`
}

func main() {

	for _, v := range sites {
		page, _ := strconv.Atoi(v[2])
		cr := &CollectRange{
			Domain: v[0],
			Tag:    v[1],
			Page:   page,
		}

		for i := 2; i <= cr.Page; i++ {
			domain := fmt.Sprintf(cr.Domain, i)

			doTask(domain)

		}
	}

	// wg.Wait()

}

func doTask(domain string) {
	// wg.Add(1)
	// go func() {
	// defer wg.Done()

	domains, e := site.GetArticleList(domain)
	if e != nil {
		return
	}

	fmt.Println(domains)

	for _, v := range domains {
		art, err := site.Collect(v)
		if err != nil {
			fmt.Println(domain, err)
		}

		if art.Title == "" {
			fmt.Println("title is null", domain)
		}

		art.Purge()
	}

	// }()

}
