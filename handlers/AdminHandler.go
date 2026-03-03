package handlers

import (
"encoding/json"
"net/http"

"github.com/office-sports/ttapp-api/data"
)

func AdminGetAllReschedules(w http.ResponseWriter, r *http.Request) {
SetHeaders(w)
claims := r.Context().Value("claims").(*AuthClaims)
if claims == nil || !claims.IsAdmin {
http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
return
}
requests, err := data.GetAllRescheduleRequests()
if err != nil {
http.Error(w, `{"error":"Failed to get reschedule requests"}`, http.StatusInternalServerError)
return
}
json.NewEncoder(w).Encode(requests)
}
