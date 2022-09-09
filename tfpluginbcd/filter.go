package tfpluginbcd

import (
	"context"
	"encoding/json"

	"github.com/open-policy-agent/opa/rego"
)

func Filter(ctx context.Context, changes []Change, regoModule string) ([]Change, error) {
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
		rego.Module("rules", regoModule))

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
