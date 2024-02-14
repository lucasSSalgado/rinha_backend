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

	if transaction.Tipo != "c" && transaction.Tipo != "d" {
		return 0, 0, errors.New("invalid op")
	}

	limite, saldo, err := cs.repo.AddTransaction(idInt, transaction)
	if err != nil {
		return 0, 0, err
	}

	return limite, saldo, nil
}

func (cs *ClientService) GetHistorico(id string) (models.Historico, error) {
	idUint, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return models.Historico{}, nil
	}

	saldo, err := cs.repo.GetSaldo(uint16(idUint))
	if err != nil {
		return models.Historico{}, nil
	}

	transacoes, err := cs.repo.GetTransacoes(uint16(idUint))
	if err != nil {
		return models.Historico{}, nil
	}

	return models.Historico{
		Saldo:      saldo,
		Transacoes: transacoes,
	}, nil
}
