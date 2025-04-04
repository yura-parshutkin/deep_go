package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

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
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	funcConstructs map[string]func(c *Container) any
	singletons     map[string]any
}

func NewContainer() *Container {
	// need to implement
	return &Container{
		funcConstructs: make(map[string]func(c *Container) any),
		singletons:     make(map[string]any),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
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

func (c *Container) Resolve(name string) (interface{}, error) {
	funcConstruct, ok := c.funcConstructs[name]
	if ok {
		return funcConstruct(c), nil
	}
	singleton, ok := c.singletons[name]
	if ok {
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
