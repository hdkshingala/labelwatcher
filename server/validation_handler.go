package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	admv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
)

func (server *Server) ServeDeploymentCreatorValidation(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %+v.\n", r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error fetching the http.Request object. Error: %s.\n", err.Error())
		responsewriters.InternalError(w, r, err)
	}

	gvk := admv1.SchemeGroupVersion.WithKind("AdmissionReview")
	var review admv1.AdmissionReview
	_, _, err = codec.UniversalDeserializer().Decode(body, &gvk, &review)
	if err != nil {
		log.Printf("Error converting the http.Request object to AdmissionReview object. Error: %s.\n", err.Error())
		responsewriters.InternalError(w, r, err)
	}

	var oldPod corev1.Pod
	var newPod corev1.Pod
	gvk = corev1.SchemeGroupVersion.WithKind("Pod")
	_, _, err = codec.UniversalDeserializer().Decode(review.Request.OldObject.Raw, &gvk, &oldPod)
	if err != nil {
		log.Printf("Error converting the AdmissionReview object to Pod object. Error: %s.\n", err.Error())
		responsewriters.InternalError(w, r, err)
	}

	_, _, err = codec.UniversalDeserializer().Decode(review.Request.Object.Raw, &gvk, &newPod)
	if err != nil {
		log.Printf("Error converting the AdmissionReview object to Pod object. Error: %s.\n", err.Error())
		responsewriters.InternalError(w, r, err)
	}
	var response admv1.AdmissionResponse
	msg, allow := server.Controller.ValidatePodLabels(oldPod.Labels, newPod.Labels, newPod.GetNamespace())
	if !allow {
		response = admv1.AdmissionResponse{
			UID:     review.Request.UID,
			Allowed: allow,
			Result: &metav1.Status{
				Message: msg,
			},
		}
	} else {
		response = admv1.AdmissionResponse{
			UID:     review.Request.UID,
			Allowed: allow,
		}
	}

	review.Response = &response
	res, err := json.Marshal(review)
	if err != nil {
		log.Printf("Error converting the AdmissionReview object to JSON object. Error: %s.\n", err.Error())
		responsewriters.InternalError(w, r, err)
	}

	_, err = w.Write(res)
	if err != nil {
		log.Printf("Error sending the response '%+v'. Error: %s.\n", res, err.Error())
		responsewriters.InternalError(w, r, err)
	}
}
