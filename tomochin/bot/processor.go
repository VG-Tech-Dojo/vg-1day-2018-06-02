package bot

import (
	"regexp"
	"strings"

	"fmt"

	"net/url"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/tomochin/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/tomochin/model"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	talkAPIEndpoint     = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
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

	// GachaProcessor は"SSレア", "Sレア", "レア", "ノーマル"のいずれかをランダムで作るprocessorの構造体です
	GachaProcessor struct{}

	TalkProcessor struct{}
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

func (p *GachaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	gacha := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	result := gacha[randIntn(len(gacha))]
	return &model.Message{
		Body: result,
		User: "gachabot",
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
	query := matchedStrings[1]

	reqBody := make(url.Values)
	reqBody.Set("apikey", env.TalkAPIAppID)
	reqBody.Set("query", query)

	type (
		Result struct {
			Perplexity float32 `json:"perplexity"`
			Reply      string  `json:reply`
		}
		Response struct {
			Status  int      `json:"status"`
			Message string   `json:"message"`
			Results []Result `json:"results"`
		}
	)

	var res Response
	if err := post(talkAPIEndpoint, reqBody, &res); err != nil {
		return nil, err
	}

	if len(res.Results) == 0 {
		return nil, fmt.Errorf("no reply")
	}

	var bestReply Result
	for _, r := range res.Results {
		if r.Perplexity > bestReply.Perplexity {
			bestReply = r
		}
	}

	return &model.Message{
		Body: bestReply.Reply,
		User: "り○なちゃん2号",
	}, nil
}
