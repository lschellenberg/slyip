package cmd

type SetNameRequest struct {
	Name string `json:"name"`
}

type SetNameResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ItemsResponse []Items

type Items struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ToggleFemale bool   `json:"toggle_female"`
	ToggleMale   bool   `json:"toggle_male"`
}

type GetUsersResponse struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Total int    `json:"total"`
	Data  []User `json:"data"`
}

type Choice struct {
	ItemId   string `json:"item_id"`
	LikeType int    `json:"like_type"`
	UserId   string `json:"user_id"`
}

type GetChoices struct {
	Choices []Choice `json:"choices"`
}
