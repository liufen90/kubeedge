package admissioncontroller

import (
	"fmt"
	"net/http"
	"strings"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	rulesv1 "github.com/kubeedge/kubeedge/cloud/pkg/apis/rules/v1"
)

func admitRule(review admissionv1beta1.AdmissionReview) *admissionv1beta1.AdmissionResponse {
	reviewResponse := admissionv1beta1.AdmissionResponse{}
	var msg string
	switch review.Request.Operation {
	case admissionv1beta1.Create:
		raw := review.Request.Object.Raw
		rule := rulesv1.Rule{}
		deserializer := codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(raw, nil, &rule); err != nil {
			klog.Errorf("validation failed with error: %v", err)
			msg = err.Error()
			break
		}
		err := validateTargetRuleEndpoint(&rule)
		if err != nil {
			msg = err.Error()
			break
		}
		err = validateSourceRuleEndpoint(&rule)
		if err != nil {
			msg = err.Error()
			break
		}
		reviewResponse.Allowed = true
	case admissionv1beta1.Delete, admissionv1beta1.Connect:
		//no rule defined for above operations, greenlight for all of above.
		reviewResponse.Allowed = true
		klog.Info("admission validation passed!")
	default:
		klog.Infof("Unsupported webhook operation %v", review.Request.Operation)
		msg = msg + "Unsupported webhook operation!"
	}
	if !reviewResponse.Allowed {
		reviewResponse.Result = &metav1.Status{Message: strings.TrimSpace(msg)}
	}
	return &reviewResponse
}

func validateSourceRuleEndpoint(rule *rulesv1.Rule) error {
	sourceKey := fmt.Sprintf("%s/%s", rule.Namespace, rule.Spec.Source)
	endpoint, err := controller.getRuleEndpoint(rule.Namespace, rule.Spec.Source)
	if err != nil {
		return fmt.Errorf("cant get source ruleEndpoint %s. reason: %s", sourceKey, err.Error())
	} else if endpoint == nil {
		return fmt.Errorf("source ruleEndpoint %s has not been created. ", sourceKey)
	}
	res := rule.Spec.SourceResource
	switch endpoint.Spec.RuleEndpointType {
	case "rest":
		_, exist := res["path"]
		if !exist {
			return fmt.Errorf("source properties do not find \"path\". ")
		}
		rules, err := controller.listRule(rule.Namespace)
		if err != nil {
			return err
		}
		for _, r := range rules {
			if res["path"] == r.Spec.SourceResource["path"] {
				return fmt.Errorf("source properties exist. path: %s", res["path"])
			}
		}
	case "eventbus":
		_, exist := res["topic"]
		if !exist {
			return fmt.Errorf("source properties do not find \"topic\". ")
		}
		_, exist = res["node_name"]
		if !exist {
			return fmt.Errorf("source properties do not find \"node_name\". ")
		}
		rules, err := controller.listRule(rule.Namespace)
		if err != nil {
			return err
		}
		for _, r := range rules {
			if res["topic"] == r.Spec.SourceResource["topic"] && res["node_name"] == r.Spec.SourceResource["node_name"] {
				return fmt.Errorf("source properties exist. node_name: %s, topic: %s", res["node_name"], res["topic"])
			}
		}
	}
	return nil
}

func validateTargetRuleEndpoint(rule *rulesv1.Rule) error {
	targetKey := fmt.Sprintf("%s/%s", rule.Namespace, rule.Spec.Target)
	endpoint, err := controller.getRuleEndpoint(rule.Namespace, rule.Spec.Target)
	if err != nil {
		return fmt.Errorf("cant get target ruleEndpoint %s. reason: %s", targetKey, err.Error())
	} else if endpoint == nil {
		return fmt.Errorf("target ruleEndpoint %s has not been created. ", targetKey)
	}
	res := rule.Spec.TargetResource
	switch endpoint.Spec.RuleEndpointType {
	case "rest":
		_, exist := res["resource"]
		if !exist {
			return fmt.Errorf("target properties do not find \"resource\". ")
		}
	case "eventbus":
		_, exist := res["topic"]
		if !exist {
			return fmt.Errorf("target properties do not find \"topic\". ")
		}
	}
	return nil
}

func serveRule(w http.ResponseWriter, r *http.Request) {
	serve(w, r, admitRule)
}
