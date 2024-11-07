package protocols

// DatabaseIsConnected is a protocol that checks if the database is connected
type DatabaseIsConnected interface {
	IsConnected() (bool, error)
}
