package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EnderecoCompleto representa os dados do endereço retornados pelas APIs.
type EnderecoCompleto struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
}

const (
	cep            = "55296-530"
	brasilAPIURL   = "https://brasilapi.com.br/api/cep/v1/"
	viaCEPURL      = "http://viacep.com.br/ws/"
	requestTimeout = 1 * time.Second
)

func main() {
	addressChannel := make(chan string)

	go fetchFromAPI("BrasilAPI", brasilAPIURL+cep, addressChannel)
	go fetchFromAPI("ViaCEP", fmt.Sprintf("%s%s/json/", viaCEPURL, cep), addressChannel)

	select {
	case result := <-addressChannel:
		fmt.Println(result)
	case <-time.After(requestTimeout):
		fmt.Println("Timeout: nenhuma resposta recebida em 1 segundo.")
	}
}

// Realiza a requisição para uma API e envia o resultado para o canal.
func fetchFromAPI(apiName, url string, resultChannel chan<- string) {
	address, err := getAddressFromAPI(url)
	if err != nil {
		resultChannel <- fmt.Sprintf("Erro ao chamar %s: %v", apiName, err)
		return
	}
	resultChannel <- formatAddress(apiName, address)
}

// Faz a requisição HTTP e decodifica a resposta JSON.
func getAddressFromAPI(url string) (EnderecoCompleto, error) {
	var address EnderecoCompleto

	resp, err := http.Get(url)
	if err != nil {
		return address, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return address, fmt.Errorf("Retorno inválido: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return address, err
	}

	return address, nil
}

// Formata o resultado em uma string para exibição.
func formatAddress(apiName string, address EnderecoCompleto) string {
	return fmt.Sprintf("Resposta de %s: %+v", apiName, address)
}
