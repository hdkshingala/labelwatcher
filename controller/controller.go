package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Controller struct {
	clientSet kubernetes.Interface
}

func (cont *Controller) ValidatePodLabels(oldLables map[string]string, newLabels map[string]string, ns string) (string, bool) {
	log.Printf("Received labels are '%v' and '%v'.\n", oldLables, newLabels)
	if reflect.DeepEqual(oldLables, newLabels) {
		return "", true
	}
	changedLables := make(map[string]string)
	for oldKey, oldValue := range oldLables {
		if newLabels[oldKey] != oldValue {
			changedLables[oldKey] = oldValue
		}
	}

	netPolicies, err := cont.clientSet.NetworkingV1().NetworkPolicies(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Sprintf("Failed to fetch NetworkPolicies due to error: '%s'", err.Error()), false
	}

	for _, netPolicy := range netPolicies.Items {
		for key, value := range changedLables {
			if netPolicy.Spec.PodSelector.MatchLabels[key] == value {
				return fmt.Sprintf(
					"Cannot modify the label with key '%s' and value '%s' as that is used as selector in NetworkPolicy with name '%s'.",
					key, value, netPolicy.Name), false
			}
		}
	}

	return "", true
}

func NewController() (*Controller, error) {
	client, err := prepareClientSet()
	if err != nil {
		return nil, err
	}

	cont := &Controller{
		clientSet: client,
	}

	return cont, nil
}

func prepareClientSet() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		err = errors.New(fmt.Sprintf("Error received while creating config from InCluster, error: %s.\n", err.Error()))
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error received while creating client set, error: %s.\n", err.Error()))
		return nil, err
	}
	return clientset, nil
}
