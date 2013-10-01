package gomws

import "net/url"
import "testing"
import "strings"
import "time"
import "encoding/json"
import "os"
import "net/http"
import "io/ioutil"

func TestCanonicalizing(t *testing.T) {
	values := make(url.Values)
	values.Add("b", " ")
	values.Add("a", "$")
	result := CanonicalizedQueryString(values)
	if strings.Index(result, "a=") > strings.Index(result, "b=") {
		t.Fatalf("strings not Sorted")
	}
	if strings.ContainsAny(result, "+") {
		t.Fatalf("spaces improperly encoded")
	}
	if !strings.EqualFold(XMLTimestamp(time.Time{}), "0001-01-01T00:00:00Z") {
		t.Fatalf("Bad format")
	}
	fname := os.Getenv("GOMWS_CONFIG")
	if len(fname) == 0 {
		fname = ".gomws_config"
	}
	f, err := os.Open(fname)
	if err != nil {
		t.Fatal("Couldn't find an aws config file. you need one.")

	}
	var c Creds
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		t.Fatal(err)
	}
	mws := NewClient(Client{
		Creds:       c,
		CompanyName: "mws unit test",
		Method:      "POST",
		Action:      "GetReportList",
	})
	client := http.Client{}
	r, err := mws.Request()
	if err != nil {
		t.Fatal(err)
	}
	if len(os.Getenv("GOMWS_INTEGRATION")) == 0 {
		return
	}
	resp, err := client.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		t.Fatal("Bad status")
	}
	os.Stdout.Write(body)
}
