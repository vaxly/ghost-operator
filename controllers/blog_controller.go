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
	"time"

	logger "github.com/sirupsen/logrus"

	networkingv1beta1 "k8s.io/api/networking/v1beta1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ghostv1alpha1 "github.com/vaxly/ghost-operator/api/v1alpha1"
)

// BlogReconciler reconciles a Blog object
type BlogReconciler struct {
	client.Client
	Logger *logger.Entry
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ghost.vaxly.io,resources=blogs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ghost.vaxly.io,resources=blogs/status,verbs=get;update;patch

func (r *BlogReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	instance := &ghostv1alpha1.Blog{}
	if err := r.Client.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if err := r.MysqlCreateOrUpdatePersistentVolumeClaim(instance); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.MyqlsCreateOrUpdateService(instance); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.MysqlCreateOrUpdateDeployment(instance); err != nil {
		return reconcile.Result{}, err
	}

	if !r.ISMysqlUp(instance) {
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	if err := r.CreateOrUpdateConfigMap(instance); err != nil {
		return reconcile.Result{}, err
	}

	if instance.Spec.Persistent.Enabled {
		if err := r.CreateOrUpdatePersistentVolumeClaim(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	if err := r.CreateOrUpdateDeployment(instance); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.CreateOrUpdateService(instance); err != nil {
		return reconcile.Result{}, err
	}

	if instance.Spec.Ingress.Enabled {
		if err := r.CreateOrUpdateIngress(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	instance.Status.Replicas = *instance.Spec.Replicas
	if err := r.Client.Status().Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *BlogReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ghostv1alpha1.Blog{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&networkingv1beta1.Ingress{}).
		Complete(r)
}
