package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ContentAncestor struct {
	ID string `json:"id"`
}

type Content struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type"`
	Status string `json:"status,omitempty"`
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
	Space struct {
		Key string `json:"key"`
	} `json:"space,omitempty"`
	Ancestors []ContentAncestor `json:"ancestors,omitempty"`
}

type ChildrenResults struct {
	GenericResults
	Results []Content `json:"results"`
}

func (w *Wiki) existingContentEndpoint(contentID string) (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/content/" + contentID)
}

func (w *Wiki) newContentEndpoint() (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/content")
}

func (w *Wiki) contentChildrenPagesEndpoint(contentID string) (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/content/" + contentID + "/child/page")
}

func (w *Wiki) DeleteContent(contentID string) error {
	contentEndPoint, err := w.existingContentEndpoint(contentID)
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
	contentEndPoint, err := w.existingContentEndpoint(contentID)
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

func (w *Wiki) UpdateContent(content *Content) (*Content, []byte, error) {
	contentEndPoint, err := w.existingContentEndpoint(content.ID)
	if err != nil {
		return nil, nil, err
	}
	return w.internalCreateOrUpdateContent(content, contentEndPoint, "PUT")
}

func (w *Wiki) CreateContent(content *Content) (*Content, []byte, error) {
	contentEndPoint, err := w.newContentEndpoint()
	if err != nil {
		return nil, nil, err
	}
	return w.internalCreateOrUpdateContent(content, contentEndPoint, "POST")
}

func (w *Wiki) internalCreateOrUpdateContent(content *Content, contentEndPoint *url.URL, method string) (*Content, []byte, error) {
	jsonBody, err := json.Marshal(content)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("sending to %s: %v\n", contentEndPoint.String(), string(jsonBody))

	req, err := http.NewRequest(method, contentEndPoint.String(), bytes.NewReader(jsonBody))
	req.Header.Add("Content-Type", "application/json")

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, res, err
	}

	var newContent Content
	err = json.Unmarshal(res, &newContent)
	if err != nil {
		return nil, res, err
	}

	return &newContent, res, nil
}

func (w *Wiki) GetContentChildrenPages(contentID string, expand []string) (*ChildrenResults, error) {
	contentEndPoint, err := w.contentChildrenPagesEndpoint(contentID)
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

	var content ChildrenResults
	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}
