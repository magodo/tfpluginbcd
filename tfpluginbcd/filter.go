package tfpluginbcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/open-policy-agent/opa/rego"
)

func buildRegoModule(content string) string {
	return fmt.Sprintf(`package provider

import future.keywords.in

%s
`, content)
}

func buildRule(content string) string {
	return fmt.Sprintf(`
breaking_change[i] {
    some i, c in input.changes
	%s
}
`, content)
}

func Filter(ctx context.Context, changes []Change, regoRules []string) ([]Change, error) {
	if len(regoRules) == 0 {
		return changes, nil
	}

	// Turn the changes from an array to an object, as rego only process on json object as input.
	type ChangeSet struct {
		Changes []Change `json:"changes"`
	}
	cs := ChangeSet{
		Changes: changes,
	}

	// Marshal and unmarshal back the change set to a Go map (default), which will then be able to be processd by rego.
	b, err := json.Marshal(cs)
	if err != nil {
		return nil, err
	}
	var input interface{}
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, err
	}

	r := rego.New(
		rego.Query("data.provider.breaking_change"),
		rego.Module("rules", buildRegoModule(strings.Join(regoRules, "\n"))))

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}
	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}

	var output []Change
	for _, idx := range rs[0].Expressions[0].Value.([]interface{}) {
		idx, _ := idx.(json.Number).Int64()
		output = append(output, changes[int(idx)])
	}

	return output, nil
}
