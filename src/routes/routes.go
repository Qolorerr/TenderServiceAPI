package routes

import (
	"github.com/gorilla/mux"
	"zadanie_6105/src/handlers"
	"zadanie_6105/src/services"
)

func RegisterRoutes(tenderService *services.TenderService) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/ping", handlers.Ping).Methods("GET")
	r.HandleFunc("/api/tenders", handlers.GetTenders(tenderService)).Methods("GET")
	r.HandleFunc("/api/tenders/new", handlers.CreateTender(tenderService)).Methods("POST")
	r.HandleFunc("/api/tenders/my", handlers.GetTendersByUser(tenderService)).Methods("GET")
	r.HandleFunc("/api/tenders/{tenderId}/status", handlers.GetTenderStatus(tenderService)).Methods("GET")
	r.HandleFunc("/api/tenders/{tenderId}/status", handlers.UpdateTenderStatus(tenderService)).Methods("PUT")
	r.HandleFunc("/api/tenders/{tenderId}/edit", handlers.UpdateTender(tenderService)).Methods("PATCH")
	r.HandleFunc("/api/tenders/{tenderId}/rollback/{version}", handlers.RollbackTender(tenderService)).Methods("PUT")

	return r
}
