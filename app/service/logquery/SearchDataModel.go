package logquery

// SearchData : query condition
type SearchData struct {
	StartDate int    `json:"startdate" validate:"gte=1"`
	EndDate   int    `json:"enddate" validate:"gte=1"`
	QueryText string `json:"querytext"`
}
