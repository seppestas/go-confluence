package confluence

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Contents struct {
	Page struct {
		Results []struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Status     string `json:"status"`
			Title      string `json:"title"`
			Extensions struct {
				Position interface{} `json:"position"`
			} `json:"extensions"`
			Links struct {
				Webui  string `json:"webui"`
				Edit   string `json:"edit"`
				Tinyui string `json:"tinyui"`
				Self   string `json:"self"`
			} `json:"_links"`
			Expandable struct {
				Container    string `json:"container"`
				Metadata     string `json:"metadata"`
				Operations   string `json:"operations"`
				Children     string `json:"children"`
				Restrictions string `json:"restrictions"`
				History      string `json:"history"`
				Ancestors    string `json:"ancestors"`
				Body         string `json:"body"`
				Version      string `json:"version"`
				Descendants  string `json:"descendants"`
				Space        string `json:"space"`
			} `json:"_expandable"`
		} `json:"results"`
		Start int `json:"start"`
		Limit int `json:"limit"`
		Size  int `json:"size"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"page"`
	Blogpost struct {
		Results []interface{} `json:"results"`
		Start   int           `json:"start"`
		Limit   int           `json:"limit"`
		Size    int           `json:"size"`
		Links   struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"blogpost"`
	Links struct {
		Base    string `json:"base"`
		Context string `json:"context"`
	} `json:"_links"`
}

type Space struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Links struct {
		Webui string `json:"webui"`
		Self  string `json:"self"`
	} `json:"_links"`
	Expandable struct {
		Metadata    string `json:"metadata"`
		Icon        string `json:"icon"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"_expandable"`
}

type Spaces struct {
	Results []Space `json:"results"`
	Start   int     `json:"start"`
	Limit   int     `json:"limit"`
	Size    int     `json:"size"`
	Links   struct {
		Self    string `json:"self"`
		Base    string `json:"base"`
		Context string `json:"context"`
	} `json:"_links"`
}

func (w *Wiki) spaceEndpoint() (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/space")
}

func (w *Wiki) spaceContentEndpoint(space string) (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/space/" + space + "/content")
}

func (w *Wiki) GetSpaces() ([]Space, error) {
	spaceEndPoint, err := w.spaceEndpoint()
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	spaceEndPoint.RawQuery = data.Encode()

	req, err := http.NewRequest("GET", spaceEndPoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var spaces = new(Spaces)
	err = json.Unmarshal(res, &spaces)
	if err != nil {
		return nil, err
	}

	return spaces.Results, nil
}

func (w *Wiki) GetSpaceContent(space string) (*Contents, error) {
	spaceContentEndpoint, err := w.spaceContentEndpoint(space)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	spaceContentEndpoint.RawQuery = data.Encode()

	req, err := http.NewRequest("GET", spaceContentEndpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var spacecontents = new(Contents)
	err = json.Unmarshal(res, &spacecontents)
	if err != nil {
		return nil, err
	}

	return spacecontents, nil
}
