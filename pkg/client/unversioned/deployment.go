/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package unversioned

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

// DeploymentsNamespacer has methods to work with Deployment resources in a namespace
type DeploymentsNamespacer interface {
	Deployments(namespace string) DeploymentInterface
}

// DeploymentInterface has methods to work with Deployment resources.
type DeploymentInterface interface {
	List(label labels.Selector, field fields.Selector) (*extensions.DeploymentList, error)
	Get(name string) (*extensions.Deployment, error)
	Delete(name string, options *api.DeleteOptions) error
	Create(*extensions.Deployment) (*extensions.Deployment, error)
	Update(*extensions.Deployment) (*extensions.Deployment, error)
	UpdateStatus(*extensions.Deployment) (*extensions.Deployment, error)
	Watch(label labels.Selector, field fields.Selector, opts unversioned.ListOptions) (watch.Interface, error)
}

// deployments implements DeploymentInterface
type deployments struct {
	client *ExtensionsClient
	ns     string
}

// newDeployments returns a Deployments
func newDeployments(c *ExtensionsClient, namespace string) *deployments {
	return &deployments{
		client: c,
		ns:     namespace,
	}
}

// List takes label and field selectors, and returns the list of Deployments that match those selectors.
func (c *deployments) List(label labels.Selector, field fields.Selector) (result *extensions.DeploymentList, err error) {
	result = &extensions.DeploymentList{}
	err = c.client.Get().Namespace(c.ns).Resource("deployments").LabelsSelectorParam(label).FieldsSelectorParam(field).Do().Into(result)
	return
}

// Get takes name of the deployment, and returns the corresponding deployment object, and an error if there is any.
func (c *deployments) Get(name string) (result *extensions.Deployment, err error) {
	result = &extensions.Deployment{}
	err = c.client.Get().Namespace(c.ns).Resource("deployments").Name(name).Do().Into(result)
	return
}

// Delete takes name of the deployment and deletes it. Returns an error if one occurs.
func (c *deployments) Delete(name string, options *api.DeleteOptions) error {
	if options == nil {
		return c.client.Delete().Namespace(c.ns).Resource("deployments").Name(name).Do().Error()
	}
	body, err := api.Scheme.EncodeToVersion(options, c.client.APIVersion())
	if err != nil {
		return err
	}
	return c.client.Delete().Namespace(c.ns).Resource("deployments").Name(name).Body(body).Do().Error()
}

// Create takes the representation of a deployment and creates it.  Returns the server's representation of the deployment, and an error, if there is any.
func (c *deployments) Create(deployment *extensions.Deployment) (result *extensions.Deployment, err error) {
	result = &extensions.Deployment{}
	err = c.client.Post().Namespace(c.ns).Resource("deployments").Body(deployment).Do().Into(result)
	return
}

// Update takes the representation of a deployment and updates it. Returns the server's representation of the deployment, and an error, if there is any.
func (c *deployments) Update(deployment *extensions.Deployment) (result *extensions.Deployment, err error) {
	result = &extensions.Deployment{}
	err = c.client.Put().Namespace(c.ns).Resource("deployments").Name(deployment.Name).Body(deployment).Do().Into(result)
	return
}

func (c *deployments) UpdateStatus(deployment *extensions.Deployment) (result *extensions.Deployment, err error) {
	result = &extensions.Deployment{}
	err = c.client.Put().Namespace(c.ns).Resource("deployments").Name(deployment.Name).SubResource("status").Body(deployment).Do().Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested deployments.
func (c *deployments) Watch(label labels.Selector, field fields.Selector, opts unversioned.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("deployments").
		VersionedParams(&opts, api.Scheme).
		LabelsSelectorParam(label).
		FieldsSelectorParam(field).
		Watch()
}
