// BELUM FINAL
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
	Input1 string `json:"Input1"`
	Input2 string `json:"Input2"`
	Output string `json:"Output"`
}

func Scrape() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var recipes []Recipe

	doc.Find("table.article-table tbody tr").Each(func(i int, s *goquery.Selection) {
		cols := s.Find("td")
		if cols.Length() < 3 {
			return
		}

		output := strings.TrimSpace(cols.Eq(0).Text())
		input1 := strings.TrimSpace(cols.Eq(1).Text())
		input2 := strings.TrimSpace(cols.Eq(2).Text())

		if input1 != "" && input2 != "" && output != "" {
			recipes = append(recipes, Recipe{
				Input1: input1,
				Input2: input2,
				Output: output,
			})
		}
	})

	// Simpan ke file JSON
	file, err := os.Create("recipes.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // biar rapi
	if err := encoder.Encode(recipes); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Resep berhasil disimpan ke recipes.json")
}
