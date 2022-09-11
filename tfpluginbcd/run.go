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
	Rules             []string
	CustomRuleContent string
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

	bcs, err := run(ctx, osch, nsch, opt)
	if err != nil {
		return "", err
	}

	var output []string
	for _, c := range bcs {
		output = append(output, c.String())
	}
	return strings.Join(output, "\n"), nil
}

func run(ctx context.Context, osch, nsch schema.ProviderSchema, opt Opt) ([]Change, error) {
	var regoRules []string
	for _, name := range opt.Rules {
		rule, ok := Rules[name]
		if !ok {
			return nil, fmt.Errorf("undefined rule: %s", name)
		}
		regoRules = append(regoRules, buildRule(rule.Expr))
	}
	if opt.CustomRuleContent != "" {
		regoRules = append(regoRules, buildRule(opt.CustomRuleContent))
	}
	return Filter(ctx, Compare(&osch, &nsch), regoRules)
}
