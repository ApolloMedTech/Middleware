package requestManager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PerformGETHTTPRequest(c *gin.Context, url string, requestBody interface{}, responseData interface{}) error {
	return performHTTPRequest(c, "GET", url, requestBody, responseData)
}

func PerformPOSTHTTPRequest(c *gin.Context, url string, requestBody interface{}, responseData interface{}) error {
	return performHTTPRequest(c, "POST", url, requestBody, responseData)
}

func BindRequestData(c *gin.Context, requestData interface{}) error {
	// Bind the request data
	if err := c.ShouldBind(requestData); err != nil {
		return err
	}
	return nil
} 

func performHTTPRequest(c *gin.Context, method, url string, requestBody interface{}, responseData interface{}) error {
	var body io.Reader

	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Microservice returned status: %d", resp.StatusCode)
		return fmt.Errorf("status not OK: %s", errMsg)
	}

	if responseData != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseData); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return nil
}
