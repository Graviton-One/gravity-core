package models

type Actor struct {
	Name      string
	InitScore float32
}

type Vote struct {
	Target string
	Score  float32
}
