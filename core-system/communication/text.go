package communication

import "net/http"

func ProcessText(text string) (string, error) {

  // Call Python service
  resp, err := http.Post("http://pythonservice/text", text)
  if err != nil {
    return "", err
  }

  // Handle response

  return resp.Body, nil
}