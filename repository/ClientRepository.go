package repository

import (
	"context"
	"errors"
	"rinha/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (c *ClientRepository) CheckIfClientExist(id uint16) error {
	var exists uint16

	rows := c.db.QueryRow(
		context.Background(),
		"SELECT user_id FROM cliente WHERE user_id = $1",
		id,
	)
	err := rows.Scan(&exists)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientRepository) AddTransaction(id uint64, transaction *models.TransacaoRequDto) (*models.TransacaoRespDto, error) {
	tx, err := c.db.Begin(context.Background())
	if err != nil {
		return &models.TransacaoRespDto{}, err
	}
	defer tx.Rollback(context.Background())

	var resp models.TransacaoRespDto
	err = tx.QueryRow(
		context.Background(),
		"SELECT limite, saldo FROM cliente WHERE user_id = $1 FOR UPDATE",
		id,
	).Scan(&resp.Limite, &resp.Saldo)
	if err != nil {
		return &models.TransacaoRespDto{}, err
	}

	var newSaldo int64
	if transaction.Tipo == "c" {
		newSaldo = int64(transaction.Valor) + resp.Saldo
	} else {
		newSaldo = resp.Saldo - int64(transaction.Valor)
	}

	if (newSaldo + int64(resp.Limite)) < 0 {
		return &models.TransacaoRespDto{}, errors.New("422")
	}

	batch := &pgx.Batch{}
	batch.Queue(
		"UPDATE cliente SET saldo = $1 WHERE user_id = $2",
		newSaldo, id,
	)
	batch.Queue(
		"INSERT INTO transacoes (user_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)",
		id, transaction.Valor, transaction.Tipo, transaction.Descricao, time.Now().UTC(),
	)

	s := tx.SendBatch(
		context.Background(),
		batch,
	)
	if err := s.Close(); err != nil {
		return &models.TransacaoRespDto{}, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return &models.TransacaoRespDto{}, err
	}
	return &models.TransacaoRespDto{Limite: resp.Limite, Saldo: resp.Saldo}, nil
}

func (c *ClientRepository) GetTransacoes(id uint16) ([]models.UltTransacoes, error) {
	var transacoes []models.UltTransacoes

	rows, err := c.db.Query(
		context.Background(),
		"SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE user_id = $1 ORDER BY realizada_em DESC LIMIT 10",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transacoes = make([]models.UltTransacoes, 0, 10)

	for rows.Next() {
		var transacao models.UltTransacoes
		err := rows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadoEm)
		if err != nil {
			return nil, err
		}
		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}

func (c *ClientRepository) GetSaldo(id uint16) (models.Saldo, error) {
	var resp models.Saldo
	err := c.db.QueryRow(
		context.Background(),
		"SELECT saldo, limite FROM cliente WHERE user_id = $1",
		id).Scan(&resp.Total, &resp.Limite)

	if err != nil {
		return models.Saldo{}, err
	}

	resp.DataExtrato = time.Now().UTC()

	return resp, nil
}
