package addlabels

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// admissionResponseError is a helper function to create an AdmissionResponse
// with an embedded error
func admissionResponseError(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Status:  "Failure",
			Message: "Admission label webhook error: " + err.Error(),
			Code:    500,
		},
	}
}

//AppendLabelsHandler append lables patch
func AppendLabelsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Println("contentType=" + contentType + ", expect application/json")
		return
	}

	var body []byte
	var err error
	if r.Body != nil {
		body, err = ioutil.ReadAll(r.Body)
	}

	if len(body) == 0 {
		log.Print("Body is empty")
		return
	}

	// The AdmissionReview that was sent to the webhook
	requestedAdmissionReview := v1beta1.AdmissionReview{}

	// The AdmissionReview that will be returned
	responseAdmissionReview := v1beta1.AdmissionReview{}

	var scheme = runtime.NewScheme()
	// Codecs provides access to encoding and decoding for the scheme.
	var Codecs = serializer.NewCodecFactory(scheme)
	deserializer := Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		log.Print(err)
		responseAdmissionReview.Response = admissionResponseError(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = prepareLabelPatch(requestedAdmissionReview)
	}
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		log.Println(err.Error())
	}
	if _, err := w.Write(respBytes); err != nil {
		log.Println(err)
	}
}

func prepareLabelPatch(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	var pod corev1.Pod
	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &pod)
	if err != nil {
		return admissionResponseError(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	if len(pod.ObjectMeta.Labels) == 0 {
		reviewResponse.Patch = []byte(`[{ "op": "add", "path": "/metadata/labels", "value": {"awesome-label": "webhook"}}]`)
	} else {
		reviewResponse.Patch = []byte(`[{ "op": "add", "path": "/metadata/labels/awesome-label", "value": "webhook" }]`)
	}

	pt := v1beta1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt
	return &reviewResponse
}
