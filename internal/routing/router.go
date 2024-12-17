package routing

import (
	"log"
	"sync"

	"github.com/matrixjnr/wave-message/internal/service"
	"github.com/matrixjnr/wave-message/pkg/message"
)

// Listener represents a channel for receiving messages
type Listener chan *message.Message

// Router manages message routing between channels and their listeners
type Router struct {
	subscriptions map[string][]Listener // ChannelID -> List of listeners
	mu            sync.RWMutex          // Mutex for thread-safe operations
	msgService    *service.MessageService
}

// NewRouter initializes and returns a new Router
func NewRouter() *Router {
	return &Router{
		subscriptions: make(map[string][]Listener),
		msgService:    service.NewMessageService(),
	}
}

// Subscribe registers a new listener to a specific channel
func (r *Router) Subscribe(channelID string) Listener {
	r.mu.Lock()
	defer r.mu.Unlock()

	listener := make(Listener, 10) // Buffered channel for message delivery
	r.subscriptions[channelID] = append(r.subscriptions[channelID], listener)
	log.Printf("Listener subscribed to channel: %s", channelID)
	return listener
}

// Unsubscribe removes a listener from a specific channel
func (r *Router) Unsubscribe(channelID string, listener Listener) {
	r.mu.Lock()
	defer r.mu.Unlock()

	listeners, exists := r.subscriptions[channelID]
	if !exists {
		log.Printf("Attempted to unsubscribe from non-existent channel: %s", channelID)
		return
	}

	// Find and remove the listener
	for i, l := range listeners {
		if l == listener {
			r.subscriptions[channelID] = append(listeners[:i], listeners[i+1:]...)
			close(listener)
			log.Printf("Listener unsubscribed from channel: %s", channelID)
			return
		}
	}
	log.Printf("Listener not found in channel: %s", channelID)
}

// Publish broadcasts a message to all listeners on a specific channel
func (r *Router) Publish(channelID, senderID string, payload interface{}, isPersistent bool) error {
	// Create a new message
	msg, err := r.msgService.CreateMessage(channelID, senderID, payload, isPersistent)
	if err != nil {
		log.Printf("Failed to create message: %v", err)
		return err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	listeners, exists := r.subscriptions[channelID]
	if !exists || len(listeners) == 0 {
		log.Printf("No listeners for channel: %s", channelID)
		return nil
	}

	// Broadcast the message to all listeners
	for i, listener := range listeners {
		select {
		case listener <- msg:
			log.Printf("Message published to channel: %s by sender: %s", channelID, senderID)
		default:
			log.Printf("Listener %d buffer full. Skipping message for channel: %s", i, channelID)
		}
	}

	return nil
}

// CleanupStaleListeners removes channels with no active listeners
func (r *Router) CleanupStaleListeners() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for channelID, listeners := range r.subscriptions {
		activeListeners := []Listener{}
		for _, l := range listeners {
			select {
			case <-l:
				log.Printf("Removed stale listener from channel: %s", channelID)
			default:
				activeListeners = append(activeListeners, l)
			}
		}
		if len(activeListeners) == 0 {
			delete(r.subscriptions, channelID)
			log.Printf("Removed empty channel: %s", channelID)
		} else {
			r.subscriptions[channelID] = activeListeners
		}
	}
}
