package controller

import (
	"encoding/json"
	"net/http"
	"shorters/domain"
	"shorters/service"
	"shorters/service/dto"

	"github.com/go-chi/chi"
)

type LinkController interface {
	Redirect(w http.ResponseWriter, r *http.Request)
	Store(w http.ResponseWriter, r *http.Request)
}

type linkController struct {
	linkService service.LinkService
}

func NewLinkController(linkService service.LinkService) LinkController {
	return &linkController{linkService: linkService}
}

func (c *linkController) Redirect(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	url, err := c.linkService.Find(key)
	if err != nil {
		switch err.(type) {
		case domain.LinkNotFoundError:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (c *linkController) Store(w http.ResponseWriter, r *http.Request) {
	var req dto.StoreRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, _ := r.Context().Value("user").(domain.User)
	res, err := c.linkService.Store(req.URL, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}
