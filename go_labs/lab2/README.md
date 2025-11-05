# Lab 2 - File Converter with Creational Design Patterns

An extensible file conversion system that converts CSV data through JSON and XML formats to YAML output, demonstrating three creational design patterns: Factory Method, Builder, and Object Pool.

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

The application converts the sample CSV data (`input_sample.csv`) through a pipeline: CSV → JSON → XML → YAML, producing `output_final.yaml`.

### Example Output

```
Creational Design Patterns Demo: CSV → JSON → XML → YAML
Processed 3 conversion steps in 1 ms
  Step 1: csv → json (1.2 KB)
  Step 2: json → xml (2.1 KB)
  Step 3: xml → yaml (1.8 KB)
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
│   │   ├── pipeline_builder.go     # Builder + Pipeline Executor
│   │   ├── csv_json_converter.go   # CSV to JSON converter
│   │   ├── json_xml_converter.go   # JSON to XML converter
│   │   └── xml_yaml_converter.go   # XML to YAML converter
│   └── models/          # Domain models
│       ├── converter.go # Converter interface and types
│       └── pipeline.go  # Pipeline and execution types
├── input_sample.csv     # Sample input data
└── output_final.yaml    # Generated output
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

type ConversionStep struct {
    From FileFormat
    To   FileFormat
}

type PipelineResult struct {
    Success  bool
    Results  []*ConversionResult
    Error    error
    Duration int64
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

Constructs complex conversion pipelines with a fluent API and executes them:

```go
// Building a pipeline
pipeline, err := factory.NewPipelineBuilder().
    WithInputPath("input_sample.csv").
    WithOutputPath("output_final.yaml").
    WithIndent().
    WithPrettyPrint().
    AddCSVToJSON().
    AddJSONToXML().
    AddXMLToYAML().
    Build()

// Executing the pipeline
executor := factory.NewPipelineExecutor(pool)
result := executor.Execute(pipeline)
```

**Key Methods**:
- **Configuration**: `WithInputPath()`, `WithOutputPath()`, `WithOptions()`
- **Formatting**: `WithIndent()`, `WithPrettyPrint()`, `WithHeaders()`
- **Pipeline Steps**: `AddConversionStep()`, `AddCSVToJSON()`, `AddJSONToXML()`, `AddXMLToYAML()`
- **Validation**: `Build()` validates required fields before creating pipeline

**Benefits**:
- **Readable**: Fluent interface makes complex construction clear
- **Flexible**: Optional parameters can be set in any order
- **Validated**: Build() ensures pipeline is properly configured

### Object Pool Pattern

Manages converter instances per type for reuse and performance optimization:

```go
type ConverterPool struct {
    pools   map[string]chan models.Converter  // Per-type pools
    factory ConverterFactory
    mu      sync.Mutex
    created map[string]int                    // Per-type counters
    maxSize int
}

func (p *ConverterPool) Get(converterType string) (models.Converter, error) {
    // Fast path: try to get existing converter from type-specific pool
    select {
    case converter := <-p.pools[converterType]:
        return converter, nil
    default:
        // Create new if under limit, otherwise fallback to temporary
        p.mu.Lock()
        if p.created[converterType] < p.maxSize {
            converter, err := p.factory.CreateConverter(converterType)
            if err == nil {
                p.created[converterType]++
            }
            p.mu.Unlock()
            return converter, err
        }
        p.mu.Unlock()
        
        // Pool exhausted, create temporary
        return p.factory.CreateConverter(converterType)
    }
}
```

**Features**:
- **Type-specific pools**: Separate pools for each converter type (csv-json, json-xml, xml-yaml)
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
    
    if len(records) == 0 {
        return &models.ConversionResult{Data: []byte("[]"), Format: models.FormatJSON}
    }
    
    headers := records[0]
    var jsonData []map[string]string
    
    for _, record := range records[1:] {
        row := make(map[string]string)
        for i, value := range record {
            if i < len(headers) {
                row[headers[i]] = value
            }
        }
        jsonData = append(jsonData, row)
    }
    
    data, err := json.MarshalIndent(jsonData, "", "  ")
    return &models.ConversionResult{Data: data, Format: models.FormatJSON}
}
```

**JSON to XML Converter**:
```go
func (j *JSONToXMLConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
    jsonData, err := io.ReadAll(input)
    var data interface{}
    json.Unmarshal(jsonData, &data)
    
    // Convert to XML using mxj library
    mv := mxj.Map{"root": data}
    xmlData, err := mv.XmlIndent("", "  ")
    
    return &models.ConversionResult{Data: xmlData, Format: models.FormatXML}
}
```

**XML to YAML Converter**:
```go
func (x *XMLToYAMLConverter) Convert(input io.Reader, from, to models.FileFormat) *models.ConversionResult {
    xmlData, err := io.ReadAll(input)
    
    // Parse XML using mxj library
    mv, err := mxj.NewMapXml(xmlData)
    
    // Convert map to YAML using gopkg.in/yaml.v3
    yamlData, err := yaml.Marshal(mv.Old())
    
    return &models.ConversionResult{Data: yamlData, Format: models.FormatYAML}
}
```

**Supported Conversions**:
- **CSV → JSON**: Tabular data to structured objects using headers as keys
- **JSON → XML**: Structured data to markup format using mxj library
- **XML → YAML**: Markup to human-readable format using yaml.v3

**Dependencies**:
- `github.com/clbanning/mxj/v2` for JSON/XML conversions
- `gopkg.in/yaml.v3` for YAML marshaling

## Open-Closed Principle Demonstration

The factory demonstrates the **Open-Closed Principle** - it's open for extension but closed for modification:

**Self-Registration**: Each converter registers itself during initialization:

```go
// csv_json_converter.go
func init() {
    RegisterConverter("csv-json", func() models.Converter {
        return &CSVToJSONConverter{}
    })
}

// json_xml_converter.go
func init() {
    RegisterConverter("json-xml", func() models.Converter {
        return &JSONToXMLConverter{}
    })
}

// xml_yaml_converter.go
func init() {
    RegisterConverter("xml-yaml", func() models.Converter {
        return &XMLToYAMLConverter{}
    })
}
```

**Adding New Converter** (Extension without modification):
```go
// new_converter.go - NEW FILE
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
- **Memory Control**: Bounded pools per type prevent unlimited resource consumption  
- **Concurrency**: Thread-safe access supports concurrent converter usage
- **Type Isolation**: Separate pools for each converter type provide better resource management

## Sample Data

**Input** (`input_sample.csv`):
```csv
name,age,city,occupation
Alice Johnson,28,New York,Software Engineer
Bob Smith,35,San Francisco,Data Scientist  
Carol Davis,42,Austin,Product Manager
David Wilson,31,Seattle,DevOps Engineer
Emma Brown,29,Boston,UX Designer
Frank Garcia,38,Chicago,Backend Developer
Grace Lee,33,Denver,Frontend Developer
Henry Chen,27,Portland,Mobile Developer
```

**Output** (`output_final.yaml`):
```yaml
doc:
    root:
        - age: "28"
          city: New York
          name: Alice Johnson
          occupation: Software Engineer
        - age: "35"
          city: San Francisco
          name: Bob Smith
          occupation: Data Scientist
        # ... remaining records
```

## Conclusion

This implementation demonstrates how creational patterns solve common object instantiation challenges in a real file conversion pipeline. The **Factory Method** provides extensible object creation through self-registration, the **Builder** simplifies complex pipeline construction with a fluent API, and the **Object Pool** optimizes resource management with type-specific pools. Together, they create a robust, maintainable file conversion system that processes CSV data through JSON and XML formats to produce YAML output.

The self-registering factory pattern particularly showcases how the Open-Closed Principle enables extension without modification, making the system truly modular and maintainable. The actual working pipeline demonstrates these patterns in action, converting real data through multiple format transformations.