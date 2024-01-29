package dtos

type SearchColumn struct {
	Column   string `json:"column"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}
