package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CalculationRequest struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

type CalculationResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/calculate", calculateHandler)
	fmt.Println("Starting server on port 80...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp := performCalculation(req)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func performCalculation(req CalculationRequest) CalculationResponse {
	var result float64
	var err error
	switch req.Operation {
	case "add":
		result = req.A + req.B
	case "subtract":
		result = req.A - req.B
	case "multiply":
		result = req.A * req.B
	case "divide":
		if req.B == 0 {
			err = fmt.Errorf("cannot divide by zero")
		} else {
			result = req.A / req.B
		}
	default:
		err = fmt.Errorf("unsupported operation: %s", req.Operation)
	}

	if err != nil {
		return CalculationResponse{Error: err.Error()}
	}
	return CalculationResponse{Result: result}
}
