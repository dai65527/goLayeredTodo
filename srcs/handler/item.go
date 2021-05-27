package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"todoapi/domain/model"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// ルートパラメータの取得（例: `/items/1/done` -> ["items", "1", "done"]）
	params := getRouteParams(r)
	if len(params) == 2 {
		updateItem(params[1], w, r, handler.usecase)
	} else if len(params) == 3 && params[2] == "done" {
		updateDone(params[1], w, r, handler.usecase)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func updateItem(id model.ID, w http.ResponseWriter, r *http.Request, usecase usecase.ItemUseCase) {
	switch r.Method {
	case "DELETE":
		deleteItem(id, w, r, usecase)
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")    // Content-Typeヘッダの使用を許可する
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS") // pre-flightリクエストに対応する
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func updateDone(id model.ID, w http.ResponseWriter, r *http.Request, usecase usecase.ItemUseCase) {
	switch r.Method {
	case "PUT":
		doneItem(id, w, r, usecase)
	case "DELETE":
		unDoneItem(id, w, r, usecase)
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")         // Content-Typeヘッダの使用を許可する
		w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, OPTIONS") // pre-flightリクエストに対応する
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
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

func deleteItem(id model.ID, w http.ResponseWriter, r *http.Request, usecase usecase.ItemUseCase) {
	err := usecase.Delete(id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func doneItem(id model.ID, w http.ResponseWriter, r *http.Request, usecase usecase.ItemUseCase) {
	err := usecase.Done(id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func unDoneItem(id model.ID, w http.ResponseWriter, r *http.Request, usecase usecase.ItemUseCase) {
	err := usecase.UnDone(id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
