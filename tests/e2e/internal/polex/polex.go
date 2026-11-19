package polex

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func New() *unstructured.Unstructured {

	u := unstructured.Unstructured{}

	u.Object = map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      "cilium-test",
			"namespace": "kube-system",
			"labels": map[string]interface{}{
				"policy.giantswarm.io/migrated-to-cel": "true",
			},
		},
		"spec": map[string]interface{}{
			"match": map[string]interface{}{
				"any": []map[string]interface{}{
					{
						"resources": map[string]interface{}{
							"namespaces": []string{
								"cilium-test*",
							},
						},
					},
				},
			},
			"exceptions": []map[string]interface{}{
				{
					"policyName": "disallow-capabilities",
					"ruleNames": []string{
						"adding-capabilities",
						"autogen-adding-capabilities",
					},
				},

				{
					"policyName": "disallow-capabilities-strict",
					"ruleNames": []string{
						"adding-capabilities-strict",
						"autogen-adding-capabilities-strict",
						"require-drop-all",
						"autogen-require-drop-all",
					},
				},
				{
					"policyName": "disallow-host-ports",
					"ruleNames": []string{
						"host-ports-none",
						"autogen-host-ports-none",
					},
				},
				{
					"policyName": "disallow-host-path",
					"ruleNames": []string{
						"host-path",
						"autogen-host-path",
					},
				},
				{
					"policyName": "disallow-privilege-escalation",
					"ruleNames": []string{
						"privilege-escalation",
						"autogen-privilege-escalation",
					},
				},
				{
					"policyName": "disallow-run-as-nonroot",
					"ruleNames": []string{
						"run-as-non-root",
						"autogen-run-as-non-root",
					},
				},
				{
					"policyName": "require-run-as-nonroot",
					"ruleNames": []string{
						"run-as-non-root",
						"autogen-run-as-non-root",
					},
				},
				{
					"policyName": "restrict-seccomp",
					"ruleNames": []string{
						"check-seccomp",
						"autogen-check-seccomp",
					},
				},
				{
					"policyName": "restrict-seccomp-strict",
					"ruleNames": []string{
						"check-seccomp-strict",
						"autogen-check-seccomp-strict",
					},
				},
				{
					"policyName": "restrict-volume-types",
					"ruleNames": []string{
						"restricted-volumes",
						"autogen-restricted-volumes",
					},
				},
				{
					"policyName": "disallow-host-namespaces",
					"ruleNames": []string{
						"host-namespaces",
						"autogen-host-namespaces",
					},
				},
			},
		},
	}

	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "kyverno.io",
		Kind:    "PolicyException",
		Version: "v2",
	})

	return &u
}

// NewCEL returns the policies.kyverno.io/v1 counterpart of New. A kyverno.io/v2
// PolicyException only applies to ClusterPolicies, so clusters that run the
// CEL (ValidatingPolicy) versions of the policies need this one.
func NewCEL() *unstructured.Unstructured {

	u := unstructured.Unstructured{}

	u.Object = map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      "cilium-test",
			"namespace": "kube-system",
		},
		"spec": map[string]interface{}{
			"policyRefs": []map[string]interface{}{
				{"name": "disallow-capabilities", "kind": "ValidatingPolicy"},
				{"name": "disallow-capabilities-strict", "kind": "ValidatingPolicy"},
				{"name": "disallow-host-ports", "kind": "ValidatingPolicy"},
				{"name": "disallow-host-path", "kind": "ValidatingPolicy"},
				{"name": "disallow-privilege-escalation", "kind": "ValidatingPolicy"},
				{"name": "disallow-run-as-nonroot", "kind": "ValidatingPolicy"},
				{"name": "require-run-as-nonroot", "kind": "ValidatingPolicy"},
				{"name": "restrict-seccomp", "kind": "ValidatingPolicy"},
				{"name": "restrict-seccomp-strict", "kind": "ValidatingPolicy"},
				{"name": "restrict-volume-types", "kind": "ValidatingPolicy"},
				{"name": "disallow-host-namespaces", "kind": "ValidatingPolicy"},
				// The following are supplemental policies and usually not enabled.
				// They are excepted because the connectivity test workloads fail
				// them when the policies are enforced.
				{"name": "require-container-requests-and-limits", "kind": "ValidatingPolicy"},
				{"name": "require-pod-probes", "kind": "ValidatingPolicy"},
				{"name": "require-ro-rootfs", "kind": "ValidatingPolicy"},
				{"name": "restrict-sa-automount-sa-token", "kind": "ValidatingPolicy"},
			},
			"matchConditions": []map[string]interface{}{
				{
					"name":       "cilium-test-namespaces",
					"expression": "has(request.namespace) && (request.namespace == 'cilium-test' || request.namespace.startsWith('cilium-test-'))",
				},
			},
		},
	}

	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "policies.kyverno.io",
		Kind:    "PolicyException",
		Version: "v1",
	})

	return &u
}
