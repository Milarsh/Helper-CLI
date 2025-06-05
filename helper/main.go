package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"helper/feeds"
	"helper/links_client"
	"helper/text"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)

	fmt.Println("Команды:  rss   |   stats   |   md2html   |   help")
	fmt.Print("Введите команду: ")
	if !reader.Scan() {
		return
	}
	cmd := strings.ToLower(strings.TrimSpace(reader.Text()))

	switch cmd {
	case "rss":
		runRSS(reader)
	case "stats":
		runStats(reader)
	case "md2html":
		runMd2HTML(reader)
	default:
		printHelp()
	}
}

func runRSS(reader *bufio.Scanner) {
	links, err := links_client.List()
	if err != nil {
		fmt.Println("links list error:", err)
		return
	}
	if len(links) == 0 {
		fmt.Println("Нет хранящихся лент")
		return
	}

	fmt.Println("Доступные ленты:")
	for _, l := range links {
		fmt.Printf("[%d] %s (%s)\n", l.ID, l.Label, l.URL)
	}

	fmt.Print("Выберите ID ленты: ")
	if !reader.Scan() {
		return
	}
	id, err := strconv.ParseInt(strings.TrimSpace(reader.Text()), 10, 64)
	if err != nil {
		fmt.Println("invalid id")
		return
	}

	link, err := links_client.Get(id)
	if err != nil {
		fmt.Println("get link error:", err)
		return
	}

	items, err := feeds.ParseRSS(link.URL, 5)
	if err != nil {
		fmt.Println("rss error:", err)
		return
	}

	for _, it := range items {
		fmt.Println("===================================")
		fmt.Println(it.Title)
		fmt.Println(it.Link)
		fmt.Println(it.Published)
		fmt.Println("-----------------------------------")
		fmt.Printf("Слов: %d | Стоп-слов: %d\n",
			text.CountWords(it.Title),
			text.CountStopWords(it.Title))
	}
}

func runStats(reader *bufio.Scanner) {
	fmt.Println("Введите текст (завершите через Ctrl + D):")
	var lines []string
	for reader.Scan() {
		line := reader.Text()
		lines = append(lines, line)
	}
	input := strings.Join(lines, "\n")

	fmt.Printf("\nСимволов   : %d\n", text.CountSymbols(input))
	fmt.Printf("Слов     : %d\n", text.CountWords(input))
	fmt.Printf("Стоп-слов: %d\n", text.CountStopWords(input))
}

func runMd2HTML(reader *bufio.Scanner) {
	fmt.Print("Введите текст в формате Markdown (завершите через Ctrl + D): ")
	var lines []string
	for reader.Scan() {
		line := reader.Text()
		lines = append(lines, line)
	}
	input := strings.Join(lines, "\n")

	html := text.MDtoHTML([]byte(input))
	fmt.Println(string(html))
}

func printHelp() {
	fmt.Println(`
		rss      – показывает последние новости из выбранной ленты
		stats    – подсчитывает количество символов, слов и стоп-слов в введённом тексте
		md2html  – конвертирует Markdown в HTML
		help     – показывает эту справку`)
}
