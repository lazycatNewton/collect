package main

import (
	"fmt"
	"go.x2ox.com/sorbifolia/httputils"
	"strconv"
	"sync"
	"time"

	"collect/site"
)

var (
	wg sync.WaitGroup

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

			tags := []string{cr.Tag}

			doTask(domain, tags)

		}
	}

	wg.Wait()
}

func doTask(domain string, tags []string) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		domains, e := site.GetArticleList(domain)
		if e != nil {
			return
		}

		for _, v := range domains {
			art, err := site.Collect(v)
			if err != nil {
				fmt.Println(domain, err)
			}

			if art.Title == "" {
				fmt.Println("title is null", domain)
			}

			art.Purge()

			art.Tags = tags

			if err = httputils.Post("http://127.0.0.1:8080/api/v1/admin/article/").
				SetBodyWithEncoder(httputils.JSON(), art).
				SetContentType(httputils.AppJSON).
				SetHeader("Authorization", "Bearer eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODE5MDk0MTIsImRhdGEiOnsia2V5IjoiSyJ9fQ.PMqEloPVhN8Nghs0j6MGBXvbr_Cvoe1Cnc4zIDoIOFd8zChWYKTIixyzg87MkciOdfSEZnvxA_44y6CKKmvOBg").
				Request(3, nil, 60*time.Second).DoRelease(); err != nil {
				fmt.Println("Upload err ", v, err)
			}

			fmt.Println("Done uploading", v)

		}

	}()

}
