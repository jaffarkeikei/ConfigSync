package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	// Import controllers
	// "github.com/configsync/configsync/pkg/controller"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = log.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	// Add custom API types to scheme
	// _ = configsyncv1alpha1.AddToScheme(scheme)
}

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
		gitRepo              string
		configDir            string
		developmentMode      bool
	)

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&gitRepo, "git-repo", "", "URL to the Git repository containing configurations.")
	flag.StringVar(&configDir, "config-dir", "configs", "Directory within the Git repository containing environment configurations.")
	flag.BoolVar(&developmentMode, "dev", false, "Run in development mode.")

	opts := zap.Options{
		Development: developmentMode,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	log.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	restConfig, err := config.GetConfig()
	if err != nil {
		setupLog.Error(err, "unable to get Kubernetes config")
		os.Exit(1)
	}

	mgr, err := manager.New(restConfig, manager.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "configsync-leader-election",
	})
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		os.Exit(1)
	}

	// Register controllers
	// if err := controller.SetupControllers(mgr, gitRepo, configDir); err != nil {
	// 	setupLog.Error(err, "unable to register controllers")
	// 	os.Exit(1)
	// }

	setupLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager", "error", err.Error())
		os.Exit(1)
	}
}
