package scheme

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	unstructuredscheme "github.com/clusterpedia-io/clusterpedia/pkg/scheme/unstructured"
)

var (
	// LegacyResourceScheme 在 import_known_versions.go 中通过引入包的方式，将所有的 k8s 标准
	// 资源注册进了这里
	LegacyResourceScheme         = legacyscheme.Scheme
	LegacyResourceCodecs         = legacyscheme.Codecs
	LegacyResourceParameterCodec = legacyscheme.ParameterCodec

	UnstructuredScheme = unstructuredscheme.NewScheme()
	UnstructuredCodecs = unstructured.UnstructuredJSONScheme
)
