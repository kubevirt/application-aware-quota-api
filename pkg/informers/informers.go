package informers

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	k6tv1 "kubevirt.io/api/core/v1"
	"kubevirt.io/applications-aware-quota/pkg/client"
	"kubevirt.io/applications-aware-quota/pkg/util"
	v1alpha13 "kubevirt.io/applications-aware-quota/staging/src/kubevirt.io/applications-aware-quota-api/pkg/apis/core/v1alpha1"
	"time"
)

const LauncherLabel = "virt-launcher"

func GetMigrationInformer(aaqCli client.AAQClient) cache.SharedIndexInformer {
	listWatcher := NewListWatchFromClient(aaqCli.KubevirtClient().KubevirtV1().RESTClient(), "virtualmachineinstancemigrations", metav1.NamespaceAll, fields.Everything(), labels.Everything())
	return cache.NewSharedIndexInformer(listWatcher, &k6tv1.VirtualMachineInstanceMigration{}, 1*time.Hour, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func GetApplicationsResourceQuotaInformer(aaqCli client.AAQClient) cache.SharedIndexInformer {
	listWatcher := NewListWatchFromClient(aaqCli.RestClient(), "applicationsresourcequotas", metav1.NamespaceAll, fields.Everything(), labels.Everything())
	return cache.NewSharedIndexInformer(listWatcher, &v1alpha13.ApplicationsResourceQuota{}, 1*time.Hour, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func GetAAQJobQueueConfig(aaqCli client.AAQClient) cache.SharedIndexInformer {
	listWatcher := NewListWatchFromClient(aaqCli.RestClient(), "aaqjobqueueconfigs", metav1.NamespaceAll, fields.Everything(), labels.Everything())
	return cache.NewSharedIndexInformer(listWatcher, &v1alpha13.AAQJobQueueConfig{}, 1*time.Hour, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func GetAAQInformer(aaqCli client.AAQClient) cache.SharedIndexInformer {
	listWatcher := NewListWatchFromClient(aaqCli.RestClient(), "aaqs", metav1.NamespaceAll, fields.Everything(), labels.Everything())
	return cache.NewSharedIndexInformer(listWatcher, &v1alpha13.AAQ{}, 1*time.Hour, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func GetPodInformer(aaqCli client.AAQClient) cache.SharedIndexInformer {
	listWatcher := NewListWatchFromClient(aaqCli.CoreV1().RESTClient(), "pods", metav1.NamespaceAll, fields.Everything(), labels.Everything())
	return cache.NewSharedIndexInformer(listWatcher, &v1.Pod{}, 1*time.Hour, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func GetSecretInformer(aaqCli client.AAQClient, ns string) cache.SharedIndexInformer {
	listWatcher := NewListWatchFromClient(aaqCli.CoreV1().RESTClient(), "secrets", ns, fields.Everything(), labels.Everything())
	return cache.NewSharedIndexInformer(listWatcher, &v1.Secret{}, 1*time.Hour, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func GetVMIInformer(aaqCli client.AAQClient) cache.SharedIndexInformer {
	listWatcher := NewListWatchFromClient(aaqCli.KubevirtClient().KubevirtV1().RESTClient(), "virtualmachineinstances", metav1.NamespaceAll, fields.Everything(), labels.Everything())
	return cache.NewSharedIndexInformer(listWatcher, &k6tv1.VirtualMachineInstance{}, 1*time.Hour, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func GetResourceQuotaInformer(aaqCli client.AAQClient) cache.SharedIndexInformer {
	labelSelector, err := labels.Parse(util.AAQLabel)
	if err != nil {
		panic(err)
	}
	listWatcher := NewListWatchFromClient(aaqCli.CoreV1().RESTClient(), "resourcequotas", metav1.NamespaceAll, fields.Everything(), labelSelector)
	return cache.NewSharedIndexInformer(listWatcher, &v1.ResourceQuota{}, 1*time.Hour, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

// NewListWatchFromClient creates a new ListWatch from the specified client, resource, kubevirtNamespace and field selector.
func NewListWatchFromClient(c cache.Getter, resource string, namespace string, fieldSelector fields.Selector, labelSelector labels.Selector) *cache.ListWatch {
	listFunc := func(options metav1.ListOptions) (runtime.Object, error) {
		options.FieldSelector = fieldSelector.String()
		options.LabelSelector = labelSelector.String()
		return c.Get().
			Namespace(namespace).
			Resource(resource).
			VersionedParams(&options, metav1.ParameterCodec).
			Do(context.Background()).
			Get()
	}
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		options.FieldSelector = fieldSelector.String()
		options.LabelSelector = labelSelector.String()
		options.Watch = true
		return c.Get().
			Namespace(namespace).
			Resource(resource).
			VersionedParams(&options, metav1.ParameterCodec).
			Watch(context.Background())
	}
	return &cache.ListWatch{ListFunc: listFunc, WatchFunc: watchFunc}
}

type fakePodSharedIndexInformer struct {
	indexer cache.Indexer
}

func NewFakePodSharedIndexInformer(podObjs []metav1.Object) cache.SharedIndexInformer {
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	for _, podObj := range podObjs {
		indexer.Add(podObj)
	}
	return fakePodSharedIndexInformer{indexer: indexer}
}

func (i fakePodSharedIndexInformer) RemoveEventHandler(handle cache.ResourceEventHandlerRegistration) error {
	//TODO implement me
	panic("implement me")
}

func (i fakePodSharedIndexInformer) IsStopped() bool {
	panic("implement me")
}

func (i fakePodSharedIndexInformer) AddIndexers(indexers cache.Indexers) error { return nil }
func (i fakePodSharedIndexInformer) GetIndexer() cache.Indexer {
	return i.indexer

}
func (i fakePodSharedIndexInformer) AddEventHandler(handler cache.ResourceEventHandler) (cache.ResourceEventHandlerRegistration, error) {
	return nil, nil
}
func (i fakePodSharedIndexInformer) AddEventHandlerWithResyncPeriod(handler cache.ResourceEventHandler, resyncPeriod time.Duration) (cache.ResourceEventHandlerRegistration, error) {
	return nil, nil
}
func (i fakePodSharedIndexInformer) GetStore() cache.Store           { return nil }
func (i fakePodSharedIndexInformer) GetController() cache.Controller { return nil }
func (i fakePodSharedIndexInformer) Run(stopCh <-chan struct{})      {}
func (i fakePodSharedIndexInformer) HasSynced() bool                 { panic("implement me") }
func (i fakePodSharedIndexInformer) LastSyncResourceVersion() string { return "" }
func (i fakePodSharedIndexInformer) SetWatchErrorHandler(handler cache.WatchErrorHandler) error {
	return nil
}
func (i fakePodSharedIndexInformer) SetTransform(f cache.TransformFunc) error {
	return nil
}