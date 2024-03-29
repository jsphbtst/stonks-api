package types

type Companies struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	About    string `json:"about"`
	Sector   string `json:"sector"`
	Industry string `json:"industry"`
	Mission  string `json:"mission"`
	Vision   string `json:"vision"`
}
