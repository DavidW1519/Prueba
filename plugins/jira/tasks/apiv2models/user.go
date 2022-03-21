package apiv2models

type User struct {
	Self         string `json:"self"`
	Key          string `json:"key"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	AccountId    string `json:"accountId"`
	AvatarUrls   struct {
		Four8X48  string `json:"48x48"`
		Two4X24   string `json:"24x24"`
		One6X16   string `json:"16x16"`
		Three2X32 string `json:"32x32"`
	} `json:"avatarUrls"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	Deleted     bool   `json:"deleted"`
	TimeZone    string `json:"timeZone"`
	Locale      string `json:"locale"`
}

func (u *User) getAccountId() string {
	if u == nil {
		return ""
	}
	if u.AccountId != "" {
		return u.AccountId
	}
	return u.EmailAddress
}
