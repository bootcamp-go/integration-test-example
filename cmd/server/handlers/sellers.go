package handlers

import (
	"app/internal/sellers"
	"app/internal/sellers/repository"
	"app/pkg/web/request"
	"app/pkg/web/response"
	"errors"
	"net/http"
	"strconv"
)

// NewControllerSellers creates a new controller for sellers
func NewControllerSellers(rp repository.RepositorySellers) *ControllerSellers {
	return &ControllerSellers{rp: rp}
}

// ControllerSellers is a controller for sellers to generate handlers
type ControllerSellers struct {
	rp repository.RepositorySellers
}

// GetById returns a seller by id
type ResponseBodyGetById struct {
	Message string `json:"message"`
	Data    *struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	} `json:"data"`
	Error bool `json:"error"`
}
func (c *ControllerSellers) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		idParam, err := request.PathLastParam(r)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBodyGetById{
				Message: "bad request",
				Data:    nil,
				Error:   true,
			}

			response.JSON(w, code, body)
			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBodyGetById{
				Message: "bad request",
				Data:    nil,
				Error:   true,
			}

			response.JSON(w, code, body)
			return
		}

		// process
		s, err := c.rp.GetById(id)
		if err != nil {
			var code int; var body *ResponseBodyGetById
			switch {
			case errors.Is(err, repository.ErrRepositorySellersNotFound):
				code = http.StatusNotFound
				body = &ResponseBodyGetById{
					Message: "seller not found",
					Data:    nil,
					Error:   true,
				}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBodyGetById{
					Message: "internal error",
					Data:    nil,
					Error:   true,
				}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBodyGetById{
			Message: "success",
			Data: &struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Email     string `json:"email"`
			}{
				ID:        s.ID,
				FirstName: s.FirstName,
				LastName:  s.LastName,
				Email:     s.Email,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}

// Save saves a seller
type RequestBodySave struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
type ResponseBodySave struct {
	Message string `json:"message"`
	Data    *struct {
		ID		  int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	} `json:"data"`
	Error bool `json:"error"`
}
func (c *ControllerSellers) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var req RequestBodySave
		if err := request.JSON(r, &req); err != nil {
			code := http.StatusBadRequest
			body := &ResponseBodySave{
				Message: "bad request",
				Data:    nil,
				Error:   true,
			}

			response.JSON(w, code, body)
			return
		}

		// process
		s := &sellers.Seller{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		}
		err := c.rp.Save(s)
		if err != nil {
			code := http.StatusInternalServerError
			body := &ResponseBodySave{
				Message: "internal error",
				Data:    nil,
				Error:   true,
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusCreated
		body := &ResponseBodySave{
			Message: "success",
			Data: &struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Email     string `json:"email"`
			}{
				ID:        s.ID,
				FirstName: s.FirstName,
				LastName:  s.LastName,
				Email:     s.Email,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}