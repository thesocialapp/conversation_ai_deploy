// models/models.go

type PythonAPIRequest struct {
	Text string `json:"text"` 
  }
  
  type PythonAPIResponse struct {
	Response string `json:"response"`
  }

  // Request and response envelopes

type RequestEnvelope struct {
	Request PythonAPIRequest `json:"request"`
  }
  
  type ResponseEnvelope struct {
	Response PythonAPIResponse `json:"response"`
  }
  
  // Standarized error structure
  
  type APIError struct {
	Code int    `json:"code"`
	Message string `json:"message"` 
  }
  
  // Generic API response
  
  type APIResponse struct {
	Data interface{} `json:"data,omitempty"`
	Error APIError   `json:"error,omitempty"`
  }
  