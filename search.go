package confluence

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type ContentResult struct {
	Content               Content `json:"content"`
	Title                 string  `json:"title"`
	Excerpt               string  `json:"excerpt"`
	URL                   string  `json:"url"`
	ResultGlobalContainer struct {
		Title      string `json:"title"`
		DisplayURL string `json:"displayUrl"`
	}
	// breadcrumbs is ignored
	EntityType           string `json:"entityType"`
	IconCSSClass         string `json:"iconCssClass"`
	LastModified         string `json:"lastModified"`
	FriendlyLastModified string `json:"friendlyLastModified"`
}

type GenericResults struct {
	Start int `json:"size"`
	Limit int `json:"limit"`
	Size  int `json:"size"`
}

type SearchResults struct {
	GenericResults
	Results        []ContentResult `json:"results"`
	TotalSize      int             `json:"totalSize"`
	CqlQuery       string          `json:"cqlQuery"`
	SearchDuration int             `json:"SearchDuration"`
	// links are ignored
}

func (w *Wiki) searchEndpoint() (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/search")
}

func (w *Wiki) Search(cql, cqlContext string, expand []string, limit int) (*SearchResults, error) {
	searchEndPoint, err := w.searchEndpoint()
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("expand", strings.Join(expand, ","))
	data.Set("cqlcontext", cqlContext)
	data.Set("cql", cql)
	searchEndPoint.RawQuery = data.Encode()

	req, err := http.NewRequest("GET", searchEndPoint.String(), nil)
	if err != nil {
		return nil, err

	}
	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var results SearchResults
	err = json.Unmarshal(res, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}
