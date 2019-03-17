package pdptw

type objective interface {
	get(*Solution) int
	isProfitable(s *Solution, i, j int, spans ...int) bool
}

type spanTime struct{}
type totalTime struct{}
type totalTimeA struct{}
