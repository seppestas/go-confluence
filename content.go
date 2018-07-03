package confluence

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Content struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Body   struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
	} `json:"body"`
	Version struct {
		Number int `json:"number"`
	} `json:"version"`
}

func (w *Wiki) contentEndpoint(contentID string) (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/content/" + contentID)
}

func (w *Wiki) DeleteContent(contentID string) error {
	contentEndPoint, err := w.contentEndpoint(contentID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", contentEndPoint.String(), nil)
	if err != nil {
		return err
	}

	_, err = w.sendRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (w *Wiki) GetContent(contentID string, expand []string) (*Content, error) {
	contentEndPoint, err := w.contentEndpoint(contentID)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("expand", strings.Join(expand, ","))
	contentEndPoint.RawQuery = data.Encode()

	req, err := http.NewRequest("GET", contentEndPoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var content Content
	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (w *Wiki) UpdateContent(content *Content) (*Content, error) {
	jsonbody, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	contentEndPoint, err := w.contentEndpoint(content.Id)
	req, err := http.NewRequest("PUT", contentEndPoint.String(), strings.NewReader(string(jsonbody)))
	req.Header.Add("Content-Type", "application/json")

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	var newContent Content
	err = json.Unmarshal(res, &newContent)
	if err != nil {
		return nil, err
	}

	return &newContent, nil
}

func (w *Wiki) GetPageIDByTitle(title, space string) (string, error) {
	type GetPageID struct {
		Results []struct {
			ID    string `json:"id,omitempty"`
			Title string `json:"title,omitempty"`
		} `json:"results,omitempty"`
	}
	contentEndPoint, err := w.contentEndpoint("")
	if err != nil {
		log.Println(err)
	}
	data := url.Values{}
	data.Set("expand", "history")
	data.Set("title", title)
	data.Set("spacekey", "space")
	contentEndPoint.RawQuery = data.Encode()

	req, err := http.NewRequest("GET", contentEndPoint.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := w.sendRequest(req)
	if err != nil {
		return "", err
	}

	gpid := GetPageID{}
	err = json.Unmarshal(res, &gpid)
	if err != nil {
		return "", err
	}
	var pageID, pageTitle string
	if len(gpid.Results) == 1 {
		for _, v := range gpid.Results {
			pageID = v.ID
			pageTitle = v.Title
		}
	}
	// In case we need to return the title at some point
	_ = pageTitle
	return pageID, nil
}
