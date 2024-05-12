package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	customappv1alpha1 "github.com/HarshitDawar55/kubernetes-basic-operator/api/v1alpha1"
)

// CustomServiceReconciler reconciles a CustomService object
type CustomServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=customapp.harshitdawar.com,resources=customservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=customapp.harshitdawar.com,resources=customservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=customapp.harshitdawar.com,resources=customservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *CustomServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the CustomService instance
	var customService customappv1alpha1.CustomService
	if err := r.Get(ctx, req.NamespacedName, &customService); err != nil {
		log.Error(err, "unable to fetch CustomService")
		return ctrl.Result{}, client.IgnoreNotFound(err) // Ignore not-found errors
	}

	// Define the desired Deployment object
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      customService.Name,
			Namespace: customService.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &customService.Spec.Replicas, // Assuming Replicas field in spec
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": customService.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": customService.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "nginx",
						Image: "nginx:latest", // Replace with your image
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
						}},
					}},
				},
			},
		},
	}

	// Set CustomService instance as the owner and controller
	if err := ctrl.SetControllerReference(&customService, deployment, r.Scheme); err != nil {
		log.Error(err, "unable to set controller reference")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&customappv1alpha1.CustomService{}).
		Complete(r)
}

