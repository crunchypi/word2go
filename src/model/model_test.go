package model

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

var (
	pathToModelSmall = "../../data/examplemodel.txt"
	pathToModelLarge = "../../data/glove/glove.6B.100d.txt"

	model *Model

	// keys: queries & vals: expected query results.
	lookupCases = map[string][]string{
		"dog":  {"cat", "dogs", "horse", "puppy"},
		"lamp": {"lamps", "bulb", "incandescent", "halogen"},
		"sofa": {"couches", "comfy", "cushions", "sofas"},
	}
	knownScores = map[string]QueryResult{
		"dog": {Word: "cat", SimiScore: 0.921801},
	}
	// What not to delete when doing pruning. This is all the
	// words in the test cases above.
	pruneInclude = []string{
		"dog", "cat", "dogs", "horse", "puppy",
		"lamp", "lamps", "bulb", "incandescent", "halogen",
		"sofa", "couches", "comfy", "cushions", "sofas",
	}
)

func reset() {
	m, err := Load(pathToModelSmall, true, true)
	if err != nil {
		panic(err.Error)
	}
	model = m
}

func init() {
	reset()
}

func checkCorrectness(m *Model) error {
	// # Compare test cases and lookup results.
	for k, v := range lookupCases {
		res, ok := m.Lookup(k, len(v))

		if !ok {
			return errors.New("failed lookup on key: " + k)
		}
		// # Check order consistency.
		for i := 0; i < len(v); i++ {
			if v[i] != (*res)[i].Word {
				msg := fmt.Sprintf("lookup err on key '%s'. Wanted '%s', got '%s'",
					k, v[i], (*res)[i].Word)
				return errors.New(msg)

			}
		}
	}
	return nil
}

func TestLookup(t *testing.T) {
	if err := checkCorrectness(model); err != nil {
		t.Error(err)
	}
}

func TestCompare(t *testing.T) {
	for k, v := range knownScores {
		r, ok := model.Compare(k, v.Word)
		if !ok {
			t.Errorf("lookup failed for %s \n", k)
		}
		// # This is a hack - float64 comparison can be tricky..
		scoreA := fmt.Sprintf("%f", r)
		scoreB := fmt.Sprintf("%f", v.SimiScore)

		if scoreA != scoreB {
			t.Errorf("wrong res for words %s and %s: %s (want %s)\n",
				k, v.Word, scoreA, scoreB)
		}
	}
}

func TestPrune(t *testing.T) {

	if ok := model.Prune(&pruneInclude); !ok {
		t.Error("pruning task failed")
	}

	l := 0
	for k := range model.Data {
		// # Just to use the k var...
		if k == k {
			l++
		}
	}
	if l != len(pruneInclude) {
		t.Error("data size incorrect after pruning")
	}

	if err := checkCorrectness(model); err != nil {
		t.Error(err)
	}
}

func TestSave(t *testing.T) {
	tempPath := "tmp.txt"
	if err := model.Save(tempPath); err != nil {
		t.Error(err)
	}

	m, err := Load(tempPath, false, true)
	if err != nil {
		t.Error(err)
	}
	if err := checkCorrectness(m); err != nil {
		t.Error(err)
	}
	os.Remove(tempPath)
}
