/*
Copyright 2025.

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

package v1alpha1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	cachev1alpha1 "popgo.domain/memcached-operator/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var memcachedlog = logf.Log.WithName("memcached-resource")

// SetupMemcachedWebhookWithManager registers the webhook for Memcached in the manager.
func SetupMemcachedWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&cachev1alpha1.Memcached{}).
		WithValidator(&MemcachedCustomValidator{}).
		WithDefaulter(&MemcachedCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-cache-popgo-domain-v1alpha1-memcached,mutating=true,failurePolicy=fail,sideEffects=None,groups=cache.popgo.domain,resources=memcacheds,verbs=create;update,versions=v1alpha1,name=mmemcached-v1alpha1.kb.io,admissionReviewVersions=v1

// MemcachedCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Memcached when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type MemcachedCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &MemcachedCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Memcached.
func (d *MemcachedCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	memcached, ok := obj.(*cachev1alpha1.Memcached)

	if !ok {
		return fmt.Errorf("expected an Memcached object but got %T", obj)
	}
	memcachedlog.Info("Defaulting for Memcached", "name", memcached.GetName())

	// case 1: set default value for spec.size
	if memcached.Spec.Size == nil {
		// Size *int32 string
		defaultSize := int32(3)
		memcached.Spec.Size = &defaultSize
		memcachedlog.Info("Defaulted spec.size to 3", "name", memcached.GetName())
	}

	// case 2: set default value for metadata.annotations
	if memcached.GetAnnotations() == nil {
		memcached.SetAnnotations(make(map[string]string))
	}
	// make sure deployment-strategy exists alway
	if _, ok := memcached.GetAnnotations()["cache.popgo.domain/deployment-strategy"]; !ok {
		memcached.GetAnnotations()["cache.popgo.domain/deployment-strategy"] = "rolling-update"
		memcachedlog.Info("Defaulted deployment-strategy annotation", "name", memcached.GetName())
	}

	// TODO(user): fill in your defaulting logic.

	return nil
}

// validateMemcachedSpec contains shared validation logic for Create and Update
func validateMemcachedSpec(memcached *cachev1alpha1.Memcached) error {
	// validate that the value of Size is between 1 and 3
	if memcached.Spec.Size == nil {
		return fmt.Errorf("spec.size is required and must not be nil")
	}

	size := *memcached.Spec.Size
	const minSize = 1
	const maxSize = 3

	if size < minSize || size > maxSize {
		return fmt.Errorf("spec.size must be between %d and %d (inclusive), but got %d", minSize, maxSize, size)
	}

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-cache-popgo-domain-v1alpha1-memcached,mutating=false,failurePolicy=fail,sideEffects=None,groups=cache.popgo.domain,resources=memcacheds,verbs=create;update;delete,versions=v1alpha1,name=vmemcached-v1alpha1.kb.io,admissionReviewVersions=v1

// MemcachedCustomValidator struct is responsible for validating the Memcached resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type MemcachedCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &MemcachedCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Memcached.
func (v *MemcachedCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	memcached, ok := obj.(*cachev1alpha1.Memcached)
	if !ok {
		return nil, fmt.Errorf("expected a Memcached object but got %T", obj)
	}
	memcachedlog.Info("Validation for Memcached upon creation", "name", memcached.GetName())

	//
	if err := validateMemcachedSpec(memcached); err != nil {
		return nil, err
	}

	// 2. case A-1
	if memcached.GetName() == "forbidden-name" {
		return nil, fmt.Errorf("the Memcached name '%s' is explicitly forbidden by policy", memcached.GetName())
	}

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Memcached.
func (v *MemcachedCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	newMemcached, ok := newObj.(*cachev1alpha1.Memcached)
	if !ok {
		return nil, fmt.Errorf("expected a Memcached object for the newObj but got %T", newObj)
	}
	oldMemcached, ok := oldObj.(*cachev1alpha1.Memcached)
	if !ok {
		return nil, fmt.Errorf("expected a Memcached object for the oldObj but got %T", oldObj)
	}

	memcachedlog.Info("Validation for Memcached upon update", "name", newMemcached.GetName())

	//
	if err := validateMemcachedSpec(newMemcached); err != nil {
		return nil, err
	}

	// 2. case B-1
	if *newMemcached.Spec.Size < *oldMemcached.Spec.Size {
		return nil, fmt.Errorf(
			"scaling down Memcached is not allowed to prevent potential data loss. Attempted scale down from %d to %d",
			*oldMemcached.Spec.Size,
			*newMemcached.Spec.Size,
		)
	}

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Memcached.
func (v *MemcachedCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	memcached, ok := obj.(*cachev1alpha1.Memcached)
	if !ok {
		return nil, fmt.Errorf("expected a Memcached object but got %T", obj)
	}
	memcachedlog.Info("Validation for Memcached upon deletion", "name", memcached.GetName())

	// case C-1
	annotations := memcached.GetAnnotations()
	if annotations != nil {
		if val, exists := annotations["popgo.domain/deletion-protection"]; exists && val == "true" {
			//
			return nil, fmt.Errorf("deletion of Memcached %s is blocked by annotation 'popgo.domain/deletion-protection: true'. Please remove the annotation to proceed", memcached.GetName())
		}
	}

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
