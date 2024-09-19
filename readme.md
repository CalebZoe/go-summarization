# Summarize API - README

## Overview

The `Summarize API` is a simple HTTP-based microservice built using the Go programming language and the Gin web framework. It takes a block of text as input and returns a summarized version of the text by calling OpenAI's GPT-based API (e.g., `text-davinci-002`).

This API is designed for use in environments where automatic summarization of text is required. The service accepts JSON input and returns JSON output.

## Prerequisites

Before running the application, ensure you have the following installed on your machine:

- [Go](https://golang.org/dl/) (version 1.16 or higher)
- [Gin web framework](https://github.com/gin-gonic/gin)
- An [OpenAI API key](https://beta.openai.com/signup/)

## Setup

### 1. Clone the Repository


### 2. Install Dependencies

You will need to install the dependencies specified in the code. You can do this using `go get`:

```bash
go get github.com/gin-gonic/gin
```

### 3. Set Environment Variables

Ensure you set your OpenAI API key as an environment variable, as the service requires it to make requests to the OpenAI API.

```bash
export OPENAI_API_KEY=your_openai_api_key
```

Alternatively, you can add this to your `.bashrc`, `.zshrc`, or other shell configuration files to make it permanent.

### 4. Run the Application

Run the Go application using the following command:

```bash
go run main.go
```

By default, the application runs on `http://localhost:8080`.

## API Usage

### Endpoint: `/summarize` (POST)

This endpoint accepts a JSON payload with a block of text and returns a summarized version.

#### Request

- **Method**: `POST`
- **URL**: `http://localhost:8080/summarize`
- **Headers**: 
  - `Content-Type: application/json`
- **Body**:
  ```json
  {
      "text": "Your text to be summarized here."
  }
  ```

#### Response

- **Status Code**: `200 OK`
- **Content-Type**: `application/json`
- **Body**:
  ```json
  {
      "summary": "This is a summarized version of the text."
  }
  ```

#### Error Responses

- **400 Bad Request**: Returned if the input JSON is invalid.
  ```json
  {
      "error": "Invalid JSON format or missing 'text' field."
  }
  ```
  
- **500 Internal Server Error**: Returned if there is an error interacting with the OpenAI API.
  ```json
  {
      "error": "OPENAI_API_KEY environment variable is not set"
  }
  ```

## Code Structure

### `main.go`

- **`SummarizeRequest`**: This struct defines the expected input payload. It contains a `text` field, which holds the text to be summarized.
  
- **`SummarizeResponse`**: This struct defines the output payload, which contains the summarized text under the `summary` field.

- **`main()`**: Initializes the Gin server and sets up the `/summarize` route.

- **`summarizeHandler()`**: The request handler for the `/summarize` endpoint. It handles parsing the input JSON, calling the `summarizeText()` function, and returning the response.

- **`summarizeText()`**: This function interacts with the OpenAI API to generate a summary for the provided text. It uses an environment variable (`OPENAI_API_KEY`) to authenticate with OpenAI.

### OpenAI API Interaction

- The API request to OpenAI includes:
  - Model: `"text-davinci-002"`
  - Prompt: A custom text prompt requesting a summary of the provided text.
  - Max Tokens: 250 tokens to limit the length of the summary.
  - Temperature: A temperature of `0.5` for balanced creativity.

## Testing the API

You can test the API using tools like `curl` or [Postman](https://www.postman.com/).

### Example using `curl`:

```bash
curl -X POST http://localhost:8080/summarize \
     -H "Content-Type: application/json" \
     -d '{"text": "Insert a long piece of text here."}'
```

1. **Invalid Input**: If the input payload is not valid JSON or the required `text` field is missing, a `400 Bad Request` response is returned.
2. **Missing API Key**: If the `OPENAI_API_KEY` environment variable is not
