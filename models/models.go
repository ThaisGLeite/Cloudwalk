package models

type Player struct {
	Name  string `json:"name"`
	Kills int    `json:"kills"`
}

type Game struct {
	GameNumber   int                `json:"game_number"`
	TotalKills   int                `json:"total_kills"`
	Players      map[string]*Player `json:"players"`
	KillsByMeans map[string]int     `json:"kills_by_means"`
}
