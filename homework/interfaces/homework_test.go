package main

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ServiceWithNestedDependencies struct {
	FooDependency *FooDependency
}

type FooDependency struct {
	BarDependency *BarDependency
}

type BarDependency struct {
	YourPropertyCanBeHere string
}

type UserService struct {
	NotEmptyStruct bool
}

type MessageService struct {
	NotEmptyStruct bool
}

type Container struct {
	mu             sync.RWMutex
	funcConstructs map[string]func(c *Container) any
	singletons     map[string]any
}

func NewContainer() *Container {
	return &Container{
		funcConstructs: make(map[string]func(c *Container) any),
		singletons:     make(map[string]any),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.alreadyRegistered(name) {
		// it's better to return error here, but let's save compatibility with the original code
		panic("service already registered")
	}
	switch constr := constructor.(type) {
	case func() any:
		// if you want to use function without arguments, you need to pass container as a parameter
		c.funcConstructs[name] = func(_ *Container) any {
			return constr()
		}
	case func(c *Container) any:
		c.funcConstructs[name] = constr
	default:
		panic("unexpected constructor type")
	}
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.alreadyRegistered(name) {
		// it's better to return error here, but let's save compatibility with the original code
		panic("service already registered")
	}
	switch constr := constructor.(type) {
	case func() any:
		c.singletons[name] = constr()
	case func(c *Container) any:
		c.singletons[name] = constr(c)
	default:
		panic("unexpected constructor type")
	}
}

func (c *Container) alreadyRegistered(name string) bool {
	_, funcOk := c.funcConstructs[name]
	_, singleOk := c.singletons[name]
	return funcOk || singleOk
}

func (c *Container) Resolve(name string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if funcConstruct, ok := c.funcConstructs[name]; ok {
		return funcConstruct(c), nil
	}
	if singleton, ok := c.singletons[name]; ok {
		return singleton, nil
	}
	return nil, errors.New("service is not registered")
}

// MustResolve is a helper function to resolve a service and panic if it fails
func (c *Container) MustResolve(name string) any {
	res, err := c.Resolve(name)
	if err != nil {
		panic(err)
	}
	return res
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}

func TestDIContainerWithNestedDependencies(t *testing.T) {
	container := NewContainer()
	container.RegisterType("BarDependency", func(c *Container) interface{} {
		return &BarDependency{
			YourPropertyCanBeHere: "test",
		}
	})
	container.RegisterType("FooDependency", func(c *Container) interface{} {
		return &FooDependency{
			BarDependency: c.MustResolve("BarDependency").(*BarDependency),
		}
	})
	container.RegisterType("ServiceWithNestedDependencies", func(c *Container) interface{} {
		return &ServiceWithNestedDependencies{
			FooDependency: c.MustResolve("FooDependency").(*FooDependency),
		}
	})

	service1, err := container.Resolve("ServiceWithNestedDependencies")
	assert.NoError(t, err)
	service2, err := container.Resolve("ServiceWithNestedDependencies")
	assert.NoError(t, err)

	s1, ok := service1.(*ServiceWithNestedDependencies)
	assert.True(t, ok)
	s2, ok := service2.(*ServiceWithNestedDependencies)
	assert.True(t, ok)

	assert.False(t, s1 == s2)
}

func TestDIContainerSingleton(t *testing.T) {
	container := NewContainer()
	container.RegisterSingletonType("UserService", func() interface{} {
		return &UserService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.True(t, u1 == u2)
}

func TestDIContainerConcurrent(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	var wg sync.WaitGroup
	const numGoroutines = 100

	// Test concurrent resolution of UserService
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			userService, err := container.Resolve("UserService")
			assert.NoError(t, err)
			assert.NotNil(t, userService)
		}()
	}
	wg.Wait()

	// Test concurrent resolution of MessageService
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			messageService, err := container.Resolve("MessageService")
			assert.NoError(t, err)
			assert.NotNil(t, messageService)
		}()
	}
	wg.Wait()
}

func TestDIContainerConcurrentSingleton(t *testing.T) {
	container := NewContainer()
	container.RegisterSingletonType("UserService", func() interface{} {
		return &UserService{}
	})

	var wg sync.WaitGroup
	const numGoroutines = 100

	// Test concurrent resolution of singleton UserService
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			userService, err := container.Resolve("UserService")
			assert.NoError(t, err)
			assert.NotNil(t, userService)
		}()
	}
	wg.Wait()
}
