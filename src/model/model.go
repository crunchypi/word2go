package model

type Model struct {
	// Data (word2vec model) is a map where keys
	// are model vocabulary and vals are vecs.
	// Chose a map because it greatly improves
	// queries, which need to check if the query
	// word exists even (alternative is lookups
	// done with linear/binar/etc searches).
	// Cost seems to be aprox 5% slower model
	// import time.
	Data map[string][]float64
	// Dim is how large vectors are in the model
	// (keys of Data). Consistency is not guaranteed
	// unless Model.ValidateDimensions is used. This
	// is not checked implicitly for performance
	// reasons.
	Dim int
}

// QueryResult is the result item of a word2vec lookup.
type QueryResult struct {
	Word      string
	SimiScore float64
}

// Load tries to load a word2vec model from a given path.
// Arguments in order: Path to model; progress printout;
// whether to fail on any issue with the model.
func Load(path string, verbose, strict bool) (*Model, error) {
	return fileToModel(path, verbose, strict)
}

// Save tries to save the model as a text file at the given path.
func (m *Model) Save(path string) error {
	return modelToFile(m, path)
}

// Checks for vector length consistency in self-contained model.
func (m *Model) ValidateDimensions() bool {
	for _, v := range m.Data {
		if len(v) != m.Dim {
			return false
		}
	}
	return true
}

// Lookup will do a word2vec lookup of 'n' neighs for the specified word.
func (m *Model) Lookup(word string, n int) (*[]QueryResult, bool) {
	// # Check if word even exists in the model.
	wordV, ok := m.Data[word]
	if !ok {
		return nil, false
	}
	res := make([]QueryResult, n)
	// # Keep track of lowest score. If a result member candidate
	// # does not have a score that is higher than the worst/lowest
	// # score, then there's no point in doing anything with it.
	low := 0.0

	// # All members have to be checked for score.
	for k, v := range m.Data {
		// # No point in getting the queried word.
		if k == word {
			continue
		}

		score := score(wordV, v)
		if score > low {

			// # This is sort-of a partial merge algorithm -
			// # length of result will always be the same and
			// # the members are in sorted order (by score).
			// # The 'insertee' variable is a carry/temp.
			insertee := QueryResult{Word: k, SimiScore: score}
			for i := 0; i < len(res); i++ {
				if insertee.SimiScore > res[i].SimiScore {
					insertee, res[i] = res[i], insertee
				}
			}
			// # Res is aways sorted.
			low = res[len(res)-1].SimiScore
			//low = insertee.SimiScore
		}
	}
	return &res, true
}

// Prune removes all words in the model that are not members
// in the specified slice.
func (m *Model) Prune(include *[]string) bool {
	if include == nil {
		return false
	}

	// # Set included vocab into a map for performance.
	// # Keys are insignificant.
	includeMap := make(map[string]bool, len(*include))
	for i := 0; i < len(*include); i++ {
		includeMap[(*include)[i]] = false
	}

	for k := range m.Data {
		// # Drop entire entry if it is not to be included.
		_, exists := includeMap[k]
		if !exists {
			delete(m.Data, k)
		}
	}
	return true
}
