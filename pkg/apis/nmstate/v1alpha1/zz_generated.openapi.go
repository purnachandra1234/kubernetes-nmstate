// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1.NodeNetworkState":       schema_pkg_apis_nmstate_v1alpha1_NodeNetworkState(ref),
		"github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1.NodeNetworkStateSpec":   schema_pkg_apis_nmstate_v1alpha1_NodeNetworkStateSpec(ref),
		"github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1.NodeNetworkStateStatus": schema_pkg_apis_nmstate_v1alpha1_NodeNetworkStateStatus(ref),
	}
}

func schema_pkg_apis_nmstate_v1alpha1_NodeNetworkState(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NodeNetworkState is the Schema for the nodenetworkstates API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1.NodeNetworkStateSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1.NodeNetworkStateStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1.NodeNetworkStateSpec", "github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1.NodeNetworkStateStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_nmstate_v1alpha1_NodeNetworkStateSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NodeNetworkStateSpec defines the desired state of NodeNetworkState",
				Properties: map[string]spec.Schema{
					"managed": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"boolean"},
							Format: "",
						},
					},
					"nodeName": {
						SchemaProps: spec.SchemaProps{
							Description: "Name of the node reporting this state",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"desiredState": {
						SchemaProps: spec.SchemaProps{
							Description: "The desired configuration for the node",
							Type:        []string{"object"},
						},
					},
				},
				Required: []string{"managed", "nodeName"},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_nmstate_v1alpha1_NodeNetworkStateStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NodeNetworkStateStatus is the status of the NodeNetworkState of a specific node",
				Properties: map[string]spec.Schema{
					"currentState": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}