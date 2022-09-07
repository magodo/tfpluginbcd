package tfpluginbcd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/magodo/tfpluginschema/schema"
)

func Run(opath, npath string) (string, error) {
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
	changes := Compare(&osch, &nsch)
	var cl []string
	for _, change := range changes {
		cl = append(cl, change.String())
	}
	return strings.Join(cl, "\n"), nil
}
