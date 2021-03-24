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

type ApiResponse struct {
	Value []PTAX `json:"value"`
}

func NewPTAX(ctx context.Context) (PTAX, error) {
	var days int
	hc := http.DefaultClient
	for {
		url := buildURL(days)
		var resp ApiResponse
		if err := getJSON(ctx, hc, url, &resp); err != nil {
			return PTAX{}, err
		}

		if len(resp.Value) > 0 {
			return resp.Value[0], nil
		}

		days++
	}
}

func buildURL(days int) string {
	date := time.Now().AddDate(0, 0, -days).Format("01-02-2006")
	return fmt.Sprintf(baseUrl, date)
}

func getJSON(ctx context.Context, hc *http.Client, url string, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(data)
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	t.Time, err = time.Parse("2006-01-02 15:04:05", s)
	return err
}

type PTAX struct {
	SellingRate float64 `json:"cotacaoVenda"`
	Timestamp   Time    `json:"dataHoraCotacao"`
}
