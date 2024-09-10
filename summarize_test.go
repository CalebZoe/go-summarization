// summarize_test.go

package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "os"
    "reflect"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestSummarizeHandler(t *testing.T) {
    // Create a test server
    r := gin.Default()
    r.POST("/summarize", summarizeHandler)

    // Create a test request
    req, err := http.NewRequest("POST", "/summarize", bytes.NewBuffer([]byte(`{"text": "This is a test text"}`)))
    if err!= nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Create a test response recorder
    w := httptest.NewRecorder()

    // Call the handler
    r.ServeHTTP(w, req)

    // Check the response status code
    if w.Code!= http.StatusOK {
        t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
    }

    // Check the response body
    var resp SummarizeResponse
    err = json.Unmarshal(w.Body.Bytes(), &resp)
    if err!= nil {
        t.Fatal(err)
    }

    // Check the summary
    if resp.Summary == "" {
        t.Errorf("expected a non-empty summary but got an empty string")
    }
}

func TestSummarizeText(t *testing.T) {
    // Set the OPENAI_API_KEY environment variable
    os.Setenv("OPENAI_API_KEY", "test-api-key")

    // Test with a valid text
    text := "This is a test text"
    summary, err := summarizeText(text)
    if err!= nil {
        t.Fatal(err)
    }

    if summary == "" {
        t.Errorf("expected a non-empty summary but got an empty string")
    }

    // Test with an empty text
    text = ""
    summary, err = summarizeText(text)
    if err!= nil {
        t.Fatal(err)
    }

    if summary!= "" {
        t.Errorf("expected an empty summary but got a non-empty string")
    }

    // Test with a nil OPENAI_API_KEY environment variable
    os.Unsetenv("OPENAI_API_KEY")
    text = "This is a test text"
    summary, err = summarizeText(text)
    if err == nil {
        t.Errorf("expected an error but got nil")
    }

    if summary!= "" {
        t.Errorf("expected an empty summary but got a non-empty string")
    }
}

func TestSummarizeText_MockAPIResponse(t *testing.T) {
    // Set the OPENAI_API_KEY environment variable
    os.Setenv("OPENAI_API_KEY", "test-api-key")

    // Create a test server to mock the OpenAI API
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, `{"choices": [{"text": "This is a test summary"}]}`)
    }))
    defer ts.Close()

    // Set the API endpoint URL to the test server URL
    originalAPIEndpoint := "https://api.openai.com/v1/completions"
    apiEndpoint = ts.URL

    // Test with a valid text
    text := "This is a test text"
    summary, err := summarizeText(text)
    if err!= nil {
        t.Fatal(err)
    }

    if summary!= "This is a test summary" {
        t.Errorf("expected summary 'This is a test summary' but got '%s'", summary)
    }

    // Restore the original API endpoint URL
    apiEndpoint = originalAPIEndpoint
}
