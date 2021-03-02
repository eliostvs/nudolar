# NUDOLAR

Simule o valor de uma compra internacional na sua fatura do Nubanco.

## Instalação

`go get github.com/eliostvs/nudolar`

## Uso

```
Uso: nudolar [OPCOES] [DOLAR]
  Simule o valor de uma compra internacional na sua fatura do Nubanco

Exemplo:
  nudolar 10.66

Opções:
  -t duration
        client timeout (default 10s)
```

### Exemplo

```
nudolar 10.66
Cotação do Dólar do dia 03/02/2021 13:31: R$5,68
Subtotal: R$60,59
Spread 4%: R$2,42
IOF 6.38%: R$4,02
Total: R$67,04
```