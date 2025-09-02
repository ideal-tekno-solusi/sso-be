package entity

type Response struct {
	RedirectUri string `json:"redirect_uri"`
	Code        string `json:"code"`
	State       string `json:"state"`
}
