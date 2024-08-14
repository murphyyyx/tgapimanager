package tgapimanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// BotAPI allows you to interact with the Telegram Bot API.
type TestBotAPI struct {
	Token  string `json:"token"`
	Debug  bool   `json:"debug"`
	Buffer int    `json:"buffer"`

	Self            User       `json:"-"`
	Client          HTTPClient `json:"-"`
	shutdownChannel chan interface{}

	apiEndpoint string
}

// NewBotAPIWithClient creates a new BotAPI instance
// and allows you to pass a http.Client.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func TestNewBotAPIWithClient(token, apiEndpoint string, client HTTPClient) (*BotAPI, error) {
	bot := &BotAPI{
		Token:           token,
		Client:          client,
		Buffer:          100,
		shutdownChannel: make(chan interface{}),

		apiEndpoint: apiEndpoint,
	}

	self, err := bot.GetMe()
	if err != nil {
		return nil, err
	}

	bot.Self = self

	return bot, nil
}

func (bot *BotAPI) TestMakeRequest(endpoint string, params Params) (*APIResponse, error) {
	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v\n", endpoint, params)
	}

	method := fmt.Sprintf(bot.apiEndpoint, bot.Token, endpoint)

	values := buildParams(params)

	req, err := http.NewRequest("POST", method, strings.NewReader(values.Encode()))
	if err != nil {
		return &APIResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := bot.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	bytes, err := bot.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return &apiResp, err
	}

	if bot.Debug {
		log.Printf("Endpoint: %s, response: %s\n", endpoint, string(bytes))
	}

	if !apiResp.Ok {
		var parameters ResponseParameters

		if apiResp.Parameters != nil {
			parameters = *apiResp.Parameters
		}

		return &apiResp, &Error{
			Code:               apiResp.ErrorCode,
			Message:            apiResp.Description,
			ResponseParameters: parameters,
		}
	}

	return &apiResp, nil
}

func (bot *BotAPI) TestingMakeRequest(endpoint string, params interface{}) (*Response, error) {
	// Эмулируем запрос, здесь можно создать различные ответы для тестирования
	switch endpoint {
	case "getMe":
		user := User{
			ID:        123,
			Username:  "testbot",
			FirstName: "Test",
		}
		userJSON, _ := json.Marshal(user)
		return &Response{Result: userJSON}, nil
	default:
		return nil, errors.New("unsupported endpoint")
	}
}

type Response struct {
	Result []byte
}

func (bot *BotAPI) TestGetMe() (User, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(resp.Result, &user)

	return user, err
}
