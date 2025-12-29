package students

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/OmarAkhter/golan-project/internal/types"
	response "github.com/OmarAkhter/golan-project/internal/utils"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("request body cannot be empty")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		response.WriteJSON(w, http.StatusCreated, response.Response{
			Status:  response.StatusSuccess,
			Message: "Student created successfully",
			Data: map[string]string{
				"id": student.ID,
			},
		})
	}
}
