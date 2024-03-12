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
// 给定一个conditions和updates, 如果updates中的condition与conditions中某些条件类型相同, 则进行更新
func MergeCondition(conditions []metav1.Condition, updates ...metav1.Condition) []metav1.Condition {

	var additions []metav1.Condition

	for i, update := range updates {
		add := true
		for j, condition := range conditions {
			if condition.Reason == update.Reason {
				add = false
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

	return additions
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
		(a.Reason != b.Reason) ||
		(a.Message != b.Message) ||
		(a.ObservedGeneration != b.ObservedGeneration)
}
