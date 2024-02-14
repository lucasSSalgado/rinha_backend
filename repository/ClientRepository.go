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

func (c *ClientRepository) AddTransaction(id uint64, transaction *models.TransacaoRequDto) (uint64, int64, error) {
	tx, err := c.db.Begin(context.Background())
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(context.Background())

	var limite uint64
	var saldo int64
	err = tx.QueryRow(
		context.Background(),
		"SELECT limite, saldo FROM cliente WHERE user_id = $1 FOR UPDATE",
		id,
	).Scan(&limite, &saldo)
	if err != nil {
		return 0, 0, err
	}

	var newSaldo int64
	if transaction.Tipo == "c" {
		newSaldo = int64(transaction.Valor) + saldo
	} else {
		newSaldo = saldo - int64(transaction.Valor)
	}

	if (newSaldo + int64(limite)) < 0 {
		return 0, 0, errors.New("422")
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
		return 0, 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, 0, err
	}
	return limite, newSaldo, nil
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

	for rows.Next() {
		var valor int
		var tipo string
		var descricao string
		var realizada_em time.Time

		err = rows.Scan(&valor, &tipo, &descricao, &realizada_em)
		if err != nil {
			return nil, err
		}
		transacao := models.UltTransacoes{
			Valor:       valor,
			Tipo:        tipo,
			Descricao:   descricao,
			RealizadoEm: realizada_em,
		}

		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}

func (c *ClientRepository) GetSaldo(id uint16) (models.Saldo, error) {
	var saldo int64
	var limite uint64
	err := c.db.QueryRow(context.Background(), "SELECT saldo, limite FROM cliente WHERE user_id = $1", id).Scan(&saldo, &limite)

	if err != nil {
		return models.Saldo{}, err
	}

	now := time.Now().UTC()

	return models.Saldo{
		Total:       saldo,
		DataExtrato: now,
		Limite:      int64(limite),
	}, nil
}
