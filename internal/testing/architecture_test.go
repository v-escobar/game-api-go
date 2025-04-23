package testing

import (
	archgo "github.com/arch-go/arch-go/api"
	config "github.com/arch-go/arch-go/api/configuration"

	"testing"
)

func TestArchitecture(t *testing.T) {
	configuration := config.Config{
		DependenciesRules: []*config.DependenciesRule{
			{
				Package: "**.interfaces.**",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"**.application.**",
						"**.interfaces.**",
					},
				},
			},
			{
				Package: "**.application.**",
				ShouldOnlyDependsOn: &config.Dependencies{
					Internal: []string{
						"**.domain.**",
						"**.application.errors",
						"**.interfaces.dto.**",
					},
				},
			},
		},
	}

	moduleInfo := config.Load("github.com/v-escobar/game-api-go")
	result := archgo.CheckArchitecture(moduleInfo, configuration)

	if !result.Pass {
		for _, rule := range result.DependenciesRuleResult.Results {
			if !rule.Passes {
				for _, verification := range rule.Verifications {
					if !verification.Passes {
						t.Errorf("Package: %s, Details: %v", verification.Package, verification.Details)
					}
				}
			}
		}
	}
}
