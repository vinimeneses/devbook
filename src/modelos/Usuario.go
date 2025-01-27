package modelos

import (
	"api/src/seguranca"
	"fmt"
	"github.com/badoux/checkmail"
	"strings"
	"time"
)

type Usuario struct {
	ID       uint64    `json:"id,omitempty"`
	Nome     string    `json:"nome,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Email    string    `json:"email,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	CriadoEm time.Time `json:"CriadoEm,omitempty"`
}

func (usuario *Usuario) Preparar(etapa string) error {
	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}
	return nil
}

func (usuario *Usuario) validar(etapa string) error {
	campos := map[string]*string{
		"nome":  &usuario.Nome,
		"nick":  &usuario.Nick,
		"email": &usuario.Email,
		"senha": &usuario.Senha,
	}

	for campo, valor := range campos {
		if *valor == "" && (etapa == "cadastro" || (etapa == "edicao" && campo != "senha")) {
			return fmt.Errorf("O campo %s é obrigatório e não pode estar em branco", campo)
		}
	}
	if err := checkmail.ValidateFormat(usuario.Email); err != nil {
		return fmt.Errorf("o campo email está em um formato inválido")
	}

	return nil
}

func (usuario *Usuario) formatar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}

		usuario.Senha = string(senhaComHash)
	}
	return nil
}
