package bot

import (
	"regexp"
	"strings"

	"fmt"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/ishiguro/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/ishiguro/model"
	"net/url"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	talkAPIURLFormat = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk?appid=%s&query=%s"
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

	GachaProcesser struct{}

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

func (p *GachaProcesser) Process(msgIn *model.Message) (*model.Message, error) {
	reality := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	result := reality[randIntn(len(reality))]
	return &model.Message{
		Body: result,
		Username: "Bot",
	}, nil
}

func (p *TalkProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Atalk (.+)\\z")
	regMatches := r.FindStringSubmatch(msgIn.Body)
	if len(regMatches) != 2 {
		return nil, fmt.Errorf("bad message: '%s'", msgIn.Body)
	}
	matchedString := regMatches[1]

	reqBody := make(url.Values)
	reqBody.Set("apikey", env.TalkAPIAppID)
	reqBody.Set("query", matchedString)

	type (
		Result struct {
			Perplexity float32 `json:"perplexity"`
			Reply      string  `json:"reply"`
		}

		Response struct {
			Status  int      `json:"status"`
			Message string   `json:"message"`
			Results []Result `json:"results"`
		}
	)

	var res Response
	if err := post(talkAPIURLFormat, reqBody, &res); err != nil {
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
		Body:     bestReply.Reply,
		Username: "talkbot",
	}, nil
}
