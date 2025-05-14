package entity

type LoginEncrypt struct {
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
	AuthorizationCode   string `json:"authorization_code"`
}
