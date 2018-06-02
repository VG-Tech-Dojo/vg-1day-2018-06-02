package bot

import (
	"regexp"
	"strings"

	"fmt"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/yumeto/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/yumeto/model"
	"net/url"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	TextsearchAPIURLFormat = "https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s"
	NearbyAPIURLFormat = "https://maps.googleapis.com/maps/api/place/nearbysearch/json?language=ja&location=%f,%f&keyword=%s&rankby=distance&key=%s"
)

type (
	// Processor はmessageを受け取り、投稿用messageを作るインターフェースです
	Processor interface {
		Process(message *model.Message) (*model.Message, error)
	}

	// HelloWorldProcessor は"hello, world!"メッセージを作るprocessorの構造体です
	HelloWorldProcessor struct{}

	// OmikujiProcessor は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかをランダムで作るprocessorの構造体です
	OmikujiProcessor struct{}

	// KeywordProcessor はメッセージ本文からキーワードを抽出するprocessorの構造体です
	KeywordProcessor struct{}

	// GachaProcessor はガチャを行うprocessorの構造体です
	GachaProcessor struct{}

	// JiroProcessor
	JiroProcessor struct{}
)

// Process は"hello, world!"というbodyがセットされたメッセージのポインタを返します
func (p *HelloWorldProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	return &model.Message{
		Body: msgIn.Body + ", world!",
	}, nil
}

// Process は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかがbodyにセットされたメッセージへのポインタを返します
func (p *OmikujiProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	fortunes := []string{
		"大吉",
		"吉",
		"中吉",
		"小吉",
		"末吉",
		"凶",
	}
	result := fortunes[randIntn(len(fortunes))]
	return &model.Message{
		Body: result,
	}, nil
}

// Process はメッセージ本文からキーワードを抽出します
func (p *KeywordProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Akeyword (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	url := fmt.Sprintf(keywordAPIURLFormat, env.KeywordAPIAppID, url.QueryEscape(text))

	type keywordAPIResponse map[string]interface{}
	var json keywordAPIResponse
	get(url, &json)

	keywords := []string{}
	for k, v := range json {
		if k == "Error" {
			return nil, fmt.Errorf("%#v", v)
		}
		keywords = append(keywords, k)
	}

	return &model.Message{
		Body: "キーワード：" + strings.Join(keywords, ", "),
	}, nil
}

func (p *GachaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	fortunes := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	resInt := randIntn(100)
	var resIndex int
	switch {
	case resInt <= 0: // SSレア 1%
		resIndex = 0
	case resInt <= 3: // Sレア 3%
		resIndex = 1
	case resInt <= 9: // レア 6%
		resIndex = 2
	default: // ノーマル 90%
		resIndex = 3
	}
	return &model.Message{
		Body: fortunes[resIndex],
	}, nil
}

func (p *JiroProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Ajiro (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	placeName := matchedStrings[1]

	textSearchAPIURL := fmt.Sprintf(TextsearchAPIURLFormat, url.QueryEscape(placeName), env.PlacesAPIKey)

	type PlacesAPIResponse struct {
		HTMLAttributions []interface{} `json:"html_attributions"`
		NextPageToken    string        `json:"next_page_token"`
		Results          []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Viewport struct {
					Northeast struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"northeast"`
					Southwest struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"southwest"`
				} `json:"viewport"`
			} `json:"geometry"`
			Icon         string `json:"icon"`
			ID           string `json:"id"`
			Name         string `json:"name"`
			OpeningHours struct {
				OpenNow     bool          `json:"open_now"`
				WeekdayText []interface{} `json:"weekday_text"`
			} `json:"opening_hours"`
			Photos []struct {
				Height           int      `json:"height"`
				HTMLAttributions []string `json:"html_attributions"`
				PhotoReference   string   `json:"photo_reference"`
				Width            int      `json:"width"`
			} `json:"photos"`
			PlaceID   string   `json:"place_id"`
			Rating    float64  `json:"rating"`
			Reference string   `json:"reference"`
			Scope     string   `json:"scope"`
			Types     []string `json:"types"`
			Vicinity  string   `json:"vicinity"`
		} `json:"results"`
		Status string `json:"status"`
	}
	var textSearchRes PlacesAPIResponse

	get(textSearchAPIURL, &textSearchRes)

	lat := textSearchRes.Results[0].Geometry.Location.Lat
	lng := textSearchRes.Results[0].Geometry.Location.Lng

	nearbyAPIURL := fmt.Sprintf(NearbyAPIURLFormat, lat, lng, url.QueryEscape("\"ラーメン二郎\""), env.PlacesAPIKey)
	var nearbyRes PlacesAPIResponse

	get(nearbyAPIURL, &nearbyRes)
	isJiro := regexp.MustCompile(".*ラーメン二郎.*")

	branchName := "近くにありません"
	vicinity := "none"
	for _, result := range nearbyRes.Results {
		if isJiro.MatchString(result.Name) {
			branchName = result.Name
			vicinity = result.Vicinity
			break
		}
	}


	return &model.Message{
		Body: fmt.Sprintf("%s (%s)", branchName, vicinity),
	}, nil
}
