# CodeRelay API Documentation

Base URL: `http://localhost:8080`

---

## Health Check

### GET /healthz

Check if the API server is running.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-12-21T22:30:00Z"
}
```

---

## Problems

### GET /api/problems

List all available problems.

**Response:**
```json
[
  {
    "id": 1,
    "title": "Two Sum",
    "description": "## Problem\n\nGiven an array...",
    "time_limit_ms": 2000,
    "memory_limit_mb": 256,
    "created_at": "2024-12-21T22:30:00Z"
  }
]
```

---

### GET /api/problems/{id}

Get a specific problem with sample test cases.

**Parameters:**
- `id` (path) - Problem ID

**Response:**
```json
{
  "id": 1,
  "title": "Two Sum",
  "description": "## Problem\n\nGiven an array...",
  "time_limit_ms": 2000,
  "memory_limit_mb": 256,
  "created_at": "2024-12-21T22:30:00Z",
  "sample_cases": [
    {
      "id": 1,
      "problem_id": 1,
      "input": "2 7 11 15\n9",
      "expected_output": "0 1",
      "is_sample": true
    }
  ]
}
```

**Errors:**
- `404 Not Found` - Problem does not exist

---

## Submissions

### POST /api/submissions

Submit code for a problem.

**Request Body:**
```json
{
  "user_id": 1,
  "problem_id": 1,
  "code": "nums = list(map(int, input().split()))\ntarget = int(input())\n...",
  "language": "python"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "user_id": 1,
  "problem_id": 1,
  "code": "...",
  "language": "python",
  "status": "queued",
  "created_at": "2024-12-21T22:30:00Z"
}
```

**Errors:**
- `400 Bad Request` - Missing required fields

---

### GET /api/submissions/{id}

Get submission status.

**Parameters:**
- `id` (path) - Submission ID

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "problem_id": 1,
  "code": "...",
  "language": "python",
  "status": "AC",
  "runtime_ms": 150,
  "created_at": "2024-12-21T22:30:00Z"
}
```

**Status Values:**
- `queued` - Waiting to be processed
- `running` - Currently being evaluated
- `AC` - Accepted (all tests passed)
- `WA` - Wrong Answer
- `TLE` - Time Limit Exceeded

**Errors:**
- `404 Not Found` - Submission does not exist

---

## Running the API

```bash
# Start the server
make backend-run

# Or with Docker
make compose-up
```

The API will be available at `http://localhost:8080`.
