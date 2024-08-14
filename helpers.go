package tgapimanager

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

func NewUpdate(offset int) UpdateConfig {
	return UpdateConfig{
		Offset:  offset,
		Limit:   0,
		Timeout: 0,
	}
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

// NewKeyboardButton creates a regular keyboard button.
func NewKeyboardButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
	}
}

// NewKeyboardButtonRow creates a row of keyboard buttons.
func NewKeyboardButtonRow(buttons ...KeyboardButton) []KeyboardButton {
	var row []KeyboardButton

	row = append(row, buttons...)

	return row
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
