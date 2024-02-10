package service

import (
	"errors"
	"rinha/models"
	"rinha/repository"
	"strconv"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{
		repo: repo,
	}
}

func (cs *ClientService) FindClientById(id string) error {
	idInt, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return err
	}

	if err := cs.repo.CheckIfClientExist(uint16(idInt)); err != nil {
		return err
	}
	return nil
}

func (cs *ClientService) LidarComTransacao(transaction *models.TransacaoRequDto, id string) (uint64, int64, error) {
	idInt, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return 0, 0, err
	}

	if transaction.Tipo == "d" {
		saldo, limite, err := cs.repo.GetLimitAndSaldoByClientId(uint16(idInt), transaction.Valor)
		if err != nil {
			return 0, 0, err
		}

		if (saldo - int64(transaction.Valor)) < (int64(limite))*-1 {
			return 0, 0, errors.New("422")
		}

		newValue := saldo - int64(transaction.Valor)
		err = cs.repo.Debitar(uint16(idInt), newValue)
		if err != nil {
			return 0, 0, err
		}
		return limite, newValue, nil
	}
	if transaction.Tipo == "c" {
		saldo, limite, err := cs.repo.GetLimitAndSaldoByClientId(uint16(idInt), transaction.Valor)
		if err != nil {
			return 0, 0, err
		}

		newValue := saldo + int64(transaction.Valor)
		err = cs.repo.Creditar(uint16(idInt), newValue)
		if err != nil {
			return 0, 0, err
		}

		return limite, newValue, nil
	}

	return 0, 0, errors.New("invalid type")
}
