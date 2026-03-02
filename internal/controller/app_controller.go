package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	portrav1 "github.com/portra-run/portra-operator/api/v1"
)

type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.portra.run,resources=apps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.portra.run,resources=apps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.portra.run,resources=apps/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	app := &portrav1.App{}
	err := r.Get(ctx, req.NamespacedName, app)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if err := r.reconcileDeployment(ctx, app); err != nil {
		log.Error(err, "Failed to reconcile Deployment")
		return ctrl.Result{}, err
	}

	if err := r.reconcileService(ctx, app); err != nil {
		log.Error(err, "Failed to reconcile Service")
		return ctrl.Result{}, err
	}

	if err := r.reconcileIngress(ctx, app); err != nil {
		log.Error(err, "Failed to reconcile Ingress")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *AppReconciler) reconcileDeployment(ctx context.Context, app *portrav1.App) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}

	_, err := ctrl.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		replicas := int32(1)
		if app.Spec.Replicas != nil {
			replicas = *app.Spec.Replicas
		}

		labels := map[string]string{"app": app.Name}
		deployment.Spec = appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "app",
						Image: app.Spec.Image,
						Ports: []corev1.ContainerPort{{
							ContainerPort: app.Spec.ContainerPort,
						}},
						Env:       app.Spec.Env,
						Resources: app.Spec.Resources,
					}},
				},
			},
		}

		if err := ctrl.SetControllerReference(app, deployment, r.Scheme); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *AppReconciler) reconcileService(ctx context.Context, app *portrav1.App) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}

	_, err := ctrl.CreateOrUpdate(ctx, r.Client, service, func() error {
		service.Spec = corev1.ServiceSpec{
			Selector: map[string]string{"app": app.Name},
			Ports: []corev1.ServicePort{{
				Port:       app.Spec.ContainerPort,
				TargetPort: intstr.FromInt32(app.Spec.ContainerPort),
			}},
			Type: corev1.ServiceTypeClusterIP,
		}

		if err := ctrl.SetControllerReference(app, service, r.Scheme); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *AppReconciler) reconcileIngress(ctx context.Context, app *portrav1.App) error {
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}

	if len(app.Spec.Domains) == 0 {
		err := r.Delete(ctx, ingress)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	_, err := ctrl.CreateOrUpdate(ctx, r.Client, ingress, func() error {
		if ingress.Annotations == nil {
			ingress.Annotations = make(map[string]string)
		}

		var rules []networkingv1.IngressRule
		for _, domain := range app.Spec.Domains {
			rules = append(rules, networkingv1.IngressRule{
				Host: domain,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{{
							Path:     "/",
							PathType: func() *networkingv1.PathType { pt := networkingv1.PathTypePrefix; return &pt }(),
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: app.Name,
									Port: networkingv1.ServiceBackendPort{
										Number: app.Spec.ContainerPort,
									},
								},
							},
						}},
					},
				},
			})
		}

		ingress.Spec.Rules = rules
		ingress.Spec.IngressClassName = func() *string { s := "nginx"; return &s }()

		if app.Spec.TLS {
			ingress.Annotations["cert-manager.io/cluster-issuer"] = "letsencrypt-prod"
			var tls []networkingv1.IngressTLS
			tls = append(tls, networkingv1.IngressTLS{
				Hosts:      app.Spec.Domains,
				SecretName: app.Name + "-tls",
			})
			ingress.Spec.TLS = tls
		} else {
			delete(ingress.Annotations, "cert-manager.io/cluster-issuer")
			ingress.Spec.TLS = nil
		}

		if err := ctrl.SetControllerReference(app, ingress, r.Scheme); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&portrav1.App{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Named("app").
		Complete(r)
}
