package quic

import (
	"context"
	"crypto/tls"
	"github.com/matrixjnr/wave/internal/channels"
	"github.com/matrixjnr/wave/internal/routing"
	"log"

	"github.com/quic-go/quic-go"
)

// Server represents the QUIC server
type Server struct {
	addr    string
	tlsConf *tls.Config
}

// NewServer initializes a new QUIC server
func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		tlsConf: generateTLSConfig(),
	}
}

// ListenAndServe starts the QUIC server
func (s *Server) ListenAndServe(router *routing.Router) error {
	listener, err := quic.ListenAddr(s.addr, s.tlsConf, nil)
	if err != nil {
		return err
	}
	log.Printf("QUIC server listening on %s", s.addr)

	for {
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("Failed to accept session: %v", err)
			continue
		}
		log.Println("Accepted new session")

		// Pass session to the channel handler
		go channels.HandleSession(session, router)
	}
}

// generateTLSConfig generates TLS configuration
func generateTLSConfig() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{generateSelfSignedCert()},
		NextProtos:   []string{"wave-protocol"},
	}
}

// generateSelfSignedCert loads placeholder certificates
func generateSelfSignedCert() tls.Certificate {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Failed to load certificates: %v", err)
	}
	return cert
}
