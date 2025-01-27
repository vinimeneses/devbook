package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Erro ao fechar o banco de dados", err)
		}
	}(db)

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		log.Fatal(erro)
	}

	respostas.JSON(w, http.StatusCreated, usuario)
}
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}
	defer func(db *sql.DB) {
		if erro := db.Close(); erro != nil {
			log.Fatal("Erro ao fechar o banco de dados", erro)
		}
	}(db)

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOuNick)

	respostas.JSON(w, http.StatusOK, usuarios)
}
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.Atoi(parametros["usuarioId"])
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer func(db *sql.DB) {
		if erro = db.Close(); erro != nil {
			log.Fatal(erro)
		}
	}(db)

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}

	respostas.JSON(w, http.StatusOK, usuario)
}
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametro := mux.Vars(r)

	usuarioId, erro := strconv.Atoi(parametro["usuarioId"])
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
	}

	if erro := usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}
	defer func(db *sql.DB) {
		if erro := db.Close(); erro != nil {
			log.Fatal(erro)
		}
	}(db)

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	erro = repositorio.Atualizar(usuarioId, usuario)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.Atoi(parametros["usuarioId"])
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
