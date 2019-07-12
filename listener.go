package linodego

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// InstanceFuture represents a future Instance value.
type InstanceFuture struct {
	filterFunc func(*Instance) bool
	outCh      chan *Instance
}

// InstanceChan returns a channel that will emit the future Instance value.
func (s InstanceFuture) InstanceChan() <-chan *Instance {
	return s.outCh
}

// Listener represents a listener on the Linode API for various future events.
type Listener struct {
	instanceFutures map[int]InstanceFuture
	client          *Client
	mu              sync.Mutex
}

// NewListener creates a new Listener for the given client.
func NewListener(client *Client) *Listener {
	instanceFutures := make(map[int]InstanceFuture)

	return &Listener{
		instanceFutures: instanceFutures,
		client:          client,
	}
}

// AddInstance creates a future for the instance with the given ID, which will resolve to a concrete
// Instance value when the provided filter function evaluates to true.
func (s *Listener) AddInstanceFuture(id int, f func(*Instance) bool) InstanceFuture {
	outCh := make(chan *Instance)

	instanceFuture := InstanceFuture{
		filterFunc: f,
		outCh:      outCh,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.instanceFutures[id] = instanceFuture

	return instanceFuture
}

// RemoveInstanceFuture removes the given instance future from the main polling loop. This function
// also closes the output channel for the given future.
func (s *Listener) RemoveInstanceFuture(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	instanceFuture, ok := s.instanceFutures[id]
	if !ok {
		return fmt.Errorf("instance future for instance %d not found", id)
	}

	close(instanceFuture.outCh)
	delete(s.instanceFutures, id)

	return nil
}

// Poll polls the API for results for the listener's futures on the given interval. If an error
// occurs for any given polling operation, that error will be emitted on the returned error
// channel.
func (s *Listener) Poll(ctx context.Context, interval time.Duration) <-chan error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	doneCh := ctx.Done()
	errCh := make(chan error, 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				if errGet := s.getResults(ctx); errGet != nil {
					errCh <- errGet
				}
			case <-doneCh:
				errCh <- ctx.Err()
				return
			}
		}
	}()

	return errCh
}

// makeFilter makes a filter clause; a little ugly but it means we don't need reflection to
// generate JSON for the filters.
func (s *Listener) makeFilter(clauses []string) string {
	fmtString := `{"+or": [%s]}`

	clausesStr := strings.Join(clauses, ", ")

	return fmt.Sprintf(fmtString, clausesStr)
}

// instancesFilter creates a list of filter clauses for instances.
func (s *Listener) instancesFilter() string {
	fmtString := `{"id": "%s"}`

	clauses := make([]string, 0, len(s.instanceFutures))

	for instanceID := range s.instanceFutures {
		clauses = append(clauses, fmt.Sprintf(fmtString, instanceID))
	}

	return s.makeFilter(clauses)
}

// getResults gets the results of the current polling operation in the given context. If it
// succeeds, all results will be sent on their respective future channels if the associated
// future's filter function evaluates to true.
func (s *Listener) getResults(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	instanceFilter := s.instancesFilter()

	listOptions := NewListOptions(1, instanceFilter)

	instances, errInstances := s.client.ListInstances(ctx, listOptions)
	if errInstances != nil {
		return errInstances
	}

	// We could wrap each send in a goroutine, but that might just make it harder to track down
	// resource leaks later.
	for _, instance := range instances {

		// Bail if we somehow find an instance with no associated future.
		future, ok := s.instanceFutures[instance.ID]
		if !ok {
			return fmt.Errorf("unexpected instance %d found", instance.ID)
		}

		// Only emit a value for that future if it passes that future's filter function.
		if future.filterFunc(&instance) {
			future.outCh <- &instance
		}
	}

	return nil
}
