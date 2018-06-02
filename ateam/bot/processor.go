package bot

import (
	"regexp"
	"strings"

	"fmt"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/ateam/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-02/ateam/model"
	"net/url"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	talkAPIURL = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
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

	GachaProcessor struct{}

	ChatProcessor struct{}

	QuizProcessor struct{}
)

// Process は"hello, world!"というbodyがセットされたメッセージのポインタを返します
func (p *HelloWorldProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	return &model.Message{
		Body: "正解",
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
	rarity := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	result := rarity[randIntn(len(rarity))]
	return &model.Message{
		Body: result,
		MessageType: 0,
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
		MessageType: 0,
	}, nil
}

func (p *ChatProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Atalk (.*)\\z")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	text := matchedStrings[1]

	request := make(url.Values)
	request.Set("apikey", env.TalkAPIAppID)
	request.Set("query", text)

	type (
		Result struct {
			Perplexity float64 `json:"perplexity"`
			Reply string `json:"reply"`
		}

		Replies struct {
			Status int `json:"status"`
			Message string `json:"message"`
			Results []Result `json:"results"`
		}
	
	)
	var json Replies

	if err := post(talkAPIURL, request, &json); err != nil {
		return nil, err
	}

	if len(json.Results) == 0 {
		return nil, fmt.Errorf("no reply")
	}

	var bestReply Result
	for _, r := range json.Results {
		if r.Perplexity > bestReply.Perplexity {
			bestReply = r
		}
	}

	return &model.Message{
		Body: bestReply.Reply,
		UserName: "ChatBot",
		MessageType: 0,
	}, nil
}

func (p *QuizProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	quizes := []string{
		"「森のバター」と呼ばれている果物は？ A.アボカド B.マンゴー C.梨",
		"板チョコにある溝は何のためにあるのか？ A.苦みを閉じ込めるため B.チョコを速く固めるため C.手で割りやすくするため",
		"鏡の水垢を取る為に有効なものは？ A.クエン酸 B.塩 C.小麦粉",
	}
	quiz := quizes[randIntn(len(quizes))]
	return &model.Message{
		Body: quiz,
		UserName: "QuizBot",
		MessageType: 1,
	}, nil
}
