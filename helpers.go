package main

import (
	"time"

	"github.com/google/uuid"
)

func getObj(i string, batchCorrelationId string) PubsubMsg {
	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	return PubsubMsg{
		Key:                   i,
		ID:                    now,
		Details:               map[string]interface{}{"missing_labels": []interface{}{"test"}},
		EventType:             "violation_audited",
		Group:                 "constraints.gatekeeper.sh",
		Version:               "v1beta1",
		Kind:                  "K8sRequiredLabels",
		Name:                  "pod-must-have-test",
		Namespace:             "",
		Message:               "you must provide labels: {\"test\"}",
		EnforcementAction:     "deny",
		ConstraintAnnotations: map[string]string(nil),
		ResourceGroup:         "",
		ResourceAPIVersion:    "v1",
		ResourceKind:          "Pod",
		ResourceNamespace:     "nginx",
		ResourceName:          "dywuperf-deployment-10kpods-69bd64c867-h2wdx",
		ResourceLabels:        map[string]string{"app": "dywuperf-app-100kpods", "pod-template-hash": "69bd64c867"},
		Timestamp:             now,
		BrokerName:            BROKERNAME,
		UserAgent:             "Publisher",
		CorrelationId:         uuid.New().String(),
		BatchCorrelationId:    batchCorrelationId}
}