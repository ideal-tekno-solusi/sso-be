package entity

type LoginEncrypt struct {
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
	AuthorizationCode   string `json:"authorizationCode"`
}
