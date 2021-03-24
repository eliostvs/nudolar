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

 Total       = R$61,21

 $10.00      = R$57,54
 IOF (6.38%) = R$3,67

 Conversão baseada na cotação de 24/03/2021 13:06:
 $1.00       = R$5,75
```