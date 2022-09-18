package functions

import (
	"encoding/json"
	"github.com/jasonlvhit/gocron"
	"github.com/reujab/wallpaper"
	"github.com/rodkranz/fetch"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

func GetWallpaper() string {
	file, readErr := os.ReadFile("config.yaml")
	if readErr != nil {
		panic(readErr)
	}
	data := make(map[interface{}]interface{})

	marshalErr := yaml.Unmarshal(file, &data)
	if marshalErr != nil {
		panic(marshalErr)
	}

	type Wallpaper struct {
		Hdurl      string `json:"hdurl"`
		Url        string `json:"url"`
		Media_type string `json:"media_type"`
	}

	client := fetch.NewDefault()
	res, getErr := client.Get("https://api.nasa.gov/planetary/apod?api_key="+data["apiKey"].(string), nil)
	if getErr != nil {
		panic(getErr)
	}

	body, StringErr := res.ToString()
	if StringErr != nil {
		panic(StringErr)
	}

	var wallpaper Wallpaper
	jsonErr := json.Unmarshal([]byte(body), &wallpaper)

	if jsonErr != nil {
		panic(jsonErr)
	}

	if wallpaper.Media_type == "video" {
		wallpaper.Url = wallpaper.Url[30 : len(wallpaper.Url)-6]

		return "https://img.youtube.com/vi/" + wallpaper.Url + "/0.jpg"
	}

	return wallpaper.Hdurl
}

func SetWallpaper() {
	err := wallpaper.SetFromURL(GetWallpaper())
	if err != nil {
		panic(err)
	}
}

func StartWallpaper() {
	SetWallpaper()
	times := time.Now()
	t := time.Date(times.Year(), times.Month(), times.Day(), 4, 50, times.Second(), times.Nanosecond(), time.UTC)
	err := gocron.Every(1).Day().From(&t).Do(SetWallpaper)
	if err != nil {
		panic(err)
	}
}