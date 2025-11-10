package types

type RankEntry struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Rank     int    `json:"rank"`
    XP       int    `json:"xp"`
    RankId   int    `json:"rank_id"`
}

type RankingRepo interface {

    GetGlobalRanking() ([]RankEntry, error)
    GetUserGlobalRank(userID int) (*RankEntry, error)

    GetLocalRanking(country string) ([]RankEntry, error)
    GetUserLocalRank(userID int, country string) (*RankEntry, error)
}

