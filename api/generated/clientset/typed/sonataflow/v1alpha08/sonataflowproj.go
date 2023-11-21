// Copyright 2023 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by client-gen. DO NOT EDIT.

package v1alpha08

import (
	"context"
	"time"

	scheme "github.com/apache/incubator-kie-kogito-serverless-operator/api/generated/clientset/scheme"
	v1alpha08 "github.com/apache/incubator-kie-kogito-serverless-operator/api/sonataflow/v1alpha08"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SonataFlowProjsGetter has a method to return a SonataFlowProjInterface.
// A group's client should implement this interface.
type SonataFlowProjsGetter interface {
	SonataFlowProjs(namespace string) SonataFlowProjInterface
}

// SonataFlowProjInterface has methods to work with SonataFlowProj resources.
type SonataFlowProjInterface interface {
	Create(ctx context.Context, sonataFlowProj *v1alpha08.SonataFlowProj, opts v1.CreateOptions) (*v1alpha08.SonataFlowProj, error)
	Update(ctx context.Context, sonataFlowProj *v1alpha08.SonataFlowProj, opts v1.UpdateOptions) (*v1alpha08.SonataFlowProj, error)
	UpdateStatus(ctx context.Context, sonataFlowProj *v1alpha08.SonataFlowProj, opts v1.UpdateOptions) (*v1alpha08.SonataFlowProj, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha08.SonataFlowProj, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha08.SonataFlowProjList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha08.SonataFlowProj, err error)
	SonataFlowProjExpansion
}

// sonataFlowProjs implements SonataFlowProjInterface
type sonataFlowProjs struct {
	client rest.Interface
	ns     string
}

// newSonataFlowProjs returns a SonataFlowProjs
func newSonataFlowProjs(c *SonataflowV1alpha08Client, namespace string) *sonataFlowProjs {
	return &sonataFlowProjs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the sonataFlowProj, and returns the corresponding sonataFlowProj object, and an error if there is any.
func (c *sonataFlowProjs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha08.SonataFlowProj, err error) {
	result = &v1alpha08.SonataFlowProj{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sonataflowprojs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SonataFlowProjs that match those selectors.
func (c *sonataFlowProjs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha08.SonataFlowProjList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha08.SonataFlowProjList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sonataflowprojs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sonataFlowProjs.
func (c *sonataFlowProjs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("sonataflowprojs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a sonataFlowProj and creates it.  Returns the server's representation of the sonataFlowProj, and an error, if there is any.
func (c *sonataFlowProjs) Create(ctx context.Context, sonataFlowProj *v1alpha08.SonataFlowProj, opts v1.CreateOptions) (result *v1alpha08.SonataFlowProj, err error) {
	result = &v1alpha08.SonataFlowProj{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("sonataflowprojs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sonataFlowProj).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a sonataFlowProj and updates it. Returns the server's representation of the sonataFlowProj, and an error, if there is any.
func (c *sonataFlowProjs) Update(ctx context.Context, sonataFlowProj *v1alpha08.SonataFlowProj, opts v1.UpdateOptions) (result *v1alpha08.SonataFlowProj, err error) {
	result = &v1alpha08.SonataFlowProj{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sonataflowprojs").
		Name(sonataFlowProj.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sonataFlowProj).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *sonataFlowProjs) UpdateStatus(ctx context.Context, sonataFlowProj *v1alpha08.SonataFlowProj, opts v1.UpdateOptions) (result *v1alpha08.SonataFlowProj, err error) {
	result = &v1alpha08.SonataFlowProj{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sonataflowprojs").
		Name(sonataFlowProj.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(sonataFlowProj).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the sonataFlowProj and deletes it. Returns an error if one occurs.
func (c *sonataFlowProjs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sonataflowprojs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sonataFlowProjs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sonataflowprojs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched sonataFlowProj.
func (c *sonataFlowProjs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha08.SonataFlowProj, err error) {
	result = &v1alpha08.SonataFlowProj{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("sonataflowprojs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
