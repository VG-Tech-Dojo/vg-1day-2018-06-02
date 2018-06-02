package bot

import (
	"regexp"
	"strings"

	"fmt"

	"net/url"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/shira/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/shira/model"
)

const (
	keywordAPIURLFormat  = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	talkAPIURL           = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
	freewordAPIURLFormat = "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/?key=%s&keyword=%s&middle_area=Y030&format=json"
)

type (
	// Processor はmessageを受け取り、投稿用messageを作るインターフェースです
	Processor interface {
		Process(message *model.Message) (*model.Message, error)
	}

	// HelloWorldProcessor は"hello, world!"メッセージを作るprocessorの構造体です
	HelloWorldProcessor struct{}

	GachaProcessor struct{}

	// OmikujiProcessor は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかをランダムで作るprocessorの構造体です
	OmikujiProcessor struct{}

	// KeywordProcessor はメッセージ本文からキーワードを抽出するprocessorの構造体です
	KeywordProcessor struct{}

	TalkProcessor struct{}

	EliteProcessor struct{}

	TimeProcessor struct{}
)

// Process は"hello, world!"というbodyがセットされたメッセージのポインタを返します
func (p *HelloWorldProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	return &model.Message{
		Body: msgIn.Body + ", world!",
	}, nil
}

func (p *GachaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	rares := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	result := rares[randIntn(len(rares))]
	return &model.Message{
		Body: result,
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

func (p *TalkProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Atalk (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	values := url.Values{}
	values.Add("apikey", env.TalkAPIKey)
	values.Add("query", url.QueryEscape(text))

	type talkAPIResponse map[string]interface{}
	var json talkAPIResponse
	post(talkAPIURL, values, &json)

	fmt.Println(json["results"])
	return &model.Message{
		Body: "OK",
	}, nil
}

func (p *EliteProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\A(.*)がたべたい\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	url := fmt.Sprintf(freewordAPIURLFormat, env.HotPapperAPI, url.QueryEscape(text))

	type freewordAPIResponse map[string]interface{}
	var json freewordAPIResponse
	get(url, &json)

	// var name string
	resluts := json["results"]
	name := resluts.(map[string]interface{})["shop"].([]interface{})[0].(map[string]interface{})["name"].(string)
	shopurl := resluts.(map[string]interface{})["shop"].([]interface{})[0].(map[string]interface{})["open"].(string)

	fmt.Println(shopurl)
	fmt.Println(name)

	return &model.Message{
		Body:     name,
		Username: shopurl,
	}, nil
}

func (p *TimeProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\A(.*)食\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	return &model.Message{
		Body:     "朝にはパンがおすすめ",
		Username: text,
	}, nil
}
