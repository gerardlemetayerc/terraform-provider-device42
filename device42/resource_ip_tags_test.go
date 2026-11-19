package device42

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceDevice42IpTagsRoundTrip(t *testing.T) {
	var gotIncludeCols string
	var gotFormTags string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/2.0/ips/":
			if err := r.ParseForm(); err != nil {
				t.Fatalf("POST parse form: %v", err)
			}
			gotFormTags = r.PostFormValue("tags")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"msg":["ip added or updated",50,"10.1.3.1",true,false],"code":0}`))
		case r.Method == http.MethodGet && r.URL.Path == "/2.0/ips/":
			gotIncludeCols = r.URL.Query().Get("include_cols")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"total_count":1,"ips":[{"ip":"10.1.3.1","available":"no","tags":["foo","bar"],"devices":[{"device_id":134,"name":"USNHCTIDS01"}]}]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := resty.New()
	client.SetBaseURL(srv.URL)

	d := schema.TestResourceDataRaw(t, resourceD42Ip().Schema, map[string]interface{}{
		"ip":        "10.1.3.1",
		"available": "no",
		"tags":      "foo,bar",
	})

	if err := resourceDevice42IpCreate(d, client); err != nil {
		t.Fatalf("create: %v", err)
	}
	if gotFormTags != "foo,bar" {
		t.Fatalf("POST tags = %q, want %q", gotFormTags, "foo,bar")
	}
	if d.Id() != "50" {
		t.Fatalf("create id = %q, want %q", d.Id(), "50")
	}

	if err := resourceDevice42IpRead(d, client); err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(gotIncludeCols, "tags") {
		t.Fatalf("GET include_cols = %q, want it to contain 'tags'", gotIncludeCols)
	}
	if got := d.Get("ip").(string); got != "10.1.3.1" {
		t.Fatalf("read ip = %q, want %q", got, "10.1.3.1")
	}
	if got := d.Get("tags").(string); got != "foo,bar" {
		t.Fatalf("read tags = %q, want %q", got, "foo,bar")
	}
}
