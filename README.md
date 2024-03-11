# DataWizard

DataWizard allows you to **parse**, **process** and **export** ATS (Applicant Tracking System) job data before being used by websites or other apps.

DataWizard allows you to connect to ATS via HTTP or SFTP and can parse both XML and JSON data. The processed data can be output in XML or JSON format. 

## Technical concepts
### SimpleMapping
The `SimpleMapping` property maps source data to exported data without any processing or manipulation. 

The `key` : `value` pairs correspond to `exported-key-name`:`key-in-source`.

```go
    SimpleMapping: map[string]string{
        "id":        "job_id",
        "title":     "title",
    }
```
The above configuration will result in the following output data:
```json
    JSON
        {
            "id": "...",
            "title": "...",
        }
```
```xml
    XML
    <job>
        <id>...</id>
        <title>...</title>
    </job>
```

### FunctionMapping
The `FunctionMapping` property allows you to apply processing functions to the source data before it is exported.

The following custom functions are already available and examples are provided below.

#### 1. CONCATENATE
**Use case**: Need to concatenate two tags into one separated by a custom delimiter. **Example**: "{city}, {country}".

```go
    FunctionMapping: map[string]map[string]interface{}{
        "CONCATENATE": {                                // Function name
            "TAGS":       []string{"city", "country"},  // Tags to concatenate
            "DELIMITER":  ", ",                         // Where to find elements (don't include JobElement)
            "EXPORTAS":   "location",                   // Exported key
        },
    },
```
The above configuration will result in the following output data:
```json
    JSON
        {
            "location": "Lisbon, Portugal"
        }
```
```xml
    XML
    <job>
        <location>Lisbon, Portugal</location>
    </job>
```