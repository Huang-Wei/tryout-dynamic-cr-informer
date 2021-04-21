package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	masterURL  string
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s.", err.Error())
	}

	dynClient := dynamic.NewForConfigOrDie(cfg)
	f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynClient, 0, v1.NamespaceAll, nil)

	klog.Info("Started.")

	// Setup watch for CRs: "foo" and "bar".
	stopCh := make(chan struct{})
	for _, s := range []string{"foo", "bar"} {
		gvr, _ := schema.ParseResourceArg(fmt.Sprintf("%vs.v1.example.com", s))
		i := f.ForResource(*gvr)
		go startWatching(s, stopCh, i.Informer())
	}

	sigCh := make(chan os.Signal, 0)
	signal.Notify(sigCh, os.Kill, os.Interrupt)
	<-sigCh
	close(stopCh)
	klog.Info("Stopped.")
}

func startWatching(kind string, stopCh <-chan struct{}, s cache.SharedIndexInformer) {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			klog.Infof("received %q add event!", kind)
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			klog.Infof("received %q update event!", kind)
		},
		DeleteFunc: func(obj interface{}) {
			klog.Infof("received %q delete event!", kind)
		},
	}
	s.AddEventHandler(handlers)
	s.Run(stopCh)
}
