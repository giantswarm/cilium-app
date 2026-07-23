// Package readiness provides stricter cluster-readiness wait conditions than
// the stock clustertest helpers, tailored to guard the Cilium connectivity
// test against transient worker-node churn.
package readiness

import (
	"context"

	"github.com/giantswarm/clustertest/v5/pkg/client"
	"github.com/giantswarm/clustertest/v5/pkg/logger"
	"github.com/giantswarm/clustertest/v5/pkg/wait"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// AllNodesAndDaemonSetsReady returns a WaitCondition that is satisfied only
// when every Node reports Ready and every DaemonSet has all of its pods
// scheduled AND available.
//
// The stock wait.AreAllDaemonSetsReady only compares CurrentNumberScheduled
// against DesiredNumberScheduled. That is too weak to protect the Cilium
// connectivity test: a DaemonSet pod on a node that is draining or being
// replaced can still count as "scheduled" while its kubelet has already torn
// the container down. cilium-cli then snapshots the agent pod list and, seconds
// later, execs into a pod whose backing node is gone, producing:
//
//	Internal error occurred: error executing command in container:
//	pods "ip-10-0-x-y.<region>.compute.internal" not found
//
// Gating on node readiness catches the churn directly, and gating on
// NumberAvailable/NumberReady ensures the agent pods are actually usable before
// the connectivity test picks one to exec into.
func AllNodesAndDaemonSetsReady(ctx context.Context, kubeClient *client.Client) wait.WaitCondition {
	return func() (bool, error) {
		// Every node must be Ready. A node that is cordoned, draining or
		// NotReady (e.g. mid-replacement) is exactly the condition that makes a
		// Cilium agent pod vanish out from under the connectivity test.
		nodeList := &corev1.NodeList{}
		if err := kubeClient.List(ctx, nodeList); err != nil {
			return false, err
		}
		if len(nodeList.Items) == 0 {
			logger.Log("no nodes found yet")
			return false, nil
		}
		for _, node := range nodeList.Items {
			if node.DeletionTimestamp != nil {
				logger.Log("node %s is terminating", node.Name)
				return false, nil
			}
			ready := false
			for _, cond := range node.Status.Conditions {
				if cond.Type == corev1.NodeReady {
					ready = cond.Status == corev1.ConditionTrue
					break
				}
			}
			if !ready {
				logger.Log("node %s is not Ready", node.Name)
				return false, nil
			}
		}

		// Every DaemonSet must have all pods scheduled, up to date, ready and
		// available - not merely scheduled.
		daemonSetList := &appsv1.DaemonSetList{}
		if err := kubeClient.List(ctx, daemonSetList); err != nil {
			return false, err
		}
		for _, ds := range daemonSetList.Items {
			desired := ds.Status.DesiredNumberScheduled
			status := ds.Status
			if status.CurrentNumberScheduled != desired ||
				status.UpdatedNumberScheduled != desired ||
				status.NumberReady != desired ||
				status.NumberAvailable != desired ||
				status.NumberUnavailable != 0 {
				logger.Log(
					"daemonSet %s/%s not fully available: desired=%d scheduled=%d updated=%d ready=%d available=%d unavailable=%d",
					ds.Namespace, ds.Name, desired,
					status.CurrentNumberScheduled, status.UpdatedNumberScheduled,
					status.NumberReady, status.NumberAvailable, status.NumberUnavailable,
				)
				return false, nil
			}
		}

		logger.Log("All (%d) nodes Ready and all (%d) daemonSets fully available", len(nodeList.Items), len(daemonSetList.Items))
		return true, nil
	}
}
