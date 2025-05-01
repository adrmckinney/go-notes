package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetIdFromUrlPath(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		return 0, fmt.Errorf("missing 'id' path parameter")
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid 'id' path parameter: %w", err)
	}
	return uint(id), nil
}
