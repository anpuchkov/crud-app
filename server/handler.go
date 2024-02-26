package server

import (
	"Api/domain"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Service interface {
	Create(ctx context.Context, service domain.Service) error
	GetByID(ctx context.Context, id int64) (domain.Service, error)
	GetAll(ctx context.Context) ([]domain.Service, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, service domain.UpdateService) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitHandler() *mux.Router {
	r := mux.NewRouter()
	r.Use(LoggingMiddleware)

	services := r.PathPrefix("/services").Subrouter()
	{
		services.HandleFunc("", h.createService).Methods(http.MethodPost)
		services.HandleFunc("", h.getAllServices).Methods(http.MethodGet)
		services.HandleFunc("/{id:[0-9]+}", h.getServiceByID).Methods(http.MethodGet)
		services.HandleFunc("/{id:[0-9]+}", h.deleteService).Methods(http.MethodDelete)
		services.HandleFunc("/{id:[0-9]+}", h.updateService).Methods(http.MethodPut)
	}
	return r
}

func (h *Handler) createService(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error in creating service: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var service domain.Service
	err = json.Unmarshal(reqBytes, &service)
	if err != nil {
		log.Println("error in unmarshalling data: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.Create(context.TODO(), service)
	if err != nil {
		log.Println("error in creating service: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getServiceByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	service, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching service by ID: %v", err)

		if errors.Is(err, domain.ErrServiceNotFound) {
			http.Error(w, "Service not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	jsonResponse, err := json.Marshal(service)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *Handler) getAllServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.service.GetAll(context.TODO())
	if err != nil {
		log.Println("error in getAll: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(services)
	if err != nil {
		log.Println("error in marshalling data: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)

}

func (h *Handler) deleteService(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(w, r)
	if err != nil {
		log.Println("error while getting id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.service.Delete(context.TODO(), id)
	if err != nil {
		log.Println("error while deleting service: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updateService(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(w, r)
	if err != nil {
		log.Println("error while getting id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error in creating service: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var input domain.UpdateService
	err = json.Unmarshal(reqBytes, &input)
	if err != nil {
		log.Println("error in unmarshalling data: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.service.Update(context.TODO(), id, input)
	if err != nil {
		log.Println("error while updating service: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(w http.ResponseWriter, r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Println("error in psarsing id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return 0, err
	}
	if id == 0 {
		return 0, errors.New("id cannot be 0")
	}
	return id, nil
}
