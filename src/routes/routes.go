package routes

import (
	"github.com/gorilla/mux"
	"zadanie_6105/src/handlers"
	"zadanie_6105/src/services"
)

func RegisterRoutes(service *services.Service) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/ping", handlers.Ping).Methods("GET")

	r.HandleFunc("/api/tenders", handlers.GetTenders(service)).Methods("GET")
	r.HandleFunc("/api/tenders/new", handlers.CreateTender(service)).Methods("POST")
	r.HandleFunc("/api/tenders/my", handlers.GetTendersByUser(service)).Methods("GET")
	r.HandleFunc("/api/tenders/{tenderId}/status", handlers.GetTenderStatus(service)).Methods("GET")
	r.HandleFunc("/api/tenders/{tenderId}/status", handlers.UpdateTenderStatus(service)).Methods("PUT")
	r.HandleFunc("/api/tenders/{tenderId}/edit", handlers.UpdateTender(service)).Methods("PATCH")
	r.HandleFunc("/api/tenders/{tenderId}/rollback/{version}", handlers.RollbackTender(service)).Methods("PUT")

	r.HandleFunc("/api/bids/new", handlers.CreateBid(service)).Methods("POST")
	r.HandleFunc("/api/bids/my", handlers.GetBidsByUser(service)).Methods("GET")
	r.HandleFunc("/api/bids/{tenderId}/list", handlers.GetBidsByTender(service)).Methods("GET")
	r.HandleFunc("/api/bids/{bidId}/status", handlers.GetBidStatus(service)).Methods("GET")
	r.HandleFunc("/api/bids/{bidId}/status", handlers.UpdateBidStatus(service)).Methods("PUT")
	r.HandleFunc("/api/bids/{bidId}/edit", handlers.UpdateBid(service)).Methods("PATCH")
	r.HandleFunc("/api/bids/{bidId}/submit_decision", handlers.SubmitBid(service)).Methods("PUT")
	r.HandleFunc("/api/bids/{bidId}/rollback/{version}", handlers.RollbackBid(service)).Methods("PUT")
	r.HandleFunc("/api/bids/{bidId}/feedback", handlers.CreateFeedback(service)).Methods("PUT")
	r.HandleFunc("/api/bids/{tenderId}/reviews", handlers.GetFeedbacks(service)).Methods("GET")

	return r
}
