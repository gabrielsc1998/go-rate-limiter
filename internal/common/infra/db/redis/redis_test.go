package redis

import (
	"context"
	"testing"
	"time"
)

func TestNewRedisClient(t *testing.T) {
	config := RedisClientConfig{Addr: "redis:6379", Password: "", DB: 0}
	client := NewRedisClient(config)

	if client == nil {
		t.Errorf("Expected client to be not nil")
	}
}

func TestSetAndGet(t *testing.T) {
	config := RedisClientConfig{Addr: "redis:6379", Password: "", DB: 0}
	client := NewRedisClient(config)

	err := client.Set(context.Background(), Input{Key: "key", Value: "value", Expire: 0})
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	value, err := client.Get(context.Background(), "key")
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value', got '%s'", value)
	}
}

func TestSetAndGetWithExpire(t *testing.T) {
	config := RedisClientConfig{Addr: "redis:6379", Password: "", DB: 0}
	client := NewRedisClient(config)

	err := client.Set(context.Background(), Input{Key: "key", Value: "value", Expire: 1})
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	value, err := client.Get(context.Background(), "key")
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value', got '%s'", value)
	}

	// Wait for key to expire
	time.Sleep(1 * time.Second)

	value, err = client.Get(context.Background(), "key")
	if err.Error() != "redis: nil" {
		t.Errorf("Expected error to be 'redis: nil', got '%s'", err.Error())
	}
	if value != "" {
		t.Errorf("Expected value to be '', got '%s'", value)
	}
}

func TestDel(t *testing.T) {
	config := RedisClientConfig{Addr: "redis:6379", Password: "", DB: 0}
	client := NewRedisClient(config)

	err := client.Set(context.Background(), Input{Key: "key", Value: "value"})
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	value, err := client.Get(context.Background(), "key")
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value', got '%s'", value)
	}

	err = client.Del("key", context.Background())
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	value, err = client.Get(context.Background(), "key")
	if err.Error() != "redis: nil" {
		t.Errorf("Expected error to be 'redis: nil', got '%s'", err.Error())
	}
	if value != "" {
		t.Errorf("Expected value to be '', got '%s'", value)
	}
}

func TestClearAll(t *testing.T) {
	config := RedisClientConfig{Addr: "redis:6379", Password: "", DB: 0}
	client := NewRedisClient(config)

	err := client.Set(context.Background(), Input{Key: "key", Value: "value"})
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	value, err := client.Get(context.Background(), "key")
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}
	if value != "value" {
		t.Errorf("Expected value to be 'value', got '%s'", value)
	}

	err = client.ClearAll(context.Background())
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	value, err = client.Get(context.Background(), "key")
	if err.Error() != "redis: nil" {
		t.Errorf("Expected error to be 'redis: nil', got '%s'", err.Error())
	}
	if value != "" {
		t.Errorf("Expected value to be '', got '%s'", value)
	}
}
