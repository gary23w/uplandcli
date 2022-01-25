package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

// BLOCK CHAIN LABELS
// n2 = New Listings
//

// Other notes
// Data:{ "a54": "qnhmzlgkpjpp", "a45": "79534306738459", "p11": "84900.00 UPX", "p3": "0.00 FIAT", "p31": 0, "p32": 1, "p33": 1 }

func GetPropName(prop_id string) *strings.Reader {
    url := "https://play.upland.me/?prop=" + prop_id
    options := []chromedp.ExecAllocatorOption{
      chromedp.Flag("headless", true), // debug usage
      chromedp.Flag("blink-settings", "imagesEnabled=false"),
      chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
    }

    options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

    c, _ := chromedp.NewExecAllocator(context.Background(), options...)


    chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))

    chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

    timeoutCtx, cancel := context.WithTimeout(chromeCtx, 5*time.Second)
    defer cancel()
    var res string                                               
    //var screnshot []byte
    err := chromedp.Run(timeoutCtx,                                                         
      chromedp.Navigate(url),                                         
      // wait main page table to load
      //chromedp.WaitVisible(`//*[@id="table-wrapper"]/table`),
      //chromedp.CaptureScreenshot(&screnshot),
      chromedp.ActionFunc(func(timeoutCtx context.Context) error {                          
        node, err := dom.GetDocument().Do(timeoutCtx)                                       
        if err != nil {                                                              
          return err                                                               
        }                                                                            
        res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(timeoutCtx)                
        return err                                                                   
    }),                                                                            
  )                                                                                
                                                                                  
  if err != nil {                                                                  
    fmt.Println(err)                                                               
  }
  return strings.NewReader(res)
}



func BlockchainScraper() *strings.Reader {
    url := "https://bloks.io/account/playuplandme"

    options := []chromedp.ExecAllocatorOption{
      chromedp.Flag("headless", true), // debug usage
      chromedp.Flag("blink-settings", "imagesEnabled=false"),
      chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
    }
  
    options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
  
    c, _ := chromedp.NewExecAllocator(context.Background(), options...)
  
  
    chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
  
    chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
  
    timeoutCtx, cancel := context.WithTimeout(chromeCtx, 5*time.Second)
    defer cancel()
    var res string                                                                   
    var screnshot []byte
    err := chromedp.Run(timeoutCtx,                                                         
      chromedp.Navigate(url),                                         
		  // wait main page table to load
		  chromedp.WaitVisible(`//*[@id="table-wrapper"]/table`),
      chromedp.CaptureScreenshot(&screnshot),
      chromedp.ActionFunc(func(timeoutCtx context.Context) error {                          
        node, err := dom.GetDocument().Do(timeoutCtx)                                       
        if err != nil {                                                              
          return err                                                               
        }                                                                            
        res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(timeoutCtx)                
        return err                                                                   
    }),                                                                            
  )                                                                                
                                                                                   
  if err != nil {                                                                  
    fmt.Println(err)                                                               
  }
  cancel()
  return strings.NewReader(res)
}

func ProcessBlockchainCollection(res *strings.Reader) {
  doc, err := goquery.NewDocumentFromReader(res)
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("tr").Each(func(i int, q *goquery.Selection) {
    if i == 0 {
      return
    }
    var row []string
    q.Find("td").Each(func(j int, s *goquery.Selection) {
      check, err := regexp.MatchString("\\bn2\\b", s.Text())
      if err != nil {
        log.Fatal(err)
      }
      if check {
        if strings.Contains(q.Text(), "a45") {
            cha45 := regexp.MustCompile(`a45:\s(.*)\sa54`)
            cha45_match := cha45.FindStringSubmatch(q.Text())
            fmt.Printf("cha45_match: %v \n", cha45_match[1])
            idurl := strings.Replace(cha45_match[1], " ", "", -1)
            fmt.Printf("https://play.upland.me/?prop_id=%v \n", idurl)
        }
      }
      row = append(row, s.Text())
    })  
  })
}


