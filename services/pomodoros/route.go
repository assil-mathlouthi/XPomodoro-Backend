package pomodoros

import (
	"backend/types"
	"backend/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.PomodoroRepo
}

func NewHandler(store types.PomodoroRepo) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/pomodoro", h.HandleAddingPomodoro).Methods(http.MethodPost)
}

//	 HandleAddingPomodoro godoc
//
//		@Summary 			Add a new pomodoro session
//		@Description 		Create a new pomodoro session for a user
//		@Tags 				pomodoros
//		@Accept 			json
//		@Produce 			json
//		@Security 			ApiKeyAuth
//		@Success 			201 {object} types.Pomodoro
//		@Failure 			400 {object} types.ErrorResponse
//		@Failure 			401 {object} types.ErrorResponse
//		@Failure 			500 {object} types.ErrorResponse
//		@Param 				request body types.AddingPomodoroPayload true "Pomodoro request payload"
//	 @Router 			/pomodoro [post]
func (h *Handler) HandleAddingPomodoro(w http.ResponseWriter, r *http.Request) {
	// get your payload
	var payload types.AddingPomodoroPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	pomodoro, err := h.store.AddPomodoro(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, pomodoro)
}
