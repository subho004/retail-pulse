// internal/service/job_service.go
package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"retail-pulse/internal/models"
	"strings"
	"sync"
	"time"
)

var (
    jobs     = make(map[string]*models.JobStatus)
    jobsLock sync.RWMutex
    stores   = make(map[string]models.Store)
)

func init() {
    // Read store master data from CSV
    content, err := os.ReadFile("StoreMasterAssignment.csv")
    if err != nil {
        log.Fatalf("Error reading store master CSV: %v", err)
    }

    // Parse CSV
    reader := csv.NewReader(strings.NewReader(string(content)))
    
    // Skip header
    _, err = reader.Read()
    if err != nil {
        log.Fatalf("Error reading CSV header: %v", err)
    }

    // Read and store all records
    stores = make(map[string]models.Store)
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("Error reading CSV record: %v", err)
        }

        // Record format: AreaCode, StoreName, StoreID
        stores[record[2]] = models.Store{
            ID:       record[2],
            Name:     record[1],
            AreaCode: record[0],
        }
    }

    log.Printf("Loaded %d stores from CSV", len(stores))
}

func ProcessJob(req models.JobRequest) (string, error) {
    jobID := generateJobID()
    
    // Initialize job status
    jobsLock.Lock()
    jobs[jobID] = &models.JobStatus{
        Status: "ongoing",
        JobID:  jobID,
    }
    jobsLock.Unlock()
    
    // Process job asynchronously
    go processImages(jobID, req)
    
    return jobID, nil
}

func processImages(jobID string, req models.JobRequest) {
    var wg sync.WaitGroup
    errorsChan := make(chan models.JobError, len(req.Visits))
    
    for _, visit := range req.Visits {
        wg.Add(1)
        go func(v models.Visit) {
            defer wg.Done()
            
            // Validate store exists
            if _, exists := stores[v.StoreID]; !exists {
                errorsChan <- models.JobError{
                    StoreID: v.StoreID,
                    Error:   "store not found",
                }
                return
            }
            
            // Process images
            for _, url := range v.ImageURLs {
                if err := processImage(url); err != nil {
                    errorsChan <- models.JobError{
                        StoreID: v.StoreID,
                        Error:   fmt.Sprintf("failed to process image: %v", err),
                    }
                    return
                }
            }
        }(visit)
    }
    
    wg.Wait()
    close(errorsChan)
    
    // Collect errors
    var errors []models.JobError
    for err := range errorsChan {
        errors = append(errors, err)
    }
    
    // Update job status
    jobsLock.Lock()
    if len(errors) > 0 {
        jobs[jobID].Status = "failed"
        jobs[jobID].Errors = errors
    } else {
        jobs[jobID].Status = "completed"
    }
    jobsLock.Unlock()
}

func processImage(url string) error {
    // Download image
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Read image data
    _, err = io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    // Simulate processing time
    time.Sleep(time.Duration(100+rand.Intn(300)) * time.Millisecond)
    
    return nil
}

func GetJobStatus(jobID string) (*models.JobStatus, error) {
    jobsLock.RLock()
    defer jobsLock.RUnlock()
    
    status, exists := jobs[jobID]
    if !exists {
        return nil, errors.New("job not found")
    }
    
    return status, nil
}

func generateJobID() string {
    return fmt.Sprintf("job_%d", time.Now().UnixNano())
}