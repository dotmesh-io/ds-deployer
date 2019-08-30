package controller

import (
	"time"

	// "github.com/rusenask/client/pkg/workgroup"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/dotmesh-io/ds-deployer/pkg/workgroup"
)

// WatchServices creates a SharedInformer for v1.Services and registers it with g.
func WatchServices(g *workgroup.Group, client *kubernetes.Clientset, l *zap.SugaredLogger, rs ...cache.ResourceEventHandler) {
	watch(g, client.CoreV1().RESTClient(), l, "services", new(v1.Service), rs...)
}

// WatchDEployments creates a SharedInformer for appsv1.Deployments and registers it with g.
func WatchDeployments(g *workgroup.Group, client *kubernetes.Clientset, l *zap.SugaredLogger, rs ...cache.ResourceEventHandler) {
	watch(g, client.AppsV1().RESTClient(), l, "deployments", new(appsv1.Deployment), rs...)
}

// WatchIngress creates a SharedInformer for v1beta1.Ingress and registers it with g.
func WatchIngress(g *workgroup.Group, client *kubernetes.Clientset, l *zap.SugaredLogger, rs ...cache.ResourceEventHandler) {
	watch(g, client.ExtensionsV1beta1().RESTClient(), l, "ingresses", new(v1beta1.Ingress), rs...)
}

func watch(g *workgroup.Group, c cache.Getter, l *zap.SugaredLogger, resource string, objType runtime.Object, rs ...cache.ResourceEventHandler) {
	lw := cache.NewListWatchFromClient(c, resource, v1.NamespaceAll, fields.Everything())
	sw := cache.NewSharedInformer(lw, objType, 30*time.Minute)
	for _, r := range rs {
		sw.AddEventHandler(r)
	}

	g.Add(func(stop <-chan struct{}) error {
		l := l.With("watch", resource)

		l.Infow("started")
		defer l.Infow("stopped")
		sw.Run(stop)

		return nil
	})
}

type buffer struct {
	ev chan interface{}
	l  *zap.SugaredLogger
	rh cache.ResourceEventHandler
}

type addEvent struct {
	obj interface{}
}

type updateEvent struct {
	oldObj, newObj interface{}
}

type deleteEvent struct {
	obj interface{}
}

// NewBuffer returns a ResourceEventHandler which buffers and serialises ResourceEventHandler events.
func NewBuffer(g *workgroup.Group, rh cache.ResourceEventHandler, l *zap.SugaredLogger, size int) cache.ResourceEventHandler {
	buf := &buffer{
		ev: make(chan interface{}, size),
		l:  l,
		rh: rh,
	}
	g.Add(buf.loop)
	return buf
}

func (b *buffer) loop(stop <-chan struct{}) error {
	log := b.l.With("module", "buffer")
	log.Infof("started")
	defer log.Infof("stopped")

	for ev := range b.ev {
		switch ev := ev.(type) {
		case *addEvent:
			b.rh.OnAdd(ev.obj)
		case *updateEvent:
			b.rh.OnUpdate(ev.oldObj, ev.newObj)
		case *deleteEvent:
			b.rh.OnDelete(ev.obj)
		default:
			log.Errorf("unhandled event type: %T: %v", ev, ev)
		}
	}
	return nil
}

func (b *buffer) OnAdd(obj interface{}) {
	b.send(&addEvent{obj})
}

func (b *buffer) OnUpdate(oldObj, newObj interface{}) {
	b.send(&updateEvent{oldObj, newObj})
}

func (b *buffer) OnDelete(obj interface{}) {
	b.send(&deleteEvent{obj})
}

func (b *buffer) send(ev interface{}) {
	select {
	case b.ev <- ev:
		// nothing to do
	default:
		b.l.Infof("event channel is full, len: %v, cap: %v", len(b.ev), cap(b.ev))
		b.ev <- ev
	}
}
