package meta_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/services/meta"
	"github.com/yitsushi/go-misskey/test"
)

func ExampleService_Announcements() {
	client := misskey.NewClient("https://slippy.xyz", os.Getenv("MISSKEY_TOKEN"))

	announcements, err := client.Meta().Announcements(
		&meta.AnnouncementOptions{
			WithUnreads: true,
			SinceID:     "",
			UntilID:     "",
		},
	)
	if err != nil {
		log.Printf("[Announcements] Error happened: %s", err)
		return
	}

	for _, announcement := range announcements {
		log.Printf("[Announcements] %s", core.StringValue(announcement.Title))
	}
}

func TestService_Announcements(t *testing.T) {
	mockClient := test.NewMockHTTPClient()
	mockClient.MockRequest("/api/announcements", func(request *http.Request) (*http.Response, error) {
		defer request.Body.Close()
		body, _ := ioutil.ReadAll(request.Body)

		var statsRequest map[string]interface{}

		err := json.Unmarshal(body, &statsRequest)
		if err != nil {
			t.Errorf("Unable to parse request: %s", err)
			return test.NewMockResponse(http.StatusInternalServerError, []byte{}, err)
		}

		if statsRequest["i"] != "thisistoken" {
			t.Errorf("expected api token = thisistoken; got = %s", statsRequest["i"])
		}

		return test.NewMockResponse(http.StatusOK, []byte(`[
			{
				"id": "firstid",
				"createdAt": "2020-10-08T08:38:24.749Z",
				"updatedAt": null,
				"title": "First Title",
				"text": "This is the content",
				"imageUrl": null
			},
			{
				"id": "secondid",
				"createdAt": "2020-10-10T11:23:23.144Z",
				"updatedAt": "2020-10-10T18:12:13.199Z",
				"title": "Second Title",
				"text": "This is 'not' the content",
				"imageUrl": null
			}
		]`), nil)
	})

	client := misskey.NewClient("https://localhost", "thisistoken")
	client.HTTPClient = mockClient

	announcements, err := client.Meta().Announcements(&meta.AnnouncementOptions{
		WithUnreads: true,
	})
	if err != nil {
		t.Errorf("Unexpected error = %s", err)
		return
	}

	if len(announcements) != 2 {
		t.Errorf("Expected 2 announcements; got = %d", len(announcements))
		return
	}

	first := announcements[0]
	if core.StringValue(first.Title) != "First Title" {
		t.Errorf("Expected title = First Title; got = %s", core.StringValue(first.Title))
	}
	
	expectedDate := time.Date(2020, time.October, 8, 8, 38, 24, 749, time.Now().UTC().Location())
	if first.CreatedAt.UTC().Equal(expectedDate) {
		t.Errorf("Expected creation time = %s; got = %s", expectedDate.UTC(), first.CreatedAt)
	}
}
