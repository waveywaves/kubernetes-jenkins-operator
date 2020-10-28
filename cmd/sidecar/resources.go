package main

import (
	"io/ioutil"
	"k8s.io/klog"
)

var (
	currentNamespace = ""
)

func getCurrentNamespace() string {
	if currentNamespace == "" {
		currentNamespaceInBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			klog.Fatal(err)
		}
		currentNamespace = string(currentNamespaceInBytes)
	}
	return currentNamespace
}
