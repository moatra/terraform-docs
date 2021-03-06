package doc

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/moatra/terraform-docs/internal/pkg/settings"
	"sort"
	"strings"
)

type Doc struct {
	Variables []Variable `json:"variables"`
	Outputs []Output `json:"outputs"`
	Providers []Provider `json:"providers"`
}

// TODO: verify that the side effects to tracker stick
func discoverAliases(tracker map[string]Provider, versionLookup map[string][]string, resources map[string]*tfconfig.Resource) {
	for _, resource := range resources {
		key := fmt.Sprintf("%s.%s", resource.Provider.Name, resource.Provider.Alias)
		var version = ""
		if requiredVersion, ok := versionLookup[resource.Provider.Name]; ok {
			version = strings.Join(requiredVersion, " ")
		}
		tracker[key] = Provider{
			Name:    resource.Provider.Name,
			Alias:   resource.Provider.Alias,
			Version: version,
		}
	}
}

func Create(module *tfconfig.Module, printSettings settings.Settings) (*Doc, error) {
	var variables = make([]Variable, 0, len(module.Variables))
	for _, variable := range module.Variables {
		var defaultValue string
		if variable.Default != nil {
			marshaled, err := json.MarshalIndent(variable.Default, "", "  ")
			if err != nil {
				return nil, err
			}
			defaultValue = string(marshaled)
		}
		variables = append(variables, Variable{
			Name:        variable.Name,
			Type:        variable.Type,
			Description: variable.Description,
			Default:	 defaultValue,
		})
	}

	var outputs = make([]Output, 0, len(module.Outputs))
	for _, output := range module.Outputs {
		outputs = append(outputs, Output{
			Name:        output.Name,
			Description: output.Description,
		})
	}

	var providerSet = make(map[string]Provider)
	discoverAliases(providerSet, module.RequiredProviders, module.DataResources)
	discoverAliases(providerSet, module.RequiredProviders, module.ManagedResources)
	var providers = make([]Provider, 0, len(providerSet))
	for _, provider := range providerSet {
		providers = append(providers, provider)
	}

	if printSettings.Has(settings.WithSortVariablesByRequired) {
		sort.Sort(variablesSortedByRequired(variables))
	} else {
		sort.Sort(variablesSortedByName(variables))
	}
	sort.Sort(outputsSortedByName(outputs))
	sort.Sort(providersSortedByRequired(providers))

	doc := &Doc{
		Variables: variables,
		Outputs:   outputs,
		Providers: providers,
	}
	return doc, nil

}