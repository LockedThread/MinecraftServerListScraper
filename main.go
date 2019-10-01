package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	pages, err := strconv.ParseInt(askQuestion(reader, "How many pages would you like to scan?"), 0, 64)

	checkErr(err)

	var allIps []string

	for i := 1; i <= int(pages); i++ {
		go func(index int) {
			allIps = append(allIps, scrape(index)...)
		}(i)
	}
	time.Sleep(time.Second)
	for index := range allIps {
		fmt.Println(allIps[index])
	}
	fmt.Print("Press 'Enter' to exit...")
	_, _ = reader.ReadBytes('\n')
}

func askQuestion(reader *bufio.Reader, question string) string {
	fmt.Print(question + " ")
	response, _ := reader.ReadString('\n')
	return strings.Replace(response, "\r\n", "", 1)
}

func scrape(page int) []string {
	var serverAddresses []string

	pageUrl := fmt.Sprintf("https://minecraft-server-list.com/page/%d/", page)

	fmt.Println("Preparing to scrape " + pageUrl)

	document, err := goquery.NewDocument(pageUrl)
	checkErr(err)

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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
