package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetForecastDoc(url string) *goquery.Document {
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
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func GetForecastEmojiText(doc *goquery.Document) (string, string) {
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

	selection := doc.Find("section.today-weather")
	iconImg := selection.Find("div.weather-icon img")
	iconUrl, _ := iconImg.Attr("src")
	iconTitle, _ := iconImg.Attr("title")
	splittedIconUrl := strings.Split(iconUrl, "/")
	iconFileName := splittedIconUrl[len(splittedIconUrl) - 1]
	iconName := strings.TrimSuffix(iconFileName, ".png")
	iconName = strings.TrimSuffix(iconName, "_n") // on night

	emoji, ok := weatherSymbolMap[iconName]

	var text string
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

	return emoji, text
}

func GetNoSmokingCount() {

}

func main() {
	doc := GetForecastDoc("https://tenki.jp/forecast/3/16/4410/13109/")
	emoji, text:= GetForecastEmojiText(doc)
	text += fmt.Sprintf(", 取得: %s", time.Now().Format("15:04"))

	fmt.Println(emoji, text)
}