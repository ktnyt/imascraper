package main

import (
	"log"

	"github.com/ktnyt/imascraper/cg"
)

func run() (err error) {
	c := new(cg.Card)
	return c.Scrape("28", "eb2571c3f125aa3fcadb1468b6a4dbee")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
