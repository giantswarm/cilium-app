package basic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/v5/pkg/state"
	"github.com/giantswarm/apptest-framework/v5/pkg/suite"
	"github.com/giantswarm/clustertest/v5/pkg/logger"
	"github.com/giantswarm/clustertest/v5/pkg/wait"

	"github.com/giantswarm/cilium-app/tests/e2e/internal/connectivity"
	"github.com/giantswarm/cilium-app/tests/e2e/internal/metrics"
	"github.com/giantswarm/cilium-app/tests/e2e/internal/polex"
	"github.com/giantswarm/cilium-app/tests/e2e/internal/readiness"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	isUpgrade = false
)

func TestBasic(t *testing.T) {
	const (
		appReadyTimeout  = 10 * time.Minute
		appReadyInterval = 5 * time.Second
	)
	suite.New().
		WithInstallNamespace("kube-system").
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		AfterClusterReady(func() {
			// no
		}).
		BeforeUpgrade(func() {
			// E.g. ensure that the initial install has completed and has settled before upgrading
		}).
		Tests(func() {
			It("should deploy the HelmRelease", func() {

				Eventually(func() (bool, error) {
					appNamespace := state.GetCluster().Organization.GetNamespace()
					appName := fmt.Sprintf("%s-cilium", state.GetCluster().Name)

					mcKubeClient := state.GetFramework().MC()

					logger.Log("HelmRelease: %s/%s", appNamespace, appName)

					release := &helmv2.HelmRelease{}
					err := mcKubeClient.Get(state.GetContext(), types.NamespacedName{Name: appName, Namespace: appNamespace}, release)
					if err != nil {
						return false, err
					}

					for _, c := range release.Status.Conditions {
						if c.Type == "Ready" {
							if c.Status == "True" {
								return true, nil
							} else {
								return false, errors.New(fmt.Sprintf("HelmRelease not ready [%s]: %s", c.Reason, c.Message))
							}
						}
					}

					return false, errors.New("HelmRelease not ready")
				}).
					WithTimeout(5 * time.Minute).
					WithPolling(15 * time.Second).
					Should(BeTrue())
			})

			It("should pass cilium connectivity test", func() {

				wcNamespace := state.GetCluster().Organization.GetNamespace()
				wcName := state.GetCluster().Name
				wcClient, _ := state.GetFramework().WC(wcName)

				By("Waiting for all Nodes and DaemonSets to be ready and stable")
				// A stricter gate than wait.AreAllDaemonSetsReady: it requires every
				// Node to be Ready and every DaemonSet pod to be available (not just
				// scheduled), held consistently over a ~60s window. This guards against
				// transient CAPA worker-node churn that would otherwise delete a Cilium
				// agent pod out from under the connectivity test.
				Eventually(
					wait.ConsistentWaitCondition(
						readiness.AllNodesAndDaemonSetsReady(state.GetContext(), wcClient),
						30,
						2*time.Second,
					)).
					WithTimeout(15 * time.Minute).
					WithPolling(wait.DefaultInterval).
					Should(Succeed())

				By("Applying policy exceptions")
				p := polex.New()
				err := wcClient.Create(context.Background(), p)
				Expect(err).ShouldNot(HaveOccurred())

				By("Running connectivity tests")
				// cilium-cli snapshots the agent pod list once at init and execs into
				// those pods seconds later. If a node vanishes in that window the run
				// fails with `pods "<node>" not found`. That staleness is self-healing on
				// a fresh attempt, so retry the whole run a few times before failing.
				Eventually(func() error {
					return connectivity.Run(wcNamespace, wcName)
				}).
					WithTimeout(15 * time.Minute).
					WithPolling(30 * time.Second).
					Should(Succeed())
			})

			It("should create the hubble-generate-certs CronJob and renew the Hubble certificates", func() {
				const (
					certgenNamespace = "kube-system"
					cronJobName      = "hubble-generate-certs"
					serverCertSecret = "hubble-server-certs"
					expectedSchedule = "*/1 * * * *"

					// Generous headroom for a certgen Job to be scheduled, pull its
					// image and run to completion.
					waitTimeout  = 6 * time.Minute
					pollInterval = 15 * time.Second
				)

				ctx := state.GetContext()
				wcName := state.GetCluster().Name
				wcClient, err := state.GetFramework().WC(wcName)
				Expect(err).ShouldNot(HaveOccurred())

				By("Ensuring the hubble-generate-certs CronJob is created")
				cronJob := &batchv1.CronJob{}
				Eventually(func() error {
					return wcClient.Get(ctx, types.NamespacedName{Name: cronJobName, Namespace: certgenNamespace}, cronJob)
				}).
					WithTimeout(waitTimeout).
					WithPolling(pollInterval).
					Should(Succeed(), "expected the hubble-generate-certs CronJob to be created")
				Expect(cronJob.Spec.Schedule).To(Equal(expectedSchedule))

				By("Waiting for a CronJob-triggered certgen Job to complete successfully")
				Eventually(func() (bool, error) {
					jobs := &batchv1.JobList{}
					if listErr := wcClient.List(ctx, jobs, cr.InNamespace(certgenNamespace)); listErr != nil {
						return false, listErr
					}
					for _, job := range jobs.Items {
						// Only count Jobs owned by the CronJob (scheduled runs), not the
						// one-shot install Job that Helm creates at release time (that Job
						// has no CronJob owner reference). Note the k8s-app label lives on
						// the pod template, not the Job object, so we match on the owner.
						for _, owner := range job.OwnerReferences {
							if owner.Kind == "CronJob" && owner.Name == cronJobName && job.Status.Succeeded > 0 {
								logger.Log("Certgen Job %q (owned by CronJob %q) succeeded", job.Name, owner.Name)
								return true, nil
							}
						}
					}
					return false, nil
				}).
					WithTimeout(waitTimeout).
					WithPolling(pollInterval).
					Should(BeTrue(), "expected the hubble-generate-certs CronJob to spawn a successful certgen Job")

				By("Recording the current hubble-server-certs certificate")
				var previousCert []byte
				Eventually(func() (bool, error) {
					secret := &corev1.Secret{}
					if getErr := wcClient.Get(ctx, types.NamespacedName{Name: serverCertSecret, Namespace: certgenNamespace}, secret); getErr != nil {
						return false, getErr
					}
					crt, ok := secret.Data["tls.crt"]
					if !ok || len(crt) == 0 {
						return false, nil
					}
					previousCert = crt
					return true, nil
				}).
					WithTimeout(waitTimeout).
					WithPolling(pollInterval).
					Should(BeTrue(), "expected the hubble-server-certs Secret to hold a certificate")

				// certgen always generates fresh key material and upserts the Secret on
				// every run, so any completed run rotates the certificate. Rather than
				// wait on the CronJob's schedule, we deterministically trigger a run from
				// the CronJob's jobTemplate (the equivalent of
				// `kubectl create job --from=cronjob/hubble-generate-certs`).
				By("Triggering a certgen Job from the CronJob's jobTemplate")
				triggeredJob := &batchv1.Job{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: cronJobName + "-e2e-renew-",
						Namespace:    certgenNamespace,
						Labels:       cronJob.Spec.JobTemplate.Labels,
						Annotations:  cronJob.Spec.JobTemplate.Annotations,
					},
					Spec: cronJob.Spec.JobTemplate.Spec,
				}
				Expect(wcClient.Create(ctx, triggeredJob)).To(Succeed())
				logger.Log("Created certgen Job %q", triggeredJob.Name)
				DeferCleanup(func() {
					_ = wcClient.Delete(ctx, triggeredJob, cr.PropagationPolicy(metav1.DeletePropagationBackground))
				})

				By("Waiting for the triggered certgen Job to complete successfully")
				Eventually(func() (bool, error) {
					job := &batchv1.Job{}
					if getErr := wcClient.Get(ctx, types.NamespacedName{Name: triggeredJob.Name, Namespace: certgenNamespace}, job); getErr != nil {
						return false, getErr
					}
					return job.Status.Succeeded > 0, nil
				}).
					WithTimeout(waitTimeout).
					WithPolling(pollInterval).
					Should(BeTrue(), "expected the triggered certgen Job to succeed")

				By("Verifying the hubble-server-certs certificate has been renewed")
				Eventually(func() (bool, error) {
					secret := &corev1.Secret{}
					if getErr := wcClient.Get(ctx, types.NamespacedName{Name: serverCertSecret, Namespace: certgenNamespace}, secret); getErr != nil {
						return false, getErr
					}
					crt, ok := secret.Data["tls.crt"]
					if !ok || len(crt) == 0 {
						return false, nil
					}
					// A freshly generated cert uses a new key pair, so it must differ
					// from the one recorded before the run.
					return !bytes.Equal(crt, previousCert), nil
				}).
					WithTimeout(waitTimeout).
					WithPolling(pollInterval).
					Should(BeTrue(), "expected the certgen Job to rotate the hubble server certificate")
			})

			It("ensure key metrics are available on mimir", func() {
				const mimirUrl = "mimir-gateway.mimir.svc:80/prometheus"
				mcClient := state.GetFramework().MC()
				expectedMetrics := []string{
					// Cilium Agent metrics
					"cilium_version",

					// Cilium Operator metrics
					"cilium_operator_version",
				}

				By("Creating a test pod")
				// Run a pod with alpine in the default namespace of the MC.
				testPodName := fmt.Sprintf("%s-metrics-test", state.GetCluster().Name)
				testPodNamespace := "default"

				err := metrics.Run(mcClient, testPodName, testPodNamespace)
				Expect(err).NotTo(HaveOccurred())

				By("ensuring that metrics are present in Mimir")
				for _, metric := range expectedMetrics {
					Eventually(metrics.Check(mcClient, metric, mimirUrl, testPodName, testPodNamespace)).
						WithTimeout(10 * time.Minute).
						WithPolling(10 * time.Second).
						Should(Succeed())
				}

				By("Cleaning up test pod")
				err = metrics.Cleanup(mcClient, testPodName, testPodNamespace)
				Expect(err).NotTo(HaveOccurred())
			})
		}).
		AfterSuite(func() {
			wcName := state.GetCluster().Name
			wcClient, _ := state.GetFramework().WC(wcName)

			By("Deleting cilium-test-1 namespace")
			p := polex.New()
			err := wcClient.Delete(context.Background(), p)
			Expect(err).ShouldNot(HaveOccurred())

			By("Deleting cilium-test PolicyException")
			testNamespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "cilium-test-1",
				},
			}
			err = wcClient.Delete(context.Background(), testNamespace)
			Expect(err).ShouldNot(HaveOccurred())
		}).
		Run(t, "Basic Test")
}
