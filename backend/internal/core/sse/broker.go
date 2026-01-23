package sse

type Subscription struct {
	topic  string
	client chan string
}

type Event struct {
	Data string
	Type string
}

type Broker struct {
	subscribe   chan Subscription
	unsubscribe chan Subscription
	notifier    chan Event
	subscribers map[string]map[chan string]struct{}
}

func NewBroker() *Broker {
	b := &Broker{
		subscribe:   make(chan Subscription),
		unsubscribe: make(chan Subscription),
		notifier:    make(chan Event),
		subscribers: make(map[string]map[chan string]struct{}),
	}

	go b.listen()
	return b
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.subscribe:
			if b.subscribers[s.topic] == nil {
				b.subscribers[s.topic] = make(map[chan string]struct{})
			}
			b.subscribers[s.topic][s.client] = struct{}{}

		case s := <-b.unsubscribe:
			if clients, ok := b.subscribers[s.topic]; ok {
				delete(clients, s.client)
				close(s.client)
			}

		case s := <-b.notifier:
			if clients, ok := b.subscribers[s.Type]; ok {
				for client := range clients {
					client <- s.Data
				}
			}
		}
	}
}

func (b *Broker) Subscribe(topic string, client chan string) {
	b.subscribe <- Subscription{
		topic:  topic,
		client: client,
	}
}

func (b *Broker) Unsubscribe(topic string, client chan string) {
	b.unsubscribe <- Subscription{
		topic:  topic,
		client: client,
	}
}

func (b *Broker) Notify(data string, eventType string) {
	b.notifier <- Event{
		Data: data,
		Type: eventType,
	}
}
