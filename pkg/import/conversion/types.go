// Copyright 2020-Present VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

// Package conversion provides types and functions to convert older API versions
// of dependency descriptors (v1alpha1, v1alpha3) to the current v1 format.
package conversion

import (
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
)

// API version constants
const (
	APIVersionV1Alpha1 = "kp.kpack.io/v1alpha1"
	APIVersionV1Alpha3 = "kp.kpack.io/v1alpha3"
)

// Lifecycle represents a single lifecycle image in v1alpha3 format
type Lifecycle struct {
	Image string `yaml:"image" json:"image"`
}

// ClusterBuilderV1Alpha1 represents a ClusterBuilder in v1alpha1 format
type ClusterBuilderV1Alpha1 struct {
	Name  string                    `yaml:"name"`
	Stack string                    `yaml:"stack"`
	Store string                    `yaml:"store"`
	Order []corev1alpha1.OrderEntry `yaml:"order"`
}

// Source represents an image source
type Source struct {
	Image string `yaml:"image"`
}

// ClusterStore represents a ClusterStore in the descriptor
type ClusterStore struct {
	Name    string   `yaml:"name" json:"name"`
	Sources []Source `yaml:"sources" json:"sources"`
}

// ClusterStack represents a ClusterStack in the descriptor
type ClusterStack struct {
	Name       string `yaml:"name" json:"name"`
	BuildImage Source `yaml:"buildImage" json:"buildImage"`
	RunImage   Source `yaml:"runImage" json:"runImage"`
}

// ClusterBuilder represents a ClusterBuilder in the descriptor (v1alpha3+)
type ClusterBuilder struct {
	Name         string                       `yaml:"name" json:"name"`
	ClusterStack string                       `yaml:"clusterStack" json:"clusterStack"`
	ClusterStore string                       `yaml:"clusterStore" json:"clusterStore"`
	Order        []v1alpha2.BuilderOrderEntry `yaml:"order" json:"order"`
}

// ClusterLifecycle represents a ClusterLifecycle in the v1 descriptor
type ClusterLifecycle struct {
	Name  string `yaml:"name" json:"name"`
	Image string `yaml:"image" json:"image"`
}

// ClusterBuildpack represents a ClusterBuildpack in the v1 descriptor
type ClusterBuildpack struct {
	Name  string `yaml:"name" json:"name"`
	Image string `yaml:"image" json:"image"`
}

// DependencyDescriptor represents the target v1 format that all conversions produce
type DependencyDescriptor struct {
	APIVersion            string             `yaml:"apiVersion" json:"apiVersion"`
	Kind                  string             `yaml:"kind" json:"kind"`
	DefaultClusterStack   string             `yaml:"defaultClusterStack" json:"defaultClusterStack"`
	DefaultClusterBuilder string             `yaml:"defaultClusterBuilder" json:"defaultClusterBuilder"`
	ClusterLifecycles     []ClusterLifecycle `yaml:"clusterLifecycles" json:"clusterLifecycles"`
	ClusterBuildpacks     []ClusterBuildpack `yaml:"clusterBuildpacks" json:"clusterBuildpacks"`
	ClusterStores         []ClusterStore     `yaml:"clusterStores" json:"clusterStores"`
	ClusterStacks         []ClusterStack     `yaml:"clusterStacks" json:"clusterStacks"`
	ClusterBuilders       []ClusterBuilder   `yaml:"clusterBuilders" json:"clusterBuilders"`
}
