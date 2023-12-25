package container

import (
	"fmt"
	"sync"
)

// Lifecycle interface for dependencies that require initialization and cleanup
type Lifecycle interface {
	Init() error
	Close() error
}

// Container struct to hold and manage dependencies
type Container struct {
	dependencies map[string]interface{}
	lock         sync.RWMutex
}

// NewContainer creates a new instance of Container
func NewContainer() *Container {
	return &Container{
		dependencies: make(map[string]interface{}),
	}
}

// Register a dependency with a given name
func (c *Container) Register(name string, dep interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.dependencies[name]; exists {
		return fmt.Errorf("dependency already registered: %s", name)
	}

	c.dependencies[name] = dep

	if lifecycleDep, ok := dep.(Lifecycle); ok {
		return lifecycleDep.Init()
	}

	return nil
}

// Resolve retrieves a dependency by its name
func (c *Container) Resolve(name string) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	dep, ok := c.dependencies[name]
	if !ok {
		return nil, fmt.Errorf("dependency not found: %s", name)
	}

	return dep, nil
}

// Close all dependencies that implement the Lifecycle interface
func (c *Container) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, dep := range c.dependencies {
		if lifecycleDep, ok := dep.(Lifecycle); ok {
			if err := lifecycleDep.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
