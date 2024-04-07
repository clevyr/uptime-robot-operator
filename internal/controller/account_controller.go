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
	"fmt"

	"github.com/clevyr/uptime-robot-operator/internal/uptimerobot"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	uptimerobotv1 "github.com/clevyr/uptime-robot-operator/api/v1"
)

var ClusterResourceNamespace string

// AccountReconciler reconciles a Account object
type AccountReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var ErrKeyNotFound = errors.New("secret key not found")

//+kubebuilder:rbac:groups=uptime-robot.clevyr.com,resources=accounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptime-robot.clevyr.com,resources=accounts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptime-robot.clevyr.com,resources=accounts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *AccountReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	account := &uptimerobotv1.Account{}
	if err := r.Client.Get(ctx, req.NamespacedName, account); err != nil {
		return ctrl.Result{}, err
	}

	apiKey, err := GetApiKey(ctx, r.Client, account, req.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	urclient := uptimerobot.NewClient(apiKey)
	if err := urclient.GetAccountDetails(ctx); err != nil {
		account.Status.Ready = false
		if err := r.Status().Update(ctx, account); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, err
	}

	account.Status.Ready = true
	if err := r.Status().Update(ctx, account); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AccountReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &uptimerobotv1.Account{}, "spec.isDefault", func(rawObj client.Object) []string {
		account := rawObj.(*uptimerobotv1.Account)
		if !account.Spec.IsDefault {
			return nil
		}
		return []string{"true"}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&uptimerobotv1.Account{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

var (
	ErrNoDefaultAccount       = errors.New("no default account")
	ErrMultipleDefaultAccount = errors.New("more than 1 default account found")
)

func GetAccount(ctx context.Context, c client.Client, account *uptimerobotv1.Account, name string) error {
	if name != "" {
		return c.Get(ctx, client.ObjectKey{Namespace: name}, account)
	}

	list := &uptimerobotv1.AccountList{}
	err := c.List(ctx, list, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.isDefault", "true"),
	})
	if err != nil {
		return err
	}
	if len(list.Items) == 0 {
		return ErrNoDefaultAccount
	}
	if len(list.Items) > 1 {
		return ErrMultipleDefaultAccount
	}

	*account = list.Items[0]
	return nil
}

func GetApiKey(ctx context.Context, c client.Client, account *uptimerobotv1.Account, name string) (string, error) {
	if account == nil {
		account = &uptimerobotv1.Account{}
		if err := GetAccount(ctx, c, account, name); err != nil {
			return "", err
		}
	}

	secret := &corev1.Secret{}
	err := c.Get(ctx, client.ObjectKey{
		Namespace: ClusterResourceNamespace,
		Name:      account.Spec.ApiKeySecretRef.Name,
	}, secret)
	if err != nil {
		return "", err
	}

	apiKey, ok := secret.Data[account.Spec.ApiKeySecretRef.Key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrKeyNotFound, account.Spec.ApiKeySecretRef.Key)
	}

	return string(apiKey), nil
}
