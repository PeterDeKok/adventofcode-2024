package info

type PartInfo struct {
	Summary      *PartInfoSummary
	LastSolution *PartInfoLastSolution
	// Solutions
}

type PartInfoSummary struct {
	BestStatus     string
	SolutionCount  string // Nr of solutions
	FastestRuntime string // Fastest solution timing
	CorrectAnswer  string
	FinishedAt     string
}

type PartInfoLastSolution struct {
	Status     string // OK, FAILED, ERROR, REVIEW/DEFERRED, BUILD
	Error      string
	Runtime    string
	FinishedAt string
	Answer     string
	Expected   string
}

func NewPartInfo() *PartInfo {
	return &PartInfo{
		Summary:      &PartInfoSummary{},
		LastSolution: &PartInfoLastSolution{},
	}
}

func (pis *PartInfo) Done() bool {
	return true
}
