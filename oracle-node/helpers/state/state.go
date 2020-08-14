package state

type States []State
type State struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

func (states States) Map() map[string]State {
	stateMap := make(map[string]State)
	for _, v := range states {
		stateMap[v.Key] = v
	}
	return stateMap
}
