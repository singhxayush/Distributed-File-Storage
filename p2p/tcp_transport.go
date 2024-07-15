package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCP-Peer represents the remode node over a TCP established connection
type TCPPeer struct {
	conn net.Conn // underlying connection of the peer
	// dial and retrieve a connection -> outbound == true
	// Accept and retrieve a connection -> outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

// TCPTransport represents a transport layer implementation using TCP.
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex      // Mutex for protecting concurrent access to peers map.
	peers map[net.Addr]Peer // A map storing connected peers, keyed by their network address.
}

// NewTCPTransport creates a new TCPTransport instance.
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

// ListenAndAccept starts listening for incoming connections and accepts them.
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	// Listen on the specified TCP address.
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	// Start a goroutine to handle incoming connections.
	go t.StartAcceptLoop()

	return nil
}

// startAcceptLoop continuously accepts incoming connections.
func (t *TCPTransport) StartAcceptLoop() {
	for {
		conn, err := t.listener.Accept() // Accept a new connection.
		if err != nil {
			fmt.Printf("TCP accept ERROR: %s\n", err)
		}

		fmt.Printf("new incoming connection: %+v\n", conn)

		// Handle the connection in a separate goroutine.
		go t.handleConn(conn)
	}
}

// handleConn han`dles a new incoming connection.
func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil { // Checking if there occurs any error in handshake to drop the connection accordingly
		conn.Close()
		fmt.Printf("TCP handshake Error: %s\n", err)
		return
	}

	// Read loop
	rpc := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			fmt.Printf("TCP  Error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()

		fmt.Printf("message: %+v\n", rpc)
	}
}
