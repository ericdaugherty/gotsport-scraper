# GotSport.com HTML Screen Scraper
This is a GoLang library that parses schedule information from the
www.gotsport.com website.  It currently parses the schedule page, ie:
http://events.gotsport.com/events/schedule.aspx and provides a list of games
in a JSON enabled GoLang struct.

## Usage

Here is a simple program that queries for a tournement schedule and pretty-prints the results as JSON.

```
package main

import (
	"encoding/json"
	"fmt"

	"github.com/ericdaugherty/gotsport-scraper"
)

func main() {

	params := map[string]string{
		"EventID": "15267",
		"GroupID": "166875",
		"Gender":  "Boys",
		"Age":     "12",
	}

	schedule, err := scraper.GetSchedule(params)
	if err != nil {
		fmt.Println("Error Occured:", err)
	}
	bytes, err := json.MarshalIndent(schedule, "", "  ")
	if err != nil {
		fmt.Println("Error Occured:", err)
	}

	fmt.Println("Schedule:", string(bytes))
}
```

## License
This code is available for use under the MIT License.  See the LICENSE file.
