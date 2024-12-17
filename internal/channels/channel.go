package channels

import (
	"context"
	"log"

	"github.com/matrixjnr/wave-message/internal/service"
	"github.com/matrixjnr/wave-message/pkg/message"
	"github.com/quic-go/quic-go"
	"wave-protocol/internal/routing"
	"wave-protocol/pkg/modulation"
)

// Channel represents a communication channel
type Channel struct {
	connection quic.Connection
	router     *routing.Router
	msgService *service.MessageService
}

// HandleSession handles a new QUIC connection session
func HandleSession(conn quic.Connection, router *routing.Router) {
	ch := &Channel{
		connection: conn,
		router:     router,
		msgService: service.NewMessageService(),
	}
	ch.handleStreams()
}

// handleStreams listens for incoming streams
func (ch *Channel) handleStreams() {
	ctx := context.Background()
	for {
		stream, err := ch.connection.AcceptStream(ctx)
		if err != nil {
			log.Printf("Failed to accept stream: %v", err)
			return
		}
		log.Println("Accepted new stream")

		// Process each stream
		go ch.handleStream(stream)
	}
}

// handleStream reads data, decodes it, validates the message, and routes it
func (ch *Channel) handleStream(stream quic.Stream) {
	buffer := make([]byte, 1024)

	for {
		n, err := stream.Read(buffer)
		if err != nil {
			log.Printf("Error reading from stream: %v", err)
			return
		}

		// QPSK Decode: Simulate decoding QPSK symbols (placeholder for real encoding)
		symbols := modulation.QPSKDecode(ch.convertToSymbols(buffer[:n]))
		log.Printf("QPSK Decoded: %v", symbols)

		// Deserialize the message
		data := symbolsToBytes(symbols)
		deserialized := ch.msgService.DeserializeData(data)
		msg, ok := deserialized.(message.Message)
		if !ok {
			log.Println("Invalid message format")
			continue
		}

		// Validate message
		if err := ch.msgService.ValidateMessage(&msg); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}

		// Route the message to appropriate listeners
		log.Printf("Routing message on channel %s from sender %s", msg.ChannelId, msg.SenderId)
		ch.router.Publish(msg.ChannelId, msg.SenderId, msg.Payload, msg.IsPersistent)
	}
}

// Helper: Convert byte slice to QPSK symbols
func (ch *Channel) convertToSymbols(data []byte) []complex64 {
	var symbols []complex64
	for _, b := range data {
		symbols = append(symbols, complex(float32(b), 0))
	}
	return symbols
}

// Helper: Convert QPSK symbols to byte slice
func symbolsToBytes(symbols []complex64) []byte {
	var bytes []byte
	for _, s := range symbols {
		bytes = append(bytes, byte(real(s)))
	}
	return bytes
}
