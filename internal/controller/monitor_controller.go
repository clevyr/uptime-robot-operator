/*
Copyright 2024.

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

package controller

import (
	"context"

	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	uptimerobotv1 "github.com/clevyr/uptime-robot-operator/api/v1"
)

// MonitorReconciler reconciles a Monitor object
type MonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=uptime-robot.clevyr.com,resources=monitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptime-robot.clevyr.com,resources=monitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptime-robot.clevyr.com,resources=monitors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *MonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	monitor := &uptimerobotv1.Monitor{}
	if err := r.Client.Get(ctx, req.NamespacedName, monitor); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	myFinalizerName := "uptime-robot.clevyr.com/finalizer"

	if !monitor.DeletionTimestamp.IsZero() {
		// Object is being deleted
		if controllerutil.ContainsFinalizer(monitor, myFinalizerName) {
			if monitor.Spec.Prune && monitor.Status.Created {
				urclient := uptimerobot.NewClient()
				if err := urclient.DeleteMonitor(ctx, monitor.Status.ID); err != nil {
					return ctrl.Result{}, err
				}
			}

			controllerutil.RemoveFinalizer(monitor, myFinalizerName)
			if err := r.Update(ctx, monitor); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	urclient := uptimerobot.NewClient()

	if monitor.Status.Created && monitor.Status.Type != monitor.Spec.Monitor.Type {
		// Type change requires recreate
		if err := urclient.DeleteMonitor(ctx, monitor.Status.ID); err != nil {
			return ctrl.Result{}, err
		}
		monitor.Status.Created = false
	}

	if !monitor.Status.Created {
		id, err := urclient.CreateMonitor(ctx, monitor.Spec.Monitor)
		if err != nil {
			return ctrl.Result{}, err
		}

		monitor.Status.Created = true
		monitor.Status.ID = id
		monitor.Status.Type = monitor.Spec.Monitor.Type
		if err := r.Status().Update(ctx, monitor); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		id, err := urclient.EditMonitor(ctx, monitor.Status.ID, monitor.Spec.Monitor)
		if err != nil {
			return ctrl.Result{}, err
		}

		if id != monitor.Status.ID {
			monitor.Status.ID = id
		}
		if err := r.Status().Update(ctx, monitor); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !controllerutil.ContainsFinalizer(monitor, myFinalizerName) {
		controllerutil.AddFinalizer(monitor, myFinalizerName)
		if err := r.Update(ctx, monitor); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uptimerobotv1.Monitor{}).
		Complete(r)
}
