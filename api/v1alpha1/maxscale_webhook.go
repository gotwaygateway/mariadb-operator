package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var maxscaleLogger = logf.Log.WithName("maxscale")

func (r *MaxScale) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-mariadb-mmontes-io-v1alpha1-maxscale,mutating=false,failurePolicy=fail,sideEffects=None,groups=mariadb.mmontes.io,resources=maxscales,verbs=create;update,versions=v1alpha1,name=vmaxscale.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &MaxScale{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MaxScale) ValidateCreate() (admission.Warnings, error) {
	maxscaleLogger.V(1).Info("Validate create", "name", r.Name)
	return nil, r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MaxScale) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	maxscaleLogger.V(1).Info("Validate update", "name", r.Name)
	oldMaxScale := old.(*MaxScale)
	if err := inmutableWebhook.ValidateUpdate(r, oldMaxScale); err != nil {
		return nil, err
	}
	return nil, r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MaxScale) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (r *MaxScale) validate() error {
	validateFns := []func() error{
		r.validateServers,
		r.validateServices,
		r.validatePodDisruptionBudget,
	}
	for _, fn := range validateFns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (r *MaxScale) validateServers() error {
	idx := r.ServerIndex()
	if len(idx) != len(r.Spec.Servers) {
		return field.Invalid(
			field.NewPath("spec").Child("servers"),
			r.Spec.Servers,
			"server names must be unique",
		)
	}
	addresses := make(map[string]struct{})
	for _, srv := range r.Spec.Servers {
		addresses[srv.Address] = struct{}{}
	}
	if len(addresses) != len(r.Spec.Servers) {
		return field.Invalid(
			field.NewPath("spec").Child("servers"),
			r.Spec.Servers,
			"server addresses must be unique",
		)
	}
	return nil
}

func (r *MaxScale) validateServices() error {
	idx := r.ServiceIndex()
	if len(idx) != len(r.Spec.Services) {
		return field.Invalid(
			field.NewPath("spec").Child("services"),
			r.Spec.Services,
			"service names must be unique",
		)
	}
	ports := make(map[int]struct{})
	for _, svc := range r.Spec.Services {
		ports[int(svc.Listener.Port)] = struct{}{}
	}
	if len(ports) != len(r.Spec.Services) {
		return field.Invalid(
			field.NewPath("spec").Child("services"),
			r.Spec.Services,
			"service listener ports must be unique",
		)
	}
	return nil
}

func (r *MaxScale) validatePodDisruptionBudget() error {
	if r.Spec.PodDisruptionBudget == nil {
		return nil
	}
	if err := r.Spec.PodDisruptionBudget.Validate(); err != nil {
		return field.Invalid(
			field.NewPath("spec").Child("podDisruptionBudget"),
			r.Spec.PodDisruptionBudget,
			err.Error(),
		)
	}
	return nil
}