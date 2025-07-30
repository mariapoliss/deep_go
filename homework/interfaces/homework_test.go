package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Constructor = func() interface{}

type Container struct {
	types map[string]Constructor
}

func NewContainer() *Container {
	return &Container{
		types: make(map[string]Constructor),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	cnstr, ok := constructor.(Constructor)
	if !ok {
		return
	}
	c.types[name] = cnstr
}

func (c *Container) Resolve(name string) (interface{}, error) {
	constructor, found := c.types[name]
	if !found {
		return nil, fmt.Errorf("Not found")
	}
	return constructor(), nil
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
