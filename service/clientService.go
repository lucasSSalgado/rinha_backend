package service

import (
	"errors"
	"rinha/models"
	"rinha/repository"
	"strconv"

	"github.com/patrickmn/go-cache"
)

type ClientService struct {
	repo  *repository.ClientRepository
	cache *cache.Cache
}

func NewClientService(repo *repository.ClientRepository, c *cache.Cache) *ClientService {
	return &ClientService{
		repo:  repo,
		cache: c,
	}
}

func (cs *ClientService) FindClientById(id string) error {
	val, exists := cs.cache.Get(id)
	if exists {
		if val.(bool) {
			return nil
		} else {
			return errors.New("")
		}
	}

	idInt, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return err
	}

	if err := cs.repo.CheckIfClientExist(uint16(idInt)); err != nil {
		cs.cache.Set(id, false, cache.NoExpiration)
		return err
	}

	cs.cache.Set(id, true, cache.NoExpiration)
	return nil
}

func (cs *ClientService) LidarComTransacao(transaction *models.TransacaoRequDto, id string) (*models.TransacaoRespDto, error) {
	idInt, err := strconv.ParseUint(id, 10, 16)
	if err != nil {
		return &models.TransacaoRespDto{}, err
	}

	if transaction.Tipo != "c" && transaction.Tipo != "d" {
		return &models.TransacaoRespDto{}, errors.New("invalid op")
	}

	resp, err := cs.repo.AddTransaction(idInt, transaction)
	if err != nil {
		return &models.TransacaoRespDto{}, err
	}

	return resp, nil
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
