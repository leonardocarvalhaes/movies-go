package dtos

type ImportMovies struct {
	PivotWords string `json:"pivot_words"`
	OtherWords string `json:"other_words"`
	Rating     string `json:"rating"`
	Votes      string `json:"votes"`
	Year       string `json:"year"`
}
