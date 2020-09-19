package smugmug

import (
	"fmt"
	"github.com/gorilla/pat"
	"github.com/mrjones/oauth"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetImage(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mockSmugmugImage(t, func(ts *httptest.Server) {

		basePath := fmt.Sprintf("%s/api/v2/", ts.URL)
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
		s.Images = NewImagesService(s)

		res, err := s.Images.Get("SD5BL92-1").Do()
		if err != nil {
			log.Fatal(err)
		}

		a.Equal("BaldyProfilePic", res.Image.Keywords)
		a.Equal("/api/v2/image/SD5BL92-1", res.Image.URI)
	})
}

func mockSmugmugImage(t *testing.T, f func(*httptest.Server)) {
	p := pat.New()

	p.Get("/api/v2/image/SD5BL92-1", func(res http.ResponseWriter, req *http.Request) {
		var s = `
		{
			"Response": {
				"Uri": "/api/v2/image/SD5BL92-1?_pretty=&_shorturis=",
				"Locator": "Image",
				"LocatorType": "Object",
				"Image": {
					"Title": "",
					"Caption": "",
					"Keywords": "BaldyProfilePic",
					"KeywordArray": [
						"BaldyProfilePic"
					],
					"Watermark": "No",
					"Latitude": "0.00000000000000",
					"Longitude": "0.00000000000000",
					"Altitude": 0,
					"Hidden": false,
					"ThumbnailUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/Th/i-SD5BL92-Th.jpg",
					"FileName": "BaldyProfilePic.jpg",
					"Processing": false,
					"UploadKey": "1251180257",
					"Date": "2011-04-13T23:08:56+00:00",
					"DateTimeUploaded": "2011-04-13T23:08:56+00:00",
					"DateTimeOriginal": "2010-12-10T22:30:58+00:00",
					"Format": "JPG",
					"OriginalHeight": 1895,
					"OriginalWidth": 1942,
					"OriginalSize": 1728655,
					"LastUpdated": "2012-09-14T06:16:40+00:00",
					"Collectable": true,
					"IsArchive": false,
					"IsVideo": false,
					"CanEdit": false,
					"CanBuy": true,
					"Protected": false,
					"ImageKey": "SD5BL92",
					"Serial": 1,
					"ArchivedUri": "https://photos.smugmug.com/Family/Photos/Baldys-beautiful-self/i-SD5BL92/1/49156caf/D/BaldyProfilePic-D.jpg",
					"ArchivedSize": 1728655,
					"ArchivedMD5": "b794d20764a53e6499d489612406ab76",
					"CanShare": true,
					"Comments": true,
					"ShowKeywords": true,
					"FormattedValues": {
						"Caption": {
							"html": "",
							"text": ""
						},
						"FileName": {
							"html": "BaldyProfilePic.jpg",
							"text": "BaldyProfilePic.jpg"
						}
					},
					"Uri": "/api/v2/image/SD5BL92-1",
					"Uris": {
						"LargestImage": "/api/v2/image/SD5BL92-1!largestimage?_shorturis=",
						"ImageSizes": "/api/v2/image/SD5BL92-1!sizes?_shorturis=",
						"ImageSizeDetails": "/api/v2/image/SD5BL92-1!sizedetails?_shorturis=",
						"PointOfInterest": "/api/v2/image/SD5BL92!pointofinterest",
						"PointOfInterestCrops": "/api/v2/image/SD5BL92!poicrops",
						"Regions": "/api/v2/image/SD5BL92!regions",
						"ImageComments": "/api/v2/image/SD5BL92!comments",
						"ImageMetadata": "/api/v2/image/SD5BL92!metadata?_shorturis=",
						"ImagePrices": "/api/v2/image/SD5BL92!prices?_shorturis=",
						"ImagePricelistExclusions": "/api/v2/image/SD5BL92!pricelistexclusions"
					}
				}
			},
			"Code": 200,
			"Message": "Ok",
			"Expansions": {
				"/api/v2/image/SD5BL92-1!largestimage?_shorturis=": {
					"Uri": "/api/v2/image/SD5BL92-1!largestimage?_shorturis=",
					"Locator": "LargestImage",
					"LocatorType": "Object",
					"LargestImage": {
						"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/O/i-SD5BL92-O.jpg",
						"Size": 1728655,
						"Height": 1895,
						"Width": 1942,
						"Usable": true,
						"Ext": "jpg",
						"Watermarked": false,
						"Uri": "/api/v2/image/SD5BL92-1!largestimage",
						"Uris": {
							"ImageSizeOriginal": "/api/v2/image/SD5BL92-1!sizeoriginal"
						}
					}
				},
				"/api/v2/image/SD5BL92-1!sizes?_shorturis=": {
					"Uri": "/api/v2/image/SD5BL92-1!sizes?_shorturis=",
					"Locator": "ImageSizes",
					"LocatorType": "Object",
					"ImageSizes": {
						"TinyImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/Ti/i-SD5BL92-Ti.jpg",
						"ThumbImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/Th/i-SD5BL92-Th.jpg",
						"SmallImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/S/i-SD5BL92-S.jpg",
						"MediumImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/M/i-SD5BL92-M.jpg",
						"LargeImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/L/i-SD5BL92-L.jpg",
						"XLargeImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/XL/i-SD5BL92-XL.jpg",
						"X2LargeImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/X2/i-SD5BL92-X2.jpg",
						"X3LargeImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/X3/i-SD5BL92-X3.jpg",
						"OriginalImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/O/i-SD5BL92-O.jpg",
						"LargestImageUrl": "https://photos.smugmug.com/photos/i-SD5BL92/1/O/i-SD5BL92-O.jpg",
						"Uri": "/api/v2/image/SD5BL92-1!sizes",
						"Uris": {
							"ImageSizeTiny": "/api/v2/image/SD5BL92-1!sizetiny",
							"ImageSizeThumb": "/api/v2/image/SD5BL92-1!sizethumb",
							"ImageSizeSmall": "/api/v2/image/SD5BL92-1!sizesmall",
							"ImageSizeMedium": "/api/v2/image/SD5BL92-1!sizemedium",
							"ImageSizeLarge": "/api/v2/image/SD5BL92-1!sizelarge",
							"ImageSizeXLarge": "/api/v2/image/SD5BL92-1!sizexlarge",
							"ImageSizeX2Large": "/api/v2/image/SD5BL92-1!sizex2large",
							"ImageSizeX3Large": "/api/v2/image/SD5BL92-1!sizex3large",
							"ImageSizeOriginal": "/api/v2/image/SD5BL92-1!sizeoriginal",
							"ImageSizeCustom": "/api/v2/image/SD5BL92-1!sizecustom",
							"LargestImage": "/api/v2/image/SD5BL92-1!largestimage"
						}
					}
				},
				"/api/v2/image/SD5BL92-1!sizedetails?_shorturis=": {
					"Uri": "/api/v2/image/SD5BL92-1!sizedetails?_shorturis=",
					"Locator": "ImageSizeDetails",
					"LocatorType": "Object",
					"ImageSizeDetails": {
						"ImageUrlTemplate": "https://photos.smugmug.com/photos/i-SD5BL92/1/#size#/i-SD5BL92-#size#.jpg",
						"UsableSizes": [
							"ImageSizeTiny",
							"ImageSizeThumb",
							"ImageSizeSmall",
							"ImageSizeMedium",
							"ImageSizeLarge",
							"ImageSizeXLarge",
							"ImageSizeX2Large",
							"ImageSizeX3Large",
							"ImageSizeOriginal"
						],
						"ImageSizeTiny": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/Ti/i-SD5BL92-Ti.jpg",
							"Ext": "jpg",
							"Height": 100,
							"Width": 100,
							"Size": 15079
						},
						"ImageSizeThumb": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/Th/i-SD5BL92-Th.jpg",
							"Ext": "jpg",
							"Height": 150,
							"Width": 150,
							"Size": 22180
						},
						"ImageSizeSmall": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/S/i-SD5BL92-S.jpg",
							"Ext": "jpg",
							"Height": 300,
							"Width": 307,
							"Size": 55259
						},
						"ImageSizeMedium": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/M/i-SD5BL92-M.jpg",
							"Ext": "jpg",
							"Height": 450,
							"Width": 461,
							"Size": 95952
						},
						"ImageSizeLarge": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/L/i-SD5BL92-L.jpg",
							"Ext": "jpg",
							"Height": 600,
							"Width": 615,
							"Size": 142554
						},
						"ImageSizeXLarge": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/XL/i-SD5BL92-XL.jpg",
							"Ext": "jpg",
							"Height": 768,
							"Width": 787,
							"Size": 209610
						},
						"ImageSizeX2Large": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/X2/i-SD5BL92-X2.jpg",
							"Ext": "jpg",
							"Height": 960,
							"Width": 984,
							"Size": 291376
						},
						"ImageSizeX3Large": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/X3/i-SD5BL92-X3.jpg",
							"Ext": "jpg",
							"Height": 1200,
							"Width": 1230,
							"Size": 414650
						},
						"ImageSizeOriginal": {
							"Url": "https://photos.smugmug.com/photos/i-SD5BL92/1/O/i-SD5BL92-O.jpg",
							"Ext": "jpg",
							"Height": 1895,
							"Width": 1942,
							"Size": 1728655,
							"Watermarked": false
						},
						"Uri": "/api/v2/image/SD5BL92-1!sizedetails"
					}
				},
				"/api/v2/image/SD5BL92!metadata?_shorturis=": {
					"Uri": "/api/v2/image/SD5BL92!metadata?_shorturis=",
					"Locator": "ImageMetadata",
					"LocatorType": "Object",
					"ImageMetadata": {
						"Title": "",
						"Caption": "",
						"UserComment": "",
						"Keywords": "",
						"Author": "",
						"Copyright": "",
						"CopyrightUrl": "",
						"CopyrightFlag": "",
						"Source": "",
						"Credit": "",
						"City": "",
						"State": "",
						"Country": "",
						"Rating": "5",
						"CreatorContactInfo": "",
						"Category": "",
						"SupplementalCategories": "",
						"SpecialInstructions": "",
						"AuthorTitle": "",
						"CountryCode": "",
						"TransmissionReference": "",
						"Headline": "",
						"WriterEditor": "",
						"Lens": "Canon EF 24mm f/1.4L II USM",
						"Make": "Canon",
						"Model": "EOS-1D Mark IV",
						"Aperture": "2",
						"DateTimeModified": "2011-04-13T11:08:09",
						"DateTimeCreated": "2010-12-10T14:30:58",
						"DateCreated": "2010-12-10",
						"TimeCreated": "14:30:58+00:00",
						"MicroDateTimeCreated": "2010-12-10T14:30:58.68",
						"MicroDateTimeDigitized": "2010-12-10T14:30:58.68",
						"DateDigitized": "2010-12-10T14:30:58-08:00",
						"Exposure": "1/5000",
						"ISO": 2000,
						"FocalLength": "24.0 mm",
						"FocalLength35mm": "26.6 mm",
						"CompressedBitsPerPixel": "",
						"Flash": "No Flash",
						"Metering": "Center-weighted average",
						"ExposureProgram": "Aperture-priority AE",
						"ExposureCompensation": "+1/3",
						"ExposureMode": "Auto",
						"LightSource": "",
						"WhiteBalance": "Auto",
						"DigitalZoomRatio": "",
						"Contrast": "",
						"Saturation": "",
						"Sharpness": "",
						"SubjectDistance": "0.88 m",
						"SubjectRange": "",
						"SensingMethod": "",
						"ColorSpace": "sRGB",
						"Brightness": "",
						"LatitudeReference": "",
						"LongitudeReference": "",
						"Latitude": 0,
						"Longitude": 0,
						"AltitudeReference": "",
						"Altitude": 0,
						"SceneCaptureType": "Standard",
						"GainControl": "",
						"ScaleFactor": "",
						"CircleOfConfusion": "0.027 mm",
						"FieldOfView": "68.2 deg",
						"DepthOfField": "0.14 m (0.81 - 0.96 m)",
						"HyperfocalDistance": "10.62 m",
						"NormalizedLightValue": "10",
						"Duration": "",
						"AudioCodec": "",
						"VideoCodec": "",
						"Software": "Adobe Photoshop CS5 Macintosh",
						"SerialNumber": "220100268",
						"LensSerialNumber": "",
						"Uri": "/api/v2/image/SD5BL92!metadata"
					}
				},
				"/api/v2/image/SD5BL92!prices?_shorturis=": {
					"Uri": "/api/v2/image/SD5BL92!prices?_shorturis=",
					"Locator": "CatalogSkuPrice",
					"LocatorType": "Objects",
					"CatalogSkuPrice": [
						{
							"Currency": "USD",
							"Price": 0.99,
							"Uri": "/api/v2/catalog/sku/100405!price?Currency=&Image=%2Fapi%2Fv2%2Fimage%2FSD5BL92",
							"Uris": {
								"CatalogSku": "/api/v2/catalog/sku/100405",
								"CatalogSkuBuy": "/api/v2/catalog/sku/100405!buy"
							},
							"ResponseLevel": "Public"
						},
						{
							"Currency": "USD",
							"Price": 39.99,
							"Uri": "/api/v2/catalog/sku/10571!price?Currency=&Image=%2Fapi%2Fv2%2Fimage%2FSD5BL92",
							"Uris": {
								"CatalogSku": "/api/v2/catalog/sku/10571",
								"CatalogSkuBuy": "/api/v2/catalog/sku/10571!buy"
							},
							"ResponseLevel": "Public"
						},
						{
							"Currency": "USD",
							"Price": 24.99,
							"Uri": "/api/v2/catalog/sku/70002!price?Currency=&Image=%2Fapi%2Fv2%2Fimage%2FSD5BL92",
							"Uris": {
								"CatalogSku": "/api/v2/catalog/sku/70002",
								"CatalogSkuBuy": "/api/v2/catalog/sku/70002!buy"
							},
							"ResponseLevel": "Public"
						}
					]
				}
			}
		}`
		res.Write([]byte(s))
	})

	ts := httptest.NewServer(p)
	defer ts.Close()

	f(ts)
}
