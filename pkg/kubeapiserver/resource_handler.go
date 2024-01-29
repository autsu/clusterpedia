package kubeapiserver

import (
	"fmt"
	"net/http"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	genericrequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/warning"
	"k8s.io/klog/v2"

	clusterv1alpha2 "github.com/clusterpedia-io/api/cluster/v1alpha2"
	clusterlister "github.com/clusterpedia-io/clusterpedia/pkg/generated/listers/cluster/v1alpha2"
	"github.com/clusterpedia-io/clusterpedia/pkg/kubeapiserver/discovery"
	"github.com/clusterpedia-io/clusterpedia/pkg/utils/request"
)

type ResourceHandler struct {
	minRequestTimeout time.Duration
	delegate          http.Handler

	rest          *RESTManager
	discovery     *discovery.DiscoveryManager
	clusterLister clusterlister.PediaClusterLister
}

// 查询貌似主要是走这个函数
func (r *ResourceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// /apis/clusterpedia.io/v1beta1/resources/apis/apps/v1/deployments		-> List
	// /apis/clusterpedia.io/v1beta1/resources/apis/apps/v1/deployments/nginx  -> Get
	// 调用 RequestInfoFrom 拿到请求信息，比如可以根据请求格式，拿到 verb
	requestInfo, ok := genericrequest.RequestInfoFrom(req.Context())
	if !ok {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			Codecs, schema.GroupVersion{}, w, req,
		)
		return
	}

	// handle discovery request
	if !requestInfo.IsResourceRequest {
		r.discovery.ServeHTTP(w, req)
		return
	}

	gvr := schema.GroupVersionResource{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion, Resource: requestInfo.Resource}

	var (
		cluster *clusterv1alpha2.PediaCluster
		err     error
	)
	// When clusterName not empty, first check cluster whether exist
	// 可以看到查询的 cluster 实际是根据 PediaCluster 的名字来查的
	clusterName := request.ClusterNameValue(req.Context())
	if clusterName != "" {
		if cluster, err = r.clusterLister.Get(clusterName); err != nil {
			if !apierrors.IsNotFound(err) {
				klog.ErrorS(err, "Failed to handle resource request, not get cluster from cache", "cluster", clusterName, "resource", gvr)
				responsewriters.ErrorNegotiated(
					apierrors.NewInternalError(err),
					Codecs, gvr.GroupVersion(), w, req,
				)
				return
			}
			responsewriters.ErrorNegotiated(
				apierrors.NewBadRequest("the server could not find the requested cluster"),
				Codecs, gvr.GroupVersion(), w, req,
			)
			return
		}
	}
	if !r.discovery.ResourceEnabled(clusterName, gvr) {
		r.delegate.ServeHTTP(w, req)
		return
	}

	info := r.rest.GetRESTResourceInfo(gvr)
	if info.Empty() {
		err := fmt.Errorf("not found request scope or resource storage")
		klog.ErrorS(err, "Failed to handle resource request", "resource", gvr)
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(err),
			Codecs, gvr.GroupVersion(), w, req,
		)
		return
	}

	resource, reqScope, storage := info.APIResource, info.RequestScope, info.Storage
	if requestInfo.Namespace != "" && !resource.Namespaced {
		r.delegate.ServeHTTP(w, req)
		return
	}

	// Check the health of the cluster
	// 如果查询条件指定了 cluster，那么会查一下该 cluster 的健康状态
	// 如果不健康，会添加一个 warning
	if cluster != nil {
		var msg string
		healthyCondition := meta.FindStatusCondition(cluster.Status.Conditions, clusterv1alpha2.ClusterHealthyCondition)
		switch {
		case healthyCondition == nil:
			msg = fmt.Sprintf("%s is not ready and the resources obtained may be inaccurate.", clusterName)
		case healthyCondition.Status != metav1.ConditionTrue:
			msg = fmt.Sprintf("%s is not ready and the resources obtained may be inaccurate, reason: %s", clusterName, healthyCondition.Reason)
		}
		/*
			TODO(scyda): Determine the synchronization status of a specific resource

			for _, resource := range c.Status.Resources {
			}
		*/

		if msg != "" {
			warning.AddWarning(req.Context(), "", msg)
		}
	}

	var handler http.Handler
	switch requestInfo.Verb {
	case "get":
		if clusterName == "" {
			responsewriters.ErrorNegotiated(
				apierrors.NewBadRequest("please specify the cluster name when using the resource name to get a specific resource."),
				Codecs, gvr.GroupVersion(), w, req,
			)
			return
		}
		// 这里的 GetResource 是 k8s 标准库里面的函数，它会进一步调用 storage.Get 方法
		handler = handlers.GetResource(storage, reqScope)
	case "list":
		// 这里的 ListResource 是 k8s 标准库里面的函数，它会进一步调用 storage.List 方法
		handler = handlers.ListResource(storage, nil, reqScope, false, r.minRequestTimeout)
	case "watch":
		handler = handlers.ListResource(storage, storage, reqScope, true, r.minRequestTimeout)
	default:
		responsewriters.ErrorNegotiated(
			apierrors.NewMethodNotSupported(gvr.GroupResource(), requestInfo.Verb),
			Codecs, gvr.GroupVersion(), w, req,
		)
	}

	if handler != nil {
		handler.ServeHTTP(w, req)
	}
}
