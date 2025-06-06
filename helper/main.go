package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"helper/feeds"
	"helper/links_client"
	"helper/text"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "update" {
		updateOnce()
		return
	}

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

func updateOnce() error {
	links, err := links_client.List()
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		return err
	}
	defer db.Close()

	const schema = `
		CREATE TABLE IF NOT EXISTS articles (
			id            BIGINT PRIMARY KEY AUTO_INCREMENT,
			link_id       BIGINT NOT NULL,
			title         TEXT    NOT NULL,
			url           TEXT NOT NULL,
			published_at  DATETIME NOT NULL,
			UNIQUE KEY uq_link_url (link_id, url(191)),
			FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
		);`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	ins, _ := db.Prepare(`
    	INSERT INTO articles (link_id,title,url,published_at)
    	VALUES (?,?,?,?)
    	ON DUPLICATE KEY UPDATE id = id`)
	if err != nil {
		return err
	}
	defer ins.Close()

	for _, l := range links {
		items, err := feeds.ParseRSS(l.URL, 5)
		if err != nil {
			log.Printf("rss %d: %v", l.ID, err)
			continue
		}
		for _, it := range items {
			t := time.Now()
			if it.PublishedParsed != nil {
				t = *it.PublishedParsed
			}
			_, _ = ins.Exec(l.ID, it.Title, it.Link, t)
		}
	}
	return nil
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
