package tfpluginbcd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/magodo/tfpluginschema/schema"
)

type Opt struct {
	Rules           []string
	CustomRuleExprs []string
}

func Run(ctx context.Context, opath, npath string, opt Opt) (string, error) {
	// Reading schemas
	ob, err := os.ReadFile(opath)
	if err != nil {
		return "", fmt.Errorf("reading the old schema file %s: %v", opath, err)
	}
	var osch schema.ProviderSchema
	if err := json.Unmarshal(ob, &osch); err != nil {
		return "", fmt.Errorf("unmarshalling the old schema: %v", err)
	}
	nb, err := os.ReadFile(npath)
	if err != nil {
		return "", fmt.Errorf("reading the new schema file %s: %v", npath, err)
	}
	var nsch schema.ProviderSchema
	if err := json.Unmarshal(nb, &nsch); err != nil {
		return "", fmt.Errorf("unmarshalling the new schema: %v", err)
	}

	changes, err := run(ctx, osch, nsch, opt)
	if err != nil {
		return "", err
	}
	return strings.Join(changes, "\n"), nil
}

func run(ctx context.Context, osch, nsch schema.ProviderSchema, opt Opt) ([]string, error) {
	var rules []Rule
	for _, name := range opt.Rules {
		rule, ok := Rules[name]
		if !ok {
			return nil, fmt.Errorf("undefined rule: %s", name)
		}
		rules = append(rules, rule)
	}
	for idx, expr := range opt.CustomRuleExprs {
		rules = append(rules, Rule{
			ID:   fmt.Sprintf("CUSTOM-%d", idx),
			Expr: expr,
		})
	}
	results, err := Filter(ctx, Compare(&osch, &nsch), rules)
	if err != nil {
		return nil, fmt.Errorf("filtering: %v", err)
	}

	var output []string
	for _, res := range results {
		if res.Rule == "" {
			output = append(output, res.Change.String())
		} else {
			output = append(output, fmt.Sprintf("[%s] %s", res.Rule, res.Change.String()))
		}
	}
	return output, nil
}
