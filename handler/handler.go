package handler

import (
	"avitosegments/database"
	"encoding/json"
	"io"
	"net/http"
)

var Api *database.API

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBody, _ := io.ReadAll(r.Body)
	var segment Segment
	err := json.Unmarshal(reqBody, &segment)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("invalid request"))
		return
	}
	err = (*Api).CreateSegment(segment.Slug)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(201)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBody, _ := io.ReadAll(r.Body)
	var segment Segment
	err := json.Unmarshal(reqBody, &segment)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("invalid request"))
		return
	}
	err = (*Api).DeleteSegment(segment.Slug)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
}

func ChangeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBody, _ := io.ReadAll(r.Body)
	var changeRequest ChangeRequest
	err := json.Unmarshal(reqBody, &changeRequest)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("invalid request"))
		return
	}
	if changeRequest.TTL < 0 {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("TTL must not be negative"))
		return
	}
	err = (*Api).ChangeSegments(
		changeRequest.ToAdd.Slugs,
		changeRequest.ToDelete.Slugs,
		changeRequest.User.ID,
		changeRequest.TTL)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBody, _ := io.ReadAll(r.Body)
	var user User
	err := json.Unmarshal(reqBody, &user)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("invalid request"))
		return
	}
	slugs, err := (*Api).GetSegments(user.ID)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	respJSON, _ := json.Marshal(Segments{slugs})
	_, _ = w.Write(respJSON)
}
