// Copyright © 2015 Eric Daugherty

package scraper

import (
	"testing"
	"time"
)

func TestNeedhamTourney2011(t *testing.T) {

	params := map[string]string{
		"EventID": "15267",
		"GroupID": "166875",
		"Gender":  "Boys",
		"Age":     "12",
	}

	schedule, err := GetSchedule(params)
	if err != nil {
		t.Error(err)
	}

	games := schedule.Games
	if len(games) != 15 {
		t.Error("Expected 15 games but got", len(games))
	}

	game0 := games[0]
	if game0.Date != time.Date(2011, time.May, 28, 8, 0, 0, 0, time.UTC) ||
		game0.Number != "#148" ||
		game0.HomeTeam != "FAR POST SC U12 BOYS PREMIER (VT)" ||
		game0.HomeScore != "2" ||
		game0.AwayTeam != "WORLD CUP SOCCER OF NASHUA PREMIER (NH)" ||
		game0.AwayScore != "5" {
		t.Error("Expected 2011-05-28 08:00:00 #148 FAR POST SC U12 BOYS PREMIER (VT) 2 vs WORLD CUP SOCCER OF NASHUA PREMIER (NH) 5", game0)
	}

	game14 := games[14]
	if game14.Date != time.Date(2011, time.May, 30, 14, 10, 0, 0, time.UTC) ||
		game14.Number != "#1122" ||
		game14.HomeTeam != "MCU PORTLAND PHOENIX ELITE (ME)" ||
		game14.HomeScore != "1" ||
		game14.AwayTeam != "FC BLAZERS (MA)" ||
		game14.AwayScore != "2" {
		t.Error("Expected 2011-05-30 14:10:00 #1122 MCU PORTLAND PHOENIX ELITE (ME) 1 vs FC BLAZERS (MA) 2", game14)
	}
}
