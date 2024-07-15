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
