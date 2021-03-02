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

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	t.Time, err = time.Parse("2006-02-01 15:04:05.00", s)
	return err
}

type PTAX struct {
	SellingRate float64 `json:"cotacaoVenda"`
	Timestamp   Time    `json:"dataHoraCotacao"`
}

func (p PTAX) String() string {
	return fmt.Sprintf(
		"Cotação do Dólar do dia %s: %s",
		p.Timestamp.Format("02/01/2006 15:05"),
		formatPrice(p.SellingRate),
	)
}

func NewPTAXClient() *PTAXClient {
	return &PTAXClient{http.DefaultClient}
}

type PTAXClient struct {
	hc *http.Client
}

func (c *PTAXClient) Get(ctx context.Context) (PTAX, error) {
	var days int
	for {
		url := c.buildURL(days)
		var resp ApiResponse
		if err := c.getJSON(ctx, url, &resp); err != nil {
			return PTAX{}, err
		}

		if len(resp.Value) > 0 {
			return resp.Value[0], nil
		}

		days++
	}
}

func (c *PTAXClient) buildURL(days int) string {
	date := time.Now().AddDate(0, 0, -days).Format("01-02-2006")
	return fmt.Sprintf(baseUrl, date)
}

func (c *PTAXClient) getJSON(ctx context.Context, url string, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(data)
}
