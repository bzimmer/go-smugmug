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

func Test_GetAuthUser(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mockSmugmugUser(t, func(ts *httptest.Server) {

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
		s.Users = NewUsersService(s)

		res, err := s.Users.GetAuthUser().Do()
		if err != nil {
			log.Fatal(err)
		}

		a.Equal("cmac", res.User.NickName)
		a.Equal("/api/v2/user/cmac", res.User.URI)
	})
}

func mockSmugmugUser(t *testing.T, f func(*httptest.Server)) {
	p := pat.New()

	p.Get("/api/v2!authuser", func(res http.ResponseWriter, req *http.Request) {
		var s = `
		{
		    "Response": {
		        "Uri": "/api/v2!authuser",
		        "Locator": "User",
		        "LocatorType": "Object",
		        "User": {
		            "AccountStatus": "Active",
		            "FirstName": "Bob",
		            "FriendsView": true,
		            "ImageCount": 288726,
		            "IsTrial": false,
		            "LastName": "cmac",
		            "NickName": "cmac",
		            "SortBy": "LastUpdated",
		            "ViewPassHint": "",
		            "ViewPassword": "",
		            "Domain": "",
		            "DomainOnly": "",
		            "RefTag": "AyHUS2jglU3Ths",
		            "Name": "cmac",
		            "Plan": "Employee",
		            "QuickShare": true,
		            "Uri": "/api/v2/user/cmac",
		            "WebUri": "https://cmac.smugmug.com",
		            "Uris": {
		                "BioImage": "/api/v2/user/cmac!bioimage",
		                "CoverImage": "/api/v2/user/cmac!coverimage",
		                "UserProfile": "/api/v2/user/cmac!profile",
		                "Node": "/api/v2/node/hdxDH",
		                "Folder": "/api/v2/folder/user/cmac",
		                "Features": "/api/v2/user/cmac!features",
		                "SiteSettings": "/api/v2/user/cmac!sitesettings",
		                "UserAlbums": "/api/v2/user/cmac!albums",
		                "UserGeoMedia": "/api/v2/user/cmac!geomedia",
		                "UserPopularMedia": "/api/v2/user/cmac!popularmedia",
		                "UserFeaturedAlbums": "/api/v2/user/cmac!featuredalbums",
		                "UserRecentImages": "/api/v2/user/cmac!recentimages",
		                "UserImageSearch": "/api/v2/user/cmac!imagesearch",
		                "UserTopKeywords": "/api/v2/user/cmac!topkeywords",
		                "UrlPathLookup": "/api/v2/user/cmac!urlpathlookup",
		                "UserAlbumTemplates": "/api/v2/user/cmac!albumtemplates",
		                "SortUserFeaturedAlbums": "/api/v2/user/cmac!sortfeaturedalbums",
		                "UserTasks": "/api/v2/user/cmac!tasks",
		                "UserWatermarks": "/api/v2/user/cmac!watermarks",
		                "UserPrintmarks": "/api/v2/user/cmac!printmarks",
		                "UserUploadLimits": "/api/v2/user/cmac!uploadlimits",
		                "UserAssetsAlbum": "/api/v2/user/cmac!assetsalbum",
		                "PhotoInvites": "/api/v2/user/cmac!photoinvites",
		                "UserCoupons": "/api/v2/user/cmac!coupons",
		                "UserGuideStates": "/api/v2/user/cmac!guides",
		                "UserHideGuides": "/api/v2/user/cmac!hideguides",
		                "UserGrants": "/api/v2/user/cmac!grants",
		                "DuplicateImageSearch": "/api/v2/user/cmac!duplicateimagesearch",
		                "UserDeletedAlbums": "/api/v2/user/cmac!deletedalbums",
		                "UserDeletedFolders": "/api/v2/user/cmac!deletedfolders",
		                "UserDeletedPages": "/api/v2/user/cmac!deletedpages",
		                "UserContacts": "/api/v2/user/cmac!contacts"
		            },
		            "ResponseLevel": "Full"
		        }
		    },
		    "Code": 200,
		    "Message": "Ok"
		}`
		res.Write([]byte(s))
	})

	ts := httptest.NewServer(p)
	defer ts.Close()

	f(ts)
}
