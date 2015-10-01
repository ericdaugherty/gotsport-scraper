// Copyright © 2015 Eric Daugherty

package scraper

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const baseURL = "http://events.gotsport.com/events/schedule.aspx"

// Schedule represents the schedule data parsed from the
// http://events.gotsport.com/events/schedule.aspx page.
type Schedule struct {
	Games []Game `json:"games"`
}

// Game represnts the data for a given game.
type Game struct {
	Date      time.Time `json:"date"`
	Number    string    `json:"gameNum"`
	HomeTeam  string    `json:"homeTeam"`
	HomeScore string    `json:"homeScore"`
	AwayTeam  string    `json:"awayTeam"`
	AwayScore string    `json:"awayScore"`
}

// GetSchedule queries the http://events.gotsport.com/events/schedule.aspx page
// and parses the HTML into a Schedule struct.
// Parameters should match the URL parameters to pass.  Ex:
// params := map[string]string{
// "EventID": "123",
// "GroupID": "123",
// "Gender":  "Boys",
// "Age":     "10",
// }
func GetSchedule(params map[string]string) (*Schedule, error) {

	t, err := query(params)
	if err != nil {
		return nil, err
	}
	schedule, err := parse(t)
	if err != nil {
		return nil, err
	}

	return schedule, err
}

func query(params map[string]string) (*html.Tokenizer, error) {

	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	values := url.Query()
	for key := range params {
		values.Add(key, params[key])
	}
	url.RawQuery = values.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Expected Response Code 200, got:" + strconv.FormatInt(int64(resp.StatusCode), 10))
	}

	return html.NewTokenizer(resp.Body), nil
}

func parse(z *html.Tokenizer) (*Schedule, error) {

	for {
		t := advanceToStartTagWithAttr("table", "width", "98%", z)
		if t == nil {
			return nil, errors.New("Unable to identifiy first table element.")
		}
		return parse2(z)
	}
}

func parse2(z *html.Tokenizer) (*Schedule, error) {

	schedule := &Schedule{}
	currentDate := ""

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return schedule, nil
		case html.StartTagToken:
			t := z.Token()
			if isTokenTagWithAttr("font", "class", "PageHeading", &t, z) {
				z.Next()
				currentDate = z.Token().Data
			} else if isTokenTagWithAttr("tr", "bgcolor", "#ffffff", &t, z) || isTokenTagWithAttr("tr", "bgcolor", "#f5f5f5", &t, z) {
				game, err := parseGame(currentDate, z)
				if err != nil {
					return nil, err
				}
				schedule.Games = append(schedule.Games, game)
			}
		}
	}
}

func parseGame(date string, z *html.Tokenizer) (Game, error) {
	var game Game
	td := advanceToStartTag("td", z)
	if td == nil {
		return game, errors.New("Unable to find Game Number")
	}
	z.Next()
	gameNum := strings.TrimSpace(z.Token().Data)

	td = advanceToStartTag("td", z)
	if td == nil {
		return game, errors.New("Unable to find Game Time")
	}
	td = advanceToStartTag("div", z)
	if td == nil {
		return game, errors.New("Unable to find Game Time")
	}
	z.Next()
	gameTime := strings.TrimSpace(z.Token().Data)
	if gameTime == "" {
		t := advanceToTextToken(z)
		gameTime = strings.TrimSpace(t.Data)
	}

	var homeTeam, homeScore, awayTeam, awayScore string

	skipAwayScore := false

	homeTeam = parseTeamName(z)
	homeScore = parseScore(z)
	if len(homeScore) > 3 {
		awayTeam = homeScore
		homeScore = ""
		skipAwayScore = true
	} else {
		awayTeam = parseTeamName(z)
	}
	if !skipAwayScore {
		awayScore = parseScore(z)
	} else {
		awayScore = ""
	}

	gameDate, err := time.Parse("1/2/2006 3:04 PM", date+" "+gameTime)
	if err != nil {
		return game, err
	}

	return Game{gameDate, gameNum, homeTeam, homeScore, awayTeam, awayScore}, nil
}

func isTokenTagWithAttr(tagName string, attrName string, attrValue string, t *html.Token, z *html.Tokenizer) bool {
	if t.Data == tagName {
		for _, attr := range t.Attr {
			if attr.Key == attrName && attr.Val == attrValue {
				return true
			}
		}
	}
	return false
}

func advanceToStartTag(tagName string, z *html.Tokenizer) *html.Token {
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return nil
		case html.StartTagToken:
			t := z.Token()
			if t.Data == tagName {
				return &t
			}
		}
	}
}

func advanceToStartTagWithAttr(tagName string, attrName string, attrValue string, z *html.Tokenizer) *html.Token {
	for {
		t := advanceToStartTag(tagName, z)
		if t == nil {
			return nil
		}
		for _, attr := range t.Attr {
			if attr.Key == attrName && attr.Val == attrValue {
				return t
			}
		}
	}
}

func advanceToTextToken(z *html.Tokenizer) *html.Token {
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return nil
		case html.TextToken:
			t := z.Token()
			return &t
		}
	}
}

func parseTeamName(z *html.Tokenizer) string {
	td := advanceToStartTag("td", z)
	if td == nil {
		return ""
	}
	td = advanceToStartTag("a", z)
	if td == nil {
		return ""
	}
	td = advanceToTextToken(z)
	return strings.TrimSpace(td.String())
}

func parseScore(z *html.Tokenizer) string {
	td := advanceToStartTag("td", z)
	if td == nil {
		return ""
	}
	td = advanceToTextToken(z)
	return strings.TrimSpace(td.Data)
}
