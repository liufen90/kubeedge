package utils

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	rulesv1 "github.com/kubeedge/kubeedge/cloud/pkg/apis/rules/v1"
)

const (
	RestType     = "rest"
	EventbusType = "eventbus"
)

func NewRule() *rulesv1.Rule {
	rule := rulesv1.Rule{
		TypeMeta: v1.TypeMeta{
			Kind:       "Rule",
			APIVersion: "rules.kubeedge.io/v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "rule-rest-eventbus-test",
			Namespace: Namespace,
		},
		Spec: rulesv1.RuleSpec{
			Source: "rest-test",
			SourceResource: map[string]string{
				"path": "/ccc",
			},
			Target: "eventbus-test",
			TargetResource: map[string]string{
				"topic": "topic-test",
			},
		},
	}
	return &rule
}

func NewRuleEndpoint(endpointType string) *rulesv1.RuleEndpoint {
	switch endpointType {
	case RestType:
		return newRestRuleEndpoint()
	case EventbusType:
		return newEventBusRuleEndpoint()
	}
	return newRestRuleEndpoint()
}

func newRestRuleEndpoint() *rulesv1.RuleEndpoint {
	restRuleEndpoint := rulesv1.RuleEndpoint{
		TypeMeta: v1.TypeMeta{
			Kind:       "RuleEndpoint",
			APIVersion: "rules.kubeedge.io/v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "rest-test",
			Namespace: Namespace,
		},
		Spec: rulesv1.RuleEndpointSpec{
			RuleEndpointType: RestType,
		},
	}
	return &restRuleEndpoint
}

func newEventBusRuleEndpoint() *rulesv1.RuleEndpoint {
	eventbusRuleEndpoint := rulesv1.RuleEndpoint{
		TypeMeta: v1.TypeMeta{
			Kind:       "RuleEndpoint",
			APIVersion: "rules.kubeedge.io/v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "eventbus-test",
			Namespace: Namespace,
		},
		Spec: rulesv1.RuleEndpointSpec{
			RuleEndpointType: EventbusType,
		},
	}
	return &eventbusRuleEndpoint
}
