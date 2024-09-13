package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
	"zadanie_6105/src/models"
	"zadanie_6105/src/services"
)

func formatBidToExport(bid *models.Bid) map[string]interface{} {
	result := map[string]interface{}{
		"id":          bid.ID.String(),
		"name":        bid.Name,
		"description": bid.Description,
		"status":      bid.Status,
		"tenderId":    bid.TenderId.String(),
		"authorType":  bid.AuthorType,
		"authorId":    bid.AuthorId.String(),
		"version":     bid.Version,
		"createdAt":   bid.CreatedAt.Format(time.RFC3339),
	}
	return result
}

func formatBidsToExport(bids *[]models.Bid) []map[string]interface{} {
	result := make([]map[string]interface{}, len(*bids))
	for i, bid := range *bids {
		result[i] = formatBidToExport(&bid)
	}
	return result
}

func formatFeedbackToExport(feedback *models.BidFeedback) map[string]interface{} {
	result := map[string]interface{}{
		"id":          feedback.ID.String(),
		"description": feedback.Description,
		"createdAt":   feedback.CreatedAt.Format(time.RFC3339),
	}
	return result
}

func formatFeedbacksToExport(feedbacks *[]models.BidFeedback) []map[string]interface{} {
	result := make([]map[string]interface{}, len(*feedbacks))
	for i, feedback := range *feedbacks {
		result[i] = formatFeedbackToExport(&feedback)
	}
	return result
}

func checkBidStatus(s string) bool {
	list := []string{"Created", "Published", "Canceled"}
	return checkParam(s, &list)
}

func checkAuthorType(s string) bool {
	list := []string{"Organization", "User"}
	return checkParam(s, &list)
}

func checkDecision(s string) bool {
	list := []string{"Approved", "Rejected"}
	return checkParam(s, &list)
}

func CreateBid(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bid models.Bid
		err := json.NewDecoder(r.Body).Decode(&bid)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if !checkAuthorType(bid.AuthorType) {
			http.Error(w, `Author type can be only "Organization", "User"`, http.StatusBadRequest)
			return
		}

		isPublished, err := service.CheckIfTenderPublished(bid.TenderId.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if !isPublished {
			http.Error(w, `Tender is not published`, http.StatusNotFound)
			return
		}

		isExist, err := service.CheckEmployeeExistence(bid.AuthorId.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isExist {
			http.Error(w, `User is not exist`, http.StatusUnauthorized)
			return
		}

		err = service.CreateBid(&bid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := formatBidToExport(&bid)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetBidsByUser(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		username := r.URL.Query().Get("username")
		limit := 5
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			} else {
				http.Error(w, "Invalid paginationLimit", http.StatusBadRequest)
				return
			}
		}
		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil {
				offset = o
			} else {
				http.Error(w, "Invalid paginationOffset", http.StatusBadRequest)
				return
			}
		}
		if username == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}

		bids, err := service.GetBidsByUser(username, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		response := formatBidsToExport(bids)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetBidsByTender(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tenderID := vars["tenderId"]
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		username := r.URL.Query().Get("username")
		limit := 5
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			} else {
				http.Error(w, "Invalid paginationLimit", http.StatusBadRequest)
				return
			}
		}
		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil {
				offset = o
			} else {
				http.Error(w, "Invalid paginationOffset", http.StatusBadRequest)
				return
			}
		}

		if username == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}
		isResponsible, err := service.CheckIfUserIsResponsibleForTender(username, tenderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		bids, err := service.GetBidsByTender(tenderID, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		response := formatBidsToExport(bids)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetBidStatus(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bidID := vars["bidId"]
		username := r.URL.Query().Get("username")

		if username == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}

		isResponsible, err := service.CheckIfUserIsResponsibleForBid(username, bidID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		status, err := service.GetBidStatus(bidID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		_, err = w.Write([]byte(status))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateBidStatus(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bidID := vars["bidId"]
		username := r.URL.Query().Get("username")
		status := r.URL.Query().Get("status")

		if !checkBidStatus(status) {
			http.Error(w, `Status can be only "Created", "Published", "Canceled"`, http.StatusBadRequest)
		}

		if username == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}

		isResponsible, err := service.CheckIfUserIsResponsibleForBid(username, bidID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		bid, err := service.UpdateBidStatus(bidID, status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response := formatBidToExport(bid)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func UpdateBid(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		bidID := vars["bidId"]
		username := r.URL.Query().Get("username")

		var edit models.BidEdit
		err := json.NewDecoder(r.Body).Decode(&edit)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check username
		if username == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}

		isResponsible, err := service.CheckIfUserIsResponsibleForBid(username, bidID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		bid, err := service.UpdateBid(bidID, &edit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := formatBidToExport(bid)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func SubmitBid(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bidID := vars["bidId"]
		username := r.URL.Query().Get("username")
		decisionStr := r.URL.Query().Get("decision")

		if !checkDecision(decisionStr) {
			http.Error(w, `Status can be only "Approved", "Rejected"`, http.StatusBadRequest)
		}
		decision := decisionStr == "Approved"

		if username == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}

		isResponsible, err := service.CheckIfUserIsResponsibleForTenderByBidID(username, bidID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		bid, err := service.SubmitBid(bidID, decision)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response := formatBidToExport(bid)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func RollbackBid(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		bidID := vars["bidId"]
		versionStr := vars["version"]
		username := r.URL.Query().Get("username")
		var version int32

		// Check version
		if versionStr != "" {
			if v, err := strconv.Atoi(versionStr); err == nil {
				version = int32(v)
			} else {
				http.Error(w, "Invalid version", http.StatusBadRequest)
				return
			}
		}

		// Check username
		if username == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}

		isResponsible, err := service.CheckIfUserIsResponsibleForBid(username, bidID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		bid, err := service.RollbackBid(bidID, version)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := formatBidToExport(bid)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func CreateFeedback(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bidIDStr := vars["bidId"]
		username := r.URL.Query().Get("username")
		description := r.URL.Query().Get("bidFeedback")

		if username == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}

		isResponsible, err := service.CheckIfUserIsResponsibleForTenderByBidID(username, bidIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		bidID, err := uuid.Parse(bidIDStr)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		feedback := models.BidFeedback{
			Description: description,
			BidId:       bidID,
		}

		bid, err := service.CreateFeedback(&feedback)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response := formatBidToExport(bid)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetFeedbacks(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tenderID := vars["tenderId"]
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		author := r.URL.Query().Get("authorUsername")
		requester := r.URL.Query().Get("requesterUsername")
		limit := 5
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			} else {
				http.Error(w, "Invalid paginationLimit", http.StatusBadRequest)
				return
			}
		}
		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil {
				offset = o
			} else {
				http.Error(w, "Invalid paginationOffset", http.StatusBadRequest)
				return
			}
		}

		if requester == "" {
			http.Error(w, "Username required", http.StatusUnauthorized)
			return
		}
		isResponsible, err := service.CheckIfUserIsResponsibleForTender(requester, tenderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		isAuthor, err := service.CheckIfBidByUserExist(author, tenderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !isAuthor {
			http.Error(w, "This user is not author", http.StatusBadRequest)
			return
		}

		feedbacks, err := service.GetFeedbacks(author, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		response := formatFeedbacksToExport(feedbacks)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
