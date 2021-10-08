package core

type Plugin interface {
	Description() string
	Init()
	Execute(options map[string]interface{}, taskId uint64, progress chan<- float32)
	// PkgPath information lost when compiled as plugin(.so)
	RootPkgPath() string
}
