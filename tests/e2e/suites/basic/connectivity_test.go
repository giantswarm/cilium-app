package basic

import (
	"os"
	"regexp"

	"github.com/cilium/cilium/cilium-cli/api"
	"github.com/cilium/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium/cilium-cli/k8s"
	"github.com/cilium/cilium/cilium-cli/sysdump"
)

func buildTestParams() check.Parameters {

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
		CurlImage:                     defaults.ConnectivityCheckAlpineCurlImage,
		CurlInsecure:                  false,
		DNSTestServerImage:            defaults.ConnectivityDNSTestServerImage,
		Debug:                         false,
		EchoServerHostPort:            4000,
		ExpectedDropReasons:           defaults.ExpectedDropReasons,
		ExpectedXFRMErrors:            defaults.ExpectedXFRMErrors,
		ExternalCIDR:                  "1.0.0.0/8",
		ExternalDeploymentPort:        8080,
		ExternalIP:                    "1.1.1.1",
		ExternalOtherIP:               "1.0.0.1",
		ExternalTarget:                "one.one.one.one.",
		ExternalTargetCAName:          "cabundle",
		ExternalTargetCANamespace:     "cilium-test-1",
		FRRImage:                      defaults.ConnectivityTestFRRImage,
		FlowValidation:                check.FlowValidationModeWarning,
		FlushCT:                       false,
		ForceDeploy:                   false,
		HelmChartDirectory:            "",
		HelmValuesSecretName:          defaults.HelmValuesSecretName,
		Hubble:                        true,
		HubbleServer:                  "localhost:4245",
		IncludeConnDisruptTest:        false,
		IncludeUnsafeTests:            false,
		JSONMockImage:                 defaults.ConnectivityCheckJSONMockImage,
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
		TestConnDisruptImage:          defaults.ConnectivityTestConnDisruptImage,
		TestNamespace:                 "cilium-test-1",
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

	// Use this regex to filter which tests to run.
	rgx, _ := regexp.Compile("no-policies")
	params.RunTests = append(params.RunTests, rgx)

	return params

}

func newConnectivityTests(client *k8s.Client, p check.Parameters, logger *check.ConcurrentLogger) ([]*check.ConnectivityTest, error) {
	hooks := &api.NopHooks{}
	cc, err := check.NewConnectivityTest(client, p, hooks, logger)
	if err != nil {
		return nil, err
	}

	connTests := make([]*check.ConnectivityTest, 0, p.TestConcurrency)
	connTests = append(connTests, cc)
	return connTests, nil
}
