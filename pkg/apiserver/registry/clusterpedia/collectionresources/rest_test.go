package collectionresources

import (
	internal "github.com/clusterpedia-io/api/clusterpedia"
	"k8s.io/apimachinery/pkg/runtime"
	"net/url"
	"testing"

	"github.com/clusterpedia-io/api/clusterpedia/scheme"
	"github.com/clusterpedia-io/api/clusterpedia/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Struct struct {
	metav1.TypeMeta

	LabelSelector string
}

func (l *Struct) DeepCopyObject() runtime.Object {
	return &Struct{LabelSelector: l.LabelSelector}
}

func TestConv(t *testing.T) {
	query := url.Values{"a": []string{"b"}}
	var opts internal.ListOptions
	//var opts Struct
	if err := scheme.ParameterCodec.DecodeParameters(query, v1beta1.SchemeGroupVersion, &opts); err != nil {
		t.Fatal(err)
	}
}
