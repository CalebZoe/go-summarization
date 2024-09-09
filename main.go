package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

type SummarizeRequest struct {
    Text string `json:"text"`
}

type SummarizeResponse struct {
    Summary string `json:"summary"`
}

func main() {
    r := gin.Default()
    r.POST("/summarize", summarizeHandler)
    r.Run()
}

func summarizeHandler(c *gin.Context) {
    var req SummarizeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    summary, err := summarizeText(req.Text)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, SummarizeResponse{
        Summary: summary,
    })
}

func summarizeText(text string) (string, error) {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
    }

    requestBody, err := json.Marshal(map[string]string{
        "model":  "text-davinci-002",
        "prompt": fmt.Sprintf("Summarize the following text:\n\n%s", text),
        "max_tokens": "250",
        "temperature": "0.5",
    })
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var respData map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
        return "", err
    }

    choices, ok := respData["choices"].([]interface{})
    if !ok || len(choices) == 0 {
        return "", fmt.Errorf("no summary generated")
    }

    choiceData, ok := choices[0].(map[string]interface{})
    if !ok {
        return "", fmt.Errorf("invalid response format")
    }

    summary, ok := choiceData["text"].(string)
    if !ok {
        return "", fmt.Errorf("summary not found in response")
    }

    return summary, nil
}
