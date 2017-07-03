package main

import (

	"net/http"
	//"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"sync"
	"strings"
	"strconv"
	"os"
	"fmt"
	"encoding/json"
)



var(

	resultList [] AmazonResult
	urls [] string
	pageMax = 400
	page = 1

)

type AmazonResult struct {

	Info string `json:"info"`

}

func respGen(urls ...string) <-chan *http.Response {
	var wg sync.WaitGroup
	out := make(chan *http.Response)
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			out <- resp
			wg.Done()
		}(url)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func rootGen(in <-chan *http.Response) <-chan *html.Node {
	var wg sync.WaitGroup
	out := make(chan *html.Node)
	for resp := range in {
		wg.Add(1)
		go func(resp *http.Response) {
			root, err := html.Parse(resp.Body)
			if err != nil {
				panic(err)
			}
			out <- root
			wg.Done()
		}(resp)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}






func resultNodeGen(in <-chan *html.Node) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)
	for root := range in {
		wg.Add(1)
			resultMatcher := func(n *html.Node) bool {
				if n.DataAtom == atom.Div && n != nil {
					return scrape.Attr(n, "class") == "a-fixed-left-grid-col a-col-right"
				}
				return false
			}

			var results = scrape.FindAll(root, resultMatcher)
		var wg1 sync.WaitGroup

		for _, result := range results{
			wg1.Add(1)
			//out := make(chan *html.Node)
					go func(n *html.Node) {
						defer wg1.Done()
						if strings.Contains(scrape.Text(n),"Sponsored P.when"){
						preString := strings.SplitAfter(scrape.Text(n),"? Leave ad feedback ")
						out <- preString[1]

					}else{
						out <- scrape.Text(n)

					}
				}(result)


				}
		go func() {
			wg1.Wait()
			wg.Done()
		}()





	}
		go func() {
			wg.Wait()
			close(out)
		}()
		return out
	}




func urlScraper(){
	for i := page; i < pageMax+1; i++ {
		urls = append(urls,"https://www.amazon.com/s/ref=sr_pg_"+strconv.Itoa(i)+"?fst=as%3Aon&rh=k%3Adell%2Cn%3A172282%2Cn%3A541966%2Cn%3A13896617011&page="+strconv.Itoa(i)+"&keywords=dell&ie=UTF8&qid=1498953958&spIA=B012PTEA0O,B01MS6TKUA,B01CH72880,B00UWD90FQ")
	}
	for result := range resultNodeGen(rootGen(respGen(urls...))) {
		resultList = append(resultList, AmazonResult{result})
	}
	jsonData, err := json.Marshal(resultList)

	if err != nil {
		panic(err)
	}

	// sanity check

	// write to JSON file

	jsonFile, err := os.Create("./DellResults.json")

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())

}







func main() {


	urlScraper()


}
