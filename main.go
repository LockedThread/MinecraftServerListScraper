package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"time"
)

const pages = 10

func main() {
	var allIps []string

	for i := 1; i <= pages; i++ {
		go func(index int) {
			ips := scrape(index)
			allIps = append(allIps, ips...)
		}(i)
	}

	time.Sleep(time.Second)

	for index := range allIps {
		fmt.Println(allIps[index])
	}

	fmt.Print("Press 'Enter' to exit...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func scrape(page int) []string {
	var serverAddresses []string

	pageUrl := fmt.Sprintf("https://minecraft-server-list.com/page/%d/", page)

	fmt.Println("Preparing to scrape " + pageUrl)

	document, err := goquery.NewDocument(pageUrl)
	if err != nil {
		log.Fatal(err)
	}

	document.Find("tr td").Each(func(index int, selection *goquery.Selection) {
		for nodeIndex := range selection.Nodes {
			node := selection.Nodes[nodeIndex]
			for attrIndex := range selection.Nodes[nodeIndex].Attr {
				attribute := node.Attr[attrIndex]
				if attribute.Key == "id" {
					serverAddresses = append(serverAddresses, attribute.Val)
				}
			}
		}
	})
	return serverAddresses
}
