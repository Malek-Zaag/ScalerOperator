package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SUCCESS = "Success"
	FAILED  = "Failed"
)

type ScalerSpec struct {
	Start       int              `json:"start"`
	End         int              `json:"end"`
	Replicas    int32            `json:"replicas"`
	Deployments []NamesapcedName `json:"deployments"`
}

type NamesapcedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// ScalerStatus defines the observed state of Scaler
type ScalerStatus struct {
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Scaler is the Schema for the scalers API
type Scaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScalerSpec   `json:"spec,omitempty"`
	Status ScalerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ScalerList contains a list of Scaler
type ScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Scaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Scaler{}, &ScalerList{})
}
