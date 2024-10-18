package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

type CoinData struct {
	CoinName  string
	Coin      string
	Price     string
	MarketCap string
}

func main() {

	// url to scrape
	scrapeUrl := "https://www.coindesk.com/price"

	//slice to hold coin objects that hold coin data
	CoinObjects := []CoinData{}
	var keywords = []string{}

	// accept comma separated keywords from terminal command
	KeywordStr := flag.String("keywords", "", "comma-separated list of coins to scrape for")
	flag.Parse()

	HandleKeywords(KeywordStr, &keywords)

	//Launch browser (headless mode)
	path, _ := launcher.New().Headless(true).Launch()

	// Connect to browser
	browser := rod.New().ControlURL(path).MustConnect()

	//to avoid zombie processes running..
	defer browser.MustClose()

	// Navigate to the target URL
	page := browser.MustPage(scrapeUrl)

	// Wait for the page to load (dynamic content)
	page.MustWaitLoad()

	// You can now proceed to scrape data or interact with the page
	log.Println("Page loaded!")
	PageFullyLoad(page)

	fmt.Println("resulting keywords are: ")
	fmt.Println(keywords)

	//get search bar element
	search_bar, err := page.Element(`input[data-slot="input"]`)
	if err != nil {
		fmt.Println("search bar access failed: ", err)
		return
	}

	//loop through keywords slice
	for _, kw := range keywords {

		//click on search bar & enter curr keyword/s
		search_bar.Input(kw)
		time.Sleep(2 * time.Second)

		table, err := page.Element("table tbody")
		if err != nil {
			fmt.Println("current table not found: ", err)
			return
		}

		// wait for modal rows to load, pass table elem to
		// MoveThroughCurrentRows, then move through
		// all rows to find match w/ curr keyword
		MoveThroughCurrentRows(page, table, kw, &CoinObjects)

		time.Sleep(500 * time.Millisecond)

		//after current rows finished, clearing text from search bar
		if _, err := search_bar.Eval(`function () {this.value = "";}`); err != nil {
			fmt.Println("Failed to clear text from search bar: ", err)
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

	if len(CoinObjects) == 0 {
		fmt.Println("Sorry, no coins were found..")
		return

	}

	// if the length of CoinObjects is the same as coins wanted,
	// that means all coins were found
	if len(CoinObjects) == len(keywords) {
		fmt.Println("All coins found!")
	}

	//now after all keywords have been gone through, print resulting coin objects
	fmt.Println("Logging resulting coin data: ")
	EncodeAndPrint(CoinObjects)

}

func HandleKeywords(kw *string, keywords *[]string) {

	if *kw == "" {
		fmt.Println("No keywords added. Default coin list will be used.")
		// here keywords will hold the default vals
		*(keywords) = append(*(keywords), "btc", "eth", "xrp")

	} else {
		*(keywords) = slices.Concat(*(keywords), strings.Split(*kw, ","))
		//alter all keyword entries to lower case & store in slice
		//defaults not used so need to modify keywords
		for i := range *(keywords) {
			(*keywords)[i] = strings.ToLower(strings.TrimSpace((*keywords)[i]))
			// if any of the keyword entries contain a ".", or special char
			// then throw error & terminate program
			if (*keywords)[i] == "" || !IsAllLetters((*keywords)[i]) {
				log.Fatal(errors.New("incorrect keyword found"))
			}
		}
		fmt.Println("Coin keywords will be used.")

	}

}

func IsAllLetters(str string) bool {
	for _, l := range str {
		if l == ' ' {
			continue
		}
		if !unicode.IsLetter(l) {
			return false
		}
	}
	return true
}

func EncodeAndPrint(obj []CoinData) {
	//standard output as I/O writer
	encode := json.NewEncoder(os.Stdout)
	encode.SetIndent("", " ")
	encode.Encode(obj)
}

func MoveThroughCurrentRows(page *rod.Page, table *rod.Element, keyword string, CO *[]CoinData) {
	// Wait for next page to load before scraping again
	page.MustWaitLoad()
	time.Sleep(2 * time.Second)
	var coin_type string
	var coin_name string

	//set current rows after search entered
	rows, _ := table.Elements(".tr")

	//return out of function if num of rows < 1
	if len(rows) < 1 {
		fmt.Println("No rows for current keyword to search through..")
		return

	}

	//move through current filtered rows
	for i := 0; i < len(rows); i++ {

		// Create a context with a 10-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// here defer would mean it doesn't get called until
		// MoveThroughPageRows returns
		// defer cancel()

		if !WaitForStableDataContext(page, ctx, rows[i]) {
			// timeout occurred

			//context now no longer needed
			// will cleanup context immediately
			cancel()

			fmt.Println("Row is unstable..")
			return
		}

		// now context will be cleaned up immediately and
		// NOT wait until the whole parent function is returned
		cancel()

		// row element here is stable

		row_elements, _ := rows[i].Elements(".td")

		//targeting coin name
		cn, cn_err := row_elements[1].Eval(`function() {
			const child = this.querySelector("h5");
			return child ? child.innerText : "";}`, nil)

		//targeting coin code
		ct, ct_err := row_elements[1].Eval(`function() {
			const child = this.querySelector("span");
			return child ? child.innerText : "";}`, nil)

		//targeting coin price
		price, pr_err := row_elements[3].Eval(`function() {
				const child = this.querySelector("h5");
				return child ? child.innerText : "";}`, nil)

		//targeting coin market cap
		MarketCap, mc_err := row_elements[6].Eval(`function() {
			const child = this.querySelector("h5");
			return child ? child.innerText : "";}`, nil)

		//check for errors from all wanted data

		if ct_err != nil {
			log.Fatalf("Failed to get coin type text: %v", ct_err)
		}
		if cn_err != nil {
			log.Fatalf("Failed to get coin name text: %v", cn_err)
		}
		if pr_err != nil {
			log.Fatalf("Failed to get coin price text: %v", pr_err)
		}
		if mc_err != nil {
			log.Fatalf("Failed to get coin market cap: %v", mc_err)
		}

		coin_type = strings.ToLower(ct.Value.Str())
		coin_name = strings.ToLower(cn.Value.Str())

		fmt.Println("curr coin type: ", coin_type)
		fmt.Println("curr coin name: ", coin_name)

		// if current keyword matches current row then add
		// coin data object to coin objects slice
		if keyword == coin_name || keyword == coin_type {

			CD := CoinData{}
			CD.CoinName = coin_name
			CD.Coin = coin_type
			CD.Price = price.Value.Str()
			CD.MarketCap = MarketCap.Value.Str()
			(*CO) = append((*CO), CD)

			//now that specific coin was found, exit out current coin row loop
			return

		}

	}

}

// scrolls down to bottom of page & waits for content to load
func PageFullyLoad(page *rod.Page) {
	// Scroll to the bottom to trigger lazy loading or additional content
	page.Keyboard.Press(input.End)

	// Wait for 2 seconds to give time for content to load (adjust this if necessary)
	time.Sleep(2 * time.Second)

	// Optionally, wait for the page to become stable
	page.MustWaitIdle()

	// Scroll back to the top if needed
	page.Keyboard.Press(input.Home)

	fmt.Println("Current page URL: ", page.MustInfo().URL)
}

func WaitForStableDataContext(page *rod.Page, ctx context.Context, row *rod.Element) bool {
	for {
		select {
		case <-ctx.Done():
			return false // Timeout or cancellation occurred
		default:

			time.Sleep(500 * time.Millisecond)

			row_elements, _ := row.Elements(".td")

			//check if coin name loaded
			loaded, err := row_elements[1].Eval(`function () {
				const child = this.querySelector("h5");
				return child ? child.innerText : "";}`, nil)

			if err != nil {
				fmt.Println("Failed to get inner text: ", err)
				// if there is an err, return false to ensure that after
				// parent func returns, cancel() gets called to clean up context
				return false
			}

			// if it's not loaded, then trying again
			if loaded.Value.Str() == "" {

				// grab html of row to see if data needed is there
				html, err := row.HTML()
				if err != nil {
					fmt.Println("Failed to get HTML of current row: ", err)
				}

				fmt.Println("HTML of current row is: ")
				fmt.Println(html)
				fmt.Println("text element not ready..trying again..")

				time.Sleep(2 * time.Second)
				continue

			}

			//if coin name text exists then return true
			return true

		}
	}
}
