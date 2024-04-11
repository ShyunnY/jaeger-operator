package runner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSortConditions(t *testing.T) {

	cases := []struct {
		caseName string
		condList []metav1.Condition
		expected []metav1.Condition
	}{
		{
			caseName: "normal",
			condList: []metav1.Condition{
				{
					Type:               "a",
					LastTransitionTime: metav1.NewTime(time.Now()),
				},
				{
					Type:               "c",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
				},
				{
					Type:               "b",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-2 * time.Hour)),
				},
			},
			expected: []metav1.Condition{
				{
					Type:               "a",
					LastTransitionTime: metav1.NewTime(time.Now()),
				},
				{
					Type:               "c",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
				},
				{
					Type:               "b",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-2 * time.Hour)),
				},
			},
		},
		{
			caseName: "sort",
			condList: []metav1.Condition{
				{
					Type:               "b",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-2 * time.Hour)),
				},
				{
					Type:               "a",
					LastTransitionTime: metav1.NewTime(time.Now()),
				},
				{
					Type:               "c",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
				},
			},
			expected: []metav1.Condition{
				{
					Type:               "a",
					LastTransitionTime: metav1.NewTime(time.Now()),
				},
				{
					Type:               "c",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
				},
				{
					Type:               "b",
					LastTransitionTime: metav1.NewTime(time.Now().Add(-2 * time.Hour)),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			actual := sortConditions(tc.condList)

			for i := range tc.expected {
				assert.Equal(t, tc.expected[i].Type, actual[i].Type)
				assert.Equal(t, tc.expected[i].LastTransitionTime, actual[i].LastTransitionTime)
			}
		})
	}

}
