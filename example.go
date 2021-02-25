package main

import (
	"fmt"
	"github.com/crunchypi/word2go/src/model"
	"os"
)

// Example for all features:
func main() {
	// # Path to where a word2vec model is located.
	path := "./data/examplemodel.txt"
	// # True will abort loading a model if there
	// # is any issue with the formatting of a file.
	strict := true
	// # True will print progress while loading a model
	// # (can take a while, depending on model, so it's
	// # nice to have).
	verbose := true

	m, err := model.Load(path, verbose, strict)
	if err != nil {
		panic(err)
	}

	// # The length of each vector per word in the model
	// # isn't implicitly checked since that isn't always
	// # necessary.
	if ok := m.ValidateDimensions(); !ok {
		panic("model dimensions are off :< ")
	}

	// # Normal lookup of 2 neighbours for the word 'dog'
	lookupWord := "dog"
	desiredResultLength := 2
	res, ok := m.Lookup(lookupWord, desiredResultLength)
	if !ok {
		panic("failed lookup :<")
	}
	for i, v := range *res {
		fmt.Printf("no. %d : %s %f\n", i, v.Word, v.SimiScore)
	}

	// # Pruning reduces the model, it removes any word that
	// # is not inside the specified slice.
	if ok := m.Prune(&[]string{}); !ok {
		panic("failed while pruning :<")
	}

	// # This should now yield an empty result, since
	// # everything was removed from the model (code
	// # block above).
	_, ok = m.Lookup(lookupWord, 2)
	if ok {
		panic("unexpected pruning result :<")
	}

	// # A model can also be saved (can be useful after pruning).
	savePath := "tmp.txt"
	if err := m.Save(savePath); err != nil {
		panic("err while saving: " + err.Error())
	}

	// # I'll just delete that file, it doesn't contain anything
	// # (removed everything in the model with the pruning step).
	os.Remove(savePath)

}
