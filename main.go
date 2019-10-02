package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pages, err := strconv.ParseInt(askQuestion(scanner, "How many pages would you like to scan?"), 0, 64)
	checkErr(err)

	ipArray := make([]string, pages * 40)

	for i := 1; i <= int(pages); i++ {
		go func(page int) {
			ips := scrape(page)
			for index := range ips {
				var dummyIndex int
				if page == 1 {
					dummyIndex = index
				} else {
					dummyIndex = page * 40 - index -1
				}
				ipArray[dummyIndex] = ips[index]
			}
		}(i)
	}
	var duration time.Duration

	if pages == 1 {
		duration = time.Second
	} else {
		duration = time.Second * time.Duration(pages / 2)
	}

	time.Sleep(duration)
	for index := range ipArray {
		fmt.Printf("%d - %s\n", index + 1, ipArray[index])
	}
	fmt.Print("Press 'Enter' to exit...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func askQuestion(reader *bufio.Scanner, question string) string {
	fmt.Print(question + " ")
	var response string
	for reader.Scan() {
		text := reader.Text()
		if (len(text) > 0) {
			response =  text
			break
		}
	}
	return response
}

func scrape(page int) []string {
	var serverAddresses []string

	pageURL := fmt.Sprintf("https://minecraft-server-list.com/page/%d/", page)

	fmt.Println("Preparing to scrape " + pageURL)

	document, err := goquery.NewDocument(pageURL)
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
