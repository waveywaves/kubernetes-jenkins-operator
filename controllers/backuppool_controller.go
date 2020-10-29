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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha2 "github.com/jenkinsci/kubernetes-operator/api/v1alpha2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// BackupPoolReconciler reconciles a BackupPool object
type BackupPoolReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=jenkins.io,resources=backuppools,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jenkins.io,resources=backuppools/status,verbs=get;update;patch

func (r *BackupPoolReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	logger := r.Log.WithValues("backuppool", request.NamespacedName)

	// Fetch the Jenkins backupPool
	backupPool := &v1alpha2.BackupPool{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, backupPool)
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

	logger.Info(fmt.Sprintf("Created BackupPool '%s' of type '%s'", backupPool.Name, backupPool.Spec.Type))

	pvc := r.getBackupPoolPVC(backupPool)
	if backupPool.Spec.Type == v1alpha2.UnifiedPool {
		err := r.Client.Create(context.TODO(), pvc)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	logger.Info(fmt.Sprintf("Created PVC '%s' for BackupPool '%s'", pvc.Name, backupPool.Name))
	backupPool.Status.PersistentVolumeClaim = pvc.Name

	return ctrl.Result{}, nil
}

func (r *BackupPoolReconciler) getBackupPoolPVC(backupPool *v1alpha2.BackupPool) *v1.PersistentVolumeClaim {
	pvc := &v1.PersistentVolumeClaim{}
	pvc.Name = r.getBackupPoolPVCName(backupPool)
	pvc.Namespace = backupPool.Namespace

	pvcSpec := v1.PersistentVolumeClaimSpec{}
	pvcSpec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadWriteMany}
	pvcSpec.Resources = v1.ResourceRequirements{
		Limits: nil,
		Requests: v1.ResourceList{
			v1.ResourceStorage: resource.MustParse("5Gi"),
		},
	}

	pvc.Spec = pvcSpec

	return pvc
}

func (r *BackupPoolReconciler) getBackupPoolPVCName(backupPool *v1alpha2.BackupPool) string {
	return fmt.Sprintf("jenkins-backuppool-%s-%s", backupPool.Spec.Type, backupPool.Name)
}

func (r *BackupPoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha2.BackupPool{}).
		Complete(r)
}
