package usersvc

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
	"github.com/remiges-tech/alya/examples/pg/sqlc-gen"
	"github.com/remiges-tech/alya/wscutils"

	"github.com/remiges-tech/alya/service"
)

type CreateUserParams struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"required,email"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func HandleCreateUserRequest(c *gin.Context, s *service.Service) {
	s.Logger.Log("CreateUser request received")
	// Parse request
	var createUserReq CreateUserRequest
	if err := wscutils.BindJSON(c, &createUserReq); err != nil {
		return
	}
	s.Logger.Log(fmt.Sprintf("CreateUser request parsed: %v", map[string]any{"username": createUserReq.Name}))

	// Validate request
	validationErrors := wscutils.WscValidate(createUserReq, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		wscutils.SendErrors(c, validationErrors)
		return
	}
	s.Logger.Log(fmt.Sprintf("CreateUser request validated %v", map[string]any{"username": createUserReq.Name}))

	// Call CreateUser function
	user, err := s.Database.(*sqlc.Queries).CreateUser(c.Request.Context(), sqlc.CreateUserParams{
		Name:  createUserReq.Name,
		Email: createUserReq.Email,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	s.Logger.Log(fmt.Sprintf("User created: %v", map[string]any{"username": createUserReq.Name}))

	// Send response
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(user))
	wscutils.SendData(c, user)
}
