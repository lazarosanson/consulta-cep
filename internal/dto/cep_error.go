package dto

type CepError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
