package meta_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/test"
)

func ExampleService_InstanceMeta() {
	client := misskey.NewClient("https://slippy.xyz", os.Getenv("MISSKEY_TOKEN"))

	meta, err := client.Meta().InstanceMeta(true)
	if err != nil {
		log.Printf("[Meta] Error happened: %s", err)
		return
	}

	log.Printf("[InstanceMeta/Name] %s", core.StringValue(meta.Name))

	for _, emoji := range meta.Emojis {
		log.Printf("[InstanceMeta/Emoji] %s", core.StringValue(emoji.Name))
	}

	boolStatusToString := func(v bool) string {
		if v {
			return "enabled"
		}

		return "disabled"
	}

	log.Printf("[InstanceMeta/Feature] Registration:   %s", boolStatusToString(meta.Features.Registration))
	log.Printf("[InstanceMeta/Feature] LocalTimeLine:  %s", boolStatusToString(meta.Features.LocalTimeLine))
	log.Printf("[InstanceMeta/Feature] GlobalTimeLine: %s", boolStatusToString(meta.Features.GlobalTimeLine))
	log.Printf("[InstanceMeta/Feature] Elasticsearch:  %s", boolStatusToString(meta.Features.Elasticsearch))
	log.Printf("[InstanceMeta/Feature] Hcaptcha:       %s", boolStatusToString(meta.Features.Hcaptcha))
	log.Printf("[InstanceMeta/Feature] Recaptcha:      %s", boolStatusToString(meta.Features.Recaptcha))
	log.Printf("[InstanceMeta/Feature] ObjectStorage:  %s", boolStatusToString(meta.Features.ObjectStorage))
	log.Printf("[InstanceMeta/Feature] Twitter:        %s", boolStatusToString(meta.Features.Twitter))
	log.Printf("[InstanceMeta/Feature] Github:         %s", boolStatusToString(meta.Features.Github))
	log.Printf("[InstanceMeta/Feature] Discord:        %s", boolStatusToString(meta.Features.Discord))
	log.Printf("[InstanceMeta/Feature] ServiceWorker:  %s", boolStatusToString(meta.Features.ServiceWorker))
	log.Printf("[InstanceMeta/Feature] MiAuth:         %s", boolStatusToString(meta.Features.MiAuth))
}

func TestService_InstanceMeta(t *testing.T) {
	mockClient := test.NewMockHTTPClient()
	mockClient.MockRequest("/api/meta", func(request *http.Request) (*http.Response, error) {
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
			"Name": "Instance Name"
		}`), nil)
	})

	client := misskey.NewClient("https://localhost", "thisistoken")
	client.HTTPClient = mockClient

	instanceMeta, err := client.Meta().InstanceMeta(true)
	if err != nil {
		t.Errorf("Unexpected error = %s", err)
		return
	}

	if core.StringValue(instanceMeta.Name) != "Instance Name" {
		t.Errorf("Expected name of the instance = Instance Name; got = %s", core.StringValue(instanceMeta.Name))
	}
}
