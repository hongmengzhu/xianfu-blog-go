package modApiBmBookmarkCategory

type VoAll struct {
	My   []Vo `json:"my"`
	Team []Vo `json:"team"`
	Org  []Vo `json:"org"`
}
