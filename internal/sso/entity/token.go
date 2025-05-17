package entity

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	Scope        string `json:"scope"`
	TokenType    string `json:"tokenType"`
}
