package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/JonathanGzzBen/ingenialists/api/v1/repository"
	"github.com/gin-gonic/gin"
)

type UpdateUserDTO struct {
	Name              string      `json:"name" binding:"required"`
	Birthdate         time.Time   `json:"birthdate" example:"2006-01-02T15:04:05Z"`
	Gender            string      `json:"gender"`
	ProfilePictureURL string      `json:"profilePictureUrl"`
	Description       string      `json:"description"`
	ShortDescription  string      `json:"shortDescription"`
	Role              models.Role `json:"role" example:"Reader"`
}

// GetAllUsers is the handler for GET requests to /users
// 	@ID GetAllUsers
// 	@Summary Get all users
// 	@Description Get all registered users.
// 	@Tags users
// 	@Success 200 {array} models.User
// 	@Failure 500 {object} models.APIError
// 	@Router /users [get]
func (s *Server) GetAllUsers(c *gin.Context) {
	users, err := s.UsersRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{Code: http.StatusInternalServerError, Message: "could not connect to database"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser is the handler for GET requests to /users/:id
// 	@ID GetUser
// 	@Summary Get user
// 	@Description Get user with matching ID.
// 	@Tags users
// 	@Param id path int true "User ID"
// 	@Success 200 {object} models.User
// 	@Failure 404 {object} models.APIError
// 	@Router /users/{id} [get]
func (s *Server) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: "invalid id: " + err.Error()})
		return
	}
	user, err := s.UsersRepo.GetUser(uint(id))
	if err == repository.ErrNotFound {
		c.JSON(http.StatusNotFound, models.APIError{Code: http.StatusNotFound, Message: "user with provided id not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser is the handler for PUT requests to /users/:id
// 	@ID UpdateUser
// 	@Summary Update user
// 	@Description Update matching user with provided data.
// 	@Tags users
// 	@Security AccessToken
// 	@Param id path int true "User ID"
// 	@Param user body UpdateUserDTO true "User"
// 	@Success 200 {object} models.User
// 	@Failure 400 {object} models.APIError
// 	@Router /users/{id} [put]
func (s *Server) UpdateUser(c *gin.Context) {
	at := c.GetHeader(AccessTokenName)
	au, err := s.userByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "not authenticated: " + err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "id is not a valid"})
		return
	}
	if au.ID != uint(id) && au.Role != models.RoleAdministrator {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "id does not match authenticated user"})
		return
	}
	var uu UpdateUserDTO
	if err := c.ShouldBindJSON(&uu); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "invalid update user: " + err.Error()})
		return
	}
	// If administrator is updating other user
	if au.Role == models.RoleAdministrator && uint(id) != au.ID {
		u, err := s.UsersRepo.GetUser(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: err.Error()})
			return
		}
		// Administrators can only change Role of other users
		u.Role = uu.Role
		u, err = s.UsersRepo.UpdateUser(u)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "could not update user: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, u)
		return
	}
	// User is updating his own information
	u, err := s.UsersRepo.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "not registered user"})
		return
	}
	u.Name = uu.Name
	u.Birthdate = uu.Birthdate
	u.Gender = uu.Gender
	u.ProfilePictureURL = uu.ProfilePictureURL
	u.Description = uu.Description
	u.ShortDescription = uu.ShortDescription
	u, err = s.UsersRepo.UpdateUser(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}
