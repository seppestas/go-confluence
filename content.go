package confluence

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"github.com/op/go-logging"
	"time"
	"github.com/google/go-querystring/query"
	"path"
	"bytes"
)

var log = logging.MustGetLogger("confluence")

type Content struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Body   Body   `json:"body"`
	Version Version `json:"version"`
}

type Ancestor struct{
	Id int `json:"id"`
}


type Body  struct {
	View View `json:"view"`
}

type View struct {
	Value string `json:"value"`
	Representation string `json:"representation"`
}

type Version struct {
	Number int `json:"number"`
}

type StorageBody  struct {
	View View `json:"storage"`
}

type UpdateContentRequest struct {
	Id     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	Body   StorageBody `json:"body"`
	Version Version`json:"version"`
	Ancestors []Ancestor `json:"ancestors,omitempty"`
	Space map[string]string `json:"space,omitempty"`
}

type ContentResult struct{
	Results []Content `json:"results"`
	Start int `json:"start"`
	Limit int `json:"limit"`
	Size int `json:"size"`
	Self map[string]string `json:_links`
}


func (w *Wiki) contentEndpoint(contentID string) (*url.URL, error) {
	return url.ParseRequestURI(w.endPoint.String() + "/content/" + contentID)
}

func (w *Wiki) contentAPIEndpoint()(*url.URL, error){
	endpoint:=*w.endPoint
	endpoint.Path=path.Join(endpoint.Path,"content")
	return &endpoint,nil
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
	log.Debugf("Calling with endpoint %s", contentEndPoint)
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

	log.Debugf("Content %s",res)

	var content Content
	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

type GetContentQuery struct{
	Type string `url:"type,omitempty"`
	SpaceKey string `url:"spaceKey,omitempty"`
	Title string `url:"title,omitempty"`
	Status string `url:"status,omitempty"`
	PostingDay time.Time `url:"postingDay,layout,2006-01-01,omitempty"`
	Expand []string `url:"expand,comma,omitempty"`
	Start int `url:"start,omitempty"`
	Limit int `url:"limit,omitempty"`
}


func (w *Wiki) GetDetailedContent(q GetContentQuery)(*ContentResult, error){
	v, err:=query.Values(q)
	if err!=nil{
		return nil, err
	}

	url,_:=w.contentAPIEndpoint()
	url.RawQuery=v.Encode()

	log.Debugf("Calling url %s",url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := w.sendRequest(req)
	if err != nil {
		return nil, err
	}

	log.Debugf("Content JSON %s",res)

	var content ContentResult
	err = json.Unmarshal(res, &content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (w *Wiki) UpdateContent(content *UpdateContentRequest) (*Content, error) {
	jsonbody, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	jsonbody = bytes.Replace(jsonbody, []byte("\\u003c"), []byte("<"), -1)
	jsonbody = bytes.Replace(jsonbody, []byte("\\u003e"), []byte(">"), -1)
	jsonbody = bytes.Replace(jsonbody, []byte("\\u0026"), []byte("&"), -1)

	log.Debugf("Request body %s",jsonbody)
	contentEndPoint, err := w.contentEndpoint(content.Id)
	if err!=nil{
		return nil, err
	}
	log.Debugf("Url %s", contentEndPoint)
	req, err := http.NewRequest("PUT", contentEndPoint.String(), strings.NewReader(string(jsonbody)))

	if err!=nil{
		return nil, err
	}
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

func (w *Wiki) CreateContent(content *UpdateContentRequest)(*Content,error){
	jsonbody, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	jsonbody = bytes.Replace(jsonbody, []byte("\\u003c"), []byte("<"), -1)
	jsonbody = bytes.Replace(jsonbody, []byte("\\u003e"), []byte(">"), -1)
	jsonbody = bytes.Replace(jsonbody, []byte("\\u0026"), []byte("&"), -1)

	log.Debugf("Request body %s",jsonbody)
	contentEndPoint, err := w.contentAPIEndpoint()
	if err!=nil{
		return nil, err
	}
	log.Debugf("Url %s", contentEndPoint)
	req, err := http.NewRequest("POST", contentEndPoint.String(), strings.NewReader(string(jsonbody)))

	if err!=nil{
		return nil, err
	}
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

