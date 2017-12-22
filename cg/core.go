package cg

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Idol struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Profile struct {
	Rarity     string `json:"rarity"`
	BloodType  string `json:"blood_type"`
	Height     string `json:"height"`
	Weight     string `json:"weight"`
	SizeB      string `json:"size_b"`
	SizeW      string `json:"size_w"`
	SizeH      string `json:"size_h"`
	Age        string `json:"age"`
	Birthday   string `json:"birthday"`
	Zodiac     string `json:"zodiac"`
	Birthplace string `json:"birthplace"`
	Hobby      string `json:"hobby"`
	Handedness string `json:"handedness"`
}

type Ability struct {
	Name   string `json:"name"`
	Effect string `json:"effect"`
}

type Status struct {
	MaxLevel    int     `json:"max_level"`
	InitAttack  int     `json:"init_attack"`
	InitDefense int     `json:"init_defense"`
	MaxAttack   int     `json:"max_attack"`
	MaxDefense  int     `json:"max_defense"`
	Cost        int     `json:"cost"`
	Ability     Ability `json:"ability"`
}

type Images struct {
	Frame   string `json:"frame"`
	Noframe string `json:"noframe"`
	Quest   string `json:"query"`
	LS      string `json:"ls"`
}

type Card struct {
	CardID    string  `json:"card_id"`
	CardName  string  `json:"card_name"`
	IdolID    string  `json:"idol_id"`
	Published string  `json:"published"`
	Profile   Profile `json:"profile"`
	Comment   string  `json:"comment"`
	Status    Status  `json:"status"`
	Images    Images  `json:"images"`
}

const (
	host       = "http://imas.gamedbs.jp"
	detailPath = "/cg/idol/detail/%s?h=%s"
)

func detailURL(args ...string) string {
	path := fmt.Sprintf("cg/idol/detal/%s", args[0])
	if len(args) > 1 {
		path += fmt.Sprintf("?h=%s", args[1])
	}

	return fmt.Sprintf("%s/%s", host, path)
}

func imageURL(path string) string {
	return fmt.Sprintf("%s%s", host, path)
}

func (c *Card) Scrape(idolID, cardID string) error {
	doc, err := goquery.NewDocument(detailURL(idolID, cardID))

	if err != nil {
		return fmt.Errorf("ScrapeCard: %s", err)
	}

	cardGallery := doc.Find("div#card-gallery")

	cardHeader := cardGallery.Find("h2")

	c.CardID = cardID
	c.CardName = cardHeader.Contents().First().Text()
	c.IdolID = idolID
	c.Published = cardHeader.Children().First().Text()

	cardImages := cardGallery.Find("table").Eq(0)

	cardImages.Find("img").Each(func(index int, s *goquery.Selection) {
		if path, ok := s.Attr("data-original"); ok {
			res, err := http.Get(imageURL(path))

			if err != nil {
				log.Printf("ScrapeCard: %s", err)
				return
			}

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				log.Printf("ScrapeCard: %s", err)
				return
			}

			image := "data:image/png;base64," + base64.StdEncoding.EncodeToString(body)

			prefix := strings.Split(path, "/")[4]

			switch prefix {
			case "l":
				c.Images.Frame = image
			case "l_noframe":
				c.Images.Noframe = image
			case "quest":
				c.Images.Quest = image
			case "ls":
				c.Images.LS = image
			}
		}
	})

	return nil
}

func ScrapeIdol(idolID string) ([]Card, error) {
	doc, err := goquery.NewDocument(detailURL(idolID))

	if err != nil {
		return nil, fmt.Errorf("ScrapeCards: %s", err)
	}

	return nil, nil
}
