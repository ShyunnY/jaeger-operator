package status

import (
	"time"

	jaegerv1a1 "github.com/ShyunnY/jaeger-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetJaegerCondition Add a new condition to Jaeger
func SetJaegerCondition(instance *jaegerv1a1.Jaeger, t string, status metav1.ConditionStatus, reason, msg string) {
	cond := newCondition(t, status, reason, msg, time.Now(), instance.Generation)
	if instance.Status.Conditions == nil {
		instance.Status.Conditions = []metav1.Condition{}
	}
	instance.Status.Conditions = MergeCondition(instance.Status.Conditions, cond)
}

// MergeCondition
// Given a conditions and updates, if the condition in updates is of the same type
// as some condition in conditions, it is updated
func MergeCondition(conditions []metav1.Condition, updates ...metav1.Condition) []metav1.Condition {

	var additions []metav1.Condition

	for i, update := range updates {
		add := true
		for j, condition := range conditions {
			// Reason represents different components, if the same component publishes the same condition,
			// we don't need to add additional conditions
			if condition.Reason == update.Reason {
				add = false

				// We change condition to the latest
				if conditionChange(condition, update) {
					conditions[j].Status = update.Status
					conditions[j].Reason = update.Reason
					conditions[j].Message = update.Message
					conditions[j].ObservedGeneration = update.ObservedGeneration
					conditions[j].LastTransitionTime = update.LastTransitionTime
					break
				}
			}
		}

		if add {
			additions = append(additions, updates[i])
		}
	}

	conditions = append(conditions, additions...)
	return conditions
}

func newCondition(t string, status metav1.ConditionStatus, reason, msg string, lt time.Time, og int64) metav1.Condition {
	return metav1.Condition{
		Type:               t,
		Status:             status,
		Reason:             reason,
		Message:            msg,
		LastTransitionTime: metav1.NewTime(lt),
		ObservedGeneration: og,
	}
}

func conditionChange(a, b metav1.Condition) bool {
	return (a.Status != b.Status) ||
		(a.Message != b.Message) ||
		(a.ObservedGeneration != b.ObservedGeneration)
}
