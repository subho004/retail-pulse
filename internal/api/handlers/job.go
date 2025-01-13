// internal/api/handlers/job.go
package handlers

import (
	"encoding/json"
	"net/http"

	"retail-pulse/internal/models"
	"retail-pulse/internal/service"
)

func SubmitJob(w http.ResponseWriter, r *http.Request) {
    var req models.JobRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Validate request
    if req.Count != len(req.Visits) {
        http.Error(w, "count does not match number of visits", http.StatusBadRequest)
        return
    }
    
    jobID, err := service.ProcessJob(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(models.JobResponse{JobID: jobID})
}