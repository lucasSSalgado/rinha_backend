package repository

import (
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

func (c *ClientRepository) Debitar(id uint16, value int64) error {
	_, err := c.db.Exec("UPDATE cliente SET saldo = $1 WHERE user_id = $2", value, id)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientRepository) Creditar(id uint16, value int64) error {
	_, err := c.db.Exec("UPDATE cliente SET saldo = $1 WHERE user_id = $2", value, id)
	if err != nil {
		return err
	}
	return nil
}
