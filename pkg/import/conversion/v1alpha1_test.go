// Copyright 2020-Present VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package conversion_test

import (
	"testing"

	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"

	"github.com/buildpacks-community/kpack-cli/pkg/import/conversion"
)

func TestDescriptorV1Alpha1(t *testing.T) {
	spec.Run(t, "TestDescriptorV1Alpha1", testDescriptorV1Alpha1)
}

func testDescriptorV1Alpha1(t *testing.T, when spec.G, it spec.S) {
	when("#ToV1", func() {
		descV1Alpha1 := conversion.DependencyDescriptorV1Alpha1{
			DefaultStack:          "some-stack",
			DefaultClusterBuilder: "some-ccb",
			Stores: []conversion.ClusterStore{
				{
					Name: "some-store",
					Sources: []conversion.Source{
						{
							Image: "some-store-image",
						},
					},
				},
			},
			Stacks: []conversion.ClusterStack{
				{
					Name: "some-stack",
					BuildImage: conversion.Source{
						Image: "build-image",
					},
					RunImage: conversion.Source{
						Image: "run-image",
					},
				},
			},
			ClusterBuilders: []conversion.ClusterBuilderV1Alpha1{
				{
					Name:  "some-ccb",
					Stack: "some-stack",
					Store: "some-store",
					Order: []corev1alpha1.OrderEntry{
						{
							Group: []corev1alpha1.BuildpackRef{
								{
									BuildpackInfo: corev1alpha1.BuildpackInfo{
										Id:      "some-buildpack",
										Version: "1.2.3",
									},
									Optional: false,
								},
							},
						},
					},
				},
			},
		}

		it("converts successfully", func() {
			d := descV1Alpha1.ToV1("kp.kpack.io/v1")
			require.Equal(t, "kp.kpack.io/v1", d.APIVersion)
			require.Empty(t, d.ClusterLifecycles)
			require.Empty(t, d.ClusterBuildpacks)
			require.Equal(t, "some-stack", d.DefaultClusterStack)
			require.Equal(t, "some-ccb", d.DefaultClusterBuilder)
		})
	})
}

