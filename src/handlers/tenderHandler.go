package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
	"zadanie_6105/src/models"
	"zadanie_6105/src/services"
)

func formatTenderToExport(tender *models.Tender) map[string]interface{} {
	result := map[string]interface{}{
		"id":             tender.ID.String(),
		"name":           tender.Name,
		"description":    tender.Description,
		"serviceType":    tender.ServiceType,
		"status":         tender.Status,
		"organizationId": tender.OrganizationId,
		"version":        tender.Version,
		"createdAt":      tender.CreatedAt.Format(time.RFC3339),
	}
	return result
}

func formatTendersToExport(tenders *[]models.Tender) []map[string]interface{} {
	result := make([]map[string]interface{}, len(*tenders))
	for i, tender := range *tenders {
		result[i] = formatTenderToExport(&tender)
	}
	return result
}

func checkParam(s string, list *[]string) bool {
	for _, v := range *list {
		if v == s {
			return true
		}
	}
	return false
}

func checkTenderStatus(s string) bool {
	list := []string{"Created", "Published", "Closed"}
	return checkParam(s, &list)
}

func checkServiceType(s string) bool {
	list := []string{"Construction", "Delivery", "Manufacture"}
	return checkParam(s, &list)
}

func GetTenders(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
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

		serviceTypeStr := r.URL.Query().Get("service_type")
		var serviceTypes []string
		if serviceTypeStr != "" {
			serviceTypes = strings.Split(serviceTypeStr, ",")
		}

		tenders, err := service.GetTenders(serviceTypes, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		response := formatTendersToExport(tenders)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func CreateTender(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tender models.Tender
		err := json.NewDecoder(r.Body).Decode(&tender)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if !checkServiceType(tender.ServiceType) {
			http.Error(w, `Service type can be only "Construction", "Delivery", "Manufacture"`, http.StatusBadRequest)
			return
		}

		isResponsible, err := service.CheckIfUserIsResponsible(tender.CreatorUsername, tender.OrganizationId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		err = service.CreateTender(&tender)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := formatTenderToExport(&tender)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetTendersByUser(service *services.Service) http.HandlerFunc {
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

		tenders, err := service.GetTendersByUser(username, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		response := formatTendersToExport(tenders)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetTenderStatus(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tenderID := vars["tenderId"]
		username := r.URL.Query().Get("username")

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

		status, err := service.GetTenderStatus(tenderID)
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

func UpdateTenderStatus(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tenderID := vars["tenderId"]
		username := r.URL.Query().Get("username")
		status := r.URL.Query().Get("status")

		if !checkTenderStatus(status) {
			http.Error(w, `Status can be only "Created", "Published", "Closed"`, http.StatusBadRequest)
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

		tender, err := service.UpdateTenderStatus(tenderID, status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response := formatTenderToExport(tender)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func UpdateTender(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		tenderID := vars["tenderId"]
		username := r.URL.Query().Get("username")

		var edit models.TenderEdit
		err := json.NewDecoder(r.Body).Decode(&edit)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check body
		if !checkServiceType(edit.ServiceType) {
			http.Error(w, `Service type can be only "Construction", "Delivery", "Manufacture"`, http.StatusBadRequest)
			return
		}

		// Check username
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

		tender, err := service.UpdateTender(tenderID, &edit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := formatTenderToExport(tender)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func RollbackTender(service *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		tenderID := vars["tenderId"]
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

		isResponsible, err := service.CheckIfUserIsResponsibleForTender(username, tenderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !isResponsible {
			http.Error(w, "This user is not responsible", http.StatusForbidden)
			return
		}

		tender, err := service.RollbackTender(tenderID, version)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := formatTenderToExport(tender)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
