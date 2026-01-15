// Copyright 2020-Present VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package _import

import (
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pkg/errors"

	"github.com/buildpacks-community/kpack-cli/pkg/import/conversion"
)

const CurrentAPIVersion = "kp.kpack.io/v1"

type API struct {
	Version string `yaml:"apiVersion" json:"apiVersion"`
}

// Type aliases to use types from the conversion package
type (
	DependencyDescriptor = conversion.DependencyDescriptor
	Source               = conversion.Source
	ClusterLifecycle     = conversion.ClusterLifecycle
	ClusterBuildpack     = conversion.ClusterBuildpack
	ClusterStore         = conversion.ClusterStore
	ClusterStack         = conversion.ClusterStack
	ClusterBuilder       = conversion.ClusterBuilder
)

func ValidateDescriptor(d DependencyDescriptor) error {
	lifecycleSet := map[string]interface{}{}
	for _, lifecycle := range d.ClusterLifecycles {
		if lifecycle.Name == "" {
			return errors.New("cluster lifecycle name cannot be empty")
		}
		if n, ok := lifecycleSet[lifecycle.Name]; ok {
			return errors.Errorf("duplicate cluster lifecycle name '%s'", n)
		}
		lifecycleSet[lifecycle.Name] = nil

		_, err := name.ParseReference(lifecycle.Image, name.WeakValidation)
		if err != nil {
			return err
		}
	}

	buildpackSet := map[string]interface{}{}
	for _, buildpack := range d.ClusterBuildpacks {
		if buildpack.Name == "" {
			return errors.New("cluster buildpack name cannot be empty")
		}
		if n, ok := buildpackSet[buildpack.Name]; ok {
			return errors.Errorf("duplicate cluster buildpack name '%s'", n)
		}
		buildpackSet[buildpack.Name] = nil

		_, err := name.ParseReference(buildpack.Image, name.WeakValidation)
		if err != nil {
			return err
		}
	}

	storeSet := map[string]interface{}{}
	for _, store := range d.ClusterStores {
		if n, ok := storeSet[store.Name]; ok {
			return errors.Errorf("duplicate store name '%s'", n)
		}
		storeSet[store.Name] = nil

		for _, src := range store.Sources {
			_, err := name.ParseReference(src.Image, name.WeakValidation)
			if err != nil {
				return err
			}
		}
	}

	stackSet := map[string]interface{}{}
	for _, stack := range d.ClusterStacks {
		if n, ok := stackSet[stack.Name]; ok {
			return errors.Errorf("duplicate stack name '%s'", n)
		}
		stackSet[stack.Name] = nil

		_, err := name.ParseReference(stack.BuildImage.Image, name.WeakValidation)
		if err != nil {
			return err
		}

		_, err = name.ParseReference(stack.RunImage.Image, name.WeakValidation)
		if err != nil {
			return err
		}
	}

	if _, ok := stackSet[d.DefaultClusterStack]; !ok && d.DefaultClusterStack != "" {
		return errors.Errorf("default cluster stack '%s' not found", d.DefaultClusterStack)
	}

	ccbSet := map[string]interface{}{}
	for _, ccb := range d.ClusterBuilders {
		if n, ok := ccbSet[ccb.Name]; ok {
			return errors.Errorf("duplicate cluster builder name '%s'", n)
		}
		ccbSet[ccb.Name] = nil
	}

	if _, ok := ccbSet[d.DefaultClusterBuilder]; !ok && d.DefaultClusterBuilder != "" {
		return errors.Errorf("default cluster builder '%s' not found", d.DefaultClusterBuilder)
	}

	return nil
}

func GetClusterLifecycles(d DependencyDescriptor) []ClusterLifecycle {
	return d.ClusterLifecycles
}

func GetClusterBuildpacks(d DependencyDescriptor) []ClusterBuildpack {
	return d.ClusterBuildpacks
}

func GetClusterStacks(d DependencyDescriptor) []ClusterStack {
	stacks := d.ClusterStacks
	for _, stack := range d.ClusterStacks {
		if stack.Name == d.DefaultClusterStack {
			stacks = append(stacks, ClusterStack{
				Name:       "default",
				BuildImage: stack.BuildImage,
				RunImage:   stack.RunImage,
			})
			break
		}
	}
	return stacks
}

func GetClusterBuilders(d DependencyDescriptor) []ClusterBuilder {
	builders := d.ClusterBuilders
	for _, cb := range d.ClusterBuilders {
		if cb.Name == d.DefaultClusterBuilder {
			builders = append(builders, ClusterBuilder{
				Name:         "default",
				ClusterStack: cb.ClusterStack,
				ClusterStore: cb.ClusterStore,
				Order:        cb.Order,
			})
			break
		}
	}
	return builders
}
