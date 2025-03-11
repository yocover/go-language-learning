# PostgreSQL JSON String Storage Example

This example demonstrates how to efficiently store and retrieve large JSON documents as strings in PostgreSQL using Go. The implementation focuses on simple storage and retrieval operations without complex JSON querying capabilities, making it ideal for scenarios where JSON data is treated as a complete document.

## Features

- Store large JSON documents as TEXT in PostgreSQL
- Basic CRUD operations (Create, Read, Update, Delete)
- Simple file listing functionality
- Efficient storage without JSON parsing overhead
- Suitable for large JSON documents
- Minimal database indexing for better write performance

## Project Structure

```
json_store/
├── cmd/
│   └── main.go         # Main application entry point
├── store/
│   └── store.go        # Database operations implementation
├── schema.sql          # Database schema
└── README.md          # This file
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- `github.com/lib/pq` package

## Database Schema

The database schema is designed for efficient storage of large JSON documents:

```sql
CREATE TABLE IF NOT EXISTS file_json_store (
    id SERIAL PRIMARY KEY,
    file_id VARCHAR(255) NOT NULL UNIQUE,
    file_name VARCHAR(255) NOT NULL,
    json_content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Usage

1. Set up your PostgreSQL database and update the connection string in `main.go` if needed.

2. Run the example:
   ```bash
   go run cmd/main.go
   ```

The example demonstrates:
- Creating a new JSON document
- Retrieving a stored document
- Updating document content
- Listing all stored files
- Deleting a document

## Key Benefits

- **Efficient Storage**: Stores JSON as plain text, avoiding JSONB parsing and indexing overhead
- **Better Write Performance**: No JSON validation or indexing during writes
- **Simple Implementation**: Straightforward CRUD operations without complex query logic
- **Suitable for Large Documents**: Ideal for storing and retrieving complete JSON documents
- **Lower Storage Overhead**: No additional indexing space required

## When to Use This Approach

This implementation is particularly suitable when:
- You need to store large JSON documents
- You don't need to query internal JSON fields
- Write performance is a priority
- You always operate on the complete JSON document
- Storage space efficiency is important

## Limitations

- No internal JSON field querying
- No JSON validation at the database level
- Must parse JSON in the application layer when needed
- No JSON-specific operations or functions

## License

MIT 