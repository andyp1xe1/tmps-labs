# Lab 2 - File Converter with Creational Design Patterns

A file conversion system demonstrating three creational design patterns: Factory Method, Builder, and Object Pool through an extensible converter architecture.

## Creational Design Patterns Overview

**Creational design patterns** deal with object creation mechanisms, trying to create objects in a manner suitable to the situation. They solve problems related to object creation complexity by optimizing, hiding, or controlling the instantiation process:

- **Factory Method** - Creates objects without specifying their exact classes
- **Builder** - Constructs complex objects step by step with a fluent interface
- **Object Pool** - Manages reusable object instances for performance optimization
- **Singleton** - Ensures a class has only one instance
- **Abstract Factory** - Creates families of related objects
- **Prototype** - Creates objects by cloning existing instances

This lab implements **Factory Method**, **Builder**, and **Object Pool** patterns.

## Usage

```bash
cd lab2
go run client/main.go
```

### Example Output

```
lab2 $ go run client/main.go 
File Converter - Creational Design Patterns Demo
================================================
1. Factory Method Pattern Demo
------------------------------
✓ Created csv-json converter
  Sample output: [
  {
    "age": "25",
    "name": "John"
  },
  {...
✓ Created json-xml converter
✓ Created xml-yaml converter
✓ Created json-csv converter

2. Builder Pattern Demo
-----------------------
✓ Built conversion pipeline with 3 steps
  Input: /tmp/input.csv
  Output: /tmp/output.yaml
  Options: Indent=true, PrettyPrint=true
  Step 1: csv → json
  Step 2: json → xml
  Step 3: xml → yaml

3. Object Pool Pattern Demo
---------------------------
✓ Created converter pool with max size: 3
  Got converter 1 (csv-json), pool size: 0, created: 1
  Got converter 2 (json-xml), pool size: 0, created: 2
  Got converter 3 (xml-yaml), pool size: 0, created: 3
  Got converter 4 (csv-json), pool size: 0, created: 3
  Returned converter 1, pool size: 1
  Returned converter 2, pool size: 2
  Returned converter 3, pool size: 3
  Returned converter 4, pool size: 3
✓ Pool demonstration complete, final pool size: 3
```

## Testing

```bash
cd lab2
go test ./...
```

## Architecture & Design Patterns

### Project Structure

Following the lab requirements, the project is organized into modules based on responsibilities:

```
lab2/
├── client/               # Client application
│   └── main.go
├── domain/              # Domain logic
│   ├── factory/         # Factory patterns implementation
│   │   ├── converter_factory.go    # Factory Method + Registry
│   │   ├── converter_pool.go       # Object Pool
│   │   ├── pipeline_builder.go     # Builder
│   │   └── json_csv_converter.go   # Extension example
│   └── models/          # Domain models
│       ├── converter.go
│       └── pipeline.go
└── cond.txt
```

### Domain Models

Core interfaces and data structures for the file conversion system:

```go
// converter.go
type Converter interface {
    Convert(input io.Reader, from, to FileFormat) *ConversionResult
    SupportsFormat(format FileFormat) bool
}

type ConversionResult struct {
    Data   []byte
    Format FileFormat
    Error  error
}

// pipeline.go  
type Pipeline struct {
    Steps      []ConversionStep
    Options    ConversionOptions
    InputPath  string
    OutputPath string
}
```

### Factory Method Pattern

Creates different converter types through a registration-based factory that follows the **Open-Closed Principle**:

```go
// Global registry for converter creators
var converterRegistry = make(map[string]ConverterCreator)

func RegisterConverter(formatType string, creator ConverterCreator) {
    registryMutex.Lock()
    defer registryMutex.Unlock()
    converterRegistry[formatType] = creator
}

// Factory creates converters from registry
func (f *DefaultConverterFactory) CreateConverter(formatType string) (models.Converter, error) {
    registryMutex.RLock()
    creator, exists := converterRegistry[formatType]
    registryMutex.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("unsupported converter type: %s", formatType)
    }
    
    return creator(), nil
}
```

**Self-Registration**: Each converter registers itself during initialization:

```go
// Each converter type auto-registers
func init() {
    RegisterConverter("csv-json", func() models.Converter {
        return &CSVToJSONConverter{}
    })
}
```

**Benefits**:
- **Extensible**: New converters added without modifying factory code
- **Thread-safe**: Uses `sync.RWMutex` for concurrent access
- **Decoupled**: Factory doesn't know about concrete implementations

### Builder Pattern

Constructs complex conversion pipelines with a fluent API:

```go
pipeline, err := factory.NewPipelineBuilder().
    WithInputPath("/tmp/input.csv").
    WithOutputPath("/tmp/output.yaml").
    WithIndent().
    WithPrettyPrint().
    WithHeaders([]string{"name", "age", "city"}).
    AddCSVToJSON().
    AddJSONToXML().
    AddXMLToYAML().
    Build()
```

**Key Methods**:
- **Configuration**: `WithInputPath()`, `WithOutputPath()`, `WithOptions()`
- **Formatting**: `WithIndent()`, `WithPrettyPrint()`, `WithHeaders()`
- **Pipeline Steps**: `AddConversionStep()`, `AddCSVToJSON()`, etc.
- **Validation**: `Build()` validates required fields before creating pipeline

**Benefits**:
- **Readable**: Fluent interface makes complex construction clear
- **Flexible**: Optional parameters can be set in any order
- **Validated**: Build() ensures pipeline is properly configured

### Object Pool Pattern

Manages converter instances for reuse and performance optimization:

```go
type ConverterPool struct {
    pool    chan models.Converter
    factory ConverterFactory
    mu      sync.Mutex
    created int
    maxSize int
}

func (p *ConverterPool) Get(converterType string) (models.Converter, error) {
    // Try to get existing converter (fast path)
    select {
    case converter := <-p.pool:
        return converter, nil
    default:
        // Create new if under limit, otherwise fallback
        p.mu.Lock()
        defer p.mu.Unlock()
        
        if p.created < p.maxSize {
            converter, err := p.factory.CreateConverter(converterType)
            if err != nil {
                return nil, err
            }
            p.created++
            return converter, nil
        }
        
        // Pool exhausted, create temporary or wait
        return p.factory.CreateConverter(converterType)
    }
}
```

**Features**:
- **Bounded**: Respects maximum pool size to control memory usage
- **Non-blocking**: Fast path for available objects
- **Thread-safe**: Concurrent access protected by mutex
- **Graceful degradation**: Creates temporary objects when pool is full

### Converter Implementations

**CSV to JSON Converter**:
```go
func (c *CSVToJSONConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
    reader := csv.NewReader(input)
    records, err := reader.ReadAll()
    
    var jsonData []map[string]string
    if len(records) > 0 {
        headers := records[0]
        for _, record := range records[1:] {
            row := make(map[string]string)
            for i, value := range record {
                if i < len(headers) {
                    row[headers[i]] = value
                }
            }
            jsonData = append(jsonData, row)
        }
    }
    
    data, err := json.MarshalIndent(jsonData, "", "  ")
    return &models.ConversionResult{Data: data, Format: models.FormatJSON}
}
```

**Supported Conversions**:
- **CSV → JSON**: Tabular data to structured objects
- **JSON → XML**: Structured data to markup format  
- **XML → YAML**: Markup to human-readable format
- **JSON → CSV**: Reverse conversion (extensibility example)

## Open-Closed Principle Demonstration

The factory demonstrates the **Open-Closed Principle** - it's open for extension but closed for modification:

**Adding New Converter** (Extension without modification):
```go
// json_csv_converter.go - NEW FILE
type JSONToCSVConverter struct{}

func init() {
    RegisterConverter("json-csv", func() models.Converter {
        return &JSONToCSVConverter{}
    })
}
```

The factory code remains unchanged while supporting new conversion types.

## Pattern Benefits

### Factory Method
- **Decoupling**: Client code doesn't depend on concrete converter classes
- **Extensibility**: New converters added without changing existing code
- **Thread Safety**: Registry operations are protected from race conditions

### Builder  
- **Readability**: Complex pipeline construction becomes self-documenting
- **Flexibility**: Parameters can be set in any order with reasonable defaults
- **Validation**: Ensures required fields are set before object creation

### Object Pool
- **Performance**: Reduces object allocation overhead for expensive converters
- **Memory Control**: Bounded pool prevents unlimited resource consumption  
- **Concurrency**: Thread-safe access supports concurrent converter usage

## Conclusion

This implementation demonstrates how creational patterns solve common object instantiation challenges. The **Factory Method** provides extensible object creation, the **Builder** simplifies complex object construction, and the **Object Pool** optimizes resource management. Together, they create a robust, maintainable file conversion system that exemplifies good object-oriented design principles.

The self-registering factory pattern particularly showcases how the Open-Closed Principle enables extension without modification, making the system truly modular and maintainable.