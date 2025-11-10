package user

import (
	"backend/helpers"
	"backend/services/auth"
	"backend/types"
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserRepo
}

// simulate the constructor in others languages
func NewHandler(store types.UserRepo) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router, authRouter *mux.Router) {
	// Public routes
	router.HandleFunc("/login", h.HandleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.HandleRegister).Methods(http.MethodPost)
	router.HandleFunc("/verify", h.HandleVerifyEmail).Methods(http.MethodGet)
	router.HandleFunc("/password/forgot", h.HandleForgotPassword).Methods(http.MethodPost)
	router.HandleFunc("/password/reset", h.HandleVerifyResetCode).Methods(http.MethodPost)

	// Protected routes
	authRouter.HandleFunc("/users/email", h.HandleUpdateEmail).Methods(http.MethodPut)
	authRouter.HandleFunc("/users/{id}/country", h.HandleUpdateCountry).Methods(http.MethodPatch)
}

//		HandleLogin godoc
//
//		@Summary 			Login a user
//		@Description 		Authenticate a user and return a JWT token
//		@Tags 				Auth
//		@Accept 			json
//		@Produce 			json
//		@Success 			200 {object} types.TokenResponse
//		@Failure 			400 {object} types.ErrorResponse
//		@Failure 			404 {object} types.ErrorResponse
//		@Failure 			500 {object} types.ErrorResponse
//		@Param 				request body types.AuthPayload true "Auth payload"
//	 	@Router 			/login [post]
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.AuthPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// see if the username exist
	user, err := h.store.GetUserByUsername(payload.Username)
	if err != nil {
		log.Fatal(err)
	}
	if user == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user with the specified username was not found"))
		return
	}
	// now that the user exist we need to compare password hash with the one stored in the database
	if !auth.ComparePasswords(user.PasswordHash, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid password"))
		return
	}
	// user + password => valid
	// generate jwt tocken now and return it to the client
	token, err := auth.CreateToken(payload.Username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.TokenResponse{Token: token})

}

//		HandleRegister godoc
//
//		@Summary 			Register a user
//		@Description 		Create a new user account
//		@Tags 				Auth
//		@Accept 			json
//		@Produce 			json
//		@Success 			201 {object} types.User
//		@Failure 			400 {object} types.ErrorResponse
//		@Failure 			500 {object} types.ErrorResponse
//		@Param 				request body types.AuthPayload true "Auth payload"
//	 	@Router 			/register [post]
func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.AuthPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// see if the user already exists (username must be unique)
	user, err := h.store.GetUserByUsername(payload.Username)
	if err != nil {
		log.Fatal(err)
	}
	// if exist return user exist
	if user != nil {
		// utils.WriteJSON(w, http.StatusOK, user)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exist"))
		return
	}
	// hash the password
	PasswordHash, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// by default users start with 0 xp and rank id 1 (Wood I)
	err = h.store.CreateUser(types.User{
		Username:     payload.Username,
		PasswordHash: PasswordHash,
		XP:           0,
		RankId:       1,
		Email:        nil,
		Country:      nil,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	user, err = h.store.GetUserByUsername(payload.Username)
	if err != nil {
		log.Fatal(err)
	}
	utils.WriteJSON(w, http.StatusCreated, user)
}

// HandleUpdateEmail godoc
//
// @Summary 			Request email update (sends verification email)
// @Description 		Validates new email, generates token, stores pending update, and sends verification email
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Param 				request body types.UpdateEmailPayload true "Update email payload"
// @Success 			200 {object} types.SuccessResponse
// @Failure 			400 {object} types.ErrorResponse
// @Failure 			500 {object} types.ErrorResponse
// @Router 			/users/email [put]
func (h *Handler) HandleUpdateEmail(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateEmailPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	id, err := strconv.Atoi(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.UpdateUserEmail(id, payload.NewEmail); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message: "Email send successfully",
	})
}

// HandleVerifyEmail godoc
//
// @Summary 			Verify email update via token
// @Description 		Confirms pending email update by token and updates the user's email
// @Tags 				User
// @Produce 			json
// @Security 			ApiKeyAuth
// @Param 				token query string true "Verification token"
// @Success 			200 {object} types.SuccessResponse
// @Failure 			400 {object} types.ErrorResponse
// @Router 				/verify [get]
func (h *Handler) HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	// token comes from query param: /verify?token=...
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing token"))
		return
	}
	if err := h.store.VerifyEmailUpdate(token); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.SuccessResponse{Message: "Email verified successfully"})
}

// HandleUpdateCountry godoc
//
// @Summary 			Update a user's country
// @Description 		Updates the country for a given user id
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Param 				id path string true "User ID"
// @Param 				request body types.UserInfoUpdate true "Country payload"
// @Success 			200 {object} types.User
// @Failure 			400 {object} types.ErrorResponse
// @Failure 			404 {object} types.ErrorResponse
// @Failure 			500 {object} types.ErrorResponse
// @Router 				/users/{id}/country [patch]
func (h *Handler) HandleUpdateCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var payload types.UserInfoUpdate
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.store.UpdateUserCountry(id, payload.Country)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)

}

// HandleForgotPassword godoc
//
// @Summary 			Request password reset code
// @Description 		Send a password reset code to the user's email if they have one associated with their account
// @Tags 				Auth
// @Accept 				json
// @Produce 			json
// @Param 				request body types.ForgotPasswordPayload true "Forgot password payload"
// @Success 			200 {object} types.SuccessResponse
// @Failure 			400 {object} types.ErrorResponse
// @Failure 			500 {object} types.ErrorResponse
// @Router 				/password/forgot [post]
func (h *Handler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var payload types.ForgotPasswordPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// get the user data
	user, err := h.store.GetUserByUsername(payload.Username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// see if the user has an email associated with his account
	if user.Email == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("this user account does not have an email associated"))
		return
	}
	// generate code , add entry and send email
	code := utils.GenerateRandomCode(8)
	if err := h.store.RequestPasswordReset(user.Id, code); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// send email with the code
	if err := helpers.SendPasswordResetCode(*user.Email, code); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.SuccessResponse{Message: "Code sent successfully"})

}

// HandleVerifyResetCode godoc
//
// @Summary 			Reset password with code
// @Description 		Reset user password using the provided reset code
// @Tags 				Auth
// @Accept 				json
// @Produce 			json
// @Param 				request body types.ResetPasswordPayload true "Reset password payload"
// @Success 			200 {object} types.SuccessResponse
// @Failure 			400 {object} types.ErrorResponse
// @Failure 			500 {object} types.ErrorResponse
// @Router 				/password/reset [post]
func (h *Handler) HandleVerifyResetCode(w http.ResponseWriter, r *http.Request) {
	var payload types.ResetPasswordPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.store.GetUserByUsername(payload.Username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.store.ResetPasswordWithCode(user.Id, payload.Code, payload.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.SuccessResponse{Message: "Password updated successfully"})
}
