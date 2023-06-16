package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lazarosanson/challenge-multithreading-goExpert/internal/dto"
	"github.com/lazarosanson/challenge-multithreading-goExpert/internal/utils"
)

type CepHandler struct {
}

func NewCepHandler() *CepHandler {
	return &CepHandler{}
}

func (cepHandler *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	cep, err := utils.CepMask(cep)
	if err != nil {
		errorHandler(w, http.StatusBadRequest, err)
		return
	}

	apiCepURL := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	viaCepURL := "http://viacep.com.br/ws/" + cep + "/json/"
	chResponse := make(chan *http.Response)
	chError := make(chan error)

	go func() {
		resp, err := http.Get(apiCepURL)
		if err != nil {
			chError <- err
			return
		}
		chResponse <- resp
	}()

	go func() {
		resp, err := http.Get(viaCepURL)
		if err != nil {
			chError <- err
			return
		}
		chResponse <- resp
	}()

	var result interface{}

	select {
	case response := <-chResponse:
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {

			responseBytes, err := io.ReadAll(response.Body)
			if err != nil {
				errorHandler(w, http.StatusInternalServerError, err)
				return
			}
			err = json.Unmarshal(responseBytes, &result)
			if err != nil {
				errorHandler(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			result = dto.CepError{
				StatusCode: response.StatusCode,
			}
		}

		var cepResponse dto.CepResponseDTO
		if response.Request.URL.String() == apiCepURL {
			cepResponse = dto.CepResponseDTO{
				Api:  "ApiCep",
				Data: result,
			}
		} else {
			cepResponse = dto.CepResponseDTO{
				Api:  "ViaCep",
				Data: result,
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cepResponse)
		return

	case err := <-chError:
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return

	case <-time.After(time.Second):
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}
}

func errorHandler(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(
		dto.CepError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
}
