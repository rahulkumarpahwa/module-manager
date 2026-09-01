package utils

import (
	"encoding/json"
	"os"

	"github.com/rahulkumarpahwa/module-manager/services/types"
)

type PackageDependency struct {
	PackageManger string
	FilePath      string
}

func NewPackageDependency(packageManager string, filePath string) *PackageDependency {
	return &PackageDependency{
		PackageManger: packageManager,
		FilePath:      filePath,
	}
}

func (p *PackageDependency) PackageDependencyParser() (*types.PackageJSON, error) {
	pj := &types.PackageJSON{}
	fileByte, err := os.ReadFile(p.FilePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileByte, pj)
	if err != nil {
		return nil, err
	}
	return pj, err
}

// todo: here we are only parsing the package.json for the package manager 'npm'. we will later create the parser for the bun, yarn, pnpm with their respective types.