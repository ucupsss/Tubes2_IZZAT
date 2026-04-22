package main

import (
	"fmt"
	"tubes2_izzat/dom"
	"tubes2_izzat/src/backend/scraper"
	"tubes2_izzat/src/backend/selector"
)

func main() {
	fmt.Println("Mencoba fungsi Scraper...")
	// Tes Scraper (internet has to be connected)
	htmlString, err := scraper.FetchHTML("https://kompas.com")
	if err != nil {
		fmt.Println("Error Scraper:", err)
	} else {
		fmt.Println("Berhasil download HTML! Panjang karakter:", len(htmlString))
	}

	fmt.Println("\nMencoba fungsi Matcher dengan Dummy Node...")
	//Tes Matcher (Membuat node palsu untuk testing)
	// Misal, <div id="utama" class="container text-center">
	atributPalsu := map[string]string{
		"id":    "utama",
		"class": "container text-center",
	}
	dummyNode := dom.NewNode("div", atributPalsu)

	//tes Parse Selector
	inputUser := ".container" // Anggap user nyari class container
	matcher := selector.ParseSelector(inputUser)

	//Cek apakah cocok
	cocok := matcher.IsMatch(dummyNode)
	if cocok {
		fmt.Printf("Selector '%s' COCOK dengan node tersebut!\n", inputUser)
	} else {
		fmt.Printf("Selector '%s' TIDAK COCOK dengan node tersebut.\n", inputUser)
	}
}