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
	s.LogHarbour.Log("CreateUser request received")
	// Parse request
	var createUserReq CreateUserRequest
	if err := wscutils.BindJSON(c, &createUserReq); err != nil {
		return
	}
	s.LogHarbour.Log(fmt.Sprintf("CreateUser request parsed: %v", map[string]any{"username": createUserReq.Name}))

	// Validate request
	validationErrors := wscutils.WscValidate(createUserReq, func(err validator.FieldError) []string { return []string{} })
	if len(validationErrors) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, validationErrors))
		return
	}
	s.LogHarbour.Log(fmt.Sprintf("CreateUser request validated %v", map[string]any{"username": createUserReq.Name}))

	// Call CreateUser function

	// Resolve the database from the container
	dbObj, err := s.Container.Resolve("database")
	if err != nil {
		// Handle error
		return
	}

	// Assert that the dbObj implements the *alya.Queries interface
	db, ok := dbObj.(*sqlc.Queries)
	if !ok {
		// Handle error
		return
	}

	// Use db to call the CreateUser function
	user, err := db.CreateUser(c.Request.Context(), sqlc.CreateUserParams{
		Name:  createUserReq.Name,
		Email: createUserReq.Email,
	})
	if err != nil {
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(wscutils.ErrcodeDatabaseError))
		return
	}
	s.LogHarbour.Log(fmt.Sprintf("User created: %v", map[string]any{"username": createUserReq.Name}))

	// Send response
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(user))
}
