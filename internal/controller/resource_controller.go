/*
Copyright 2026.

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

package controller

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	idpv1alpha1 "github.com/jjsiv/idp/api/v1alpha1"
	"github.com/jjsiv/idp/internal/provisioner"
)

// ResourceReconciler reconciles a Resource object
type ResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	logr.Logger
}

// +kubebuilder:rbac:groups=idp.autopay.pl,resources=resources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=idp.autopay.pl,resources=resourcetemplates,verbs=get;list;watch
// +kubebuilder:rbac:groups=idp.autopay.pl,resources=resources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=idp.autopay.pl,resources=resources/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile
func (r *ResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Logger = logf.FromContext(ctx).WithValues("name", req.Name, "namespace", req.Namespace)

	var resource idpv1alpha1.Resource
	if err := r.Client.Get(ctx, req.NamespacedName, &resource); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	templateNamespace := resource.Namespace
	if resource.Spec.TemplateRef.Namespace != "" {
		templateNamespace = resource.Spec.TemplateRef.Namespace
	}
	var template idpv1alpha1.ResourceTemplate
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name:      resource.Spec.TemplateRef.Name,
		Namespace: templateNamespace,
	}, &template); err != nil {
		r.Logger.Error(err, "failed to get ResourceTemplate")
		return ctrl.Result{}, err
	}

	provisioners, err := r.setupProvisioners(ctx, &resource, &template)
	if err != nil {
		r.Logger.Error(err, "failed to setup provisioners")
		return ctrl.Result{}, fmt.Errorf("failed to setup provisioners: %w", err)
	}

	if !resource.ObjectMeta.DeletionTimestamp.IsZero() {
		for _, p := range provisioners {
			p.Deprovision()
		}

		return ctrl.Result{}, nil
	}

	for _, p := range provisioners {
		if err := p.Provision(); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to provision resource: %w", err)
		}
	}

	patch := client.MergeFrom(resource.DeepCopy())
	resource.Status.Provisioned = true
	if err := r.Status().Patch(ctx, &resource, patch); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to patch resource status: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *ResourceReconciler) setupProvisioners(ctx context.Context, resource *idpv1alpha1.Resource, template *idpv1alpha1.ResourceTemplate) ([]provisioner.Provisioner, error) {
	var provisioners []provisioner.Provisioner

	if template.Spec.Provisioning.Git != nil {
		for _, gitProvisioner := range template.Spec.Provisioning.Git {
			secretNamespace := template.Namespace
			if gitProvisioner.Config.Auth.SecretRef.Namespace != "" {
				secretNamespace = gitProvisioner.Config.Auth.SecretRef.Namespace
			}
			var keySecret v1.Secret
			if err := r.Client.Get(ctx, types.NamespacedName{
				Name:      gitProvisioner.Config.Auth.SecretRef.Name,
				Namespace: secretNamespace,
			}, &keySecret); err != nil {
				return nil, fmt.Errorf("failed to retrieve authentication secret: %s", err)
			}

			privateKey, ok := keySecret.Data[gitProvisioner.Config.Auth.SecretRef.Key]
			if !ok {
				return nil, fmt.Errorf("key %s not found in secret %s", gitProvisioner.Config.Auth.SecretRef.Key, gitProvisioner.Config.Auth.SecretRef.Name)
			}

			repo := provisioner.GitRepository{
				URL: gitProvisioner.RepositoryURL,
				Ref: gitProvisioner.Branch,
				CommitAuthor: &provisioner.CommitAuthor{
					Name:  gitProvisioner.Config.Author.Name,
					Email: gitProvisioner.Config.Author.Email,
				},
				KeyAuth: &provisioner.GitKeyAuth{
					PrivateKey: privateKey,
				},
			}

			files, err := renderFileTemplates(resource.Spec.TemplateRef.Parameters, gitProvisioner.Templates)
			if err != nil {
				return nil, err
			}

			p, err := provisioner.NewGitProvisioner(&repo, files)
			if err != nil {
				return nil, fmt.Errorf("failed to create git provisioner: %w", err)
			}

			provisioners = append(provisioners, p)
		}
	}

	return provisioners, nil
}

func renderFileTemplates(params []idpv1alpha1.ResourceParameter, fileTemplates []idpv1alpha1.FileTemplate) ([]provisioner.GitFile, error) {
	paramMap := make(map[string]any, len(params))
	for _, param := range params {
		paramMap[param.Name] = param.Value
	}

	variables := map[string]any{
		"parameters": paramMap,
	}

	var files []provisioner.GitFile
	for _, fileTemplate := range fileTemplates {
		// TODO: there is probably a better way to do this
		filename, err := templateText(variables, fileTemplate.Filename)
		if err != nil {
			return nil, err
		}

		path, err := templateText(variables, fileTemplate.Path)
		if err != nil {
			return nil, err
		}

		content, err := templateText(variables, fileTemplate.Content)
		if err != nil {
			return nil, err
		}

		files = append(files, provisioner.GitFile{
			Path:    string(path) + "/" + string(filename),
			Content: content,
		})
	}

	return files, nil
}

func templateText(vars map[string]any, text string) ([]byte, error) {
	t, err := template.New("").Parse(text)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&idpv1alpha1.Resource{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Named("resource").
		Complete(r)
}
