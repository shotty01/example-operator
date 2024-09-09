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
	testv1beta1 "github.com/shotty01/example-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

// ExampleConfigReconciler reconciles a ExampleConfig object
type ExampleConfigReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

//+kubebuilder:rbac:groups=test.example.com,resources=exampleconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.example.com,resources=exampleconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.example.com,resources=exampleconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ExampleConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ExampleConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := ctrl.Log.WithName("controllers").WithName("ExampleConfig")

	example := &testv1beta1.ExampleConfig{}
	err := r.Get(ctx, req.NamespacedName, example)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("ExampleConfig resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ExampleConfig")
		return ctrl.Result{}, err
	}

	log.Info("ExampleConfig was found with count " + strconv.Itoa(example.Spec.EventCount) + " and message " + example.Spec.MessageContent)

	r.createEvents(example)

	return ctrl.Result{}, nil
}

func (r *ExampleConfigReconciler) createEvents(example *testv1beta1.ExampleConfig) {

	for i := 0; i < example.Spec.EventCount; i++ {
		r.Recorder.Event(example, "Warning", "example", example.Spec.MessageContent+" number "+strconv.Itoa(i))
	}

}

// SetupWithManager sets up the controller with the Manager.
func (r *ExampleConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1beta1.ExampleConfig{}).
		Complete(r)
}
