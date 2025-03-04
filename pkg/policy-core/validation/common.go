package validation

import (
	"context"

	"github.com/weaveworks/policy-agent/pkg/policy-core/domain"
)

func matchEntity(entity domain.Entity, policy domain.Policy) bool {
	var matchKind bool
	var matchNamespace bool
	var matchLabel bool

	if len(policy.Targets.Kinds) == 0 {
		matchKind = true
	} else {
		resourceKind := entity.Kind
		for _, kind := range policy.Targets.Kinds {
			if resourceKind == kind {
				matchKind = true
				break
			}
		}
	}

	if len(policy.Targets.Namespaces) == 0 {
		matchNamespace = true
	} else {
		resourceNamespace := entity.Namespace
		for _, namespace := range policy.Targets.Namespaces {
			if resourceNamespace == namespace {
				matchNamespace = true
				break
			}
		}
	}

	if len(policy.Targets.Labels) == 0 {
		matchLabel = true
	} else {
	outer:
		for _, obj := range policy.Targets.Labels {
			for key, val := range obj {
				entityVal, ok := entity.Labels[key]
				if ok {
					if val != "*" && val != entityVal {
						continue
					}
					matchLabel = true
					break outer
				}
			}
		}
	}

	return matchKind && matchNamespace && matchLabel
}

func writeToSinks(
	ctx context.Context,
	resultsSinks []domain.PolicyValidationSink,
	PolicyValidationSummary domain.PolicyValidationSummary,
	writeCompliance bool) {
	for _, resutsSink := range resultsSinks {
		if len(PolicyValidationSummary.Violations) > 0 {
			resutsSink.Write(ctx, PolicyValidationSummary.Violations)
		}
		if writeCompliance && len(PolicyValidationSummary.Compliances) > 0 {
			resutsSink.Write(ctx, PolicyValidationSummary.Compliances)
		}
	}
}
