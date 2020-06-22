package models

type Request struct {
	Actors []Actor
	Votes  map[string][]Vote
}

type Actor struct {
	Name      string
	InitScore float32
}

type Vote struct {
	Target string
	Score  float32
}
