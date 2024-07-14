package tgapilib

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int    `json:"update_id"`
	Message  string `json:"message"`
}
