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

package workflowproj

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/parser"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/kiegroup/kogito-serverless-operator/api/metadata"
	operatorapi "github.com/kiegroup/kogito-serverless-operator/api/v1alpha08"
)

var _ WorkflowProjectHandler = &workflowProjectHandler{}

// defaultResourcePath is the default resource path to add to the generated ConfigMaps
const defaultResourcePath = "specs"

// WorkflowProjectHandler is the description of the handler interface.
// A handler can generate Kubernetes manifests to deploy a new Kogito Serverless Workflow project in the cluster
type WorkflowProjectHandler interface {
	// Named overwrites the workflow ID. The handler will use this name instead to generate the manifests name.
	// Remember that together with the Namespace, the Name is the unique key of a Kubernetes object.
	Named(name string) WorkflowProjectHandler
	// WithWorkflow reader for a file or the content stream of a workflow definition.
	WithWorkflow(reader io.Reader) WorkflowProjectHandler
	// WithAppProperties reader for a file or the content stream of a workflow application properties.
	WithAppProperties(reader io.Reader) WorkflowProjectHandler
	// AddResource reader for a file or the content stream of any resource needed by the workflow. E.g. an OpenAPI specification file.
	// Name is required, should match the workflow function definition.
	AddResource(name string, reader io.Reader) WorkflowProjectHandler
	// AddResourceAt same as AddResource, but defines the path instead of using the default.
	AddResourceAt(name, path string, reader io.Reader) WorkflowProjectHandler
	// SaveAsKubernetesManifests saves the project in the given file system path in YAML format.
	SaveAsKubernetesManifests(path string) error
	// AsObjects returns a reference to the WorkflowProject holding the Kubernetes Manifests based on your files.
	AsObjects() (*WorkflowProject, error)
}

// WorkflowProject is a structure to hold every Kubernetes object generated by the given WorkflowProjectHandler handler.
type WorkflowProject struct {
	// Workflow the workflow definition
	Workflow *operatorapi.KogitoServerlessWorkflow
	// Properties the application properties for the workflow
	Properties *corev1.ConfigMap
	// Resources any resource that this workflow requires, like an OpenAPI specification file.
	Resources []*corev1.ConfigMap
}

type resource struct {
	name     string
	contents io.Reader
}

// New is the entry point for this package.
// You can create a new handler with the given namespace, meaning that every manifest generated will use this namespace.
// namespace is a required parameter.
func New(namespace string) WorkflowProjectHandler {
	s := scheme.Scheme
	utilruntime.Must(operatorapi.AddToScheme(s))
	utilruntime.Must(corev1.AddToScheme(s))
	return &workflowProjectHandler{
		scheme:       s,
		namespace:    namespace,
		rawResources: map[string][]*resource{},
	}
}

type workflowProjectHandler struct {
	name             string
	namespace        string
	scheme           *runtime.Scheme
	project          WorkflowProject
	rawWorkflow      io.Reader
	rawAppProperties io.Reader
	rawResources     map[string][]*resource
	parsed           bool
}

func (w *workflowProjectHandler) Named(name string) WorkflowProjectHandler {
	w.name = strings.ToLower(name)
	w.parsed = false
	return w
}

func (w *workflowProjectHandler) WithWorkflow(reader io.Reader) WorkflowProjectHandler {
	w.rawWorkflow = reader
	w.parsed = false
	return w
}

func (w *workflowProjectHandler) WithAppProperties(reader io.Reader) WorkflowProjectHandler {
	w.rawAppProperties = reader
	w.parsed = false
	return w
}

func (w *workflowProjectHandler) AddResource(name string, reader io.Reader) WorkflowProjectHandler {
	return w.AddResourceAt(name, defaultResourcePath, reader)
}

func (w *workflowProjectHandler) AddResourceAt(name, path string, reader io.Reader) WorkflowProjectHandler {
	for _, r := range w.rawResources[path] {
		if r.name == name {
			r.contents = reader
			return w
		}
	}
	w.rawResources[path] = append(w.rawResources[path], &resource{name: name, contents: reader})
	w.parsed = false
	return w
}

func (w *workflowProjectHandler) SaveAsKubernetesManifests(path string) error {
	if err := ensurePath(path); err != nil {
		return err
	}
	if err := w.parseRawProject(); err != nil {
		return err
	}
	fileCount := 1
	if err := saveAsKubernetesManifest(w.project.Workflow, path, fmt.Sprintf("%02d-", 1)); err != nil {
		return err
	}
	for i, r := range w.project.Resources {
		fileCount = i + 1
		if err := saveAsKubernetesManifest(r, path, fmt.Sprintf("%02d-", fileCount)); err != nil {
			return err
		}
	}
	fileCount++
	if err := saveAsKubernetesManifest(w.project.Properties, path, fmt.Sprintf("%02d-", fileCount)); err != nil {
		return err
	}
	return nil
}

func (w *workflowProjectHandler) AsObjects() (*WorkflowProject, error) {
	if err := w.parseRawProject(); err != nil {
		return nil, err
	}
	return &w.project, nil
}

func (w *workflowProjectHandler) parseRawProject() error {
	if w.parsed {
		return nil
	}
	if err := w.sanityCheck(); err != nil {
		return err
	}
	if err := w.parseRawWorkflow(); err != nil {
		return err
	}
	if err := w.parseRawAppProperties(); err != nil {
		return err
	}
	if err := w.parseRawResources(); err != nil {
		return err
	}
	w.parsed = true
	return nil
}

func (w *workflowProjectHandler) sanityCheck() error {
	if len(w.namespace) == 0 {
		return errors.New("Namespace is required when building Workflow projects")
	}
	if w.rawWorkflow == nil {
		return errors.New("A workflow reader pointer is required when building Workflow projects")
	}
	return nil
}

func (w *workflowProjectHandler) parseRawWorkflow() error {
	workflowContents, err := io.ReadAll(w.rawWorkflow)
	if err != nil {
		return err
	}
	var workflowDef *model.Workflow
	// TODO: add this to the SDK, also an input from io.Reader
	workflowDef, err = parser.FromJSONSource(workflowContents)
	if err != nil {
		workflowDef, err = parser.FromYAMLSource(workflowContents)
		if err != nil {
			return errors.Errorf("Failed to parse the workflow either as a JSON or as a YAML file: %+v", err)
		}
	}

	if len(w.name) == 0 {
		w.name = strings.ToLower(workflowDef.ID)
	}

	w.project.Workflow, err = operatorapi.FromCNCFWorkflow(workflowDef, context.TODO())
	w.project.Workflow.Name = w.name
	w.project.Workflow.Namespace = w.namespace

	SetWorkflowProfile(w.project.Workflow, metadata.DevProfile)
	SetDefaultLabels(w.project.Workflow, w.project.Workflow)
	if err = SetTypeToObject(w.project.Workflow, w.scheme); err != nil {
		return err
	}

	return nil
}

func (w *workflowProjectHandler) parseRawAppProperties() error {
	if w.rawAppProperties == nil {
		return nil
	}
	appPropsContent, err := io.ReadAll(w.rawAppProperties)
	if err != nil {
		return err
	}
	w.project.Properties = CreateNewAppPropsConfigMap(w.project.Workflow, string(appPropsContent))
	if err = SetTypeToObject(w.project.Properties, w.scheme); err != nil {
		return err
	}
	return nil
}

func (w *workflowProjectHandler) parseRawResources() error {
	if len(w.rawResources) == 0 {
		return nil
	}

	resourceCount := 1
	for path, resources := range w.rawResources {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Namespace: w.namespace, Name: fmt.Sprintf("%02d-%s-resources", resourceCount, w.name)},
			Data:       map[string]string{},
		}

		for _, r := range resources {
			contents, err := io.ReadAll(r.contents)
			if err != nil {
				return err
			}
			if len(contents) == 0 {
				return errors.Errorf("Content for the resource %s is empty. Can't add an empty resource to the workflow project", r.name)
			}
			cm.Data[r.name] = string(contents)
		}

		if err := w.addResourceConfigMapToProject(cm, path); err != nil {
			return err
		}
		resourceCount++
	}

	return nil
}

func (w *workflowProjectHandler) addResourceConfigMapToProject(cm *corev1.ConfigMap, path string) error {
	if cm.Data != nil {
		if err := SetTypeToObject(cm, w.scheme); err != nil {
			return err
		}
		w.project.Workflow.Spec.Resources.ConfigMaps = append(w.project.Workflow.Spec.Resources.ConfigMaps,
			operatorapi.ConfigMapWorkflowResource{ConfigMap: corev1.LocalObjectReference{Name: cm.Name}, WorkflowPath: path})
		w.project.Resources = append(w.project.Resources, cm)
	}
	return nil
}
