package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Target struct {
	Target string
}

func (h *Handler) Target(writer http.ResponseWriter, request *http.Request) {
	target := h.getTargetFromRequest(request)

	switch request.Method {
	case "POST":
		h.addTarget(writer, target)
	case "DELETE":
		h.deleteTarget(writer, target)
	}
}

func (h *Handler) addTarget(writer http.ResponseWriter, target Target)  {
	_, err := h.Res.Db.Exec("INSERT OR IGNORE INTO urls (url) VALUES ('" + target.Target + "'); ")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Error " + error.Error(err)))
		error.Error(err)
	}

	h.checkRows()

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Added " + target.Target))
}

func (h *Handler) deleteTarget(writer http.ResponseWriter, target Target)  {
	_, err := h.Res.Db.Exec("DELETE FROM urls where url = '" + target.Target +"';")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Error " + error.Error(err)))
		error.Error(err)
	}

	h.checkRows()
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Deleted " + target.Target))
}

func (h *Handler) getTargetFromRequest(request *http.Request) Target {
	tu := Target{}
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&tu)

	return tu
}

func (h *Handler)checkRows()  {
	rows, _:= h.Res.Db.Query("SELECT * FROM urls")
	var id int
	var url string
	for rows.Next() {
		rows.Scan(&id, &url)
		fmt.Println(id, url)
	}
}