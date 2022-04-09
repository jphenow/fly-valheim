package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type connection struct {
	clientAddr *net.UDPAddr // Address of the client
	serverConn *net.UDPConn // UDP connection to server
}

// Generate a new connection by opening a UDP connection to the server
func newConnection(srvAddr, cliAddr *net.UDPAddr) *connection {
	conn := new(connection)
	conn.clientAddr = cliAddr
	srvudp, err := net.DialUDP("udp", nil, srvAddr)
	if checkreport(1, err) {
		return nil
	}
	conn.serverConn = srvudp
	return conn
}

// Global state
// Connection used by clients as the proxy server
var proxyConn *net.UDPConn

// Address of server
var serverAddr *net.UDPAddr

// Mapping from client addresses (as host:port) to connection
var clientDict map[string]*connection = make(map[string]*connection)

// Mutex used to serialize access to the dictionary
var dmutex *sync.Mutex = new(sync.Mutex)

func setup(hostport string, port int) bool {
	// Set up Proxy
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("fly-global-services:%d", port))
	if checkreport(1, err) {
		return false
	}
	pudp, err := net.ListenUDP("udp", saddr)
	if checkreport(1, err) {
		return false
	}
	proxyConn = pudp
	verboseLogf(2, "Proxy serving on port %d\n", port)

	// Get server address
	srvaddr, err := net.ResolveUDPAddr("udp", hostport)
	if checkreport(1, err) {
		return false
	}
	serverAddr = srvaddr
	verboseLogf(2, "Connected to server at %s\n", hostport)
	return true
}

func dlock() {
	dmutex.Lock()
}

func dunlock() {
	dmutex.Unlock()
}

// Go routine which manages connection from server to single client
func runConnection(conn *connection) {
	var buffer [1500]byte
	for {
		// Read from server
		n, err := conn.serverConn.Read(buffer[0:])
		if checkreport(1, err) {
			continue
		}
		// Relay it to client
		_, err = proxyConn.WriteToUDP(buffer[0:n], conn.clientAddr)
		if checkreport(1, err) {
			continue
		}
		verboseLogf(3, "Relayed '%s' from server to %s.\n",
			string(buffer[0:n]), conn.clientAddr.String())
	}
}

// Routine to handle inputs to Proxy port
func runProxy() {
	var buffer [1500]byte
	for {
		n, cliaddr, err := proxyConn.ReadFromUDP(buffer[0:])
		if checkreport(1, err) {
			continue
		}
		verboseLogf(3, "Read '%s' from client %s\n",
			string(buffer[0:n]), cliaddr.String())
		saddr := cliaddr.String()
		dlock()
		conn, found := clientDict[saddr]
		if !found {
			conn = newConnection(serverAddr, cliaddr)
			if conn == nil {
				dunlock()
				continue
			}
			clientDict[saddr] = conn
			dunlock()
			verboseLogf(2, "Created new connection for client %s\n", saddr)
			// Fire up routine to manage new connection
			go runConnection(conn)
		} else {
			verboseLogf(5, "Found connection for client %s\n", saddr)
			dunlock()
		}
		// Relay to server
		_, err = conn.serverConn.Write(buffer[0:n])
		if checkreport(1, err) {
			continue
		}
	}
}

var verbosity int = 6

// Log result if verbosity level high enough
func verboseLogf(level int, format string, v ...interface{}) {
	if level <= verbosity {
		log.Printf(format, v...)
	}
}

// Handle errors
func checkreport(level int, err error) bool {
	if err == nil {
		return false
	}
	verboseLogf(level, "Error: %s", err.Error())
	return true
}

func main() {
	var ihelp *bool = flag.Bool("h", false, "Show help information")
	var ipport *int = flag.Int("p", 6667, "Proxy port")
	var isport *int = flag.Int("P", 6666, "Server port")
	var ishost *string = flag.String("H", "localhost", "Server address")
	var iverb *int = flag.Int("v", 1, "Verbosity (0-6)")

	flag.Parse()
	verbosity = *iverb
	if *ihelp {
		flag.Usage()
		os.Exit(0)
	}
	if flag.NArg() > 0 {
		ok := true
		fields := strings.Split(flag.Arg(0), ":")
		ok = ok && len(fields) == 2
		if ok {
			*ishost = fields[0]
			n, err := fmt.Sscanf(fields[1], "%d", isport)
			ok = ok && n == 1 && err == nil
		}
		if !ok {
			flag.Usage()
			os.Exit(0)
		}
	}
	hostport := fmt.Sprintf("%s:%d", *ishost, *isport)
	verboseLogf(3, "Proxy port = %d, Server address = %s\n",
		*ipport, hostport)
	if setup(hostport, *ipport) {
		runProxy()
	}
	os.Exit(0)
}
