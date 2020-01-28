package addlabel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestPrepareLabelPatch(t *testing.T) {
	patchType := v1beta1.PatchTypeJSONPatch
	tests := []struct {
		name     string
		in       v1beta1.AdmissionReview
		expected *v1beta1.AdmissionResponse
	}{
		{
			name:     "empty admission",
			in:       v1beta1.AdmissionReview{},
			expected: &v1beta1.AdmissionResponse{Allowed: false},
		},
		{
			name:     "empty request",
			in:       v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{}},
			expected: &v1beta1.AdmissionResponse{Allowed: false},
		},
		{
			name:     "empty pod spec object",
			in:       v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{Object: runtime.RawExtension{}}},
			expected: &v1beta1.AdmissionResponse{Allowed: false},
		},
		{
			name: "empty pod labels",
			in: v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{Object: runtime.RawExtension{
				Raw: []byte(`{"metadata":{"labels":null}}`),
			}}},
			expected: &v1beta1.AdmissionResponse{Allowed: true, Patch: []byte(`[{ "op": "add", "path": "/metadata/labels", "value": {"awesome-label": "webhook"}}]`),
				PatchType: &patchType,
			},
		},
		{
			name: "labels exists",
			in: v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{Object: runtime.RawExtension{
				Raw: []byte(`{"metadata":{"labels":{"my-other-awesome-label":"test"}}}`),
			}}},
			expected: &v1beta1.AdmissionResponse{Allowed: true, Patch: []byte(`[{ "op": "add", "path": "/metadata/labels/awesome-label", "value": "webhook" }]`),
				PatchType: &patchType,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			actual := prepareLabelPatch(test.in)
			assert.Equal(tt, test.expected.Allowed, actual.Allowed)
			assert.Equal(tt, test.expected.Result, actual.Result)
			assert.Equal(tt, test.expected.Patch, actual.Patch)
			assert.Equal(tt, test.expected.PatchType, actual.PatchType)
		})
	}
}
