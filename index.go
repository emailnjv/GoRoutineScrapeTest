package main

import (

	"net/http"
	//"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"fmt"
	"sync"
	"strings"
	"strconv"
)



var(

	resultList [] string
	urls [] string
	pageMax = 400
	page = 1
	startingUrl = "https://www.amazon.com/s/ref=sr_pg_"+string(page)+"?rh=i%3Aaps%2Ck%3ASony&page="+string(page)+"&keywords=Sony&ie=UTF8&qid=1498947694&spIA=B00J7VMRMW,B01IT8LXS2,B06Y2HG62W,B01N7S869U"
	dellSURL = "https://www.amazon.com/s/ref=sr_pg_"+string(page)+"?fst=as%3Aon&rh=k%3Adell%2Cn%3A172282%2Cn%3A541966%2Cn%3A13896617011&page="+string(page)+"&keywords=dell&ie=UTF8&qid=1498953958&spIA=B012PTEA0O,B01MS6TKUA,B01CH72880,B00UWD90FQ"

)

type AmazonResult struct {

	Title string
	SoldBy string
	Info string

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
		go func(n *html.Node) {
			resultMatcher := func(n *html.Node) bool {
				if n.DataAtom == atom.Div && n != nil {
					return scrape.Attr(n, "class") == "a-fixed-left-grid-col a-col-right"
				}
				return false
			}

			resultNodes := scrape.FindAll(n, resultMatcher)
			    for _, result := range resultNodes{
					if strings.Contains(scrape.Text(result),"Sponsored P.when"){
						preString := strings.SplitAfter(scrape.Text(result),"? Leave ad feedback ")
						out <- preString[1]
					}else{
						out <- scrape.Text(result)
					}


				}
			wg.Done()


		}(root)

	}
		go func() {
			wg.Wait()
			close(out)
		}()
		return out
	}




func urlScraper(){
	for i := page; i < pageMax+1; i++ {
		urls = append(urls,"https://www.amazon.com/s/ref=sr_pg_"+strconv.Itoa(page)+"?fst=as%3Aon&rh=k%3Adell%2Cn%3A172282%2Cn%3A541966%2Cn%3A13896617011&page="+strconv.Itoa(page)+"&keywords=dell&ie=UTF8&qid=1498953958&spIA=B012PTEA0O,B01MS6TKUA,B01CH72880,B00UWD90FQ")
	}
	for result := range resultNodeGen(rootGen(respGen(urls...))) {
		fmt.Println(result)
		fmt.Println("***************************")
	}
}







func main() {


	urlScraper()

	//resp, err := http.Get(startingUrl)
	//if err != nil {
	//	panic(err)
	//}
	//root, err := html.Parse(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// define a matcher
	//matcher := func(n *html.Node) bool {
	//	// must check for nil values
	//	if n.DataAtom == atom.Div && scrape.Attr(n, "id") ==  {
	//		return scrape.Attr(n.Parent.Parent, "class") == "athing"
	//	}
	//	return false
	//}
	//// grab all articles and print them
	//articles := scrape.FindAll(root, matcher)
	//for i, article := range articles {
	//	fmt.Printf("%2d %s (%s)\n", i, scrape.Text(article), scrape.Attr(article, "href"))
	//}
	//

}