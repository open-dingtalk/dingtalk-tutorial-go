package models

type User struct {
	UserID  string
	UnionID string
	Name    string
	Nick    string
	Avatar  string
}

type GetUserAccessRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Code         string `json:"code"`
	RefreshToken string `json:"refreshToken"`
	GrantType    string `json:"grantType"`
}

type GetUserAccessResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireIn     int64  `json:"expireIn"`
	CorpID       string `json:"corpId"`
	RequestID    string `json:"requestid"`
	ErrorCode    string `json:"code"`
	ErrorMessage string `json:"message"`
}

type GetContactUserResponse struct {
	Nick               string             `json:"nick"`
	Avatar             string             `json:"avatarUrl"`
	Mobile             string             `json:"mobile"`
	OpenID             string             `json:"openId"`
	UnionID            string             `json:"unionId"`
	Email              string             `json:"email"`
	StateCode          string             `json:"stateCode"`
	RequestID          string             `json:"requestid"`
	ErrorCode          string             `json:"code"`
	ErrorMessage       string             `json:"message"`
	AccessDeniedDetail AccessDeniedDetail `json:"accessdenieddetail"`
}
