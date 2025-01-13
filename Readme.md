# Retail Pulse Image Processing Service

A Go-based microservice that processes store visit images in parallel, calculates image perimeters, and manages job statuses.

## Description

This service provides APIs to:

- Submit image processing jobs for store visits
- Track job processing status
- Handle multiple concurrent jobs
- Process thousands of images in parallel
- Validate store information against a master store list

### Key Features

- Concurrent image processing using goroutines
- Job status tracking
- Store master data validation
- Error handling for failed downloads or invalid stores
- Docker support
- Comprehensive testing suite

## Assumptions

1. **Store Master Data**:

   - Stores are pre-loaded from a CSV file named `StoreMasterAssignment.csv`
   - CSV format: AreaCode, StoreName, StoreID
   - Store IDs are unique

2. **Image Processing**:

   - Images are accessible via HTTP(S) URLs
   - Processing time is simulated (0.1 to 0.4 seconds per image)
   - Failed image downloads should fail the entire job

3. **Performance**:
   - Multiple jobs can run concurrently
   - Each job can process multiple store visits in parallel
   - Multiple images per store visit are processed in parallel

## Installation & Setup

### Prerequisites

- Go 1.21 or later
- Docker (optional)
- jq (for running tests)

### Local Setup

1. Clone the repository:

```bash
git clone https://github.com/subho004/retail-pulse
cd retail-pulse
```

2. Initialize Go modules:

```bash
go mod init retail-pulse
go mod tidy
```

3. Place store master data:

```bash
# Ensure StoreMasterAssignment.csv is in the project root
```

4. Run the server:

```bash
go run cmd/server/main.go
```

### Docker Setup

1. Build and run using Docker Compose:

```bash
docker-compose up --build
```

## API Routes

### 1. Submit Job

Submit a new image processing job.

**Endpoint**: `POST /api/submit`

**Request Body**:

```json
{
  "count": 2,
  "visits": [
    {
      "store_id": "RP00001",
      "image_url": [
        "https://example.com/image1.jpg",
        "https://example.com/image2.jpg"
      ],
      "visit_time": "2024-01-13T12:00:00Z"
    },
    {
      "store_id": "RP00002",
      "image_url": ["https://example.com/image3.jpg"],
      "visit_time": "2024-01-13T12:30:00Z"
    }
  ]
}
```

**Success Response** (201 Created):

```json
{
  "job_id": "job_1234567890"
}
```

**Error Response** (400 Bad Request):

```json
{
  "error": "count does not match number of visits"
}
```

### 2. Get Job Status

Check the status of a submitted job.

**Endpoint**: `GET /api/status?jobid=<job_id>`

**Success Response** (200 OK):

```json
{
  "status": "completed",
  "job_id": "job_1234567890"
}
```

**Failed Job Response** (200 OK):

```json
{
  "status": "failed",
  "job_id": "job_1234567890",
  "error": [
    {
      "store_id": "RP00001",
      "error": "store not found"
    }
  ]
}
```

**Error Response** (400 Bad Request):

```json
{
  "error": "job not found"
}
```

## Testing

### Automated Testing

Run the test script:

```bash
# Make script executable
chmod +x test_api.sh

# Run tests
./test_api.sh
```

The test script checks:

1. Valid job submission
2. Job status retrieval
3. Invalid store ID handling
4. Count mismatch validation
5. Invalid job ID handling

### Manual Testing

Use curl commands:

```bash
# Submit a job
curl -X POST http://localhost:8080/api/submit \
-H "Content-Type: application/json" \
-d '{
    "count": 2,
    "visits": [
        {
            "store_id": "RP00001",
            "image_url": [
                "https://www.gstatic.com/webp/gallery/2.jpg",
                "https://www.gstatic.com/webp/gallery/3.jpg"
            ],
            "visit_time": "2024-01-13T12:00:00Z"
        },
        {
            "store_id": "RP00002",
            "image_url": [
                "https://www.gstatic.com/webp/gallery/3.jpg"
            ],
            "visit_time": "2024-01-13T12:30:00Z"
        }
    ]
}'

# Check job status (replace with actual job_id)
curl "http://localhost:8080/api/status?jobid=job_1234567890"
```

## Development Environment

- **Operating System**: macOS/Linux
- **Language**: Go 1.21
- **IDE**: VSCode with Go extension
- **Dependencies**:
  - gorilla/mux: HTTP router
  - rs/cors: CORS middleware
  - Docker & Docker Compose
  - jq: JSON processor for testing

## Future Improvements

Given more time, the following improvements could be made:

1. **Infrastructure**:

   - Add database persistence for job status and results
   - Implement redis for job queue management
   - Set up CI/CD pipeline

2. **Features**:

   - Add authentication/authorization
   - Implement rate limiting
   - Add API versioning
   - Add metrics collection

3. **Code Quality**:
   - Add comprehensive unit tests
   - Add integration tests
   - Implement input validation middleware
   - Add API documentation (Swagger/OpenAPI)
