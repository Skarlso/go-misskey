package meta_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/test"
)

func ExampleService_Stats() {
	client := misskey.NewClient("https://slippy.xyz", os.Getenv("MISSKEY_TOKEN"))

	stats, err := client.Meta().Stats()
	if err != nil {
		log.Printf("[Meta] Error happened: %s", err)
		return
	}

	log.Printf("[Stats] Instances:          %d", stats.Instances)
	log.Printf("[Stats] NotesCount:         %d", stats.NotesCount)
	log.Printf("[Stats] UsersCount:         %d", stats.UsersCount)
	log.Printf("[Stats] OriginalNotesCount: %d", stats.OriginalNotesCount)
	log.Printf("[Stats] OriginalUsersCount: %d", stats.OriginalUsersCount)
}

func TestService_Stats(t *testing.T) {
	mockClient := test.NewMockHTTPClient()
	mockClient.MockRequest("/api/stats", func(request *http.Request) (*http.Response, error) {
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

		return test.NewMockResponse(http.StatusOK, []byte(`{
			"originalNotesCount": 101,
			"notesCount": 102,
			"originalUsersCount": 103,
			"usersCount": 104,
			"instances": 100
		}`), nil)
	})

	client := misskey.NewClient("https://localhost", "thisistoken")
	client.HTTPClient = mockClient

	stats, err := client.Meta().Stats()
	if err != nil {
		t.Errorf("Unexpected error = %s", err)
		return
	}

	if stats.Instances != 100 {
		t.Errorf("Expected number of instances = %d; got = %d", 100, stats.Instances)
	}
	
	if stats.OriginalNotesCount != 101 {
		t.Errorf("Expected number of original notes = %d; got = %d", 101, stats.OriginalNotesCount)
	}

	if stats.NotesCount != 102 {
		t.Errorf("Expected number of notes = %d; got = %d", 102, stats.NotesCount)
	}
	
	if stats.OriginalUsersCount != 103 {
		t.Errorf("Expected number of original users = %d; got = %d", 103, stats.OriginalUsersCount)
	}
	
	if stats.UsersCount != 104 {
		t.Errorf("Expected number of users = %d; got = %d", 104, stats.UsersCount)
	}
}
