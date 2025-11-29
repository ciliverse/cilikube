# CiliKube API

This directory contains API definitions, specifications, and documentation for CiliKube.

## API Versions

- **v1** (`/api/v1/`) - Current stable API version

## Directory Structure

```
api/
├── README.md           # This file
└── v1/                 # API v1 definitions
    ├── README.md       # v1 API documentation
    ├── openapi.yaml    # OpenAPI 3.0 specification
    └── schemas/        # JSON schemas (if needed)
```

## API Design Principles

### RESTful Design
- Follow REST conventions for resource operations
- Use appropriate HTTP methods (GET, POST, PUT, DELETE)
- Use meaningful HTTP status codes
- Consistent URL patterns

### Consistent Response Format
All API responses follow a standard format:
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (RBAC)
- Secure token handling

### Error Handling
- Consistent error response format
- Meaningful error messages
- Appropriate HTTP status codes

### Versioning Strategy
- URL path versioning (`/api/v1/`)
- Backward compatibility within major versions
- Clear deprecation notices for breaking changes

## Implementation

The actual API implementation is located in:
- **Handlers**: `internal/handlers/` - HTTP request handlers
- **Routes**: `internal/routes/` - Route definitions and middleware
- **Models**: `internal/models/` - Data models and DTOs
- **Services**: `internal/service/` - Business logic layer

## Documentation

- **OpenAPI Specification**: Complete API specification in OpenAPI 3.0 format
- **Interactive Documentation**: Can be viewed with Swagger UI or similar tools
- **Code Examples**: Provided in language-specific documentation

## Development Workflow

1. **Design**: Define API endpoints in OpenAPI specification
2. **Review**: Review API design with team
3. **Implement**: Implement handlers, routes, and business logic
4. **Test**: Write and run API tests
5. **Document**: Update documentation and examples
6. **Deploy**: Deploy to staging/production environments

## Tools and Resources

- **OpenAPI Generator**: Generate client SDKs and documentation
- **Swagger UI**: Interactive API documentation
- **Postman**: API testing and development
- **curl**: Command-line API testing

## Contributing

When contributing to the API:

1. Follow existing patterns and conventions
2. Update OpenAPI specification for any changes
3. Add appropriate tests
4. Update documentation
5. Ensure backward compatibility