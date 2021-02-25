package model

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

var (
	pathToModel = "../../data/examplemodel.txt"

	model *Model

	// keys: queries & vals: expected query results.
	lookupCases = map[string][]string{
		"dog":  {"cat", "dogs", "horse", "puppy"},
		"lamp": {"lamps", "bulb", "incandescent", "halogen"},
		"sofa": {"couches", "comfy", "cushions", "sofas"},
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
	m, err := Load(pathToModel, true, true)
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
