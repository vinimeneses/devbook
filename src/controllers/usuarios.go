package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

	if erro = usuario.Preparar(); erro != nil {
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
	WriteResponse(w, "Buscando usuários!")
}
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, "Buscando usuário!")
}
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, "Atualizando usuário!")
}
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, "Deletando usuário!")
}
func WriteResponse(w http.ResponseWriter, message string) {
	if _, err := w.Write([]byte(message)); err != nil {
		log.Fatalf("Erro ao fechar o write: %v", err)
	}
}
