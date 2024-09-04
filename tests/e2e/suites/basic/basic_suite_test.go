package basic

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/logger"

	"k8s.io/apimachinery/pkg/types"
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

					if release.Status.LastAppliedRevision == strings.TrimPrefix(os.Getenv("E2E_APP_VERSION"), "v") {
						return true, nil
					}
					return false, errors.New("err")
				}).
					WithTimeout(1 * time.Minute).
					WithPolling(5 * time.Second).
					Should(BeTrue())
			})
		}).
		Run(t, "Basic Test")
}
