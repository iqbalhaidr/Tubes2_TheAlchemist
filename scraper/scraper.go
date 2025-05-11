package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Recipe struct {
	Output string     `json:"Output"`
	Inputs [][]string `json:"Inputs"`
	Tier   int        `json:"Tier"`
}

func Scrape() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var results []Recipe
	var tierCtr int = -2

	doc.Find("table.list-table").Each(func(i int, table *goquery.Selection) {
		tierCtr++
		table.Find("tr").Each(func(j int, s *goquery.Selection) {
			tds := s.Find("td")
			if tds.Length() >= 2 {
				aTags := tds.Eq(0).Find("a")
				if aTags.Length() < 2 {
					return
				}
				elementHasil := strings.TrimSpace(aTags.Eq(1).Text())

				var elementBahan [][]string
				tds.Eq(1).Find("li").Each(func(j int, li *goquery.Selection) {
					aFromLi := li.Find("a")
					if aFromLi.Length() >= 4 {
						combo := []string{
							strings.TrimSpace(aFromLi.Eq(1).Text()),
							strings.TrimSpace(aFromLi.Eq(3).Text()),
						}
						elementBahan = append(elementBahan, combo)
					}
				})

				if elementHasil != "" && len(elementBahan) > 0 {
					results = append(results, Recipe{
						Output: elementHasil,
						Inputs: elementBahan,
						Tier:   tierCtr,
					})
				}
			}
		})
	})

	// Simpan ke file JSON
	file, err := os.Create("./data/recipe.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	err = encoder.Encode(results)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scraping selesai! Data disimpan ke ./data/recipe.json")
}

// package main

// import "littlealchemy/scraper"

// // import "littlealchemy/scraper"

// func main() {
// 	scraper.Scrape()
// 	// scraper.Scrape() Scrape the web and save it to ./data/recipes.json
// }
