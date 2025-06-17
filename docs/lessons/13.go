package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// 1. Basic interfaces
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

// 2. Shape interface example
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.Radius
}

type Square struct {
	Side float64
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

func (s Square) Perimeter() float64 {
	return 4 * s.Side
}

// 3. Custom error interface
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// 4. Storage interface example
type Storage interface {
	Save(data []byte) error
	Load() ([]byte, error)
	Delete() error
}

// File Storage implementation
type FileStorage struct {
	path string
}

func (fs FileStorage) Save(data []byte) error {
	return ioutil.WriteFile(fs.path, data, 0644)
}

func (fs FileStorage) Load() ([]byte, error) {
	return ioutil.ReadFile(fs.path)
}

func (fs FileStorage) Delete() error {
	return os.Remove(fs.path)
}

// Memory Storage implementation
type MemoryStorage struct {
	data []byte
}

func (ms *MemoryStorage) Save(data []byte) error {
	ms.data = make([]byte, len(data))
	copy(ms.data, data)
	return nil
}

func (ms *MemoryStorage) Load() ([]byte, error) {
	return ms.data, nil
}

func (ms *MemoryStorage) Delete() error {
	ms.data = nil
	return nil
}

// 5. Plugin system example
type Plugin interface {
	Name() string
	Initialize() error
	Execute() error
	Cleanup() error
}

type LoggerPlugin struct {
	startTime time.Time
}

func (p *LoggerPlugin) Name() string {
	return "Logger"
}

func (p *LoggerPlugin) Initialize() error {
	p.startTime = time.Now()
	fmt.Println("Logger plugin initialized")
	return nil
}

func (p *LoggerPlugin) Execute() error {
	fmt.Printf("Logger running (started at: %v)\n", p.startTime)
	return nil
}

func (p *LoggerPlugin) Cleanup() error {
	fmt.Println("Logger plugin cleaned up")
	return nil
}

// Plugin Manager
type PluginManager struct {
	plugins []Plugin
}

func (pm *PluginManager) Register(p Plugin) {
	pm.plugins = append(pm.plugins, p)
}

func (pm *PluginManager) ExecuteAll() error {
	for _, p := range pm.plugins {
		if err := p.Execute(); err != nil {
			return fmt.Errorf("plugin %s failed: %v", p.Name(), err)
		}
	}
	return nil
}

func main() {
	// 1. Shape interface demo
	fmt.Println("=== Shape Interface Demo ===")
	shapes := []Shape{
		Circle{Radius: 5},
		Square{Side: 4},
	}

	for _, shape := range shapes {
		fmt.Printf("Area: %.2f, Perimeter: %.2f\n",
			shape.Area(), shape.Perimeter())
	}

	// 2. Custom error demo
	fmt.Println("\n=== Custom Error Demo ===")
	err := ValidationError{
		Field:   "email",
		Message: "invalid email format",
	}
	fmt.Println("Error:", err)

	// 3. Storage interface demo
	fmt.Println("\n=== Storage Interface Demo ===")

	// Test with memory storage
	memStorage := &MemoryStorage{}
	data := []byte("Hello, World!")

	fmt.Println("Testing MemoryStorage:")
	fmt.Println("Saving data...")
	memStorage.Save(data)

	if loaded, err := memStorage.Load(); err == nil {
		fmt.Printf("Loaded data: %s\n", loaded)
	}

	fmt.Println("Deleting data...")
	memStorage.Delete()

	if loaded, err := memStorage.Load(); err == nil {
		fmt.Printf("After delete: %v\n", loaded)
	}

	// 4. Plugin system demo
	fmt.Println("\n=== Plugin System Demo ===")
	pm := &PluginManager{}

	logger := &LoggerPlugin{}
	pm.Register(logger)

	logger.Initialize()
	pm.ExecuteAll()
	logger.Cleanup()

	// 5. Empty interface and type assertion demo
	fmt.Println("\n=== Empty Interface and Type Assertion Demo ===")

	var i interface{}
	i = "Hello"

	// Type assertion
	if str, ok := i.(string); ok {
		fmt.Printf("Value is string: %s\n", str)
	}

	i = 42

	// Type switch
	switch v := i.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	default:
		fmt.Printf("Unknown type\n")
	}
}
