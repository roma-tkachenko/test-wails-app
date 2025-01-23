package service

import "encoding/json"

type RewardCards struct {
	Cards       json.RawMessage `json:"cards"`
	IfReward    string          `json:"if_reward"`
	RewardLimit int             `json:"reward_limit"`
	StopReward  string          `json:"stop_reward"`
}

type Card struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Author     string `json:"author"`
	NewsID     string `json:"news_id"`
	Rank       string `json:"rank"`
	Image      string `json:"image"`
	Source     string `json:"source"`
	Approve    string `json:"approve"`
	Reward     string `json:"reward"`
	Comment    string `json:"comment"`
	NewAuthor  string `json:"new_author"`
	OwnerID    int    `json:"owner_id"`
	AnimeName  string `json:"anime_name"`
	AnimeLink  string `json:"anime_link"`
	CanTrade   string `json:"can_trade"`
	IsFavorite string `json:"is_favorite"`
}
