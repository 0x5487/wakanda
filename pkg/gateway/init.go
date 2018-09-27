package gateway

var (
	_manager *Manager
)

func Initialize() {
	_manager = NewManager()
}
