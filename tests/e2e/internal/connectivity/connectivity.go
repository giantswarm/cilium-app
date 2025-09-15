package connectivity

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/cilium/cilium/cilium-cli/api"
	"github.com/cilium/cilium/cilium-cli/connectivity"
	"github.com/cilium/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium/cilium-cli/utils/codeowners"
	"github.com/cilium/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium/cilium-cli/k8s"
	"github.com/cilium/cilium/cilium-cli/sysdump"

	"github.com/giantswarm/apptest-framework/pkg/state"
)

func buildParams() check.Parameters {

	params := defaultConnectivityTestParams()

	// Custom settings for E2E tests:
	params.Hubble = false
	params.TestNamespace = "cilium-test-1"

	// Use this regex to filter which tests to run.
	rgx, _ := regexp.Compile("no-policies")
	params.RunTests = append(params.RunTests, rgx)

	return params
}

func newTests(client *k8s.Client, p check.Parameters, logger *check.ConcurrentLogger) ([]*check.ConnectivityTest, error) {
	hooks := &api.NopHooks{}
	ruleset := &codeowners.Ruleset{}
	cc, err := check.NewConnectivityTest(client, p, hooks, logger, ruleset)
	if err != nil {
		return nil, err
	}

	connTests := make([]*check.ConnectivityTest, 0, p.TestConcurrency)
	connTests = append(connTests, cc)
	return connTests, nil
}

func Run(wcNamespace, wcName string) error {
	ciliumNamespace := "kube-system"
	params := buildParams()
	hooks := &api.NopHooks{}
	tmpKubeconfig := fmt.Sprintf("/tmp/kubeconfig-%s", wcName)

	mcClient := state.GetFramework().MC()
	kubeconfig, err := mcClient.GetClusterKubeConfig(context.Background(), wcName, wcNamespace)
	err = os.WriteFile(tmpKubeconfig, []byte(kubeconfig), 0644)
	if err != nil {
		return err
	}

	k8sClient, err := k8s.NewClient("", tmpKubeconfig, ciliumNamespace, "", nil)
	if err != nil {
		return err
	}

	ctx := api.SetNamespaceContextValue(context.Background(), ciliumNamespace)
	ctx = api.SetK8sClientContextValue(ctx, k8sClient)

	logger := check.NewConcurrentLogger(params.Writer)
	logger.Start()
	defer logger.Stop()

	connTests, err := newTests(k8sClient, params, logger)
	if err != nil {
		return err
	}

	err = connectivity.Run(ctx, connTests, hooks)
	if err != nil {
		return err
	}

	return nil
}

func defaultConnectivityTestParams() check.Parameters {

	// parameters taken from https://github.com/cilium/cilium/blob/main/cilium-cli/cli/connectivity.go#L101
	// Defaults - Do not replace them here. Add an override in the buildConnectivityTestParam function.
	var params = check.Parameters{
		AgentDaemonSetName:            defaults.AgentDaemonSetName,
		AgentPodSelector:              defaults.AgentPodSelector,
		AllFlows:                      false,
		AssumeCiliumVersion:           "",
		CiliumNamespace:               "kube-system",
		CiliumPodSelector:             defaults.CiliumPodSelector,
		CollectSysdumpOnFailure:       false,
		ConnDisruptDispatchInterval:   0,
		ConnDisruptTestRestartsPath:   "/tmp/cilium-conn-disrupt-restarts",
		ConnDisruptTestSetup:          false,
		ConnDisruptTestXfrmErrorsPath: "/tmp/cilium-conn-disrupt-xfrm-errors",
		ConnectTimeout:                defaults.ConnectTimeout,
		CurlImage:                     defaults.ConnectivityCheckImagesTest["ConnectivityCheckAlpineCurlImage"],
		CurlInsecure:                  false,
		DNSTestServerImage:            defaults.ConnectivityCheckImagesTest["ConnectivityDNSTestServerImage"],
		Debug:                         false,
		EchoServerHostPort:            4000,
		ExpectedDropReasons:           defaults.ExpectedDropReasons,
		ExpectedXFRMErrors:            defaults.ExpectedXFRMErrors,
		ExternalCIDRv4:                "1.0.0.0/8",
		ExternalCIDRv6:                "2606:4700:4700::/96",
		ExternalDeploymentPort:        8080,
		ExternalIPv4:                  "1.1.1.1",
		ExternalIPv6:                  "2606:4700:4700::1111",
		ExternalOtherIPv4:             "1.0.0.1",
		ExternalOtherIPv6:             "2606:4700:4700::1001",
		ExternalTarget:                "one.one.one.one.",
		ExternalTargetCAName:          "cabundle",
		ExternalTargetCANamespace:     "cilium-test-1",
		FRRImage:                      defaults.ConnectivityCheckImagesTest["ConnectivityDNSTestServerImage"],
		FlowValidation:                check.FlowValidationModeWarning,
		FlushCT:                       false,
		ForceDeploy:                   false,
		HelmChartDirectory:            "",
		HelmValuesSecretName:          defaults.HelmValuesSecretName,
		Hubble:                        true,
		HubbleServer:                  "localhost:4245",
		IncludeConnDisruptTest:        false,
		IncludeUnsafeTests:            false,
		JSONMockImage:                 defaults.ConnectivityCheckImagesTest["ConnectivityCheckJSONMockImage"],
		JunitFile:                     "",
		JunitProperties:               make(map[string]string),
		K8sLocalHostTest:              false,
		K8sVersion:                    "",
		MultiCluster:                  "",
		NodeCIDRs:                     nil,
		NodeSelector:                  make(map[string]string),
		PauseOnFail:                   false,
		PostTestSleepDuration:         0,
		PrintFlows:                    false,
		RequestTimeout:                defaults.RequestTimeout,
		Retry:                         defaults.ConnectRetry,
		RetryDelay:                    defaults.ConnectRetryDelay,
		SecondaryNetworkIface:         "",
		ServiceType:                   "NodePort",
		SingleNode:                    false,
		SkipIPCacheCheck:              true,
		TestConcurrency:               1,
		TestConnDisruptImage:          defaults.ConnectivityCheckImagesTest["ConnectivityTestConnDisruptImage"],
		TestNamespaceIndex:            0,
		Timeout:                       defaults.ConnectivityTestSuiteTimeout,
		Timestamp:                     false,
		Verbose:                       false,
		Writer:                        os.Stdout,
		SysdumpOptions: sysdump.Options{
			LargeSysdumpAbortTimeout: sysdump.DefaultLargeSysdumpAbortTimeout,
			LargeSysdumpThreshold:    sysdump.DefaultLargeSysdumpThreshold,
			Writer:                   os.Stdout,
		},
	}

	return params
}
