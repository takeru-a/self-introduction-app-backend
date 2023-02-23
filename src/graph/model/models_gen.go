// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type FetchRoom struct {
	ID string `json:"id"`
}

type FetchUser struct {
	ID string `json:"id"`
}

type LoginedUser struct {
	UserID    string `json:"userId"`
	Name      string `json:"name"`
	RoomID    string `json:"roomId"`
	RoomToken string `json:"roomToken"`
	Code      int    `json:"code"`
}

type NewRoom struct {
	HostName string `json:"host_name"`
}

type NewUser struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Room struct {
	ID      string  `json:"id"`
	Host    *User   `json:"host"`
	Token   string  `json:"token"`
	Players []*User `json:"players"`
}

type Subscriber struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Answer string `json:"answer"`
	Score  int    `json:"score"`
}
