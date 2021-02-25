# word2go


A very simple and performant word2vec _reader_ -- created because the readers I found on this site (Go) were too slow for my use-case.


Features are (examples are listed further down).
- import a word2vec model (example is found in ./data/examplemodel.txt)
- lookup a word2vec relationship.
- pruning (reducing the vocabulary in the model type)
- saving the model type


(note all these examples can be found in ./example.go)


##### Importing a model:
```
import "github.com/crunchypi/word2go/src/model"

func myFunc() {
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
  // # necessary
  if ok := m.ValidateDimensions(); !ok {
      panic("model dimensions are off :< ")
  }
  
  ...

}
```



##### Lookup with imported model:
```
import (
  "fmt"
  "github.com/crunchypi/word2go/src/model"
)

func myFunc() {
  m, _ := model.Load("./data/examplemodel.txt", true, true)
  
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
  
  ...
}
```



##### Lookup with imported model:
```
import "github.com/crunchypi/word2go/src/model"

func myFunc() {
  m, _ := model.Load("./data/examplemodel.txt", true, true)
  
  // # Stuff to ignore in the pruning (these words will not
  // # be deleted).
  preserve := []string{"dog"}
  
  if ok := m.Prune(&preserve); !ok {
    panic("failed while pruning :<")
  }
}
```


##### Saving a model:
```
import "github.com/crunchypi/word2go/src/model"

func myFunc() {
  oldPath, newPath := "./data/examplemodel.txt", "tmp.txt"
  m, _ := model.Load(oldPath, true, true)
  
  // # This will just copy the word2vec file, since nothing
  // # is pruned.
  if err := m.Save(newPath); err != nil {
    panic("err while saving :<")
  }

}

```
