package smugmug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type AlbumsService struct {
	s *Service
}

func NewAlbumsService(s *Service) *AlbumsService {
	r := &AlbumsService{s: s}
	return r
}

func (r *AlbumsService) Get(id string) *AlbumsGetCall {
	c := &AlbumsGetCall{s: r.s, urlParams: url.Values{}}
	c.id = id
	c.isAlbum = true
	return c
}

func (r *AlbumsService) GetN(id string) *AlbumsGetCall {
	c := &AlbumsGetCall{s: r.s, urlParams: url.Values{}}
	c.id = id
	c.isAlbum = false
	return c.Paginate(0, 50)
}

type AlbumsServiceResponse struct {
	Code     int
	Message  string
	Response struct {
		ServiceResponse
		Album *json.RawMessage
		Pages *json.RawMessage
	}
	Expansions map[string]*json.RawMessage `json:",omitempty"`
}

type AlbumsGetCall struct {
	id        string
	s         *Service
	urlParams url.Values
	isAlbum   bool
}

func (c *AlbumsGetCall) Expand(expansions []string) *AlbumsGetCall {
	c.urlParams.Set("_expand", strings.Join(expansions, ","))
	return c
}

func (c *AlbumsGetCall) Filter(filter []string) *AlbumsGetCall {
	c.urlParams.Set("_filter", strings.Join(filter, ","))
	return c
}

func (c *AlbumsGetCall) Paginate(start int, count int) *AlbumsGetCall {
	c.urlParams.Set("start", fmt.Sprintf("%d", start))
	c.urlParams.Set("count", fmt.Sprintf("%d", count))
	return c
}

func (c *AlbumsGetCall) doRequest() (*http.Response, error) {
	var path string
	if c.isAlbum {
		path = "album/" + c.id
	} else {
		path = fmt.Sprintf("user/%s!albums", c.id)
	}
	urls := resolveRelative(c.s.BasePath, path)
	urls += "?" + encodeURLParams(c.urlParams)
	req, _ := http.NewRequest("GET", urls, nil)
	c.s.setHeaders(req)
	debugRequest(req)
	return c.s.client.Do(req)
}

func (c *AlbumsGetCall) Do() (*AlbumsGetResponse, error) {
	res, err := c.doRequest()
	if err != nil {
		return nil, err
	}
	debugResponse(res)
	defer closeBody(res)
	if err := checkResponse(res); err != nil {
		return nil, err
	}

	ret := &AlbumsGetResponse{
		ServerResponse: ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}

	albumsRes := &AlbumsServiceResponse{}
	if err := json.NewDecoder(res.Body).Decode(&albumsRes); err != nil {
		return nil, err
	}
	if c.isAlbum {
		ret.Album = &Album{}
		if err := json.Unmarshal(*albumsRes.Response.Album, &ret.Album); err != nil {
			return nil, err
		}
		exp, err := unmarshallExpansions(ret.Album.URIs, albumsRes.Expansions)
		if err != nil {
			return nil, err
		}
		for name, v := range exp {
			switch name {
			case "Node":
				ret.Node = v.(*Node)
			case "User":
				ret.User = v.(*User)
			case "AlbumImages":
				ret.Album.Images = v.([]*AlbumImage)
			}
		}
	} else {
		ret.UserAlbums = &UserAlbums{
			Album: []*Album{},
			Pages: &Pages{},
		}
		if err := json.Unmarshal(*albumsRes.Response.Pages, &ret.UserAlbums.Pages); err != nil {
			return nil, err
		}

		if ret.UserAlbums.Pages.Count > 0 {
			if err := json.Unmarshal(*albumsRes.Response.Album, &ret.UserAlbums.Album); err != nil {
				return nil, err
			}
		}
	}
	return ret, nil
}

type AlbumsGetResponse struct {
	Album *Album

	// AlbumDownload
	// AlbumGeoMedia
	// AlbumHighlightImage // deprecated
	// AlbumImages
	// AlbumPopularMedia
	// AlbumPrices
	// AlbumShareUris
	// ApplyAlbumTemplate
	// CollectImages
	// DeleteAlbumImages
	// Folder // deprecated
	// HighlightImage
	// MoveAlbumImages
	Node *Node
	// ParentFolders // deprecated
	// SortAlbumImages
	// UploadFromUri
	User       *User
	UserAlbums *UserAlbums

	ServerResponse `json:"-"`
}

type Album struct {
	AlbumKey            string `json:",omitempty"`
	AllowDownloads      bool   `json:",omitempty"`
	Backprinting        string `json:",omitempty"`
	BoutiquePackaging   string `json:",omitempty"`
	CanRank             bool   `json:",omitempty"`
	CanShare            bool   `json:",omitempty"`
	Clean               bool   `json:",omitempty"`
	Comments            bool   `json:",omitempty"`
	Date                *Time  `json:",omitempty"`
	Description         string `json:",omitempty"`
	EXIF                bool   `json:",omitempty"`
	External            bool   `json:",omitempty"`
	FamilyEdit          bool   `json:",omitempty"`
	Filenames           bool   `json:",omitempty"`
	FriendEdit          bool   `json:",omitempty"`
	Geography           bool   `json:",omitempty"`
	HasDownloadPassword bool   `json:",omitempty"`
	Header              string `json:",omitempty"`
	HideOwner           bool   `json:",omitempty"`
	ImageCount          int    `json:",omitempty"`
	ImagesLastUpdated   string `json:",omitempty"`
	InterceptShipping   string `json:",omitempty"`
	Keywords            string `json:",omitempty"`
	LargestSize         string `json:",omitempty"`
	LastUpdated         string `json:",omitempty"`
	Name                string `json:",omitempty"`
	NiceName            string `json:",omitempty"`
	NodeID              string `json:",omitempty"`
	OriginalSizes       int    `json:",omitempty"`
	PackagingBranding   bool   `json:",omitempty"`
	Password            string `json:",omitempty"`
	PasswordHint        string `json:",omitempty"`
	Printable           bool   `json:",omitempty"`
	Privacy             string `json:",omitempty"`
	ProofDays           int    `json:",omitempty"`
	Protected           bool   `json:",omitempty"`
	SecurityType        string `json:",omitempty"`
	Share               bool   `json:",omitempty"`
	SmugSearchable      string `json:",omitempty"`
	SortDirection       string `json:",omitempty"`
	SortMethod          string `json:",omitempty"`
	SquareThumbs        bool   `json:",omitempty"`
	TemplateURI         string `json:"TemplateUri"`
	Title               string `json:",omitempty"`
	TotalSizes          int    `json:",omitempty"`
	URLName             string `json:"UrlName,omitempty"`
	URLPath             string `json:"UrlPath,omitempty"`
	Watermark           bool   `json:",omitempty"`
	WorldSearchable     bool   `json:",omitempty"`

	ResponseLevel  string `json:",omitempty"`
	URI            string `json:"Uri,omitempty"`
	URIDescription string `json:"UriDescription,omitempty"`
	URIs           *URIs  `json:"Uris,omitempty"`
	WebURI         string `json:"WebUri,omitempty"`

	Images []*Image `json:",omitempty"`
}

type UserAlbums struct {
	URI         string   `json:",omitempty"`
	Locator     string   `json:",omitempty"`
	LocatorType string   `json:",omitempty"`
	Album       []*Album `json:",omitempty"`
	Pages       *Pages   `json:",omitempty"`
}

// This alias is needed to parse the JSON from SM
//  In the API an expansion image is called an AlbumImage
//  whereas it's an Image if you ask for it directly.
type AlbumImage = Image
