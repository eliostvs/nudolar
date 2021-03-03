package nudolar

import (
	"context"
	"errors"
	"flag"
	"fmt"
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
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("nudolar", flag.ContinueOnError)
	fl.Usage = func() {
		fmt.Fprint(fl.Output(), "Uso: nudolar [OPCOES] DOLAR\n  Simule o valor de uma compra internacional na sua fatura do Nubanco\n\n")
		fmt.Fprint(fl.Output(), "Exemplo: \n  nudolar 10.66\n\n")
		fmt.Fprintln(fl.Output(), "Opções:")
		fl.PrintDefaults()
	}
	fl.DurationVar(&app.timeout, "t", 10*time.Second, "client timeout")

	if err := fl.Parse(args); err != nil {
		return err
	}

	if err := app.parsePrice(fl.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fl.Usage()
		return err
	}

	return nil
}

func (app *appEnv) run() error {
	ctx, cancel := context.WithTimeout(context.Background(), app.timeout)
	defer cancel()

	ptax := NewPTAXClient()
	if quoting, err := ptax.Get(ctx); err != nil {
		return err
	} else {
		return app.show(quoting, calcPrices(app.price, quoting.SellingRate))
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

func (app *appEnv) show(ptax PTAX, p Prices) error {
	_, err := fmt.Printf("%s\n%s", ptax, p)
	return err
}

func calcPrices(price float64, sellingPrice float64) Prices {
	subtotal := price * sellingPrice
	spread := (subtotal / 100) * Spread
	subtotalPlusSpread := subtotal + spread
	iof := (subtotalPlusSpread / 100) * IOF
	return Prices{
		subtotal: subtotal,
		spread:   spread,
		iof:      iof,
		total:    subtotalPlusSpread + iof,
	}
}

type Prices struct {
	subtotal float64
	spread   float64
	iof      float64
	total    float64
}

func (p Prices) String() string {
	return fmt.Sprintf(
		"Subtotal: %s\nSpread %d%%: %s\nIOF %.2f%%: %s\nTotal: %s\n",
		formatPrice(p.subtotal),
		Spread,
		formatPrice(p.spread),
		IOF,
		formatPrice(p.iof),
		formatPrice(p.total))
}

func formatPrice(price float64) string {
	return money.New(int64(price*100), "BRL").Display()
}
