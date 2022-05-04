package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/JonathanGzzBen/nutrity-api/api/v1/models"
	"github.com/JonathanGzzBen/nutrity-api/api/v1/repository"
	"github.com/gin-gonic/gin"
)

type UpdateUserDTO struct {
	Username          string   `json:"username"`
	Email             string   `json:"email"`
	FirstName         string   `json:"firstname"`
	LastName          string   `json:"lastname"`
	UserProfileEdited bool     `json:"userProfileEdited"`
	Calories          uint     `json:"calories"`
	Carbs             uint     `json:"carbs"`
	Day               uint     `json:"day"`
	Fats              uint     `json:"fats"`
	Proteins          uint     `json:"proteins"`
	RecipesAdded      []string `json:"recipesAdded"`
}

type UserDTO struct {
	ID                uint     `json:"id,omitempty"`
	Username          string   `json:"username"`
	Email             string   `json:"email"`
	FirstName         string   `json:"firstname"`
	LastName          string   `json:"lastname"`
	UserProfileEdited bool     `json:"userProfileEdited"`
	Calories          uint     `json:"calories"`
	Carbs             uint     `json:"carbs"`
	Day               uint     `json:"day"`
	Fats              uint     `json:"fats"`
	Proteins          uint     `json:"proteins"`
	RecipesAdded      []string `json:"recipesAdded"`
}

func userDTOFromUser(u *models.User) UserDTO {
	return UserDTO{
		ID:                u.ID,
		Username:          u.Username,
		Email:             u.Email,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		UserProfileEdited: u.UserProfileEdited,
		Calories:          u.Calories,
		Carbs:             u.Carbs,
		Day:               u.Day,
		Fats:              u.Fats,
		Proteins:          u.Proteins,
		RecipesAdded:      strings.Split(u.RecipesAdded, "^"),
	}
}

func userFromUserDTO(uDTO *UserDTO) models.User {
	// Convert array to single string separating elements using character '^'
	recipesAddedString := ""
	for _, r := range uDTO.RecipesAdded {
		recipesAddedString += r + "^"
	}
	recipesAddedString = recipesAddedString[:len(recipesAddedString)-1] // Remove trailing '^'

	return models.User{
		ID:                uDTO.ID,
		Username:          uDTO.Username,
		Email:             uDTO.Email,
		FirstName:         uDTO.FirstName,
		LastName:          uDTO.LastName,
		UserProfileEdited: uDTO.UserProfileEdited,
		Calories:          uDTO.Calories,
		Carbs:             uDTO.Carbs,
		Day:               uDTO.Day,
		Fats:              uDTO.Fats,
		Proteins:          uDTO.Proteins,
		RecipesAdded:      recipesAddedString,
	}
}

// GetAllUsers is the handler for GET requests to /users
// 	@ID GetAllUsers
// 	@Summary Get all users
// 	@Description Get all registered users.
// 	@Tags users
// 	@Success 200 {array} UserDTO
// 	@Failure 500 {object} models.APIError
// 	@Router /users [get]
func (s *Server) GetAllUsers(c *gin.Context) {
	users, err := s.UsersRepo.GetAllUsers()
	userDTOs := make([]UserDTO, len(users))
	for _, u := range users {
		userDTOs = append(userDTOs, userDTOFromUser(&u))
	}
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
// 	@Success 200 {object} UserDTO
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
	c.JSON(http.StatusOK, userDTOFromUser(user))
}

// UpdateUser is the handler for PUT requests to /users/:id
// 	@ID UpdateUser
// 	@Summary Update user
// 	@Description Update matching user with provided data.
// 	@Tags users
// 	@Security AccessToken
// 	@Param id path int true "User ID"
// 	@Param user body UpdateUserDTO true "User"
// 	@Success 200 {object} UserDTO
// 	@Failure 400 {object} models.APIError
// 	@Router /users/{id} [put]
func (s *Server) UpdateUser(c *gin.Context) {
	at := c.GetHeader(AccessTokenName)
	au, err := s.UsersRepo.GetUserByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "not authenticated: " + err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "id is not a valid"})
		return
	}
	if au.ID != uint(id) {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "id does not match authenticated user"})
		return
	}
	var uu UpdateUserDTO
	if err := c.ShouldBindJSON(&uu); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "invalid update user: " + err.Error()})
		return
	}

	// User is updating his own information
	u, err := s.UsersRepo.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "not registered user"})
		return
	}

	// Convert array to single string separating elements using character '^'
	recipesAddedString := ""
	for _, r := range uu.RecipesAdded {
		recipesAddedString += r + "^"
	}
	recipesAddedString = recipesAddedString[:len(recipesAddedString)-1] // Remove trailing '^'

	u.Username = uu.Username
	u.Email = uu.Email
	u.FirstName = uu.FirstName
	u.LastName = uu.LastName
	u.UserProfileEdited = uu.UserProfileEdited
	u.Calories = uu.Calories
	u.Carbs = uu.Carbs
	u.Day = uu.Day
	u.Fats = uu.Fats
	u.Proteins = uu.Proteins
	u.RecipesAdded = recipesAddedString

	u, err = s.UsersRepo.UpdateUser(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, userDTOFromUser(u))
}
