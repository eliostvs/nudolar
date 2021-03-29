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

 Total       = R$67,27

 $10.66      = R$63,23
 IOF (6.38%) = R$4,03

 Conversão baseada na cotação de 26/03/2021 13:02:
 $1.00       = R$5,93

 PTAX        = R$5,70
 Spread (4%) = R$0,23
```