package biological_driven_architecture

type Error struct {
	Type    string
	Message string
}

type Worker interface {
	Run() Error
	HandleError(Error)
	Stop()
}
