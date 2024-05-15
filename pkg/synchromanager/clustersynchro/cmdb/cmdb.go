package cmdb

import (
	"github.com/IBM/sarama"
	"github.com/clusterpedia-io/api/cluster/v1alpha2"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ZoneInfo struct {
	ClusterType string
	RegionID    string
}

func PutResource(pc *v1alpha2.PediaCluster, kafka sarama.SyncProducer, obj runtime.Object) error {
	kind := obj.GetObjectKind().GroupVersionKind()
	switch kind.GroupKind().Kind {
	case "Deployment":
		dep, ok := obj.(*appsv1.Deployment)
		if !ok {

		}
		if !isBusinessNamespace(dep.Namespace) {
			return nil
		}
		// 后端数据库不存在的 namespace，数据不再推送
		// TODO: 这个需要查后端库去判断，怎么做比较好？
		// 1. 调 backend 的接口（貌似还没有这个接口，不过好写）判断，需要引入 backend addr 的环境变量
		// 2. 在 clusterpedia 查表（感觉不好）

		zoneID := pc.Labels["zoneId"]
		zoneName := pc.Labels["zoneName"]
		clusterType := pc.Labels["clusterType"]
		regionID := pc.Labels["regionID"]

		kafka.SendMessage()
	case "Pod":
		
	}

}
