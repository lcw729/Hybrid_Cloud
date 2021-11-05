package util

type WatchingLevel struct {
	Levels []Level
}

type Level struct {
	Type  string
	Value []string
}
