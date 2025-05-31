package entity

type TokenRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"codeVerifier"`
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	Scope        string `json:"scope"`
	TokenType    string `json:"tokenType"`
}
