package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cwbudde/gll-tools/pkg/gll"
)

func TestAcousticJSONResponses(t *testing.T) {
	resetFlags(rootCmd)
	t.Cleanup(func() { resetFlags(rootCmd); rootCmd.SetOut(nil); rootCmd.SetErr(nil) })
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"acoustic", "../../../testdata/gll/example-vis.gll", "--json", "--responses", "--source", "0"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var sources []gll.SourceDefinitionItem
	if err := json.Unmarshal(out.Bytes(), &sources); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(sources) != 1 || sources[0].Definition == nil {
		t.Fatalf("expected one source: %+v", sources)
	}
	b := sources[0].Definition.BalloonData
	if b == nil || len(b.Responses) == 0 || len(b.Responses) != int(b.ResponseCount) {
		t.Fatal("JSON omitted loaded balloon responses")
	}
	if len(b.Responses[0].Phase) != len(b.Responses[0].Level) {
		t.Fatal("incomplete complex response")
	}
}
