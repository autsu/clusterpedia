package internalstorage

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	internal "github.com/clusterpedia-io/api/clusterpedia"
	"github.com/clusterpedia-io/clusterpedia/pkg/scheme"
)

const (
	CollectionResourceAny           = "any"
	CollectionResourceWorkloads     = "workloads"
	CollectionResourceKubeResources = "kuberesources"
)

var collectionResources = []internal.CollectionResource{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: CollectionResourceAny,
		},
		ResourceTypes: []internal.CollectionResourceType{},
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: CollectionResourceWorkloads,
		},
		ResourceTypes: []internal.CollectionResourceType{
			{
				Group:    "apps",
				Resource: "deployments",
			},
			{
				Group:    "apps",
				Resource: "daemonsets",
			},
			{
				Group:    "apps",
				Resource: "statefulsets",
			},
		},
	},
	{
		// 这里的 ResourceTypes 会在下面的 init 里面填充
		ObjectMeta: metav1.ObjectMeta{
			Name: CollectionResourceKubeResources,
		},
	},
}

func init() {
	groups := sets.NewString()
	// scheme.LegacyResourceScheme 会在 import_known_versions.go 中 register 所有的 k8s 标准资源 struct
	// 通过引入包的方式，调用这些包里面的 init 来完成注册
	for _, groupversion := range scheme.LegacyResourceScheme.PreferredVersionAllGroups() {
		groups.Insert(groupversion.Group)
	}

	// 进一步拿到所有的 group 信息
	types := make([]internal.CollectionResourceType, 0, len(groups))
	for _, group := range groups.List() {
		types = append(types, internal.CollectionResourceType{
			Group: group,
		})
	}

	for i := range collectionResources {
		if collectionResources[i].Name == CollectionResourceKubeResources {
			// 填充 kuberesources 的 ResourceTypes
			collectionResources[i].ResourceTypes = types
		}
	}
}
