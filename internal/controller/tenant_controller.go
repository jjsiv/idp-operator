/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and
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
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	yaml "github.com/goccy/go-yaml"
	idpv1alpha1 "github.com/jjsiv/idp/api/v1alpha1"
	"github.com/jjsiv/idp/internal/backstage"
	"github.com/jjsiv/idp/internal/provisioner"
	"github.com/jjsiv/idp/internal/utils"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	logr.Logger
}

// +kubebuilder:rbac:groups=idp.autopay.pl,resources=tenants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=idp.autopay.pl,resources=tenants/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=idp.autopay.pl,resources=tenants/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile
func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Logger = logf.FromContext(ctx).WithValues("tenant", req.NamespacedName)
	r.Logger.Info("Reconciling Tenant")

	var tenant idpv1alpha1.Tenant
	if err := r.Client.Get(ctx, req.NamespacedName, &tenant); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	p, err := r.setupProvisioner(ctx, &tenant)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set up provisioner: %s", err.Error())
	}

	var backstageEntities []any
	topGroup, err := backstage.NewGroup().FromTenant(&tenant)
	if err != nil {
		return ctrl.Result{}, err
	}

	backstageEntities = append(backstageEntities, topGroup)

	var roles []string
	roleSet := make(utils.Set)
	for _, member := range tenant.Spec.Members {
		for _, memberRole := range member.Roles {
			if !roleSet.Has(memberRole) {
				roleSet.Insert(memberRole)
				roles = append(roles, memberRole)
			}
		}
	}

	for _, role := range roles {
		roleGroup, err := backstage.NewGroup().FromTenant(&tenant, backstage.WithGroupRole(role))
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create role subgroup for role %s: %s", role, err.Error())
		}

		backstageEntities = append(backstageEntities, roleGroup)
	}

	system, err := backstage.NewSystem().FromTenant(&tenant)
	backstageEntities = append(backstageEntities, system)

	file, err := createCatalogInfoFile(backstageEntities, &tenant)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create catalog info file: %s", err.Error())
	}

	if !tenant.ObjectMeta.DeletionTimestamp.IsZero() {
		p.Deprovision()
		return ctrl.Result{}, nil
	}

	err = p.Provision(file)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to provision catalog info file: %s", err.Error())
	}

	r.Logger.Info("Reconciliation complete")
	return ctrl.Result{}, nil
}

func createCatalogInfoFile(entities []any, tenant *idpv1alpha1.Tenant) (*provisioner.GitFile, error) {
	var contentbuf bytes.Buffer
	encoder := yaml.NewEncoder(&contentbuf)
	for _, entity := range entities {
		err := encoder.Encode(entity)
		if err != nil {
			return nil, err
		}
	}

	tenantVarMap := make(map[string]any)
	tenantVarMap["name"] = tenant.Name
	variables := map[string]any{
		"tenant": tenantVarMap,
	}

	t, err := template.New("").Parse(tenant.Spec.BackstageCatalog.Provisioning.Git.Path)
	if err != nil {
		return nil, err
	}

	var pathbuf bytes.Buffer
	err = t.Execute(&pathbuf, variables)
	if err != nil {
		return nil, err
	}

	return &provisioner.GitFile{
		Path:    pathbuf.String(),
		Content: contentbuf.Bytes(),
	}, nil
}

func (r *TenantReconciler) setupProvisioner(ctx context.Context, tenant *idpv1alpha1.Tenant) (provisioner.Provisioner, error) {
	if tenant.Spec.BackstageCatalog.Provisioning.Git != nil {
		provisionerSpec := tenant.Spec.BackstageCatalog.Provisioning.Git

		secretNamespace := tenant.Namespace
		if provisionerSpec.Config.Auth.SecretRef.Namespace != "" {
			secretNamespace = provisionerSpec.Config.Auth.SecretRef.Namespace
		}

		var keySecret v1.Secret
		if err := r.Client.Get(ctx, types.NamespacedName{
			Name:      provisionerSpec.Config.Auth.SecretRef.Name,
			Namespace: secretNamespace,
		}, &keySecret); err != nil {
			return nil, fmt.Errorf("failed to retrieve authentication secret: %s", err)
		}

		privateKey, ok := keySecret.Data[provisionerSpec.Config.Auth.SecretRef.Key]
		if !ok {
			return nil, fmt.Errorf("key %s not found in secret %s", provisionerSpec.Config.Auth.SecretRef.Key, provisionerSpec.Config.Auth.SecretRef.Name)
		}

		repo := provisioner.GitRepository{
			URL:    provisionerSpec.RepositoryURL,
			Branch: provisionerSpec.Branch,
			CommitAuthor: &provisioner.CommitAuthor{
				Name:  provisionerSpec.Config.Author.Name,
				Email: provisionerSpec.Config.Author.Email,
			},
			KeyAuth: &provisioner.GitKeyAuth{
				PrivateKey: privateKey,
			},
		}

		gitProvisioner, err := provisioner.NewGitProvisioner(&repo)
		if err != nil {
			return nil, err
		}

		return gitProvisioner, nil
	}

	return nil, fmt.Errorf("unsupported or unspecified provisioning method")
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&idpv1alpha1.Tenant{}).
		Named("tenant").
		Complete(r)
}
