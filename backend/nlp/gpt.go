package nlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func CallOpenAI(prompt string) (string, error) {
	// Define the API endpoint for OpenAI completion
	url := "https://api.openai.com/v1/chat/completions"

	// Create the request body with the necessary parameters
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo-0613", // Specify the model you want to use
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}

	// Create a new HTTP request to the OpenAI API
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	secretKey := os.Getenv("OPENAI_SECRET_KEY")
	// Set the required headers, including the Authorization header with your API key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secretKey)

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and return the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// You might want to unmarshal and process the response body based on your needs
	// For simplicity, we're returning the raw response body as a string here
	return string(body), nil
}

func IsMessageAffirmative(message string) bool {
	// Construct a prompt that asks the model to interpret the message
	prompt := fmt.Sprintf("Is the following message affirmative or negative?\n\nMessage: \"%s\"", message)

	// Call the OpenAI API with the constructed prompt
	response, err := CallOpenAI(prompt)
	if err != nil {
		log.Fatalf("Failed to call OpenAI API: %v", err)
	}

	// Assuming the response structure, parse the API response to get the text
	// This structure and parsing might need to be adjusted based on the actual response format
	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(response), &apiResponse); err != nil {
		log.Fatalf("Failed to parse OpenAI API response: %v", err)
	}

	// Analyze the response text to determine if it's affirmative
	// This is a basic check and might need adjustments based on your specific needs
	for _, choice := range apiResponse.Choices {
		text := choice.Message.Content
		if strings.Contains(text, "affirmative") {
			return true
		}
	}

	return false
}

type RatingResponse struct {
	RatingProvided bool `json:"ratingProvided"`
	Rating         int  `json:"rating"`
}

func AnalyzeUserRating(userResponse string) (bool, int, error) {
	prompt := fmt.Sprintf("User was asked to rate the iPhone 15 on a scale of 1-5. The user's response was: \"%s\". Did the user provide a rating and if so, what was the rating? Give me an answer in the following format only: ```\n{ratingProvided: boolean, rating: number}\n```", userResponse)

	apiResponse, err := CallOpenAI(prompt)
	if err != nil {
		return false, 0, err
	}

	// Assume the structure of apiResponse and parse it to extract the rating
	// This parsing depends on the specific format of the response from OpenAI
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal([]byte(apiResponse), &response); err != nil {
		return false, 0, err
	}

	// Example of how to parse the response content to determine if a rating was provided and what it is
	// This is a simplified example. You may need to adjust the parsing logic based on the actual response format and content
	for _, choice := range response.Choices {
		content := choice.Message.Content

		// Parse the JSON content into the RatingResponse struct
		var ratingResp RatingResponse
		err := json.Unmarshal([]byte(content), &ratingResp)
		if err != nil {
			log.Printf("Error parsing JSON from content: %v", err)
			// Handle error, e.g., by continuing to the next choice or returning an error
			continue
		}

		// Check if a rating was provided and if it's within the expected range
		if ratingResp.RatingProvided && ratingResp.Rating >= 1 && ratingResp.Rating <= 5 {
			log.Printf("Rating provided: %d", ratingResp.Rating)
			return true, ratingResp.Rating, nil
		}
	}

	return false, 0, nil
}
