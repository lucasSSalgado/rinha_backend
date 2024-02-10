package repository

import (
	"rinha/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (c *ClientRepository) CheckIfClientExist(id uint16) error {
	var exists uint16

	rows := c.db.QueryRow("Select user_id FROM cliente WHERE user_id = $1", id)
	err := rows.Scan(&exists)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientRepository) GetLimitAndSaldoByClientId(id uint16, value uint32) (int64, uint64, error) {
	var saldo int64
	var limite uint64
	err := c.db.QueryRow("SELECT limite, saldo FROM cliente WHERE user_id = $1", id).Scan(&limite, &saldo)
	if err != nil {
		return 0, 0, err
	}

	return saldo, limite, nil
}

func (c *ClientRepository) Debitar(id uint16, newValue int64, debito int64, descricao string) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE cliente SET saldo = $1 WHERE user_id = $2", newValue, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO transacoes (user_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)",
		id, debito, "d", debito, time.Now(), 
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err!=nil {
		return err
	}
	return nil
}

func (c *ClientRepository) Creditar(id uint16, newValue int64, credito int64, descricao string) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE cliente SET saldo = $1 WHERE user_id = $2", newValue, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO transacoes (user_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)", 
		id, credito, "c", descricao, time.Now(), 		
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientRepository) GetTransacoes(id uint16) ([]models.UltTransacoes, error) {
	var transacoes []models.UltTransacoes

	rows, err := c.db.Queryx("SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE user_id = $1 ORDER BY realizada_em LIMIT 10", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var transacao models.UltTransacoes
		err = rows.StructScan(&transacao)
		if err != nil {
			return nil, err
		}

		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}

func (c *ClientRepository) GetSaldo(id uint16) (models.Saldo, error) {
	var saldo int64
	var limite uint64
	err := c.db.QueryRow("SELECT saldo, limite FROM cliente WHERE user_id = $1", id).Scan(&saldo, &limite)
	if err != nil {
		return models.Saldo{}, err
	}

	now := time.Now()

	return models.Saldo{
		Total:       saldo,
		DataExtrato: now,
		Limite:      int64(limite),
	}, nil
}
