package handler

import (
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
	fmt.Fprintln(w, "<h1>itemHandler</h1>")
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
