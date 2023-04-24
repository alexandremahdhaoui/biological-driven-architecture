package biological_driven_architecture

type Strategy interface {
	Init(Runtime) Error
	Run(Runtime) Error
	HandleError(Runtime, Error) Error
	Stop(Runtime) Error
}

type defaultStrategy struct{}

func (s *defaultStrategy) Init(runtime Runtime) Error {
	return runtime.Init()
}

func (s *defaultStrategy) Run(runtime Runtime) Error {
	return runtime.Run()
}

func (s *defaultStrategy) HandleError(runtime Runtime, err Error) Error {
	return runtime.HandleError(err)
}

func (s *defaultStrategy) Stop(runtime Runtime) Error {
	return runtime.Stop()
}

func DefaultStrategy() Strategy {
	return &defaultStrategy{}
}
