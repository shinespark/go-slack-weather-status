package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
)

func GetForecastDoc(url string) (doc *goquery.Document) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func GetForecastEmojiText(doc *goquery.Document) (h2 string, emoji string, text string) {
	weatherSymbolMap := map[string]string{
		"01": ":sunny:",
		"02": ":mostly_sunny:",
		"03": ":partly_sunny_lain:",
		"04": ":snow_cloud:",
		"05": ":partly_sunny:",
		"06": ":partly_sunny_lain:",
		"07": ":snow_cloud:",
		"08": ":cloud:",
		"09": ":partly_sunny:",
		"10": ":rain_cloud:",
		"11": ":snow_cloud:",
		"12": ":partly_sunny:",
		"13": ":rain_cloud:",
		"14": ":snow_cloud:",
		"15": ":umbrella:",
		"16": ":rain_cloud:",
		"17": ":umbrella:",
		"18": ":umbrella:",
		"19": ":umbrella:",
		"20": ":rain_cloud:",
		"21": ":rain_cloud:",
		"22": ":rain_cloud:",
		"23": ":snowman_without_snow:",
		"24": ":snowman_without_snow:",
		"25": ":snowman_without_snow:",
		"26": ":snowman_without_snow:",
		"27": ":snowman_without_snow:",
		"28": ":snowman_without_snow:",
		"29": ":snowman_without_snow:",
		"30": ":snowman:",
	}

	h2Text := doc.Find("h2").Text()
	h2 = strings.Split(h2Text, "の天気")[0]

	selection := doc.Find("section.today-weather")
	iconImg := selection.Find("div.weather-icon img")
	iconUrl, _ := iconImg.Attr("src")
	iconTitle, _ := iconImg.Attr("title")
	splitIconUrl := strings.Split(iconUrl, "/")
	iconFileName := splitIconUrl[len(splitIconUrl)-1]
	iconName := strings.TrimSuffix(iconFileName, ".png")
	iconName = strings.TrimSuffix(iconName, "_n") // on night

	var ok bool
	emoji, ok = weatherSymbolMap[iconName]

	if ok {
		text = iconTitle

	} else {
		// 絵文字マップ未定義
		emoji = ":question:"
		text = fmt.Sprintf("【%s】%s", iconName, iconTitle)
	}

	highTemp := selection.Find("dd.high-temp").Text()
	lowTemp := selection.Find("dd.low-temp").Text()
	text += fmt.Sprintf(" %s 〜 %s", lowTemp, highTemp)

	return
}

type Config struct {
	ForecastUrl string
	SlackToken  string
}

func UpdateSlackStatus(emoji string, text string, token string) (err error) {
	var jsonStr = []byte(fmt.Sprintf(`{"profile": {"status_emoji": "%s", "status_text": "%s"}}`, emoji, text))

	req, err := http.NewRequest(
		"POST",
		"https://slack.com/api/users.profile.set",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

func main() {
	var conf Config
	absPath, _ := filepath.Abs("./config.toml")
	_, err := toml.DecodeFile(absPath, &conf)
	if err != nil {
		panic(err)
	}

	doc := GetForecastDoc(conf.ForecastUrl)
	h2, emoji, text := GetForecastEmojiText(doc)

	statusText := fmt.Sprintf("%s: %s, 取得: %s", h2, text, time.Now().Format("15:04"))

	UpdateSlackStatus(emoji, statusText, conf.SlackToken)
}
