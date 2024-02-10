package models

import "time"

type TransacaoRequDto struct {
	Valor     uint32 `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type TransacaoRespDto struct {
	Limite uint64 `json:"limite"`
	Saldo  int64  `json:"saldo"`
}

type Historico struct {
	Saldo      Saldo           `json:"saldo"`
	Transacoes []UltTransacoes `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       int64     `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int64     `json:"limite"`
}

type UltTransacoes struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadoEm time.Time `json:"realizada_em"`
}
