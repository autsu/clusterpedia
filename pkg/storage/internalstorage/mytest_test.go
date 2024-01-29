package internalstorage

import (
	"k8s.io/apiserver/pkg/server"
	"net/http"
	"net/url"
	"testing"
)

func TestJSONQuery(t *testing.T) {
	in := JSONQuery("object", "metadata", "labels").In("123", "456")
	t.Logf("%+v\n", in)
}

func TestName(t *testing.T) {
	//urlv, _ := url.Parse("/apis/clusterpedia.io/v1beta1/resources/apis/apps/v1/deployments")
	//urlv, _ := url.Parse("/apis/clusterpedia.io/v1beta1/resources/apis")
	urlv, _ := url.Parse("/apis/group/v1/resources/")
	info, err := server.NewRequestInfoResolver(&server.Config{}).NewRequestInfo(&http.Request{URL: urlv, Method: http.MethodGet})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}
