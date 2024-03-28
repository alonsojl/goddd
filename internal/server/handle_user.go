package server

import (
	"goddd/internal/domain/user"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type UserHandler interface {
	HandleGetUsers() http.HandlerFunc
	HandleGetUserByEmail() http.HandlerFunc
	HandleCreateUser() http.HandlerFunc
}

type userHandler struct {
	logger      *logrus.Logger
	userService user.Service
}

func NewUserHandler(logger *logrus.Logger, userService user.Service) UserHandler {
	return &userHandler{
		logger:      logger,
		userService: userService,
	}
}

type userGetResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Age       int16  `json:"age"`
	Position  string `json:"position"`
}

// @Summary 	Get users.
// @Description Description...
// @Tags 		User
// @Router 		/users [get]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Success     200	{object} server.HandleGetUsers.response "Success"
// @Failure     500	{object} server.errorResponse "Error Internal Server"
func (u *userHandler) HandleGetUsers() http.HandlerFunc {
	type response struct {
		Status    string            `json:"status"`
		HTTPCode  int               `json:"http_code"`
		Datetime  string            `json:"datetime"`
		Timestamp int64             `json:"timestamp"`
		User      []userGetResponse `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		users, err := u.userService.GetUsers(r.Context())
		if err != nil {
			u.logger.Printf("error get users: %v", err)
			encodeError(w, err)
			return
		}

		var result []userGetResponse
		for _, user := range users {
			result = append(result, userGetResponse{
				ID:        user.ID,
				Name:      user.Name,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Phone:     user.Phone,
				Age:       user.Age,
				Position:  user.Position,
			})
		}

		respond(w, response{
			Status:    "success",
			HTTPCode:  http.StatusOK,
			Datetime:  time.Now().Format("2006-01-02 15:04:05"),
			Timestamp: time.Now().Unix(),
			User:      result,
		}, http.StatusOK)
	}
}

// @Summary 	Get user by email.
// @Description Description...
// @Tags 		User
// @Router 		/users/{email} [get]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param       email path string true "Email user"
// @Success     200	{object} server.HandleGetUserByEmail.response "Success"
// @Failure     500	{object} server.errorResponse "Error Internal Server"
func (u *userHandler) HandleGetUserByEmail() http.HandlerFunc {
	type response struct {
		Status    string          `json:"status"`
		HTTPCode  int             `json:"http_code"`
		Datetime  string          `json:"datetime"`
		Timestamp int64           `json:"timestamp"`
		User      userGetResponse `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		user, err := u.userService.GetUserByEmail(r.Context(), email)
		if err != nil {
			u.logger.Printf("error get user by email: %v", err)
			encodeError(w, err)
			return
		}

		respond(w, response{
			Status:    "success",
			HTTPCode:  http.StatusOK,
			Datetime:  time.Now().Format("2006-01-02 15:04:05"),
			Timestamp: time.Now().Unix(),
			User: userGetResponse{
				ID:        user.ID,
				Name:      user.Name,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Phone:     user.Phone,
				Age:       user.Age,
				Position:  user.Position,
			},
		}, http.StatusOK)
	}
}

type userAddResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// @Summary 	Create user.
// @Description Description...
// @Tags 		User
// @Router 		/users [post]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param	    request body server.HandleCreateUser.request true "User params"
// @Success     200	{object} server.HandleCreateUser.response "Success"
// @Failure     400	{object} server.errorResponse "Invalid Request"
// @Failure     500	{object} server.errorResponse "Error Internal Server"
func (u *userHandler) HandleCreateUser() http.HandlerFunc {
	type (
		request struct {
			Name      string `json:"name" example:"Jorge Luis"`
			FirstName string `json:"first_name" example:"Alonso"`
			LastName  string `json:"last_name" example:"Hdez"`
			Email     string `json:"email" example:"alonso12.dev@gmail.com"`
			Phone     string `json:"phone" example:"7713037204"`
			Age       int16  `json:"age" example:"25"`
			Position  string `json:"position" example:"Go developer"`
		}
		response struct {
			Status    string          `json:"status"`
			HTTPCode  int             `json:"http_code"`
			Datetime  string          `json:"datetime"`
			Timestamp int64           `json:"timestamp"`
			User      userAddResponse `json:"user"`
		}
	)
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req request
			err error
		)

		if err = decode(r, &req); err != nil {
			u.logger.Errorf("invalid user request: %v", err)
			encodeError(w, err)
			return
		}
		params := &user.Params{
			Name:      req.Name,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
			Age:       req.Age,
			Position:  req.Position,
		}

		user, err := u.userService.CreateUser(r.Context(), params)
		if err != nil {
			u.logger.Errorf("error created user: %v", err)
			encodeError(w, err)
			return
		}

		respond(w, response{
			Status:    "success",
			HTTPCode:  http.StatusCreated,
			Datetime:  time.Now().Format("2006-01-02 15:04:05"),
			Timestamp: time.Now().Unix(),
			User: userAddResponse{
				ID:        user.ID,
				Name:      user.Name,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
		}, http.StatusCreated)
	}
}
