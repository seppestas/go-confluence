package confluence

// https://developer.atlassian.com/server/confluence/pagination-in-the-rest-api/

type ResultPagination struct {
	Start int `json:"size"`
	Limit int `json:"limit"`
	Size  int `json:"size"`
}
