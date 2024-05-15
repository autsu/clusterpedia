package cmdb

import "strings"

const (
	jumpcloudNS = "jumpcloud"
	StressENVNS = "stress-env"
)

// 判断是否是业务 namespace
func isBusinessNamespace(namespace string) bool {
	return strings.HasPrefix(namespace, "app-") || namespace == jumpcloudNS || namespace == StressENVNS
}
