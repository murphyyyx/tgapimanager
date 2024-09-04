package tgapimanager

import (
	"io"
	"net/url"
)

const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

// BaseChat is base type for all chat config types.
type BaseChat struct {
	ChatID                   int64 // required
	ChannelUsername          string
	ReplyToMessageID         int
	ReplyMarkup              interface{}
	DisableNotification      bool
	AllowSendingWithoutReply bool
}
type MessageConfig struct {
	BaseChat
	Text                  string
	ParseMode             string
	Entities              []MessageEntity
	DisableWebPagePreview bool
}

func (chat *BaseChat) params() (Params, error) {
	params := make(Params)

	params.AddFirstValid("chat_id", chat.ChatID, chat.ChannelUsername)
	params.AddNonZero("reply_to_message_id", chat.ReplyToMessageID)
	params.AddBool("disable_notification", chat.DisableNotification)
	params.AddBool("allow_sending_without_reply", chat.AllowSendingWithoutReply)

	err := params.AddInterface("reply_markup", chat.ReplyMarkup)

	return params, err
}

func (config MessageConfig) params() (Params, error) {
	params, err := config.BaseChat.params()
	if err != nil {
		return params, err
	}

	params.AddNonEmpty("text", config.Text)
	params.AddBool("disable_web_page_preview", config.DisableWebPagePreview)
	params.AddNonEmpty("parse_mode", config.ParseMode)
	err = params.AddInterface("entities", config.Entities)

	return params, err
}

func (config MessageConfig) method() string {
	return "sendMessage"
}

// Chattable is any config type that can be sent.
type Chattable interface {
	params() (Params, error)
	method() string
}

// Fileable is any config type that can be sent that includes a file.
type Fileable interface {
	Chattable
	files() []RequestFile
}

// RequestFile represents a file associated with a field name.
type RequestFile struct {
	// The file field name.
	Name string
	// The file data to include.
	Data RequestFileData
}

// RequestFileData represents the data to be used for a file.
type RequestFileData interface {
	// NeedsUpload shows if the file needs to be uploaded.
	NeedsUpload() bool

	// UploadData gets the file name and an `io.Reader` for the file to be uploaded. This
	// must only be called when the file needs to be uploaded.
	UploadData() (string, io.Reader, error)
	// SendData gets the file data to send when a file does not need to be uploaded. This
	// must only be called when the file does not need to be uploaded.
	SendData() string
}

// UpdateConfig contains information about a GetUpdates request.
type UpdateConfig struct {
	Offset         int
	Limit          int
	Timeout        int
	AllowedUpdates []string
}

func (UpdateConfig) method() string {
	return "getUpdates"
}

func (config UpdateConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonZero("offset", config.Offset)
	params.AddNonZero("limit", config.Limit)
	params.AddNonZero("timeout", config.Timeout)
	params.AddInterface("allowed_updates", config.AllowedUpdates)

	return params, nil
}

// SetMyCommandsConfig sets a list of commands the bot understands.
type SetMyCommandsConfig struct {
	Commands     []BotCommand
	Scope        *BotCommandScope
	LanguageCode string
}

func (config SetMyCommandsConfig) method() string {
	return "setMyCommands"
}

func (config SetMyCommandsConfig) params() (Params, error) {
	params := make(Params)

	if err := params.AddInterface("commands", config.Commands); err != nil {
		return params, err
	}
	err := params.AddInterface("scope", config.Scope)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, err
}

type DeleteMyCommandsConfig struct {
	Scope        *BotCommandScope
	LanguageCode string
}

func (config DeleteMyCommandsConfig) method() string {
	return "deleteMyCommands"
}

func (config DeleteMyCommandsConfig) params() (Params, error) {
	params := make(Params)

	err := params.AddInterface("scope", config.Scope)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, err
}

// GetMyCommandsConfig gets a list of the currently registered commands.
type GetMyCommandsConfig struct {
	Scope        *BotCommandScope
	LanguageCode string
}

func (config GetMyCommandsConfig) method() string {
	return "getMyCommands"
}

func (config GetMyCommandsConfig) params() (Params, error) {
	params := make(Params)

	err := params.AddInterface("scope", config.Scope)
	params.AddNonEmpty("language_code", config.LanguageCode)

	return params, err
}

// BaseEdit is base type of all chat edits.
type BaseEdit struct {
	ChatID          int64
	ChannelUsername string
	MessageID       int
	InlineMessageID string
	ReplyMarkup     *InlineKeyboardMarkup
}

func (edit BaseEdit) params() (Params, error) {
	params := make(Params)

	if edit.InlineMessageID != "" {
		params["inline_message_id"] = edit.InlineMessageID
	} else {
		params.AddFirstValid("chat_id", edit.ChatID, edit.ChannelUsername)
		params.AddNonZero("message_id", edit.MessageID)
	}

	err := params.AddInterface("reply_markup", edit.ReplyMarkup)

	return params, err
}

// StopPollConfig allows you to stop a poll sent by the bot.
type StopPollConfig struct {
	BaseEdit
}

func (config StopPollConfig) params() (Params, error) {
	return config.BaseEdit.params()
}

func (StopPollConfig) method() string {
	return "stopPoll"
}

// LocationConfig contains information about a SendLocation request.
type LocationConfig struct {
	BaseChat
	Latitude             float64 // required
	Longitude            float64 // required
	HorizontalAccuracy   float64 // optional
	LivePeriod           int     // optional
	Heading              int     // optional
	ProximityAlertRadius int     // optional
}

func (config LocationConfig) params() (Params, error) {
	params, err := config.BaseChat.params()

	params.AddNonZeroFloat("latitude", config.Latitude)
	params.AddNonZeroFloat("longitude", config.Longitude)
	params.AddNonZeroFloat("horizontal_accuracy", config.HorizontalAccuracy)
	params.AddNonZero("live_period", config.LivePeriod)
	params.AddNonZero("heading", config.Heading)
	params.AddNonZero("proximity_alert_radius", config.ProximityAlertRadius)

	return params, err
}

func (config LocationConfig) method() string {
	return "sendLocation"
}

// EditMessageLiveLocationConfig allows you to update a live location.
type EditMessageLiveLocationConfig struct {
	BaseEdit
	Latitude             float64 // required
	Longitude            float64 // required
	HorizontalAccuracy   float64 // optional
	Heading              int     // optional
	ProximityAlertRadius int     // optional
}

func (config EditMessageLiveLocationConfig) params() (Params, error) {
	params, err := config.BaseEdit.params()

	params.AddNonZeroFloat("latitude", config.Latitude)
	params.AddNonZeroFloat("longitude", config.Longitude)
	params.AddNonZeroFloat("horizontal_accuracy", config.HorizontalAccuracy)
	params.AddNonZero("heading", config.Heading)
	params.AddNonZero("proximity_alert_radius", config.ProximityAlertRadius)

	return params, err
}

func (config EditMessageLiveLocationConfig) method() string {
	return "editMessageLiveLocation"
}

// StopMessageLiveLocationConfig stops updating a live location.
type StopMessageLiveLocationConfig struct {
	BaseEdit
}

func (config StopMessageLiveLocationConfig) params() (Params, error) {
	return config.BaseEdit.params()
}

func (config StopMessageLiveLocationConfig) method() string {
	return "stopMessageLiveLocation"
}

// VenueConfig contains information about a SendVenue request.
type VenueConfig struct {
	BaseChat
	Latitude        float64 // required
	Longitude       float64 // required
	Title           string  // required
	Address         string  // required
	FoursquareID    string
	FoursquareType  string
	GooglePlaceID   string
	GooglePlaceType string
}

func (config VenueConfig) params() (Params, error) {
	params, err := config.BaseChat.params()

	params.AddNonZeroFloat("latitude", config.Latitude)
	params.AddNonZeroFloat("longitude", config.Longitude)
	params["title"] = config.Title
	params["address"] = config.Address
	params.AddNonEmpty("foursquare_id", config.FoursquareID)
	params.AddNonEmpty("foursquare_type", config.FoursquareType)
	params.AddNonEmpty("google_place_id", config.GooglePlaceID)
	params.AddNonEmpty("google_place_type", config.GooglePlaceType)

	return params, err
}

func (config VenueConfig) method() string {
	return "sendVenue"
}

// ContactConfig allows you to send a contact.
type ContactConfig struct {
	BaseChat
	PhoneNumber string
	FirstName   string
	LastName    string
	VCard       string
}

// WebhookConfig contains information about a SetWebhook request.
type WebhookConfig struct {
	URL                *url.URL
	Certificate        RequestFileData
	IPAddress          string
	MaxConnections     int
	AllowedUpdates     []string
	DropPendingUpdates bool
}

func (config WebhookConfig) method() string {
	return "setWebhook"
}

func (config WebhookConfig) params() (Params, error) {
	params := make(Params)

	if config.URL != nil {
		params["url"] = config.URL.String()
	}

	params.AddNonEmpty("ip_address", config.IPAddress)
	params.AddNonZero("max_connections", config.MaxConnections)
	err := params.AddInterface("allowed_updates", config.AllowedUpdates)
	params.AddBool("drop_pending_updates", config.DropPendingUpdates)

	return params, err
}

func (config WebhookConfig) files() []RequestFile {
	if config.Certificate != nil {
		return []RequestFile{{
			Name: "certificate",
			Data: config.Certificate,
		}}
	}

	return nil
}

// DeleteWebhookConfig is a helper to delete a webhook.
type DeleteWebhookConfig struct {
	DropPendingUpdates bool
}

func (config DeleteWebhookConfig) method() string {
	return "deleteWebhook"
}

func (config DeleteWebhookConfig) params() (Params, error) {
	params := make(Params)

	params.AddBool("drop_pending_updates", config.DropPendingUpdates)

	return params, nil
}

// EditMessageTextConfig allows you to modify the text in a message.
type EditMessageTextConfig struct {
	BaseEdit
	Text                  string
	ParseMode             string
	Entities              []MessageEntity
	DisableWebPagePreview bool
}

func (config EditMessageTextConfig) params() (Params, error) {
	params, err := config.BaseEdit.params()
	if err != nil {
		return params, err
	}

	params["text"] = config.Text
	params.AddNonEmpty("parse_mode", config.ParseMode)
	params.AddBool("disable_web_page_preview", config.DisableWebPagePreview)
	err = params.AddInterface("entities", config.Entities)

	return params, err
}

func (config EditMessageTextConfig) method() string {
	return "editMessageText"
}

// EditMessageCaptionConfig allows you to modify the caption of a message.
type EditMessageCaptionConfig struct {
	BaseEdit
	Caption         string
	ParseMode       string
	CaptionEntities []MessageEntity
}

func (config EditMessageCaptionConfig) params() (Params, error) {
	params, err := config.BaseEdit.params()
	if err != nil {
		return params, err
	}

	params["caption"] = config.Caption
	params.AddNonEmpty("parse_mode", config.ParseMode)
	err = params.AddInterface("caption_entities", config.CaptionEntities)

	return params, err
}

func (config EditMessageCaptionConfig) method() string {
	return "editMessageCaption"
}

// EditMessageMediaConfig allows you to make an editMessageMedia request.
type EditMessageMediaConfig struct {
	BaseEdit

	Media interface{}
}

func (EditMessageMediaConfig) method() string {
	return "editMessageMedia"
}

// EditMessageReplyMarkupConfig allows you to modify the reply markup
// of a message.
type EditMessageReplyMarkupConfig struct {
	BaseEdit
}

func (config EditMessageReplyMarkupConfig) params() (Params, error) {
	return config.BaseEdit.params()
}

func (config EditMessageReplyMarkupConfig) method() string {
	return "editMessageReplyMarkup"
}
