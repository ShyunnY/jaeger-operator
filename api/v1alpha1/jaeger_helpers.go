package v1alpha1

// Jaeger Operator Components
const (
	// ReconcilerComponent Define the Reconciler component
	ReconcilerComponent = "reconciler"

	// TranslatorComponent Define the Translator component
	TranslatorComponent = "translator"

	// StatusComponent Define the Status component
	StatusComponent = "status"

	// InfrastructureComponent Define the Infrastructure component
	InfrastructureComponent = "infrastructure"
)

// Jaeger Operator Common Labels Key
const (
	OperatorByLabelKey = "jaegertracing.orange.io/operated-by"
)
