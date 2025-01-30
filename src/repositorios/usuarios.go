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
		"insert into usuarios (nome, nick, email, senha) values(?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			log.Fatal(erro)
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

// Buscar traz todos os usuários que atendem um filtro de nome ou nick
func (repositorio Usuarios) Buscar(nomeOuNick string) ([]modelos.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) // %nomeOuNick%

	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?",
		nomeOuNick, nomeOuNick,
	)

	if erro != nil {
		return nil, erro
	}
	defer func(linhas *sql.Rows) {
		err := linhas.Close()
		if err != nil {
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
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// BuscarPorID traz um usuário do banco de dados
func (repositorio Usuarios) BuscarPorID(ID uint64) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"select id, nome, nick, email, criadoEm from usuarios where id = ?",
		ID,
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer func(linhas *sql.Rows) {
		err := linhas.Close()
		if err != nil {
			log.Fatal(erro)
		}
	}(linhas)

	var usuario modelos.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Atualizar altera as informações de um usuário no banco de dados
func (repositorio Usuarios) Atualizar(ID uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"update usuarios set nome = ?, nick = ?, email = ? where id = ?",
	)
	if erro != nil {
		return erro
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			log.Fatal(erro)
		}
	}(statement)

	if _, erro = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID); erro != nil {
		return erro
	}

	return nil
}

// Deletar exclui as informações de um usuário no banco de dados
func (repositorio Usuarios) Deletar(ID uint64) error {
	statement, erro := repositorio.db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		return erro
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			log.Fatal(erro)
		}
	}(statement)

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio Usuarios) BuscarPorEmail(email string) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query("select id, senha from usuarios where email = ?", email)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer func(linha *sql.Rows) {
		err := linha.Close()
		if err != nil {
			log.Fatal(erro)
		}
	}(linha)

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil

}

func (repositorio Usuarios) Seguir(usuarioID uint64, seguidorID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"insert ignore into seguidores (usuario_id, seguidor_id) values (?, ?)",
	)
	if erro != nil {
		return erro
	}

	if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
		return erro
	}

	return nil
}

func (repositorio Usuarios) PararDeSeguir(usuarioID, SeguidorID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM seguidores WHERE usuario_id = ? AND seguidor_id = ?",
	)
	if erro != nil {
		return erro
	}

	if _, erro = statement.Exec(usuarioID, SeguidorID); erro != nil {
		return erro
	}
	return nil
}

func (repositorio Usuarios) BuscarSeguidores(id uint64) ([]modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
	SELECT 
	    u.id,
	    u.nome,
	    u.nick,
	    u.email,
	    u.criadoEm
	FROM 
	    usuarios u
	INNER JOIN 
		seguidores s
	ON
		u.id = s.seguidor_id
	WHERE
	    s.usuario_id = ?
	`, id)
	if erro != nil {
		return nil, erro
	}
	defer func(linhas *sql.Rows) {
		err := linhas.Close()
		if err != nil {
			log.Fatal(erro)
		}
	}(linhas)

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro := linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}
	return usuarios, nil
}

func (repositorio Usuarios) UsuarioSeguindo(usuarioID uint64) ([]modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
	SELECT
	    u.id,
	    u.nome,
	    u.nick,
	    u.email,
	    u.criadoEm
	FROM
	    usuarios u
	INNER JOIN 
		seguidores s
	ON
		u.id = s.usuario_id
	WHERE
	    seguidor_id = ?
`, usuarioID)
	if erro != nil {
		return nil, erro
	}

	defer func(linhas *sql.Rows) {
		err := linhas.Close()
		if err != nil {
			log.Fatal(erro)
		}
	}(linhas)

	var usuarios []modelos.Usuario
	for linhas.Next() {
		var usuario modelos.Usuario

		if erro := linhas.Scan(
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

func (repositorio Usuarios) BuscarSenha(id uint64) (string, error) {
	linha, erro := repositorio.db.Query("select senha from usuarios where id = ?", id)
	if erro != nil {
		return "", erro
	}
	defer func(linha *sql.Rows) {
		erro := linha.Close()
		if erro != nil {
			log.Fatal(erro)
		}
	}(linha)

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.Senha); erro != nil {
			return "", erro
		}
	}

	return usuario.Senha, nil
}

func (repositorio Usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statement, erro := repositorio.db.Prepare(`
	UPDATE 
	usuarios 
	SET 
	senha = ?
	WHERE
	id = ?
`)
	if erro != nil {
		return erro
	}

	_, erro = statement.Exec(senha, usuarioID)
	if erro != nil {
		return erro
	}
	return nil
}
