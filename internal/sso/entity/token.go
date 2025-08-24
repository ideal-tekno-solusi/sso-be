package entity

type TokenRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"codeVerifier"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RedirectUrl  string `json:"redirect_url"`
}
