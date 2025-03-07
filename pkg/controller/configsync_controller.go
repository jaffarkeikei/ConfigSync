package controller

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configsyncv1alpha1 "github.com/configsync/configsync/pkg/apis/configsync/v1alpha1"
	"github.com/configsync/configsync/pkg/git"
)

// ConfigSyncReconciler reconciles a ConfigSync object
type ConfigSyncReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=configsync.io,resources=configsyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=configsync.io,resources=configsyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services;configmaps;secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile handles the reconciliation of ConfigSync resources
func (r *ConfigSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("configsync", req.NamespacedName)
	log.Info("Reconciling ConfigSync")

	// Fetch the ConfigSync instance
	var configSync configsyncv1alpha1.ConfigSync
	if err := r.Get(ctx, req.NamespacedName, &configSync); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		log.Error(err, "Failed to get ConfigSync resource")
		return ctrl.Result{}, err
	}

	// Example implementation - in a real implementation, you would:
	// 1. Clone/fetch the Git repository
	// 2. Check for changes since the last sync
	// 3. Validate the configurations
	// 4. Apply the configurations to the cluster
	// 5. Update the ConfigSync status

	gitClient := git.NewClient(configSync.Spec.GitRepository, configSync.Spec.Branch)

	// Clone or pull the repository
	if err := gitClient.SyncRepository(); err != nil {
		log.Error(err, "Failed to sync Git repository")
		// Update status to reflect error
		return ctrl.Result{}, err
	}

	// Check if there are any changes from the last sync
	hasChanges, commitID, err := gitClient.HasChanges(configSync.Status.LastCommitID)
	if err != nil {
		log.Error(err, "Failed to check for changes in Git repository")
		return ctrl.Result{}, err
	}

	if !hasChanges {
		log.Info("No changes detected since last sync")
		// Requeue after the sync interval
		syncInterval, _ := time.ParseDuration(configSync.Spec.SyncInterval)
		if syncInterval == 0 {
			syncInterval = 5 * time.Minute // Default sync interval
		}
		return ctrl.Result{RequeueAfter: syncInterval}, nil
	}

	log.Info("Changes detected in Git repository", "commitID", commitID)

	// Apply the configurations
	// TODO: Implement the actual application of configurations

	// Update status
	configSync.Status.LastCommitID = commitID
	configSync.Status.LastSyncTime = &configsyncv1alpha1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, &configSync); err != nil {
		log.Error(err, "Failed to update ConfigSync status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *ConfigSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configsyncv1alpha1.ConfigSync{}).
		Complete(r)
}
