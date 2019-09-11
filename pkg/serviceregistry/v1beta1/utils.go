// Package v1beta1 of connmgrs Defines the Reconciliation logic and required setup for component operator.
// IBM Confidential
// OCO Source Materials
// 5737-E67
// (C) Copyright IBM Corporation 2019 All Rights Reserved
// The source code for this program is not published or otherwise divested of its trade secrets, irrespective of what has been deposited with the U.S. Copyright Office.
package v1beta1

import (
	"context"

	multicloudv1beta1 "github.ibm.com/IBMPrivateCloud/ibm-klusterlet-operator/pkg/apis/multicloud/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IsReady helps the other components to check whether the connmgr pod is ready
func IsReady(instance *multicloudv1beta1.Endpoint, c client.Client) (bool, error) {
	var svcregIsReady, coreDNSIsReady bool

	foundSvcRegDeployment := &extensionsv1beta1.Deployment{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-search", Namespace: instance.Namespace}, foundSvcRegDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			log.V(5).Info("Cannot find the svcreg deployment")
			return false, nil
		}
		log.Error(err, "Unexpected ERROR")
		return false, err
	}
	log.V(5).Info("Svcreg deployment exist")
	for _, condition := range foundSvcRegDeployment.Status.Conditions {
		if condition.Type == "Available" && condition.Status == "True" {
			svcregIsReady = true
		}
	}

	foundSvcRegCoreDNSDeployment := &extensionsv1beta1.Deployment{}
	err = c.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-search", Namespace: instance.Namespace}, foundSvcRegCoreDNSDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			log.V(5).Info("Cannot find the svcreg coreDNS deployment")
			return false, nil
		}
		log.Error(err, "Unexpected ERROR")
		return false, err
	}
	log.V(5).Info("Svcreg coreDNS deployment exist")
	for _, condition := range foundSvcRegCoreDNSDeployment.Status.Conditions {
		if condition.Type == "Available" && condition.Status == "True" {
			coreDNSIsReady = true
		}
	}

	return svcregIsReady && coreDNSIsReady, nil
}