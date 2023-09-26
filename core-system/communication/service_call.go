package communication

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type YourDataType struct {
    // Define the fields of the data type you're sending to Python here
}

type YourResponseType struct {
    // Define the fields of the response type you're receiving from Python here
}

func CallPythonService(data YourDataType) (*YourResponseType, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal data: %v", err)
    }

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Post("http://python_service_address/process-text", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to call Python service: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Python service returned non-200 status code: %d", resp.StatusCode)
    }

    var response YourResponseType
    err = json.NewDecoder(resp.Body).Decode(&response)
    if err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    return &response, nil
}
