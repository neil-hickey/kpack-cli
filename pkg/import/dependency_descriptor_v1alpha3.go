// Copyright 2020-Present VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package _import

const APIVersionV1Alpha3 = "kp.kpack.io/v1alpha3"

type Lifecycle struct {
	Image string `yaml:"image" json:"image"`
}

type DependencyDescriptorV1Alpha3 struct {
	APIVersion            string           `yaml:"apiVersion"`
	Kind                  string           `yaml:"kind"`
	DefaultClusterStack   string           `yaml:"defaultClusterStack"`
	DefaultClusterBuilder string           `yaml:"defaultClusterBuilder"`
	Lifecycle             Lifecycle        `yaml:"lifecycle"`
	ClusterStores         []ClusterStore   `yaml:"clusterStores"`
	ClusterStacks         []ClusterStack   `yaml:"clusterStacks"`
	ClusterBuilders       []ClusterBuilder `yaml:"clusterBuilders"`
}

func (d DependencyDescriptorV1Alpha3) ToV1() DependencyDescriptor {
	var v1 DependencyDescriptor
	v1.APIVersion = d.APIVersion
	v1.Kind = d.Kind
	v1.DefaultClusterStack = d.DefaultClusterStack
	v1.DefaultClusterBuilder = d.DefaultClusterBuilder
	v1.ClusterStores = d.ClusterStores
	v1.ClusterStacks = d.ClusterStacks
	v1.ClusterBuilders = d.ClusterBuilders

	// Convert single lifecycle to array with name "default"
	if d.Lifecycle.Image != "" {
		v1.ClusterLifecycles = []ClusterLifecycle{
			{
				Name:  "default",
				Image: d.Lifecycle.Image,
			},
		}
	}

	return v1
}
