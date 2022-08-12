package nudolar

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Rhymond/go-money"
)

const IOF = 6.38
const Spread = 4

func NewReal(p float64) *money.Money {
	return money.New(int64(p*100), "BRL")
}

type Calculation struct {
	amount,
	spread,
	exchangeRate,
	subtotal,
	iof,
	total,
	quoteAmount *money.Money
	quoteTime Time
}

func (c Calculation) String() string {
	template := `
 Total       = %s

 %-9v   = %s
 IOF (%.2f%%) = %s

 Conversão baseada na cotação de %s:
 $1.00       = %s

 PTAX        = %s
 Spread (4%%) = %s

`
	return fmt.Sprintf(
		template,
		c.total.Display(),
		c.amount.Display(),
		c.subtotal.Display(),
		IOF,
		c.iof.Display(),
		c.quoteTime,
		c.exchangeRate.Display(),
		c.quoteAmount.Display(),
		c.spread.Display(),
	)
}

func CLI(args []string) int {
	var app App

	if err := app.ParseArgs(args); err != nil {
		return 1
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro inesperado: %v\n", err)
		return 2
	}

	return 0
}

type App struct {
	amount float64
	hc     *http.Client
}

func (a *App) ParseArgs(args []string) error {
	var timeout time.Duration

	fl := flag.NewFlagSet("nudolar", flag.ContinueOnError)
	fl.Usage = func() {
		fmt.Fprint(fl.Output(), "\nSimule o valor de uma compra internacional na sua fatura do Nubanco\n\n")
		fmt.Fprint(fl.Output(), "Uso:\n  nudolar [opcoes] <dolar>\n\n")
		fmt.Fprint(fl.Output(), "Exemplo:\n  nudolar 10.66\n\n")
		fmt.Fprintln(fl.Output(), "Opções:")
		fl.PrintDefaults()
	}
	fl.DurationVar(&timeout, "timeout", 10*time.Second, "client timeout")

	if err := fl.Parse(args); err != nil {
		return fmt.Errorf("erro: %w", err)
	}

	if err := a.ParsePrice(fl.Args()); err != nil {
		fmt.Fprintf(fl.Output(), "%v\n", err)
		fl.Usage()
		return err
	}

	a.hc = &http.Client{Timeout: timeout}

	return nil
}

func (a *App) ParsePrice(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("valor não informado")
	}

	value := args[0]

	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("valor inválido: %s", value)
	}

	if amount <= 0.0 {
		return errors.New("valor deve ser maior que zero")
	}

	a.amount = amount

	return nil
}

func (a *App) Run() error {
	quote, err := GetQuote(context.Background(), a.hc)
	if err != nil {
		return err
	}

	if err := a.Show(a.calculatePrice(a.amount, quote)); err != nil {
		return err
	}

	return nil
}

func (a *App) Show(calc Calculation) error {
	_, err := fmt.Printf("%s", calc)
	if err != nil {
		return fmt.Errorf("erro ao exibir a simulação: %w", err)
	}
	return nil
}

func (a *App) calculatePrice(amount float64, quote Quote) Calculation {
	spread := (quote.Amount / 100) * Spread
	exchangeRate := spread + quote.Amount
	subtotal := exchangeRate * amount
	iof := (subtotal / 100) * IOF

	return Calculation{
		amount:       money.New(int64(amount*100), "USD"),
		quoteAmount:  NewReal(quote.Amount),
		quoteTime:    quote.Time,
		spread:       NewReal(spread),
		exchangeRate: NewReal(exchangeRate),
		subtotal:     NewReal(subtotal),
		iof:          NewReal(iof),
		total:        NewReal(subtotal + iof),
	}
}
