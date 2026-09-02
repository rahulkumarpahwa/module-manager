package types

import "time"

type FirstRunEntry struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	ModifiedDate time.Time `json:"modified_date"`
	CreationDate time.Time `json:"creation_date"`
	IsWorkspace  bool      `json:"workspace,omitempty"`
	IsProject    bool      `json:"project,omitempty"`
	Language     string    `json:"framework/language,omitempty"`
}

var FirstRunEntriesStore = []FirstRunEntry{}
