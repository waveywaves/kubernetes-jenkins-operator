/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/jenkinsci/kubernetes-operator/pkg/configuration/base/resources"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha2 "github.com/jenkinsci/kubernetes-operator/api/v1alpha2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// BackupWorkerReconciler reconciles a BackupWorker object
type BackupWorkerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=jenkins.io,resources=backupworkers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jenkins.io,resources=backupworkers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=jenkins.io,resources=jenkins,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *BackupWorkerReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	logger := r.Log.WithValues("backupworker", request.NamespacedName)

	// Fetch the Jenkins backupWorker
	backupWorker := &v1alpha2.BackupWorker{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, backupWorker)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	logger.Info("Created BackupWorker with Name " + backupWorker.Name)

	// Fetch the Jenkins Deployment associated to the backupWorker
	jenkinsDep := &appsv1.Deployment{}
	jenkinsRef := backupWorker.Spec.JenkinsRef
	namespacedName := types.NamespacedName{Namespace: request.Namespace, Name: resources.CreateJenkinsDeploymentName(jenkinsRef)}
	err = r.Client.Get(context.TODO(), namespacedName, jenkinsDep)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info(fmt.Sprintf("Couldn't find deployment with name %s ", resources.CreateJenkinsDeploymentName(jenkinsRef)))
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	logger.Info("Getting Jenkins Deployment with Name " + jenkinsDep.Name)

	return ctrl.Result{}, nil
}

func (r *BackupWorkerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha2.BackupWorker{}).
		Complete(r)
}
