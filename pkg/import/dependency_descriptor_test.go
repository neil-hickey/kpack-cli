// Copyright 2020-Present VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package _import_test

import (
	"testing"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"

	importpkg "github.com/buildpacks-community/kpack-cli/pkg/import"
)

func TestDescriptor(t *testing.T) {
	spec.Run(t, "TestDescriptor", testDescriptor)
}

func testDescriptor(t *testing.T, when spec.G, it spec.S) {
	desc := importpkg.DependencyDescriptor{
		DefaultClusterStack:   "some-stack",
		DefaultClusterBuilder: "some-cb",
		ClusterLifecycles: []importpkg.ClusterLifecycle{
			{
				Name:  "default",
				Image: "lifecycle-image",
			},
		},
		ClusterStores: []importpkg.ClusterStore{
			{
				Name: "some-store",
				Sources: []importpkg.Source{
					{
						Image: "some-store-image",
					},
				},
			},
		},
		ClusterStacks: []importpkg.ClusterStack{
			{
				Name: "some-stack",
				BuildImage: importpkg.Source{
					Image: "build-image",
				},
				RunImage: importpkg.Source{
					Image: "run-image",
				},
			},
		},
		ClusterBuilders: []importpkg.ClusterBuilder{
			{
				Name:         "some-cb",
				ClusterStack: "some-stack",
				ClusterStore: "some-store",
				Order: []v1alpha2.BuilderOrderEntry{
					{
						Group: []v1alpha2.BuilderBuildpackRef{
							{
								BuildpackRef: corev1alpha1.BuildpackRef{
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
		},
	}
	when("#ValidateDescriptor", func() {

		it("validates successfully", func() {
			require.NoError(t, importpkg.ValidateDescriptor(desc))
		})

		when("there is a duplicate lifecycle name", func() {
			desc.ClusterLifecycles = append(desc.ClusterLifecycles, importpkg.ClusterLifecycle{
				Name:  "default",
				Image: "another-image",
			})

			it("fails validation", func() {
				require.Error(t, importpkg.ValidateDescriptor(desc))
			})
		})

		when("there is a lifecycle with empty name", func() {
			it("fails validation", func() {
				descWithEmptyName := importpkg.DependencyDescriptor{
					ClusterLifecycles: []importpkg.ClusterLifecycle{
						{
							Name:  "",
							Image: "some-image",
						},
					},
				}
				err := importpkg.ValidateDescriptor(descWithEmptyName)
				require.Error(t, err)
				require.Contains(t, err.Error(), "cluster lifecycle name cannot be empty")
			})
		})

		when("there is a buildpack with empty name", func() {
			it("fails validation", func() {
				descWithEmptyName := importpkg.DependencyDescriptor{
					ClusterBuildpacks: []importpkg.ClusterBuildpack{
						{
							Name:  "",
							Image: "some-image",
						},
					},
				}
				err := importpkg.ValidateDescriptor(descWithEmptyName)
				require.Error(t, err)
				require.Contains(t, err.Error(), "cluster buildpack name cannot be empty")
			})
		})

		when("there is a duplicate buildpack name", func() {
			it("fails validation", func() {
				descWithDupe := importpkg.DependencyDescriptor{
					ClusterBuildpacks: []importpkg.ClusterBuildpack{
						{Name: "my-bp", Image: "image1"},
						{Name: "my-bp", Image: "image2"},
					},
				}
				err := importpkg.ValidateDescriptor(descWithDupe)
				require.Error(t, err)
				require.Contains(t, err.Error(), "duplicate cluster buildpack name")
			})
		})

		when("there is a duplicate store name", func() {
			desc.ClusterStores = append(desc.ClusterStores, importpkg.ClusterStore{
				Name: "some-store",
			})

			it("fails validation", func() {
				require.Error(t, importpkg.ValidateDescriptor(desc))
			})
		})

		when("there is a duplicate stack name", func() {
			desc.ClusterStacks = append(desc.ClusterStacks, importpkg.ClusterStack{
				Name: "some-stack",
			})

			it("fails validation", func() {
				require.Error(t, importpkg.ValidateDescriptor(desc))
			})
		})

		when("there is a duplicate cb name", func() {
			desc.ClusterBuilders = append(desc.ClusterBuilders, importpkg.ClusterBuilder{
				Name:         "some-cb",
				ClusterStack: "some-stack",
				ClusterStore: "some-store",
			})

			it("fails validation", func() {
				require.Error(t, importpkg.ValidateDescriptor(desc))
			})
		})

		when("the default stack does not exist", func() {
			desc.DefaultClusterStack = "does-not-exist"

			it("fails validation", func() {
				require.Error(t, importpkg.ValidateDescriptor(desc))
			})
		})

		when("there is no default clusterstack", func() {
			desc.DefaultClusterStack = ""

			it("validates successfully", func() {
				require.NoError(t, importpkg.ValidateDescriptor(desc))
			})
		})

		when("the default clusterbuilder does not exist", func() {
			desc.DefaultClusterBuilder = "does-not-exist"

			it("fails validation", func() {
				require.Error(t, importpkg.ValidateDescriptor(desc))
			})
		})

		when("there is no default clusterbuilder", func() {
			desc.DefaultClusterBuilder = ""

			it("validates successfully", func() {
				require.NoError(t, importpkg.ValidateDescriptor(desc))
			})
		})
	})

	when("GetClusterStacks", func() {
		it("returns the cluster stacks and the default cluster stack", func() {
			stacks := importpkg.GetClusterStacks(desc)
			expectedStacks := []importpkg.ClusterStack{
				{Name: "some-stack", BuildImage: importpkg.Source{Image: "build-image"}, RunImage: importpkg.Source{Image: "run-image"}},
				{Name: "default", BuildImage: importpkg.Source{Image: "build-image"}, RunImage: importpkg.Source{Image: "run-image"}}}
			require.Equal(t, expectedStacks, stacks)
		})
	})

	when("GetClusterBuilders", func() {
		it("returns the cluster builders and the default cluster builder", func() {
			builders := importpkg.GetClusterBuilders(desc)
			expectedBuilders := []importpkg.ClusterBuilder{
				{
					Name:         "some-cb",
					ClusterStack: "some-stack",
					ClusterStore: "some-store",
					Order: []v1alpha2.BuilderOrderEntry{
						{
							Group: []v1alpha2.BuilderBuildpackRef{
								{
									BuildpackRef: corev1alpha1.BuildpackRef{
										BuildpackInfo: corev1alpha1.BuildpackInfo{
											Id: "some-buildpack", Version: "1.2.3",
										},
										Optional: false,
									},
								},
							},
						},
					},
				},
				{
					Name:         "default",
					ClusterStack: "some-stack",
					ClusterStore: "some-store",
					Order: []v1alpha2.BuilderOrderEntry{
						{
							Group: []v1alpha2.BuilderBuildpackRef{
								{
									BuildpackRef: corev1alpha1.BuildpackRef{
										BuildpackInfo: corev1alpha1.BuildpackInfo{
											Id: "some-buildpack", Version: "1.2.3",
										},
										Optional: false,
									},
								},
							},
						},
					},
				},
			}
			require.Equal(t, expectedBuilders, builders)
		})
	})
}
