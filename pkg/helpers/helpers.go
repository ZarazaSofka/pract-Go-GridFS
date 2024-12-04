package helpers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func JSONSend(w http.ResponseWriter, data any, code int) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, `{"error": "response write error"}`, http.StatusInternalServerError)
	}
}

func JSONMessageSend(w http.ResponseWriter, text string, code int) {
	JSONSend(w, map[string]string{"message": text}, code)
}

func DatabaseError(w http.ResponseWriter, err error, logger *zap.SugaredLogger) {
	JSONMessageSend(w, "database error", http.StatusInternalServerError)
	logger.Errorf("%s: %w", err.Error(), err)
}
