package main

import (
	"fmt"
	"github.com/jenkinsci/kubernetes-operator/api/v1alpha2"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
	"os"
)

func onUpdate(new v1alpha2.JCasCProfile) (string, error) {
	podName := os.Getenv("POD_NAME")
	postUrl := fmt.Sprintf("http://localhost:8080/reload-configuration-as-code/?casc-reload-token=%s", podName)
	resp, err := http.Post(postUrl, "application/json", nil)
	if err != nil {
		klog.Fatal(err)
		return resp.Status, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			klog.Fatal(err)
			return resp.Status, err
		}
		bodyString := string(bodyBytes)
		klog.Infof("OnUpdate : POST Request to Jenkins for casc reload. Status: %v , Body: %s", resp.Status, bodyString)
	}

	klog.Infof("JCasCProfile has been updated to %+v", new)

	return resp.Status, nil
}
