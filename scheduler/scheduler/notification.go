package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"scheduler/cli/logger"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type AlertRequest struct {
	Dedup   string `json:"dedup"`
	Summary string `json:"summary"`
	Details string `json:"details"`
	Urgency string `json:"urgency,omitempty"`
}

func sendNotification(name, res string) {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Error("Unable to load env file", zap.Error(err))
	}
	baseURL := os.Getenv("BASE_URL")
	payload := AlertRequest{
		Dedup:   fmt.Sprintf("alert - %v", time.Now()),
		Summary: name,
		Details: res,
		Urgency: "High",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		// logger.Log.Error("Error in marshaling", zap.Error(err))
		fmt.Println(err)
	}
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// logger.Log.Error("Error in marshaling", zap.Error(err))
		fmt.Println(err)
	}
	defer resp.Body.Close()
	// bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// logger.Log.Error("Error in marshaling", zap.Error(err))
		fmt.Println(err)
	}
	// logger.Log.Info("Status of sending event to destination", zap.String("status", resp.Status), zap.Int("statusCode", resp.StatusCode), zap.String("statusBody", string(bodyBytes)))
	fmt.Println("notification send")
}
