package smugmug

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/pat"
	"github.com/mrjones/oauth"
	"github.com/stretchr/testify/assert"
)

func testAlbumsService(url string) *Service {

	basePath := fmt.Sprintf("%s/api/v2/", url)

	consumer := oauth.NewConsumer(
		"consumerKey-0",
		"consumerSecret-0",
		oauth.ServiceProvider{
			RequestTokenUrl:   basePath + "requestToken/",
			AuthorizeTokenUrl: basePath + "authorizeToken/",
			AccessTokenUrl:    basePath + "accessToken/",
		},
	)
	consumer.AdditionalAuthorizationUrlParams = map[string]string{
		"Access":      "Full",
		"Permissions": "Modify",
	}

	token := &oauth.AccessToken{
		Token:  "accessToken-1",
		Secret: "accessTokenSecret-1",
	}

	client, err := consumer.MakeHttpClient(token)
	if err != nil {
		log.Fatal(err)
	}

	s := &Service{client: client, BasePath: basePath}
	s.Albums = NewAlbumsService(s)
	return s
}

func Test_GetAlbum(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mockSmugmugAlbum(t, func(ts *httptest.Server) {
		s := testAlbumsService(ts.URL)

		res, err := s.Albums.Get("kQ3t8P").Do()
		if err != nil {
			log.Fatal(err)
		}

		a.Equal("2015-Oct-Dec", res.Album.NiceName)
		a.Equal("2019-11-26T21:08:41+00:00", res.Album.ImagesLastUpdated)
	})
}

func Test_NoMoreAlbumsAvailable(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mockSmugmugAlbumNoResults(t, func(ts *httptest.Server) {
		s := testAlbumsService(ts.URL)

		res, err := s.Albums.GetN("cmac").Do()
		if err != nil {
			log.Fatal(err)
		}

		a.Equal(0, res.UserAlbums.Pages.Count)
	})
}

func Test_GetAlbumN(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	// test default pagination
	mockSmugmugAlbumN(t, 0, 50, func(ts *httptest.Server) {
		s := testAlbumsService(ts.URL)

		res, err := s.Albums.GetN("cmac").Do()
		if err != nil {
			log.Fatal(err)
		}

		a.Equal(2, len(res.UserAlbums.Album))
		a.Equal(436, res.UserAlbums.Pages.Total)
		a.Equal(50, res.UserAlbums.Pages.RequestedCount)
	})

	// test with pagination
	mockSmugmugAlbumN(t, 3, 22, func(ts *httptest.Server) {
		s := testAlbumsService(ts.URL)

		res, err := s.Albums.GetN("cmac").Paginate(3, 22).Do()
		if err != nil {
			log.Fatal(err)
		}

		a.Equal(2, len(res.UserAlbums.Album))
		a.Equal(436, res.UserAlbums.Pages.Total)
		a.Equal(22, res.UserAlbums.Pages.RequestedCount)
	})
}

func Test_GetAlbumImages(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mockSmugmugAlbumImages(t, func(ts *httptest.Server) {
		s := testAlbumsService(ts.URL)

		res, err := s.Albums.Get("kQ3t8P").Expand([]string{"AlbumImages"}).Do()
		if err != nil {
			log.Fatal(err)
		}

		a.Equal(2, len(res.Album.Images))
	})
}

func mockSmugmugAlbumNoResults(t *testing.T, f func(*httptest.Server)) {
	p := pat.New()

	p.Get("/api/v2/user/cmac!albums", func(res http.ResponseWriter, req *http.Request) {
		var json = fmt.Sprintf(`
		{
			"Response": {
				"Uri": "/api/v2/user/cmac!albums",
				"Locator": "Album",
				"LocatorType": "Objects",
				"Pages": {
					"Total": 436,
					"Start": 400,
					"Count": 0,
					"RequestedCount": 15,
					"FirstPage": "/api/v2/user/cmac!albums?_pretty=&_shorturis=&count=15&start=1",
					"LastPage": "/api/v2/user/cmac!albums?_pretty=&_shorturis=&count=15&start=436",
					"NextPage": "/api/v2/user/cmac!albums?_pretty=&_shorturis=&count=15&start=18",
					"Misaligned": true
				}
			},
			"Code": 200,
			"Message": "Ok"
		}`)
		res.Write([]byte(json))
	})
	ts := httptest.NewServer(p)
	defer ts.Close()

	f(ts)
}

func mockSmugmugAlbumN(t *testing.T, start int, count int, f func(*httptest.Server)) {
	p := pat.New()
	a := assert.New(t)

	p.Get("/api/v2/user/cmac!albums", func(res http.ResponseWriter, req *http.Request) {
		s, _ := req.URL.Query()["start"]
		c, _ := req.URL.Query()["count"]

		a.Equal([]string([]string{fmt.Sprintf("%d", start)}), s)
		a.Equal([]string([]string{fmt.Sprintf("%d", count)}), c)

		var json = fmt.Sprintf(`
	{
		"Response": {
			"Uri": "/api/v2/user/cmac!albums?_pretty=&_shorturis=&count=15&start=3",
			"Locator": "Album",
			"LocatorType": "Objects",
			"Album": [
				{
					"NiceName": "Black-lives-matter-protest",
					"UrlName": "Black-lives-matter-protest",
					"Title": "Black lives matter protest",
					"Name": "Black lives matter protest",
					"AllowDownloads": false,
					"Description": "",
					"EXIF": true,
					"External": true,
					"Filenames": false,
					"Geography": true,
					"Keywords": "",
					"PasswordHint": "",
					"Protected": false,
					"SortDirection": "Ascending",
					"SortMethod": "Date Taken",
					"SecurityType": "None",
					"CommerceLightbox": true,
					"AlbumKey": "jbBNhR",
					"CanBuy": true,
					"CanFavorite": false,
					"LastUpdated": "2020-06-01T03:11:56+00:00",
					"ImagesLastUpdated": "2020-06-01T03:12:43+00:00",
					"NodeID": "q2qP7F",
					"ImageCount": 16,
					"UrlPath": "/Events/Black-lives-matter-protest",
					"CanShare": true,
					"HasDownloadPassword": false,
					"Packages": false,
					"Uri": "/api/v2/album/jbBNhR",
					"WebUri": "https://cmac.smugmug.com/Events/Black-lives-matter-protest",
					"Uris": {
						"AlbumShareUris": "/api/v2/album/jbBNhR!shareuris",
						"Node": "/api/v2/node/q2qP7F",
						"NodeCoverImage": "/api/v2/node/q2qP7F!cover",
						"User": "/api/v2/user/cmac",
						"Folder": "/api/v2/folder/user/cmac/Events",
						"ParentFolders": "/api/v2/folder/user/cmac/Events!parents",
						"HighlightImage": "/api/v2/highlight/node/q2qP7F",
						"AlbumHighlightImage": "/api/v2/album/jbBNhR!highlightimage",
						"AlbumImages": "/api/v2/album/jbBNhR!images?_shorturis=",
						"AlbumPopularMedia": "/api/v2/album/jbBNhR!popularmedia",
						"AlbumGeoMedia": "/api/v2/album/jbBNhR!geomedia",
						"AlbumComments": "/api/v2/album/jbBNhR!comments",
						"AlbumPrices": "/api/v2/album/jbBNhR!prices",
						"AlbumPricelistExclusions": "/api/v2/album/jbBNhR!pricelistexclusions"
					},
					"ResponseLevel": "Public"
				},
				{
					"NiceName": "Mosko-zoom",
					"UrlName": "Mosko-zoom",
					"Title": "Mosko zoom",
					"Name": "Mosko zoom",
					"AllowDownloads": true,
					"Description": "",
					"EXIF": true,
					"External": true,
					"Filenames": false,
					"Geography": true,
					"Keywords": "",
					"PasswordHint": "",
					"Protected": false,
					"SortDirection": "Ascending",
					"SortMethod": "Date Taken",
					"SecurityType": "None",
					"CommerceLightbox": true,
					"AlbumKey": "mW5sgS",
					"CanBuy": true,
					"CanFavorite": false,
					"LastUpdated": "2020-05-06T01:43:20+00:00",
					"ImagesLastUpdated": "2020-05-06T01:43:40+00:00",
					"NodeID": "9MLJRT",
					"ImageCount": 29,
					"UrlPath": "/Other/Mosko-zoom",
					"CanShare": true,
					"HasDownloadPassword": false,
					"Packages": false,
					"Uri": "/api/v2/album/mW5sgS",
					"WebUri": "https://cmac.smugmug.com/Other/Mosko-zoom",
					"Uris": {
						"AlbumShareUris": "/api/v2/album/mW5sgS!shareuris",
						"Node": "/api/v2/node/9MLJRT",
						"NodeCoverImage": "/api/v2/node/9MLJRT!cover",
						"User": "/api/v2/user/cmac",
						"Folder": "/api/v2/folder/user/cmac/Other",
						"ParentFolders": "/api/v2/folder/user/cmac/Other!parents",
						"HighlightImage": "/api/v2/highlight/node/9MLJRT",
						"AlbumHighlightImage": "/api/v2/album/mW5sgS!highlightimage",
						"AlbumImages": "/api/v2/album/mW5sgS!images?_shorturis=",
						"AlbumPopularMedia": "/api/v2/album/mW5sgS!popularmedia",
						"AlbumGeoMedia": "/api/v2/album/mW5sgS!geomedia",
						"AlbumComments": "/api/v2/album/mW5sgS!comments",
						"AlbumDownload": "/api/v2/album/mW5sgS!download",
						"AlbumPrices": "/api/v2/album/mW5sgS!prices",
						"AlbumPricelistExclusions": "/api/v2/album/mW5sgS!pricelistexclusions"
					},
					"ResponseLevel": "Public"
				}
			],
			"Pages": {
				"Total": 436,
				"Start": 3,
				"Count": 15,
				"RequestedCount": %d,
				"FirstPage": "/api/v2/user/cmac!albums?_pretty=&_shorturis=&count=15&start=1",
				"LastPage": "/api/v2/user/cmac!albums?_pretty=&_shorturis=&count=15&start=436",
				"NextPage": "/api/v2/user/cmac!albums?_pretty=&_shorturis=&count=15&start=18",
				"Misaligned": true
			}
		},
		"Code": 200,
		"Message": "Ok"
	}`, count)
		res.Write([]byte(json))
	})
	ts := httptest.NewServer(p)
	defer ts.Close()

	f(ts)
}

func mockSmugmugAlbum(t *testing.T, f func(*httptest.Server)) {
	p := pat.New()

	p.Get("/api/v2/album/kQ3t8P", func(res http.ResponseWriter, req *http.Request) {
		var json = `
	{
		"Request": {
			"Version": "v2",
			"Method": "GET",
			"Uri": "/api/v2/album/kQ3t8P"
		},
		"Options": {
		},
		"Response": {
			"Uri": "/api/v2/album/kQ3t8P",
			"Locator": "Album",
			"LocatorType": "Object",
			"Album": {
				"NiceName": "2015-Oct-Dec",
				"UrlName": "2015-Oct-Dec",
				"Title": "2015 Oct-Dec",
				"Name": "2015 Oct-Dec",
				"AllowDownloads": true,
				"Description": "",
				"EXIF": true,
				"External": true,
				"Filenames": false,
				"Geography": true,
				"Keywords": "",
				"PasswordHint": "",
				"Protected": false,
				"SortDirection": "Ascending",
				"SortMethod": "Position",
				"SecurityType": "None",
				"CommerceLightbox": false,
				"AlbumKey": "kQ3t8P",
				"CanBuy": true,
				"CanFavorite": false,
				"LastUpdated": "2019-11-26T21:08:23+00:00",
				"ImagesLastUpdated": "2019-11-26T21:08:41+00:00",
				"NodeID": "h22spN",
				"ImageCount": 183,
				"UrlPath": "/Family/Photos/2015-Oct-Dec",
				"CanShare": true,
				"HasDownloadPassword": false,
				"Packages": false,
				"Uri": "/api/v2/album/kQ3t8P",
				"WebUri": "https://cmac.smugmug.com/Family/Photos/2015-Oct-Dec",
				"UriDescription": "Album by key",
				"Uris": {
					"AlbumShareUris": {
						"Uri": "/api/v2/album/kQ3t8P!shareuris",
						"Locator": "AlbumShareUris",
						"LocatorType": "Object",
						"UriDescription": "URIs that are useful for sharing",
						"EndpointType": "AlbumShareUris"
					},
					"Node": {
						"Uri": "/api/v2/node/h22spN",
						"Locator": "Node",
						"LocatorType": "Object",
						"UriDescription": "Node with the given id.",
						"EndpointType": "Node"
					},
					"NodeCoverImage": {
						"Uri": "/api/v2/node/h22spN!cover",
						"Locator": "Image",
						"LocatorType": "Object",
						"UriDescription": "Cover image for a folder, album, or page",
						"EndpointType": "NodeCoverImage"
					},
					"User": {
						"Uri": "/api/v2/user/cmac",
						"Locator": "User",
						"LocatorType": "Object",
						"UriDescription": "User By Nickname",
						"EndpointType": "User"
					},
					"Folder": {
						"Uri": "/api/v2/folder/user/cmac/Family/Photos",
						"Locator": "Folder",
						"LocatorType": "Object",
						"UriDescription": "A folder or legacy (sub)category by UrlPath",
						"EndpointType": "Folder"
					},
					"ParentFolders": {
						"Uri": "/api/v2/folder/user/cmac/Family/Photos!parents",
						"Locator": "Folder",
						"LocatorType": "Objects",
						"UriDescription": "The sequence of parent folders, from the given folder to the root",
						"EndpointType": "ParentFolders"
					},
					"HighlightImage": {
						"Uri": "/api/v2/highlight/node/h22spN",
						"Locator": "Image",
						"LocatorType": "Object",
						"UriDescription": "Highlight image for a folder, album, or page",
						"EndpointType": "HighlightImage"
					},
					"AddSamplePhotos": {
						"Uri": "/api/v2/album/kQ3t8P!addsamplephotos",
						"UriDescription": "Add sample photos to Album",
						"EndpointType": "AddSamplePhotos"
					},
					"AlbumHighlightImage": {
						"Uri": "/api/v2/album/kQ3t8P!highlightimage",
						"Locator": "AlbumImage",
						"LocatorType": "Object",
						"UriDescription": "Highlight image for album",
						"EndpointType": "AlbumHighlightImage"
					},
					"AlbumImages": {
						"Uri": "/api/v2/album/kQ3t8P!images",
						"Locator": "AlbumImage",
						"LocatorType": "Objects",
						"UriDescription": "Images from album",
						"EndpointType": "AlbumImages"
					},
					"AlbumPopularMedia": {
						"Uri": "/api/v2/album/kQ3t8P!popularmedia",
						"Locator": "AlbumImage",
						"LocatorType": "Objects",
						"UriDescription": "Popular images from album",
						"EndpointType": "AlbumPopularMedia"
					},
					"AlbumGeoMedia": {
						"Uri": "/api/v2/album/kQ3t8P!geomedia",
						"Locator": "AlbumImage",
						"LocatorType": "Objects",
						"UriDescription": "Geotagged images from album",
						"EndpointType": "AlbumGeoMedia"
					},
					"AlbumComments": {
						"Uri": "/api/v2/album/kQ3t8P!comments",
						"Locator": "Comment",
						"LocatorType": "Objects",
						"UriDescription": "Comments on album",
						"EndpointType": "AlbumComments"
					},
					"AlbumDownload": {
						"Uri": "/api/v2/album/kQ3t8P!download",
						"Locator": "Download",
						"LocatorType": "Objects",
						"UriDescription": "Download album",
						"EndpointType": "AlbumDownload"
					},
					"AlbumPrices": {
						"Uri": "/api/v2/album/kQ3t8P!prices",
						"Locator": "CatalogSkuPrice",
						"LocatorType": "Objects",
						"UriDescription": "Purchasable Skus",
						"EndpointType": "AlbumPrices"
					},
					"AlbumPricelistExclusions": {
						"Uri": "/api/v2/album/kQ3t8P!pricelistexclusions",
						"Locator": "AlbumPricelistExclusions",
						"LocatorType": "Object",
						"UriDescription": "Pricelist information for an Album",
						"EndpointType": "AlbumPricelistExclusions"
					}
				},
				"ResponseLevel": "Public"
			},
			"UriDescription": "Album by key",
			"EndpointType": "Album",
			"DocUri": "https://api.smugmug.com/api/v2/doc/reference/album.html",
			"Timing": {
				"Total": {
					"time": 0.02675,
					"cycles": 1,
					"objects": 0
				}
			}
		},
		"Code": 200,
		"Message": "Ok"
	}`
		res.Write([]byte(json))
	})

	ts := httptest.NewServer(p)
	defer ts.Close()

	f(ts)
}

func mockSmugmugAlbumImages(t *testing.T, f func(*httptest.Server)) {
	p := pat.New()

	p.Get("/api/v2/album/kQ3t8P", func(res http.ResponseWriter, req *http.Request) {
		var json = `
	{
		"Response": {
		"Uri": "/api/v2/album/kQ3t8P?_pretty=&_shorturis=",
		"Locator": "Album",
		"LocatorType": "Object",
		"Album": {
			"NiceName": "2015-Oct-Dec",
			"UrlName": "2015-Oct-Dec",
			"Title": "2015 Oct-Dec",
			"Name": "2015 Oct-Dec",
			"AllowDownloads": true,
			"Description": "",
			"EXIF": true,
			"External": true,
			"Filenames": false,
			"Geography": true,
			"Keywords": "",
			"PasswordHint": "",
			"Protected": false,
			"SortDirection": "Ascending",
			"SortMethod": "Position",
			"SecurityType": "None",
			"CommerceLightbox": false,
			"AlbumKey": "kQ3t8P",
			"CanBuy": true,
			"CanFavorite": false,
			"LastUpdated": "2019-11-26T21:08:23+00:00",
			"ImagesLastUpdated": "2019-11-26T21:08:41+00:00",
			"NodeID": "h22spN",
			"ImageCount": 183,
			"UrlPath": "/Family/Photos/2015-Oct-Dec",
			"CanShare": true,
			"HasDownloadPassword": false,
			"Packages": false,
			"Uri": "/api/v2/album/kQ3t8P",
			"WebUri": "https://cmac.smugmug.com/Family/Photos/2015-Oct-Dec",
			"Uris": {
			"AlbumShareUris": "/api/v2/album/kQ3t8P!shareuris",
			"Node": "/api/v2/node/h22spN",
			"NodeCoverImage": "/api/v2/node/h22spN!cover",
			"User": "/api/v2/user/cmac",
			"Folder": "/api/v2/folder/user/cmac/Family/Photos",
			"ParentFolders": "/api/v2/folder/user/cmac/Family/Photos!parents",
			"HighlightImage": "/api/v2/highlight/node/h22spN",
			"AlbumHighlightImage": "/api/v2/album/kQ3t8P!highlightimage",
			"AlbumImages": "/api/v2/album/kQ3t8P!images?_shorturis=",
			"AlbumPopularMedia": "/api/v2/album/kQ3t8P!popularmedia",
			"AlbumGeoMedia": "/api/v2/album/kQ3t8P!geomedia",
			"AlbumComments": "/api/v2/album/kQ3t8P!comments",
			"AlbumDownload": "/api/v2/album/kQ3t8P!download",
			"AlbumPrices": "/api/v2/album/kQ3t8P!prices",
			"AlbumPricelistExclusions": "/api/v2/album/kQ3t8P!pricelistexclusions"
			},
			"ResponseLevel": "Public"
		}
		},
		"Code": 200,
		"Message": "Ok",
		"Expansions": {
		"/api/v2/album/kQ3t8P!images?_shorturis=": {
			"Uri": "/api/v2/album/kQ3t8P!images?_shorturis=",
			"Locator": "AlbumImage",
			"LocatorType": "Objects",
			"AlbumImage": [
			{
				"Title": "",
				"Caption": "",
				"Keywords": "",
				"KeywordArray": [],
				"Watermark": "No",
				"Latitude": "0.00000000000000",
				"Longitude": "0.00000000000000",
				"Altitude": 0,
				"Hidden": false,
				"ThumbnailUrl": "https://photos.smugmug.com/photos/i-rPZcMrk/0/Th/i-rPZcMrk-Th.jpg",
				"FileName": "_DSC6480.jpg",
				"Processing": false,
				"UploadKey": "4493986498",
				"Date": "2015-11-05T10:23:44+00:00",
				"DateTimeUploaded": "2015-11-05T10:23:44+00:00",
				"DateTimeOriginal": "2015-10-09T20:18:58+00:00",
				"Format": "JPG",
				"OriginalHeight": 3486,
				"OriginalWidth": 2925,
				"OriginalSize": 4836122,
				"LastUpdated": "2015-11-05T10:23:47+00:00",
				"Collectable": true,
				"IsArchive": false,
				"IsVideo": false,
				"CanEdit": false,
				"CanBuy": true,
				"Protected": false,
				"Watermarked": false,
				"ImageKey": "rPZcMrk",
				"Serial": 0,
				"ArchivedUri": "https://photos.smugmug.com/Family/Photos/2015-Oct-Dec/i-rPZcMrk/0/4dd3a6ef/D/_DSC6480-D.jpg",
				"ArchivedSize": 4836122,
				"ArchivedMD5": "5f6ae54060f16056eb2ca087c86be490",
				"CanShare": true,
				"Comments": true,
				"ShowKeywords": true,
				"FormattedValues": {
				"Caption": {
					"html": "",
					"text": ""
				},
				"FileName": {
					"html": "_DSC6480.jpg",
					"text": "_DSC6480.jpg"
				}
				},
				"Uri": "/api/v2/album/kQ3t8P/image/rPZcMrk-0",
				"Uris": {
				"LargestImage": "/api/v2/image/rPZcMrk-0!largestimage",
				"ImageSizes": "/api/v2/image/rPZcMrk-0!sizes",
				"ImageSizeDetails": "/api/v2/image/rPZcMrk-0!sizedetails",
				"PointOfInterest": "/api/v2/image/rPZcMrk!pointofinterest",
				"PointOfInterestCrops": "/api/v2/image/rPZcMrk!poicrops",
				"Regions": "/api/v2/image/rPZcMrk!regions",
				"ImageComments": "/api/v2/image/rPZcMrk!comments",
				"ImageMetadata": "/api/v2/image/rPZcMrk!metadata",
				"ImagePrices": "/api/v2/image/rPZcMrk!prices",
				"ImagePricelistExclusions": "/api/v2/image/rPZcMrk!pricelistexclusions",
				"Album": "/api/v2/album/kQ3t8P",
				"Image": "/api/v2/image/rPZcMrk-0",
				"AlbumImageMetadata": "/api/v2/album/kQ3t8P/image/rPZcMrk-0!metadata",
				"AlbumImageShareUris": "/api/v2/album/kQ3t8P/image/rPZcMrk-0!shareuris"
				},
				"Movable": true,
				"Origin": "Album",
				"WebUri": "https://cmac.smugmug.com/Family/Photos/2015-Oct-Dec/i-rPZcMrk"
			},
			{
				"Title": "",
				"Caption": "",
				"Keywords": "",
				"KeywordArray": [],
				"Watermark": "No",
				"Latitude": "0.00000000000000",
				"Longitude": "0.00000000000000",
				"Altitude": 0,
				"Hidden": false,
				"ThumbnailUrl": "https://photos.smugmug.com/photos/i-xr2CptT/0/Th/i-xr2CptT-Th.jpg",
				"FileName": "_DSC6498.jpg",
				"Processing": false,
				"UploadKey": "4493987525",
				"Date": "2015-11-05T10:25:07+00:00",
				"DateTimeUploaded": "2015-11-05T10:25:07+00:00",
				"DateTimeOriginal": "2015-10-09T20:21:09+00:00",
				"Format": "JPG",
				"OriginalHeight": 4912,
				"OriginalWidth": 4834,
				"OriginalSize": 10220985,
				"LastUpdated": "2015-11-05T10:25:11+00:00",
				"Collectable": true,
				"IsArchive": false,
				"IsVideo": false,
				"CanEdit": false,
				"CanBuy": true,
				"Protected": false,
				"Watermarked": false,
				"ImageKey": "xr2CptT",
				"Serial": 0,
				"ArchivedUri": "https://photos.smugmug.com/Family/Photos/2015-Oct-Dec/i-xr2CptT/0/b3c8920a/D/_DSC6498-D.jpg",
				"ArchivedSize": 10220985,
				"ArchivedMD5": "6843adeb6d63cf1db77a65c3f09b3325",
				"CanShare": true,
				"Comments": true,
				"ShowKeywords": true,
				"FormattedValues": {
				"Caption": {
					"html": "",
					"text": ""
				},
				"FileName": {
					"html": "_DSC6498.jpg",
					"text": "_DSC6498.jpg"
				}
				},
				"Uri": "/api/v2/album/kQ3t8P/image/xr2CptT-0",
				"Uris": {
				"LargestImage": "/api/v2/image/xr2CptT-0!largestimage",
				"ImageSizes": "/api/v2/image/xr2CptT-0!sizes",
				"ImageSizeDetails": "/api/v2/image/xr2CptT-0!sizedetails",
				"PointOfInterest": "/api/v2/image/xr2CptT!pointofinterest",
				"PointOfInterestCrops": "/api/v2/image/xr2CptT!poicrops",
				"Regions": "/api/v2/image/xr2CptT!regions",
				"ImageComments": "/api/v2/image/xr2CptT!comments",
				"ImageMetadata": "/api/v2/image/xr2CptT!metadata",
				"ImagePrices": "/api/v2/image/xr2CptT!prices",
				"ImagePricelistExclusions": "/api/v2/image/xr2CptT!pricelistexclusions",
				"Album": "/api/v2/album/kQ3t8P",
				"Image": "/api/v2/image/xr2CptT-0",
				"AlbumImageMetadata": "/api/v2/album/kQ3t8P/image/xr2CptT-0!metadata",
				"AlbumImageShareUris": "/api/v2/album/kQ3t8P/image/xr2CptT-0!shareuris"
				},
				"Movable": true,
				"Origin": "Album",
				"WebUri": "https://cmac.smugmug.com/Family/Photos/2015-Oct-Dec/i-xr2CptT"
			}
			],
			"Pages": {
			"Total": 183,
			"Start": 1,
			"Count": 100,
			"RequestedCount": 100,
			"FirstPage": "/api/v2/album/kQ3t8P!images?_shorturis=&start=1&count=100",
			"LastPage": "/api/v2/album/kQ3t8P!images?_shorturis=&start=101&count=100",
			"NextPage": "/api/v2/album/kQ3t8P!images?_shorturis=&start=101&count=100"
			}
		}
		}
	}`
		res.Write([]byte(json))
	})

	ts := httptest.NewServer(p)
	defer ts.Close()

	f(ts)
}
