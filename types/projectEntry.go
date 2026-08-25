package types

type ProjectEntry struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsDir     bool   `json:"is_dir"`
	TotalSize int64  `json:"total_size"`
}
