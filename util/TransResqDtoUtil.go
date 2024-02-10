package util

import (
	"errors"
	"rinha/models"
)

func CheckFields(dto *models.TransacaoRequDto) error {
	if dto.Valor == 0 || dto.Tipo == "" || dto.Descricao == "" {
		return errors.New("")
	}

	if len(dto.Descricao) > 11 {
		return errors.New("")
	}

	return nil
}
