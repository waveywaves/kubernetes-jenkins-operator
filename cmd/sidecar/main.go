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

package main

import (
	"flag"
	"fmt"
	"os"
	currentruntime "runtime"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	jenkinsv1alpha2 "github.com/jenkinsci/kubernetes-operator/api/v1alpha2"
	sidecarv1alpha2 "github.com/jenkinsci/kubernetes-operator/cmd/sidecar/api/v1alpha2"
	"github.com/jenkinsci/kubernetes-operator/controllers"
	"github.com/jenkinsci/kubernetes-operator/version"

	// sdkVersion "github.com/operator-framework/operator-sdk/version"
	kzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	watchNamespaceEnvVar = "WATCH_NAMESPACE"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(jenkinsv1alpha2.AddToScheme(scheme))
	utilruntime.Must(sidecarv1alpha2.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	parsePglags(metricsAddr, enableLeaderElection)

	setupLog.Info("Registering Components.")
	manager := initManager(metricsAddr, enableLeaderElection)

	// setup Backup controller
	setupBackupReconciler(manager)

	// start the Cmd
	setupLog.Info("Starting the Cmd.")
	runMananger(manager)
	// +kubebuilder:scaffold:builder
}

func getRestClient(debug *bool) *rest.Config {
	// get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		fatal(errors.Wrap(err, "failed to get config"), *debug)
	}

	return cfg
}

func parsePglags(metricsAddr string, enableLeaderElection bool) {

	flag.StringVar(&metricsAddr, "metrics-addr", ":8081", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()
	ctrl.SetLogger(kzap.New(kzap.UseDevMode(true)))
}

func initManager(metricsAddr string, enableLeaderElection bool) manager.Manager {
	printInfo()
	mgr, err := startManager(metricsAddr, enableLeaderElection)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	return mgr
}

func runMananger(mgr manager.Manager) {
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
	setupLog.Info("manager started")
}

func startManager(metricsAddr string, enableLeaderElection bool) (manager.Manager, error) {
	// Jenkins is listening on 8080, hence binding metrics to 8081
	metricsBindAddress := ":8081"
	options := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsBindAddress,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "9cf053ac.jenkins.io",
		Namespace:          getWatchNamespace(), // namespaced-scope when the value is not an empty string
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	return mgr, err
}

// getWatchNamespace returns the Namespace the operator should be watching for changes
func getWatchNamespace() string {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	ns, _ := os.LookupEnv(watchNamespaceEnvVar)
	return ns
}

func setupBackupReconciler(mgr manager.Manager) {
	if err := newBackupReconciler(mgr).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Jenkins")
		os.Exit(1)
	}
}

func newBackupReconciler(mgr manager.Manager) *controllers.BackupReconciler {
	return &controllers.BackupReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Backup"),
		Scheme: mgr.GetScheme(),
	}
}

func fatal(err error, debug bool) {
	if debug {
		setupLog.Error(nil, fmt.Sprintf("%+v", err))
	} else {
		setupLog.Error(nil, fmt.Sprintf("%s", err))
	}
	os.Exit(-1)
}

func printInfo() {
	setupLog.Info(fmt.Sprintf("Version: %s", version.Version))
	setupLog.Info(fmt.Sprintf("Git commit: %s", version.GitCommit))
	setupLog.Info(fmt.Sprintf("Go Version: %s", currentruntime.Version()))
	setupLog.Info(fmt.Sprintf("Go OS/Arch: %s/%s", currentruntime.GOOS, currentruntime.GOARCH))
}
