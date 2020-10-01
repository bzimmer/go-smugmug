package smugmug

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const basePath = "https://api.smugmug.com/api/v2/"

var dateFormats = []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02"}

// Time represents a datetime
type Time struct {
	time.Time
}

// UnmarshalJSON parses the string representing a Time struct
func (t *Time) UnmarshalJSON(data []byte) error {
	// data should be a string representing a time, possibly with a leading negative sign
	date := string(data)
	if strings.HasPrefix(date, `"-`) {
		return nil
	}

	for _, fmt := range dateFormats {
		r, err := time.Parse(`"`+fmt+`"`, date)
		if err == nil {
			t.Time = r
			return nil
		}
	}

	return fmt.Errorf("unable to parse '%s'", date)
}

type FormattedValues struct {
	Caption struct {
		HTML string `json:html`
		Text string `json:text`
	}
	Name struct {
		HTML string `json:html`
	}
	Description struct {
		HTML string `json:html`
		Text string `json:text`
	}
}

type URI interface{}
type URIs map[string]URI

type Pages struct {
	Total          int  `json:"Total"`
	Start          int  `json:"Start"`
	Count          int  `json:"Count"`
	RequestedCount int  `json:"RequestedCount"`
	FirstPage      *URI `json:"FirstPage"`
	LastPage       *URI `json:"LastPage"`
	NextPage       *URI `json:"NextPage"`
}

func (p *Pages) Next() int {
	return p.Start + p.Count
}

func (p *Pages) Previous() int {
	return p.Start - p.Count
}

func (p *Pages) Remaining() int {
	return p.Total - p.Start - p.Count
}

func parseURI(u URI) string {
	switch u.(type) {
	case string:
		return u.(string)
	case map[string]interface{}:
		return u.(map[string]interface{})["Uri"].(string)
	default:
		log.Fatal(u)
	}
	return ""
}

type Service struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	Albums *AlbumsService
	Images *ImagesService
	Nodes  *NodesService
	Users  *UsersService
}

func (s *Service) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", s.userAgent())
}

func (s *Service) userAgent() string {
	return "go-smugmug"
}

type ServiceResponse struct {
	URI            string `json:"Uri"`
	URIDescription string `json:"UriDescription"`
	ResponseLevel  string
	DocURI         string `json:"DocUri"`
	EndpointType   string
	Locator        string
	LocatorType    string
	Timing         struct {
		Total struct {
			time    float32
			cycles  int
			objects int
		}
	}
}

type ServerResponse struct {
	HTTPStatusCode int
	Header         http.Header
}

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.Albums = NewAlbumsService(s)
	s.Images = NewImagesService(s)
	s.Nodes = NewNodesService(s)
	s.Users = NewUsersService(s)
	return s, nil
}

func closeBody(res *http.Response) error {
	return res.Body.Close()
}

func checkResponse(res *http.Response) error {
	if res.StatusCode >= 400 {
		return fmt.Errorf("%s", res.Status)
	}
	return nil
}

func encodeURLParams(overrides url.Values) string {
	params := url.Values{}
	params.Set("_expand", "")
	params.Set("_shorturis", "")
	params.Set("_verbosity", "1")
	if debug {
		params.Set("_pretty", "")
	}
	for k := range overrides {
		params[k] = overrides[k]
	}
	ret := params.Encode()
	ret = strings.Replace(ret, "%2C", ",", -1)
	return ret
}

func resolveRelative(basestr string, relstr string) string {
	u, _ := url.Parse(basestr)
	rel, _ := url.Parse(relstr)
	u = u.ResolveReference(rel)
	us := u.String()
	us = strings.Replace(us, "%7B", "{", -1)
	us = strings.Replace(us, "%7D", "}", -1)
	return us
}

func unmarshallExpansions(uris *URIs, exp map[string]*json.RawMessage) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	for name, uri := range *uris {
		u := parseURI(uri)
		switch name {
		case "Album", "ImageAlbum":
			if value, ok := exp[u]; ok {
				res := struct{ Album *Album }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.Album
			}
		case "Node", "ParentNode":
			if value, ok := exp[u]; ok {
				res := struct{ Node *Node }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.Node
			}
		case "ChildNodes", "ParentNodes":
			if value, ok := exp[u]; ok {
				res := struct{ Node []*Node }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.Node
			}
		case "HighlightImage":
			if value, ok := exp[u]; ok {
				res := struct{ Image *Image }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.Image
			}
		case "ImageDownload":
			if value, ok := exp[u]; ok {
				res := struct{ ImageDownload *ImageDownload }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.ImageDownload
			}
		case "ImageMetadata":
			if value, ok := exp[u]; ok {
				res := struct{ ImageMetadata *ImageMetadata }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.ImageMetadata
			}
		case "User", "ImageOwner":
			if value, ok := exp[u]; ok {
				res := struct{ User *User }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.User
			}
		case "ImagePrices":
			if value, ok := exp[u]; ok {
				res := struct{ CatalogSkuPrice []*CatalogSkuPrice }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.CatalogSkuPrice
			}
		case "ImageSizeDetails":
			if value, ok := exp[u]; ok {
				res := struct{ ImageSizeDetails *ImageSizeDetails }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.ImageSizeDetails
			}
		case "ImageSizes":
			if value, ok := exp[u]; ok {
				res := struct{ ImageSizes *ImageSizes }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.ImageSizes
			}
		case "LargestImage":
			if value, ok := exp[u]; ok {
				res := struct{ LargestImage *LargestImage }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.LargestImage
			}
		case "UserAlbums":
			if value, ok := exp[u]; ok {
				res := struct{ UserAlbums *UserAlbums }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res
			}
		case "AlbumImages":
			if value, ok := exp[u]; ok {
				res := struct{ AlbumImage []*AlbumImage }{}
				if err := json.Unmarshal(*value, &res); err != nil {
					return nil, err
				}
				ret[name] = res.AlbumImage
			}
		}
	}
	return ret, nil
}
