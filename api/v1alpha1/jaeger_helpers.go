package v1alpha1

// Jaeger Operator Components
const (
	// ReconcilerComponent Define the Reconciler component
	ReconcilerComponent = "reconciler"

	// TranslatorComponent Define the Translator component
	TranslatorComponent = "translator"

	// StatusComponent Define the Status component
	StatusComponent = "status-manager"

	// InfrastructureComponent Define the Infrastructure component
	InfrastructureComponent = "infrastructure"
)

const (
	// OperatorByLabelKey Operator Common Labels Key
	OperatorByLabelKey = "jaegertracing.orange.io/operated-by"

	ServiceTargetLabelKey = "tracing.orange.io/service-target"
)

func (j *Jaeger) EnableHTTPRoute() bool {
	return j.Spec.CommonSpec.HTTPRoute != nil &&
		len(j.Spec.CommonSpec.HTTPRoute) != 0
}

func (j *Jaeger) GetDeploymentType() DeploymentType {
	return j.Spec.Type
}
