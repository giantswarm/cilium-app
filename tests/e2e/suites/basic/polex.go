package basic

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func buildPolicyException() *unstructured.Unstructured {

	u := unstructured.Unstructured{}

	u.Object = map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      "cilium-test",
			"namespace": "kube-system",
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
		Version: "v2beta1",
	})

	return &u
}
