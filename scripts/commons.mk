# Current Operator version
VERSION ?= 0.6.0
GIT_COMMIT_ID ?= $(shell git rev-parse --short HEAD)
# Default bundle image tag
BUNDLE_IMG ?= quay.io/redhat-developer/openshift-jenkins-operator-bundle:$(VERSION)
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Image URL to use all building/pushing image targets
IMG ?= quay.io/redhat-developer/openshift-jenkins-operator:$(VERSION)-$(GIT_COMMIT_ID)
# Image URL to use for building/pushing sidecar image targets
SIDECAR_IMG ?= docker.io/waveywaves/jenkins-sidecar:$(VERSION)-$(GIT_COMMIT_ID)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

XDG_CACHE_HOME := $(shell pwd)/build/_output/xdgcache
GOOS := $(shell go env GOOS)
GOARCH :=  $(shell go env GOARCH)
