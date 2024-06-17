package types

type WebAppUser struct {
	Id           int    `form:"id"`
	FirstName    string `form:"first_name"`
	LastName     string `form:"last_name"`
	Username     string `form:"username"`
	LanguageCode string `form:"language_code"`
	PhotoUrl     string `form:"photo_url"`
}

type WebAppChat struct {
	Id       int    `form:"id"`
	Type     string `form:"type"`
	Title    string `form:"title"`
	Username string `form:"username"`
	PhotoUrl string `form:"photo_url"`
}

type InitData struct {
	QueryId      string     `form:"query_id" json:"query_id"`
	User         WebAppUser `form:"user" json:"user"`
	Receiver     WebAppUser `form:"receiver" json:"receiver"`
	Chat         WebAppChat `form:"chat" json:"chat"`
	ChatType     string     `form:"chat_type" json:"chat_type"`
	ChatInstance string     `form:"chat_instance" json:"chat_instance"`
	StartParam   string     `form:"start_param" json:"start_param"`
	CanSendAfter int        `form:"can_send_after" json:"can_send_after"`
	AuthDate     int        `form:"auth_date" json:"auth_date"`
	Hash         string     `form:"hash" json:"hash"`
}
