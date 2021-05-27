package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"todoapi/usecase"
)

type ItemHandler interface {
	HandleOne(http.ResponseWriter, *http.Request)
	HandleAll(http.ResponseWriter, *http.Request)
}

func NewItemHandler(usecase usecase.ItemUseCase) ItemHandler {
	return &itemHandler{
		usecase: usecase,
	}
}

type itemHandler struct {
	usecase usecase.ItemUseCase
}

func (handler itemHandler) HandleAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // localhost:3000からのオリジン間アクセスを許可する

	switch r.Method {
	case "GET":
		getAllItems(w, r, handler.usecase) // 全てのitemの取得
	case "POST":
		addNewItem(w, r, handler.usecase) // 新しいitemの追加
	case "DELETE":
		deleteDoneItems(w, r, handler.usecase) // 実行済みitemの削除
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")               // Content-Typeヘッダの使用を許可する
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS") // pre-flightリクエストに対応する
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (handler itemHandler) HandleOne(w http.ResponseWriter, r *http.Request) {
	params := getRouteParams(r)
	if len(params) != 3 || params[2] != "done" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	id := params[1]
	fmt.Fprintf(w, "<h1>%s</h1>", id)
}

func getAllItems(w http.ResponseWriter, r *http.Request, usecase usecase.ItemUseCase) {
	items, err := usecase.GetAll()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(items); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, buf.String())
}

func addNewItem(w http.ResponseWriter, r *http.Request, usecase usecase.ItemUseCase) {
	var reqBody struct {
		Name string `json:"name"`
	}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&reqBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = usecase.AddItem(reqBody.Name)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func deleteDoneItems(w http.ResponseWriter, r *http.Request, usecase usecase.ItemUseCase) {
	err := usecase.DeleteDone()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
