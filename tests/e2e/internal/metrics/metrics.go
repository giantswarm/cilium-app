package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/giantswarm/clustertest/v3/pkg/client"
	"github.com/giantswarm/clustertest/v3/pkg/logger"
	. "github.com/onsi/gomega" //nolint:staticcheck
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client2 "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/apptest-framework/v3/pkg/state"
)

func Check(mcClient *client.Client, metric string, mimirUrl string, testPodName string, testPodNamespace string) func() error {
	return func() error {
		query := fmt.Sprintf("absent(%[1]s{cluster_id=\"%[2]s\"}) or label_replace(vector(0), \"cluster_id\", \"%[2]s\", \"\", \"\")", metric, state.GetCluster().Name)

		cmd := []string{"wget", "-O-", "-Y", "off", "--header", "X-Scope-OrgID: anonymous|giantswarm", fmt.Sprintf("%[1]s/api/v1/query?query=%[2]s", mimirUrl, url.QueryEscape(query))}
		stdout, stderr, err := mcClient.ExecInPod(context.Background(), testPodName, testPodNamespace, "test", cmd)
		if err != nil {
			return fmt.Errorf("can't exec command in pod %s: %s (stderr: %q)", testPodName, err, stderr)
		}

		// {"status":"success","data":{"resultType":"vector","result":[{"metric":{},"value":[1681718763.145,"1"]}]}}

		type result struct {
			Value []any
		}

		response := struct {
			Status string
			Data   struct {
				ResultType string
				Result     []result
			}
		}{}

		err = json.Unmarshal([]byte(stdout), &response)
		if err != nil {
			return fmt.Errorf("can't parse mimir query output: %s (output: %q)", err, stdout)
		}

		if response.Status != "success" {
			return fmt.Errorf("unexpected response status %q when running query %q (output: %q)", response.Status, query, stdout)
		}

		if response.Data.ResultType != "vector" {
			return fmt.Errorf("unexpected response type %q when running query %q (wanted vector,  output: %q)", response.Status, query, stdout)
		}

		if len(response.Data.Result) != 1 {
			return fmt.Errorf("unexpected count of results when running query %q (wanted 1, got %d; output: %q)", query, len(response.Data.Result), stdout)
		}

		// Second field of first result is the metric value. [1681718763.145,"1"] => "1"
		str, ok := (response.Data.Result[0].Value[1]).(string)
		if !ok {
			return fmt.Errorf("cannot cast result value to string when running query %q (output: %q)", query, stdout)
		}
		if str != "0" {
			return fmt.Errorf("unexpected value for query %q (wanted '0', got %q; output: %q)", query, str, stdout)
		}

		logger.Log("Metric %q was found", metric)
		return nil
	}
}

func Run(mcClient *client.Client, podName string, ns string) error {
	t := true
	f := false
	userAndGroup := int64(35)

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: ns,
		},
		Spec: corev1.PodSpec{
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    &userAndGroup,
				RunAsGroup:   &userAndGroup,
				RunAsNonRoot: &t,
				SeccompProfile: &corev1.SeccompProfile{
					Type: "RuntimeDefault",
				},
			},
			Containers: []corev1.Container{
				{
					Name:  "test",
					Image: "gsoci.azurecr.io/giantswarm/alpine:latest",
					Args:  []string{"sleep", "99999999"},
					SecurityContext: &corev1.SecurityContext{
						Capabilities: &corev1.Capabilities{
							Drop: []corev1.Capability{
								"ALL",
							},
						},
						AllowPrivilegeEscalation: &f,
					},
				},
			},
		},
	}
	// Check if pods exists already.
	create := false
	existing := corev1.Pod{}
	err := mcClient.Get(context.Background(), client2.ObjectKey{Namespace: ns, Name: podName}, &existing)
	if errors.IsNotFound(err) {
		create = true
	} else if err != nil {
		return fmt.Errorf("error ensuring test pod is deleted %s: %s", podName, err)
	}

	if !create {
		// Check if pod is running.
		if existing.Status.Phase != corev1.PodRunning {
			// Pod unhealthy, delete and recreate it.
			err := Cleanup(mcClient, podName, ns)
			if err != nil {
				return err
			}

			create = true
		}
	}

	if create {
		// Create the pod.
		err = mcClient.Create(context.Background(), &pod)
		if err != nil {
			return fmt.Errorf("can't create test pod %s: %s", podName, err)
		}
	}

	// Wait for pod to be running.
	Eventually(func() (bool, error) {
		err = mcClient.Get(context.Background(), client2.ObjectKey{Namespace: ns, Name: podName}, &existing)
		if err != nil {
			return false, fmt.Errorf("error ensuring test pod is running %s: %s", podName, err)
		}

		if existing.Status.Phase == corev1.PodRunning {
			return true, nil
		}

		return false, fmt.Errorf("waiting for pod %s to be running", podName)
	}).
		WithTimeout(5 * time.Minute).
		WithPolling(5 * time.Second).
		Should(BeTrue())

	return nil
}

func Cleanup(mcClient *client.Client, podName string, ns string) error {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: ns,
		},
	}

	// Pod unhealthy, delete and recreate it.
	err := mcClient.Delete(context.Background(), &pod)
	if errors.IsNotFound(err) {
		// Fallthrough (in case the pod was deleting already).
	} else if err != nil {
		return fmt.Errorf("error deleting test pod %s: %s", podName, err)
	}

	// Wait for pod to be deleted.
	Eventually(func() (bool, error) {
		err = mcClient.Get(context.Background(), client2.ObjectKey{Namespace: ns, Name: podName}, &pod)
		if errors.IsNotFound(err) {
			return true, nil
		} else if err != nil {
			return false, fmt.Errorf("error ensuring test pod %s is deleted: %s", podName, err)
		}

		return false, fmt.Errorf("waiting for pod %s to be deleted", podName)
	}).
		WithTimeout(5 * time.Minute).
		WithPolling(5 * time.Second).
		Should(BeTrue())

	return nil
}
