package ranking

import (
	"backend/types"
	"database/sql"
)

type RankingRepoImpl struct {
	db *sql.DB
}

func NewRankingRepoImpl(db *sql.DB) *RankingRepoImpl {
	return &RankingRepoImpl{db: db}
}

func (r *RankingRepoImpl) GetGlobalRanking() ([]types.RankEntry, error) {
	rows, err := r.db.Query(
		`SELECT  id, username, xp, rank_id,
		RANK() OVER (ORDER BY xp DESC)
		FROM users
		ORDER BY xp DESC;`,
	)
	if err != nil {
		return nil, err
	}
	list := make([]types.RankEntry, 0)
	for rows.Next() {
		var e types.RankEntry
		rows.Scan(&e.UserID, &e.Username, &e.XP, &e.RankId, &e.Rank)
		list = append(list, e)
	}
	return list, nil
}

func (r *RankingRepoImpl) GetUserGlobalRank(userID int) (*types.RankEntry, error) {
	row := r.db.QueryRow(`
        SELECT *
        FROM (
            SELECT id, username, xp, rank_id, RANK() OVER (ORDER BY xp DESC)
            FROM users
        ) AS ranked
        WHERE id = ?;`,
		userID,
	)

	var e types.RankEntry
	if err := row.Scan(&e.UserID, &e.Username, &e.XP, &e.RankId, &e.Rank); err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *RankingRepoImpl) GetLocalRanking(country string) ([]types.RankEntry, error) {
	rows, err := r.db.Query(
		`SELECT  id, username, xp, rank_id,
		RANK() OVER (ORDER BY xp DESC)
		FROM users where country = ?
		ORDER BY xp DESC;`,
		country,
	)
	if err != nil {
		return nil, err
	}
	list := make([]types.RankEntry, 0)
	for rows.Next() {
		var e types.RankEntry
		rows.Scan(&e.UserID, &e.Username, &e.XP, &e.RankId, &e.Rank)
		list = append(list, e)
	}
	return list, nil
}

func (r *RankingRepoImpl) GetUserLocalRank(userID int, country string) (*types.RankEntry, error) {
	row := r.db.QueryRow(`
        SELECT *
        FROM (
            SELECT id, username, xp, rank_id, RANK() OVER (ORDER BY xp DESC)
            FROM users where country = ?
        ) AS ranked
        WHERE id = ?;`,
		country,
		userID,
	)

	var e types.RankEntry
	if err := row.Scan(&e.UserID, &e.Username, &e.XP, &e.RankId, &e.Rank); err != nil {
		return nil, err
	}

	return &e, nil
}
