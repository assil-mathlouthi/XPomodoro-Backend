package ranking

import (
	"backend/types"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.RankingRepo
}

func NewHandler(store types.RankingRepo) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/ranking/global", h.GetGlobalRanking).Methods(http.MethodGet)
	router.HandleFunc("/ranking/global/{id}", h.GetUserGlobalRanking).Methods(http.MethodGet)
	router.HandleFunc("/ranking/{country}", h.GetLocalRanking).Methods(http.MethodGet)
	router.HandleFunc("/ranking/{country}/{id}", h.GetUserLocalRanking).Methods(http.MethodGet)
}

// GetGlobalRanking docs
//
// @Summary             Get global ranking
// @Description         Retrieve the global ranking list ordered by XP descending
// @Tags                ranking
// @Accept              json
// @Produce             json
// @Security 			ApiKeyAuth
// @Success             200 {array} types.RankEntry
// @Failure             500 {object} types.ErrorResponse
// @Router              /ranking/global [get]
func (h *Handler) GetGlobalRanking(w http.ResponseWriter, r *http.Request) {

	rankingList, err := h.store.GetGlobalRanking()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, rankingList)
}

// GetLocalRanking docs
//
// @Summary             Get local ranking by country
// @Description         Retrieve ranking for users in a specific country ordered by XP descending
// @Tags                ranking
// @Accept              json
// @Produce             json
// @Security 			ApiKeyAuth
// @Param               country path string true "Country code"
// @Success             200 {array} types.RankEntry
// @Failure             500 {object} types.ErrorResponse
// @Router              /ranking/{country} [get]
func (h *Handler) GetLocalRanking(w http.ResponseWriter, r *http.Request) {

	rankingList, err := h.store.GetLocalRanking(mux.Vars(r)["country"])
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, rankingList)
}

// GetUserGlobalRanking docs
//
// @Summary             Get a user's global rank
// @Description         Retrieve the global rank entry for a specific user ID
// @Tags                ranking
// @Accept              json
// @Produce             json
// @Security 			ApiKeyAuth
// @Param               id path int true "User ID"
// @Success             200 {object} types.RankEntry
// @Failure             500 {object} types.ErrorResponse
// @Router              /ranking/global/{id} [get]
func (h *Handler) GetUserGlobalRanking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	rankEntry, err := h.store.GetUserGlobalRank(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, rankEntry)
}

// GetUserLocalRanking docs
//
// @Summary             Get a user's local rank
// @Description         Retrieve the local rank entry for a specific user ID within a country
// @Tags                ranking
// @Accept              json
// @Produce             json
// @Security 			ApiKeyAuth
// @Param               country path string true "Country code"
// @Param               id path int true "User ID"
// @Success             200 {object} types.RankEntry
// @Failure             500 {object} types.ErrorResponse
// @Router              /ranking/{country}/{id} [get]
func (h *Handler) GetUserLocalRanking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	country := vars["country"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	rankEntry, err := h.store.GetUserLocalRank(id, country)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, rankEntry)
}
