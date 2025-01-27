package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
	"log"
)

type Usuarios struct {
	db *sql.DB
}

func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

func (repositorio Usuarios) Criar(usuario modelos.Usuario) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"insert into usuarios (nome, nick, email, senha) values (?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
		}
	}(statement)

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

func (repositorio Usuarios) Buscar(nomeOuNick string) ([]modelos.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) // %nomeOuNick%

	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?",
		nomeOuNick, nomeOuNick)
	if erro != nil {
		return nil, erro
	}
	defer func(linhas *sql.Rows) {
		if erro = linhas.Close(); erro != nil {
			log.Fatal(erro)
		}
	}(linhas)

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repositorio Usuarios) BuscarPorID(usuarioID int) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where id = ?",
		usuarioID)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer func(linhas *sql.Rows) {
		if erro := linhas.Close(); erro != nil {
			log.Fatal(erro)
		}
	}(linhas)

	var usuario modelos.Usuario
	for linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm); erro != nil {
			return modelos.Usuario{}, nil
		}
	}
	return usuario, nil
}

func (repositorio Usuarios) Atualizar(ID int, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"update usuarios set nome = ?, nick = ?, email = ? where id = ?",
	)
	if erro != nil {
		return erro
	}
	defer func(statement *sql.Stmt) {
		erro := statement.Close()
		if erro != nil {

		}
	}(statement)

	if _, erro = statement.Exec(usuario.Nome, usuario.ID, usuario.Email, ID); erro != nil {
		return erro
	}
	return nil
}

func (repositorio Usuarios) Deletar(ID int) error {
	statement, erro := repositorio.db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		return erro
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {

		}
	}(statement)

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}
	return nil
}
