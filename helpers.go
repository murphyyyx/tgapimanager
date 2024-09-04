package tgapimanager

import (
	"net/url"
)

// NewMessage creates a new Message.
//
// chatID is where to send it, text is the message text.
func NewMessage(chatID int64, text string) MessageConfig {
	return MessageConfig{
		BaseChat: BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text:                  text,
		DisableWebPagePreview: false,
	}
}

// NewLocation shares your location.
//
// chatID is where to send it, latitude and longitude are coordinates.
func NewLocation(chatID int64, latitude float64, longitude float64) LocationConfig {
	return LocationConfig{
		BaseChat: BaseChat{
			ChatID: chatID,
		},
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewVenue allows you to send a venue and its location.
func NewVenue(chatID int64, title, address string, latitude, longitude float64) VenueConfig {
	return VenueConfig{
		BaseChat: BaseChat{
			ChatID: chatID,
		},
		Title:     title,
		Address:   address,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// NewUpdate gets updates since the last Offset.
//
// offset is the last Update ID to include.
// You likely want to set this to the last Update ID plus 1.
func NewUpdate(offset int) UpdateConfig {
	return UpdateConfig{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
}

// NewWebhook creates a new webhook.
//
// link is the url parsable link you wish to get the updates.
func NewWebhook(link string) (WebhookConfig, error) {
	u, err := url.Parse(link)

	if err != nil {
		return WebhookConfig{}, err
	}

	return WebhookConfig{
		URL: u,
	}, nil
}

// NewWebhookWithCert creates a new webhook with a certificate.
//
// link is the url you wish to get webhooks,
// file contains a string to a file, FileReader, or FileBytes.
func NewWebhookWithCert(link string, file RequestFileData) (WebhookConfig, error) {
	u, err := url.Parse(link)

	if err != nil {
		return WebhookConfig{}, err
	}

	return WebhookConfig{
		URL:         u,
		Certificate: file,
	}, nil
}

// NewEditMessageText allows you to edit the text of a message.
func NewEditMessageText(chatID int64, messageID int, text string) EditMessageTextConfig {
	return EditMessageTextConfig{
		BaseEdit: BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text: text,
	}
}

// NewEditMessageTextAndMarkup allows you to edit the text and replymarkup of a message.
func NewEditMessageTextAndMarkup(chatID int64, messageID int, text string, replyMarkup InlineKeyboardMarkup) EditMessageTextConfig {
	return EditMessageTextConfig{
		BaseEdit: BaseEdit{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: &replyMarkup,
		},
		Text: text,
	}
}

// NewEditMessageCaption allows you to edit the caption of a message.
func NewEditMessageCaption(chatID int64, messageID int, caption string) EditMessageCaptionConfig {
	return EditMessageCaptionConfig{
		BaseEdit: BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Caption: caption,
	}
}

// NewEditMessageReplyMarkup allows you to edit the inline
// keyboard markup.
func NewEditMessageReplyMarkup(chatID int64, messageID int, replyMarkup InlineKeyboardMarkup) EditMessageReplyMarkupConfig {
	return EditMessageReplyMarkupConfig{
		BaseEdit: BaseEdit{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: &replyMarkup,
		},
	}
}

// NewRemoveKeyboard hides the keyboard, with the option for being selective
// or hiding for everyone.
func NewRemoveKeyboard(selective bool) ReplyKeyboardRemove {
	return ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      selective,
	}
}

// NewKeyboardButton creates a regular keyboard button.
func NewKeyboardButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
	}
}

// NewKeyboardButtonContact creates a keyboard button that requests
// user contact information upon click.
func NewKeyboardButtonContact(text string) KeyboardButton {
	return KeyboardButton{
		Text:           text,
		RequestContact: true,
	}
}

// NewKeyboardButtonLocation creates a keyboard button that requests
// user location information upon click.
func NewKeyboardButtonLocation(text string) KeyboardButton {
	return KeyboardButton{
		Text:            text,
		RequestLocation: true,
	}
}

// NewKeyboardButtonRow creates a row of keyboard buttons.
func NewKeyboardButtonRow(buttons ...KeyboardButton) []KeyboardButton {
	var row []KeyboardButton

	row = append(row, buttons...)

	return row
}

// NewReplyKeyboard creates a new regular keyboard with sane defaults.
func NewReplyKeyboard(rows ...[]KeyboardButton) ReplyKeyboardMarkup {
	var keyboard [][]KeyboardButton

	keyboard = append(keyboard, rows...)

	return ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
	}
}

// NewOneTimeReplyKeyboard creates a new one time keyboard.
func NewOneTimeReplyKeyboard(rows ...[]KeyboardButton) ReplyKeyboardMarkup {
	markup := NewReplyKeyboard(rows...)
	markup.OneTimeKeyboard = true
	return markup
}

// NewInlineKeyboardButtonData creates an inline keyboard button with text
// and data for a callback.
func NewInlineKeyboardButtonData(text, data string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:         text,
		CallbackData: &data,
	}
}

// NewInlineKeyboardButtonLoginURL creates an inline keyboard button with text
// which goes to a LoginURL.
func NewInlineKeyboardButtonLoginURL(text string, loginURL LoginURL) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:     text,
		LoginURL: &loginURL,
	}
}

// NewInlineKeyboardButtonURL creates an inline keyboard button with text
// which goes to a URL.
func NewInlineKeyboardButtonURL(text, url string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		URL:  &url,
	}
}

// NewInlineKeyboardButtonSwitch creates an inline keyboard button with
// text which allows the user to switch to a chat or return to a chat.
func NewInlineKeyboardButtonSwitch(text, sw string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:              text,
		SwitchInlineQuery: &sw,
	}
}

// NewInlineKeyboardRow creates an inline keyboard row with buttons.
func NewInlineKeyboardRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	var row []InlineKeyboardButton

	row = append(row, buttons...)

	return row
}

// NewInlineKeyboardMarkup creates a new inline keyboard.
func NewInlineKeyboardMarkup(rows ...[]InlineKeyboardButton) InlineKeyboardMarkup {
	var keyboard [][]InlineKeyboardButton

	keyboard = append(keyboard, rows...)

	return InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

// NewBotCommandScopeDefault represents the default scope of bot commands.
func NewBotCommandScopeDefault() BotCommandScope {
	return BotCommandScope{Type: "default"}
}

// NewBotCommandScopeAllPrivateChats represents the scope of bot commands,
// covering all private chats.
func NewBotCommandScopeAllPrivateChats() BotCommandScope {
	return BotCommandScope{Type: "all_private_chats"}
}

// NewBotCommandScopeAllGroupChats represents the scope of bot commands,
// covering all group and supergroup chats.
func NewBotCommandScopeAllGroupChats() BotCommandScope {
	return BotCommandScope{Type: "all_group_chats"}
}

// NewBotCommandScopeAllChatAdministrators represents the scope of bot commands,
// covering all group and supergroup chat administrators.
func NewBotCommandScopeAllChatAdministrators() BotCommandScope {
	return BotCommandScope{Type: "all_chat_administrators"}
}

// NewBotCommandScopeChat represents the scope of bot commands, covering a
// specific chat.
func NewBotCommandScopeChat(chatID int64) BotCommandScope {
	return BotCommandScope{
		Type:   "chat",
		ChatID: chatID,
	}
}

// NewBotCommandScopeChatAdministrators represents the scope of bot commands,
// covering all administrators of a specific group or supergroup chat.
func NewBotCommandScopeChatAdministrators(chatID int64) BotCommandScope {
	return BotCommandScope{
		Type:   "chat_administrators",
		ChatID: chatID,
	}
}

// NewBotCommandScopeChatMember represents the scope of bot commands, covering a
// specific member of a group or supergroup chat.
func NewBotCommandScopeChatMember(chatID, userID int64) BotCommandScope {
	return BotCommandScope{
		Type:   "chat_member",
		ChatID: chatID,
		UserID: userID,
	}
}

// NewGetMyCommandsWithScope allows you to set the registered commands for a
// given scope.
func NewGetMyCommandsWithScope(scope BotCommandScope) GetMyCommandsConfig {
	return GetMyCommandsConfig{Scope: &scope}
}

// NewGetMyCommandsWithScopeAndLanguage allows you to set the registered
// commands for a given scope and language code.
func NewGetMyCommandsWithScopeAndLanguage(scope BotCommandScope, languageCode string) GetMyCommandsConfig {
	return GetMyCommandsConfig{Scope: &scope, LanguageCode: languageCode}
}

// NewSetMyCommands allows you to set the registered commands.
func NewSetMyCommands(commands ...BotCommand) SetMyCommandsConfig {
	return SetMyCommandsConfig{Commands: commands}
}

// NewSetMyCommandsWithScope allows you to set the registered commands for a given scope.
func NewSetMyCommandsWithScope(scope BotCommandScope, commands ...BotCommand) SetMyCommandsConfig {
	return SetMyCommandsConfig{Commands: commands, Scope: &scope}
}

// NewSetMyCommandsWithScopeAndLanguage allows you to set the registered commands for a given scope
// and language code.
func NewSetMyCommandsWithScopeAndLanguage(scope BotCommandScope, languageCode string, commands ...BotCommand) SetMyCommandsConfig {
	return SetMyCommandsConfig{Commands: commands, Scope: &scope, LanguageCode: languageCode}
}

// NewDeleteMyCommands allows you to delete the registered commands.
func NewDeleteMyCommands() DeleteMyCommandsConfig {
	return DeleteMyCommandsConfig{}
}

// NewDeleteMyCommandsWithScope allows you to delete the registered commands for a given
// scope.
func NewDeleteMyCommandsWithScope(scope BotCommandScope) DeleteMyCommandsConfig {
	return DeleteMyCommandsConfig{Scope: &scope}
}

// NewDeleteMyCommandsWithScopeAndLanguage allows you to delete the registered commands for a given
// scope and language code.
func NewDeleteMyCommandsWithScopeAndLanguage(scope BotCommandScope, languageCode string) DeleteMyCommandsConfig {
	return DeleteMyCommandsConfig{Scope: &scope, LanguageCode: languageCode}
}
