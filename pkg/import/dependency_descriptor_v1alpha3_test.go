// Copyright 2020-Present VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package _import_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"

	importpkg "github.com/buildpacks-community/kpack-cli/pkg/import"
)

func TestDescriptorV1Alpha3(t *testing.T) {
	spec.Run(t, "TestDescriptorV1Alpha3", testDescriptorV1Alpha3)
}

func testDescriptorV1Alpha3(t *testing.T, when spec.G, it spec.S) {
	when("#ToV1", func() {
		descV1Alpha3 := importpkg.DependencyDescriptorV1Alpha3{
			DefaultClusterStack:   "some-stack",
			DefaultClusterBuilder: "some-ccb",
			Lifecycle: importpkg.Lifecycle{
				Image: "some-lifecycle-image",
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
					Name:         "some-ccb",
					ClusterStack: "some-stack",
					ClusterStore: "some-store",
				},
			},
		}

		it("converts successfully", func() {
			v1 := descV1Alpha3.ToV1()
			require.NoError(t, v1.Validate())
			require.Len(t, v1.ClusterLifecycles, 1)
			require.Equal(t, "default", v1.ClusterLifecycles[0].Name)
			require.Equal(t, "some-lifecycle-image", v1.ClusterLifecycles[0].Image)
		})

		it("converts with empty lifecycle", func() {
			descV1Alpha3.Lifecycle = importpkg.Lifecycle{}
			v1 := descV1Alpha3.ToV1()
			require.NoError(t, v1.Validate())
			require.Len(t, v1.ClusterLifecycles, 0)
		})
	})
}
