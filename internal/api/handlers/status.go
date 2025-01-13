// internal/api/handlers/status.go
package handlers

import (
	"encoding/json"
	"net/http"

	"retail-pulse/internal/service"
)

func GetJobStatus(w http.ResponseWriter, r *http.Request) {
    jobID := r.URL.Query().Get("jobid")
    if jobID == "" {
        http.Error(w, "jobid is required", http.StatusBadRequest)
        return
    }
    
    status, err := service.GetJobStatus(jobID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    json.NewEncoder(w).Encode(status)
}