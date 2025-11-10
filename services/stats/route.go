package stats

import (
	"backend/types"
	"backend/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.StatsRepo
}

func NewHandler(store types.StatsRepo) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/stats/heatmap", h.GetUserHeatMap).Methods(http.MethodGet)
	router.HandleFunc("/stats/heatmap", h.UpsertUserHeatmapEntry).Methods(http.MethodPut)
	router.HandleFunc("/stats/{id}", h.GetUserStats).Methods(http.MethodGet)
	router.HandleFunc("/stats", h.UpdateUserStats).Methods(http.MethodPut)
	router.HandleFunc("/stats", h.AddUserStats).Methods(http.MethodPost)
}

// GetUserStats docs
//
// @Summary 			Get user statistics
// @Description 		Get statistics for a user by ID
// @Tags 				stats
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Success 			200 {object} types.ExtendedStats
// @Failure 			401 {object} types.ErrorResponse
// @Failure 			500 {object} types.ErrorResponse
// @Param 				id path string true "User ID"
// @Router 				/stats/{id} [get]
func (h *Handler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	userCompleteStats, err := h.store.GetUserStats(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, userCompleteStats)
}

// UpdateUserStats docs
//
// @Summary 			Update user statistics
// @Description 		Update statistics for a user by ID
// @Tags 				stats
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Success 			200 {object} types.Stats
// @Failure 			400 {object} types.ErrorResponse
// @Failure 			401 {object} types.ErrorResponse
// @Failure 			500 {object} types.ErrorResponse
// @Param 				stats body types.StatsPayload true "Stats update payload"
// @Router 				/stats [put]
func (h *Handler) UpdateUserStats(w http.ResponseWriter, r *http.Request) {
	var payload types.StatsPayload
	var stats types.Stats
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// copy the payload into the stats
	stats.UserID = payload.UserID
	stats.LongestStreak = payload.LongestStreak
	stats.CurrentStreak = payload.CurrentStreak
	// set the last updated to the current time
	stats.LastUpdated = time.Now()

	_, err := h.store.UpdateUserStats(&stats)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, stats)
}

// AddUserStats docs
//
// @Summary 			Add user statistics
// @Description 		Add statistics for a user by ID
// @Tags 				stats
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Success 			201 {object} types.Stats
// @Failure 			400 {object} types.ErrorResponse
// @Failure 			401 {object} types.ErrorResponse
// @Failure 			500 {object} types.ErrorResponse
// @Param 				stats body types.Stats true "Stats update payload"
// @Router 				/stats [post]
func (h *Handler) AddUserStats(w http.ResponseWriter, r *http.Request) {
	var payload types.StatsPayload
	var stats types.Stats
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// copy values of payload into stats
	stats.UserID = payload.UserID
	stats.CurrentStreak = payload.CurrentStreak
	stats.LongestStreak = payload.LongestStreak
	stats.CreatedAt = time.Now()
	stats.LastUpdated = time.Now()
	_, err := h.store.AddUserStats(&stats)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, stats)
}

// GetUserHeatMap docs
//
// @Summary 			Get user heatmap data
// @Description 		Retrieve heatmap data for a user within a specified date range. Returns daily pomodoro counts including days with zero activity.
// @Tags 				stats
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Success 			200 {object} types.HeatMapResponse "Heatmap data with user ID and daily counts"
// @Failure 			400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 			401 {object} types.ErrorResponse "Unauthorized"
// @Failure 			500 {object} types.ErrorResponse "Internal server error"
// @Param 				payload body types.HeatMapPayload true "Heatmap query parameters including user ID and date range"
// @Router 				/stats/heatmap [get]
// @Example 			Request body: {"user_id": 1, "start_date": "2024-01-01T00:00:00Z", "end_date": "2024-01-31T23:59:59Z"}
// @Example 			Response: {"user_id": 1, "data": [{"date": "2024-01-01T00:00:00Z", "count": 5}, {"date": "2024-01-02T00:00:00Z", "count": 0}]}
func (h *Handler) GetUserHeatMap(w http.ResponseWriter, r *http.Request) {
	var payload types.HeatMapPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	rows, err := h.store.GetUserHeatmap(&payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(
		w,
		http.StatusOK,
		types.HeatMapResponse{UserID: payload.UserID, Data: rows},
	)
}

// UpsertUserHeatmapEntry docs
//
// @Summary 			Create or update heatmap entry
// @Description 		Create a new heatmap entry for a user on a specific date, or update the count if an entry already exists. Uses UPSERT operation.
// @Tags 				stats
// @Accept 				json
// @Produce 			json
// @Security 			ApiKeyAuth
// @Success 			200 {object} types.HeatMap "Updated heatmap entry"
// @Failure 			400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 			401 {object} types.ErrorResponse "Unauthorized"
// @Failure 			500 {object} types.ErrorResponse "Internal server error"
// @Param 				payload body types.HeatMap true "Heatmap entry with user ID, date, and pomodoro count"
// @Router 				/stats/heatmap [put]
// @Example 			Request body: {"user_id": 1, "date": "2024-01-15T00:00:00Z", "count": 3}
// @Example 			Response: {"user_id": 1, "date": "2024-01-15T00:00:00Z", "count": 3}
func (h *Handler) UpsertUserHeatmapEntry(w http.ResponseWriter, r *http.Request) {
	var payload types.HeatMap
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.store.UpsertUserHeatmapEntry(&payload); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, payload)

}
