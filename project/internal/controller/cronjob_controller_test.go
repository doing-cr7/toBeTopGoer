/*
Copyright 2024.

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
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	batchv1 "k8s.io/api/batch/v1"
	cronjobv1 "tutorial.kubebuilder.io/project/api/v1"
)

var _ = Describe("CronJob Controller", func() {

	const (
		CronjobName      = "test-cronjob"
		CronjobNamespace = "default"
		JobName          = "test-job"
		timeout          = time.Second * 10
		duration         = time.Second * 10
		interval         = time.Millisecond * 250
	)

	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      CronjobName,
			Namespace: CronjobNamespace, // TODO(user):Modify as needed
		}
		cronjob := &cronjobv1.CronJob{
			TypeMeta: metav1.TypeMeta{
				Kind:       "CronJob",
				APIVersion: "batch.tutorial.kubebuilder.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      CronjobName,
				Namespace: CronjobNamespace,
			},
			Spec: cronjobv1.CronJobSpec{
				Schedule: "1 * * * *",
				JobTemplate: batchv1.JobTemplateSpec{
					Spec: batchv1.JobSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name:  "test-container",
										Image: "test-image",
									},
								},
								RestartPolicy: v1.RestartPolicyOnFailure,
							},
						},
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, cronjob)).Should(Succeed())

		BeforeEach(func() {
			By("creating the custom resource for the Kind CronJob")
			err := k8sClient.Get(ctx, typeNamespacedName, cronjob)
			if err != nil && errors.IsNotFound(err) {
				resource := &batchv1.CronJob{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					// TODO(user): Specify other spec details if needed.
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &batchv1.CronJob{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance CronJob")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &CronJobReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})

	})
})

//		cronjobLookupKey := types.NamespacedName{
//			Namespace: CronjobName,
//			Name:      CronjobNamespace,
//		}
//		createdCronjob := &cronjobv1.CronJob{}
//
//		Eventually(func() bool {
//			err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
//			return err == nil
//		}, timeout, interval).Should(BeTrue())
//
//		Expect(createdCronjob.Spec.Schedule).Should(Equal("1 * * * *"))
//
//		By("By checking the CronJob has zero active Jobs")
//		Consistently(func() (int, error) {
//			err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
//			if err != nil {
//				return -1, err
//			}
//			return len(createdCronjob.Status.Active), nil
//		}, duration, interval).Should(Equal(0))
//
//		By("By creating a new Job")
//		testJob := &batchv1.Job{
//			ObjectMeta: metav1.ObjectMeta{
//				Name:      JobName,
//				Namespace: CronjobNamespace,
//			},
//			Spec: batchv1.JobSpec{
//				Template: v1.PodTemplateSpec{
//					Spec: v1.PodSpec{
//						// For simplicity, we only fill out the required fields.
//						Containers: []v1.Container{
//							{
//								Name:  "test-container",
//								Image: "test-image",
//							},
//						},
//						RestartPolicy: v1.RestartPolicyOnFailure,
//					},
//				},
//			},
//			Status: batchv1.JobStatus{
//				Active: 2,
//			},
//		}
//
//		// Note that your CronJob’s GroupVersionKind is required to set up this owner reference.
//		kind := reflect.TypeOf(cronjobv1.CronJob{}).Name()
//		gvk := cronjobv1.GroupVersion.WithKind(kind)
//
//		controllerRef := metav1.NewControllerRef(createdCronjob, gvk)
//		testJob.SetOwnerReferences([]metav1.OwnerReference{*controllerRef})
//		Expect(k8sClient.Create(ctx, testJob)).Should(Succeed())
//
//		By("By checking that the CronJob has one active Job")
//		Eventually(func() ([]string, error) {
//			err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
//			if err != nil {
//				return nil, err
//			}
//
//			names := []string{}
//			for _, job := range createdCronjob.Status.Active {
//				names = append(names, job.Name)
//			}
//			return names, nil
//		}, timeout, interval).Should(ConsistOf(JobName), "should list our active job %s in the active jobs list in status", JobName)
