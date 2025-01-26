package controllers

import (
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		log.Fatal(erro)
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		log.Fatal(erro)
	}

	db, erro := banco.Conectar()
	if erro != nil {
		log.Fatal(erro)
	}

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioID, erro := repositorio.Criar(usuario)
	if erro != nil {
		log.Fatal(erro)
	}

	WriteResponse(w, fmt.Sprintf("Id inserido %d", usuarioID))
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
