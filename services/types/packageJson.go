package types

type PackageJSON struct {
	Name        string `json:"name,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`

	License string      `json:"license,omitempty"`
	Author  interface{} `json:"author,omitempty"`

	Type   string `json:"type,omitempty"`
	Main   string `json:"main,omitempty"`
	Module string `json:"module,omitempty"`
	Types  string `json:"types,omitempty"`

	Private bool `json:"private,omitempty"`

	Engines map[string]string `json:"engines,omitempty"`
	Scripts map[string]string `json:"scripts,omitempty"`

	Dependencies         map[string]string `json:"dependencies,omitempty"`
	DevDependencies      map[string]string `json:"devDependencies,omitempty"`
	PeerDependencies     map[string]string `json:"peerDependencies,omitempty"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`

	Keywords []string `json:"keywords,omitempty"`
	Files    []string `json:"files,omitempty"`

	Workspaces     interface{} `json:"workspaces,omitempty"`
	PackageManager string      `json:"packageManager,omitempty"`
}
