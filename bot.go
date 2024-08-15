package tgapimanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient is the type needed for the bot to perform HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// BotAPI allows you to interact with the Telegram Bot API.
type BotAPI struct {
	Token  string `json:"token"`
	Debug  bool   `json:"debug"`
	Buffer int    `json:"buffer"`

	Self            User       `json:"-"`
	Client          HTTPClient `json:"-"`
	shutdownChannel chan interface{}

	apiEndpoint string
}

// NewBotAPI creates a new BotAPI instance.
//
// It requires a token, provided by @BotFather on Telegram.
func NewBotAPI(token string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, APIEndpoint, &http.Client{})
}

// NewBotAPIWithAPIEndpoint creates a new BotAPI instance
// and allows you to pass API endpoint.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewBotAPIWithAPIEndpoint(token, apiEndpoint string) (*BotAPI, error) {
	return NewBotAPIWithClient(token, apiEndpoint, &http.Client{})
}

// NewBotAPIWithClient creates a new BotAPI instance
// and allows you to pass a http.Client.
//
// It requires a token, provided by @BotFather on Telegram and API endpoint.
func NewBotAPIWithClient(token, apiEndpoint string, client HTTPClient) (*BotAPI, error) {
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

// SetAPIEndpoint changes the Telegram Bot API endpoint used by the instance.
func (bot *BotAPI) SetAPIEndpoint(apiEndpoint string) {
	bot.apiEndpoint = apiEndpoint
}

func buildParams(in Params) url.Values {
	if in == nil {
		return url.Values{}
	}

	out := url.Values{}

	for key, value := range in {
		out.Set(key, value)
	}

	return out
}

// MakeRequest makes a request to a specific endpoint with our token.
func (bot *BotAPI) MakeRequest(endpoint string, params Params) (*APIResponse, error) {
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

// decodeAPIResponse decode response and return slice of bytes if debug enabled.
// If debug disabled, just decode http.Response.Body stream to APIResponse struct
// for efficient memory usage
func (bot *BotAPI) decodeAPIResponse(responseBody io.Reader, resp *APIResponse) ([]byte, error) {
	if !bot.Debug {
		dec := json.NewDecoder(responseBody)
		err := dec.Decode(resp)
		return nil, err
	}

	// if debug, read response body
	data, err := io.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UploadFiles makes a request to the API with files.
func (bot *BotAPI) UploadFiles(endpoint string, params Params, files []RequestFile) (*APIResponse, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	// This code modified from the very helpful @HirbodBehnam
	// https://github.com/go-telegram-bot-api/telegram-bot-api/issues/354#issuecomment-663856473
	go func() {
		defer w.Close()
		defer m.Close()

		for field, value := range params {
			if err := m.WriteField(field, value); err != nil {
				w.CloseWithError(err)
				return
			}
		}

		for _, file := range files {
			if file.Data.NeedsUpload() {
				name, reader, err := file.Data.UploadData()
				if err != nil {
					w.CloseWithError(err)
					return
				}

				part, err := m.CreateFormFile(file.Name, name)
				if err != nil {
					w.CloseWithError(err)
					return
				}

				if _, err := io.Copy(part, reader); err != nil {
					w.CloseWithError(err)
					return
				}

				if closer, ok := reader.(io.ReadCloser); ok {
					if err = closer.Close(); err != nil {
						w.CloseWithError(err)
						return
					}
				}
			} else {
				value := file.Data.SendData()

				if err := m.WriteField(file.Name, value); err != nil {
					w.CloseWithError(err)
					return
				}
			}
		}
	}()

	if bot.Debug {
		log.Printf("Endpoint: %s, params: %v, with %d files\n", endpoint, params, len(files))
	}

	method := fmt.Sprintf(bot.apiEndpoint, bot.Token, endpoint)

	req, err := http.NewRequest("POST", method, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", m.FormDataContentType())

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
			Message:            apiResp.Description,
			ResponseParameters: parameters,
		}
	}

	return &apiResp, nil
}

// GetMe fetches the currently authenticated bot.
//
// This method is called upon creation to validate the token,
// and so you may get this data from BotAPI.Self without the need for
// another request.
func (bot *BotAPI) GetMe() (User, error) {
	resp, err := bot.MakeRequest("getMe", nil)
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(resp.Result, &user)

	return user, err
}
func hasFilesNeedingUpload(files []RequestFile) bool {
	for _, file := range files {
		if file.Data.NeedsUpload() {
			return true
		}
	}

	return false
}

// Request sends a Chattable to Telegram, and returns the APIResponse.
func (bot *BotAPI) Request(c Chattable) (*APIResponse, error) {
	params, err := c.params()
	if err != nil {
		return nil, err
	}

	if t, ok := c.(Fileable); ok {
		files := t.files()

		// If we have files that need to be uploaded, we should delegate the
		// request to UploadFile.
		if hasFilesNeedingUpload(files) {
			return bot.UploadFiles(t.method(), params, files)
		}

		// However, if there are no files to be uploaded, there's likely things
		// that need to be turned into params instead.
		for _, file := range files {
			params[file.Name] = file.Data.SendData()
		}
	}

	return bot.MakeRequest(c.method(), params)
}

// Send will send a Chattable item to Telegram and provides the
// returned Message.
func (bot *BotAPI) Send(c Chattable) (Message, error) {
	resp, err := bot.Request(c)
	if err != nil {
		return Message{}, err
	}

	var message Message
	err = json.Unmarshal(resp.Result, &message)

	return message, err
}

func (bot *BotAPI) GetUpdates(config UpdateConfig) ([]Update, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	err = json.Unmarshal(resp.Result, &updates)

	return updates, err
}

// GetWebhookInfo allows you to fetch information about a webhook and if
// one currently is set, along with pending update count and error messages.
func (bot *BotAPI) GetWebhookInfo() (WebhookInfo, error) {
	resp, err := bot.MakeRequest("getWebhookInfo", nil)
	if err != nil {
		return WebhookInfo{}, err
	}

	var info WebhookInfo
	err = json.Unmarshal(resp.Result, &info)

	return info, err
}

// GetUpdatesChan starts and returns a channel for getting updates.
func (bot *BotAPI) GetUpdatesChan(config UpdateConfig) UpdatesChannel {
	ch := make(chan Update, bot.Buffer)

	go func() {
		for {
			select {
			case <-bot.shutdownChannel:
				close(ch)
				return
			default:
			}

			updates, err := bot.GetUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()

	return ch
}

// StopReceivingUpdates stops the go routine which receives updates
func (bot *BotAPI) StopReceivingUpdates() {
	if bot.Debug {
		log.Println("Stopping the update receiver routine...")
	}
	close(bot.shutdownChannel)
}

// ListenForWebhook registers a http handler for a webhook.
func (bot *BotAPI) ListenForWebhook(pattern string) UpdatesChannel {
	ch := make(chan Update, bot.Buffer)

	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		update, err := bot.HandleUpdate(r)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(errMsg)
			return
		}

		ch <- *update
	})

	return ch
}

// ListenForWebhookRespReqFormat registers a http handler for a single incoming webhook.
func (bot *BotAPI) ListenForWebhookRespReqFormat(w http.ResponseWriter, r *http.Request) UpdatesChannel {
	ch := make(chan Update, bot.Buffer)

	func(w http.ResponseWriter, r *http.Request) {
		defer close(ch)

		update, err := bot.HandleUpdate(r)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(errMsg)
			return
		}

		ch <- *update
	}(w, r)

	return ch
}

// HandleUpdate parses and returns update received via webhook
func (bot *BotAPI) HandleUpdate(r *http.Request) (*Update, error) {
	if r.Method != http.MethodPost {
		err := errors.New("wrong HTTP method required POST")
		return nil, err
	}

	var update Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		return nil, err
	}

	return &update, nil
}

// WriteToHTTPResponse writes the request to the HTTP ResponseWriter.
//
// It doesn't support uploading files.
func WriteToHTTPResponse(w http.ResponseWriter, c Chattable) error {
	params, err := c.params()
	if err != nil {
		return err
	}

	if t, ok := c.(Fileable); ok {
		if hasFilesNeedingUpload(t.files()) {
			return errors.New("unable to use http response to upload files")
		}
	}

	values := buildParams(params)
	values.Set("method", c.method())

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = w.Write([]byte(values.Encode()))
	return err
}

// GetMyCommands gets the currently registered commands.
func (bot *BotAPI) GetMyCommands() ([]BotCommand, error) {
	return bot.GetMyCommandsWithConfig(GetMyCommandsConfig{})
}

// GetMyCommandsWithConfig gets the currently registered commands with a config.
func (bot *BotAPI) GetMyCommandsWithConfig(config GetMyCommandsConfig) ([]BotCommand, error) {
	resp, err := bot.Request(config)
	if err != nil {
		return nil, err
	}

	var commands []BotCommand
	err = json.Unmarshal(resp.Result, &commands)

	return commands, err
}
