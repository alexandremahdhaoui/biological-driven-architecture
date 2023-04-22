package biological_driven_architecture

type Orchestrator struct {
	WorkerPools []WorkerPool
	Strategy    Strategy
}

type Strategy interface {
}
