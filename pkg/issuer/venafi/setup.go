package venafi

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"

	"github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
)

func (v *Venafi) Setup(ctx context.Context) error {
	err := v.client.Ping()
	if err != nil {
		glog.Infof("Issuer could not connect to endpoint with provided credentials. Issuer failed to connect to endpoint\n")
		v.issuer.UpdateStatusCondition(v1alpha1.IssuerConditionReady, v1alpha1.ConditionFalse,
			"ErrorPing", fmt.Sprintf("Failed to connect to Venafi endpoint"))
		return fmt.Errorf("error verifying Venafi client: %s", err.Error())
	}

	// If it does not already have a 'ready' condition, we'll also log an event
	// to make it really clear to users that this Issuer is ready.
	if !v.issuer.HasCondition(v1alpha1.IssuerCondition{
		Type:   v1alpha1.IssuerConditionReady,
		Status: v1alpha1.ConditionTrue,
	}) {
		v.Recorder.Eventf(v.issuer, corev1.EventTypeNormal, "Ready", "Verified issuer with Venafi server")
	}

	glog.Info("Venafi issuer started")
	v.issuer.UpdateStatusCondition(v1alpha1.IssuerConditionReady, v1alpha1.ConditionTrue, "Venafi issuer started", "Venafi issuer started")

	return nil
}
