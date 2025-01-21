package basic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/logger"
	"github.com/giantswarm/clustertest/pkg/wait"

	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/cilium/cilium/cilium-cli/api"
	"github.com/cilium/cilium/cilium-cli/connectivity"
	"github.com/cilium/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium/cilium-cli/k8s"
)

const (
	isUpgrade = false
)

func TestBasic(t *testing.T) {
	const (
		appReadyTimeout  = 10 * time.Minute
		appReadyInterval = 5 * time.Second
	)
	suite.New(config.MustLoad("../../config.yaml")).
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

					release := &helmv2beta1.HelmRelease{}
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

				mcClient := state.GetFramework().MC()
				wcNamespace := state.GetCluster().Organization.GetNamespace()
				wcName := state.GetCluster().Name
				wcClient, _ := state.GetFramework().WC(wcName)

				By("Waiting for all DaemonSets to be ready")
				Eventually(
					wait.ConsistentWaitCondition(
						wait.AreAllDaemonSetsReady(state.GetContext(), wcClient),
						10,
						time.Second,
					)).
					WithTimeout(15 * time.Minute).
					WithPolling(wait.DefaultInterval).
					Should(Succeed())

				By("Applying policy exceptions")
				polex := buildPolicyException()
				err := wcClient.Create(context.Background(), polex)
				Expect(err).ShouldNot(HaveOccurred())

				By("Running connectivity tests")
				ciliumNamespace := "kube-system"
				params := buildConnectivityTestParams()
				hooks := &api.NopHooks{}
				tmpKubeconfig := fmt.Sprintf("/tmp/kubeconfig-%s", wcName)

				kubeconfig, err := mcClient.GetClusterKubeConfig(context.Background(), wcName, wcNamespace)
				err = os.WriteFile(tmpKubeconfig, []byte(kubeconfig), 0644)
				Expect(err).ShouldNot(HaveOccurred())

				k8sClient, err := k8s.NewClient("", tmpKubeconfig, ciliumNamespace, "", nil)
				Expect(err).ShouldNot(HaveOccurred())
				ctx := api.SetNamespaceContextValue(context.Background(), ciliumNamespace)
				ctx = api.SetK8sClientContextValue(ctx, k8sClient)

				logger := check.NewConcurrentLogger(params.Writer, params.TestConcurrency)
				logger.Start(ctx)
				defer logger.Stop()

				connTests, err := newConnectivityTests(k8sClient, params, logger)
				Expect(err).ShouldNot(HaveOccurred())

				err = connectivity.Run(ctx, connTests, hooks)
				Expect(err).ShouldNot(HaveOccurred())
			})
		}).
		AfterSuite(func() {
			wcName := state.GetCluster().Name
			wcClient, _ := state.GetFramework().WC(wcName)

			By("Deleting cilium-test-1 namespace")
			polex := buildPolicyException()
			err := wcClient.Delete(context.Background(), polex)
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
