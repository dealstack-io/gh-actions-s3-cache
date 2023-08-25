package main

const (
	SaveAction    = "save"
	RestoreAction = "restore"
)

type Action struct {
	Action          string
	Bucket          string
	Key             string
	Paths           []string
	FailOnCacheMiss bool
}
