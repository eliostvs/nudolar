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

func CLI(args []string) int {
	var app appEnv

	if err := app.fromArgs(args); err != nil {
		return 2
	}

	if err := app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro inesperado: %v\n", err)
		return 1
	}

	return 0
}

type appEnv struct {
	price   float64
	timeout time.Duration
	hc      *http.Client
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("nudolar", flag.ContinueOnError)
	fl.Usage = func() {
		fmt.Fprint(fl.Output(), "\nSimule o valor de uma compra internacional na sua fatura do Nubanco\n\n")
		fmt.Fprint(fl.Output(), "Uso:\n  nudolar [opcoes] <dolar>\n\n")
		fmt.Fprint(fl.Output(), "Exemplo:\n  nudolar 10.66\n\n")
		fmt.Fprintln(fl.Output(), "Opções:")
		fl.PrintDefaults()
	}
	fl.DurationVar(&app.timeout, "t", 10*time.Second, "client timeout")

	if err := fl.Parse(args); err != nil {
		return err
	}

	if err := app.parsePrice(fl.Args()); err != nil {
		fmt.Fprintf(fl.Output(), "%v\n", err)
		fl.Usage()
		return err
	}

	app.hc = http.DefaultClient

	return nil
}

func (app *appEnv) run() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.timeout)
	defer cancel()

	if ptax, err := NewPTAX(ctx, app.hc); err != nil {
		return err
	} else {
		return app.show(ptax, calcPrices(app.price, ptax.SellingRate))
	}
}

func (app *appEnv) parsePrice(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("valor não informado")
	}

	price, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return fmt.Errorf("valor inválido: %s", args[0])
	}

	if price <= 0.0 {
		return errors.New("valor deve ser maior que zero")
	}

	app.price = price

	return nil
}

func (app *appEnv) show(ptax PTAX, prices Prices) error {
	template := `
 Total       = %s

 %-9v   = %s
 IOF (%.2f%%) = %s

 Conversão baseada na cotação de %s:
 %-11v = %s

 PTAX        = %s
 Spread (4%%) = %s

`
	_, err := fmt.Printf(
		template,
		displayReal(prices.total),
		displayDollar(prices.price),
		displayReal(prices.subtotal),
		IOF,
		displayReal(prices.iof),
		ptax.Timestamp.Format("02/01/2006 15:04"),
		displayDollar(1),
		displayReal(prices.dollar),
		displayReal(ptax.SellingRate),
		displayReal(prices.spread),
	)
	return err
}

func calcPrices(price float64, dollar float64) Prices {
	spread := (dollar / 100) * Spread
	subtotal := (spread + dollar) * price
	iof := (subtotal / 100) * IOF
	return Prices{
		spread:   spread,
		price:    price,
		dollar:   dollar + spread,
		subtotal: subtotal,
		iof:      iof,
		total:    subtotal + iof,
	}
}

type Prices struct {
	price, spread, dollar, subtotal, iof, total float64
}

func displayReal(price float64) string {
	return money.New(int64(price*100), "BRL").Display()
}

func displayDollar(price float64) string {
	return money.New(int64(price*100), "USD").Display()
}
