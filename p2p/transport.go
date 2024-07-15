package p2p

// Peer is interface that represents the "remote node"
type Peer interface {
	Close() error
}

// Transport handles the inter-communication between nodes in network
// This can be of the form - TCP, UDP, WebSockets
type Transport interface {
	ListenAndAccept() error
}
