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
	"errors"

	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot"
	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot/urtypes"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	uptimerobotv1 "github.com/clevyr/uptime-robot-operator/api/v1"
)

// MonitorReconciler reconciles a Monitor object
type MonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var ErrContactMissingID = errors.New("contact missing ID")

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

	_, apiKey, err := GetApiKey(ctx, r.Client, monitor.Spec.Account.Name)
	if err != nil {
		return ctrl.Result{}, err
	}
	urclient := uptimerobot.NewClient(apiKey)

	myFinalizerName := "uptime-robot.clevyr.com/finalizer"

	if !monitor.DeletionTimestamp.IsZero() {
		// Object is being deleted
		if controllerutil.ContainsFinalizer(monitor, myFinalizerName) {
			if monitor.Spec.Prune && monitor.Status.Ready {
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

	if monitor.Status.Ready && monitor.Status.Type != monitor.Spec.Monitor.Type {
		// Type change requires recreate
		if err := urclient.DeleteMonitor(ctx, monitor.Status.ID); err != nil {
			return ctrl.Result{}, err
		}
		monitor.Status.Ready = false
	}

	contacts := make([]uptimerobot.MonitorContact, 0, len(monitor.Spec.Contacts))
	for _, ref := range monitor.Spec.Contacts {
		contact := &uptimerobotv1.Contact{}

		if err := GetContact(ctx, r.Client, contact, ref.Name); err != nil {
			return ctrl.Result{}, err
		}

		if contact.Status.ID == "" {
			return ctrl.Result{}, ErrContactMissingID
		}

		contacts = append(contacts, uptimerobot.MonitorContact{
			ID:                   contact.Status.ID,
			MonitorContactCommon: ref.MonitorContactCommon,
		})
	}

	if !monitor.Status.Ready {
		id, err := urclient.CreateMonitor(ctx, monitor.Spec.Monitor, contacts)
		if err != nil {
			return ctrl.Result{}, err
		}

		monitor.Status.Ready = true
		monitor.Status.ID = id
		monitor.Status.Type = monitor.Spec.Monitor.Type
		monitor.Status.Status = monitor.Spec.Monitor.Status
		if err := r.Status().Update(ctx, monitor); err != nil {
			return ctrl.Result{}, err
		}

		if monitor.Spec.Monitor.Status == urtypes.MonitorPaused {
			if _, err := urclient.EditMonitor(ctx, id, monitor.Spec.Monitor, contacts); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		id, err := urclient.EditMonitor(ctx, monitor.Status.ID, monitor.Spec.Monitor, contacts)
		if err != nil {
			return ctrl.Result{}, err
		}

		monitor.Status.ID = id
		monitor.Status.Status = monitor.Spec.Monitor.Status
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

	return ctrl.Result{RequeueAfter: monitor.Spec.Interval.Duration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uptimerobotv1.Monitor{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}
