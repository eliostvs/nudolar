package nudolar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const baseUrl = "https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata/CotacaoDolarDia(dataCotacao=@dataCotacao)?@dataCotacao='%s'&$top=1&$format=json&$select=cotacaoVenda,dataHoraCotacao"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	t.Time, err = time.Parse("2006-01-02 15:04:05", s)
	return err
}

func (t Time) String() string {
	return t.Format("02/01/2006 15:04")
}

type Quote struct {
	Amount float64 `json:"cotacaoVenda"`
	Time   Time    `json:"dataHoraCotacao"`
}

func GetQuote(ctx context.Context, hc *http.Client) (Quote, error) {
	var days int
	for {
		var res struct {
			Value []Quote `json:"value"`
		}

		if err := getJSON(ctx, hc, buildURL(days), &res); err != nil {
			return Quote{}, err
		}

		if len(res.Value) > 0 {
			return res.Value[0], nil
		}

		days++
	}
}

func buildURL(days int) string {
	date := time.Now().AddDate(0, 0, -days).Format("01-02-2006")
	return fmt.Sprintf(baseUrl, date)
}

func getJSON(ctx context.Context, hc *http.Client, url string, dst interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisção: %w", err)
	}

	res, err := hc.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao buscar url '%s': %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("código http inesperado: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(dst); err != nil {
		return fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return nil
}
