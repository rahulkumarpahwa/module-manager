package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProjectType string

const (
	NodeJS     ProjectType = "node.js"
	JavaScript ProjectType = "javascript"
	TypeScript ProjectType = "typescript"
	Golang     ProjectType = "golang"
	Java       ProjectType = "java"
	CPP        ProjectType = "cpp"
	C          ProjectType = "c"
	Python     ProjectType = "python"

	// React ecosystem
	ReactNative     ProjectType = "react-native"
	ExpoReactNative ProjectType = "expo-react-native"
	React           ProjectType = "react"
	NextJS          ProjectType = "next.js"

	// Vue ecosystem
	Vue     ProjectType = "vue"
	ViteVue ProjectType = "vite-vue"

	// Vite
	ViteReact ProjectType = "vite-react"
	Vite      ProjectType = "vite"

	// Angular
	Angular ProjectType = "angular"

	// HTML / Go templates
	HTML       ProjectType = "html"
	GoTemplate ProjectType = "go-template"

	Unknown ProjectType = ""
)

type ProjectInfo struct {
	Language  ProjectType
	IsProject bool
}

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

func IsProject(root, name string) bool {
	_, err := os.Stat(filepath.Join(root, name))
	return err == nil
}

func hasFile(root, pattern string) bool {
	matches, err := filepath.Glob(filepath.Join(root, pattern))
	return err == nil && len(matches) > 0
}

func readPackageJSON(root string) (*packageJSON, bool) {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return nil, false
	}

	var pkg packageJSON

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, false
	}

	return &pkg, true
}

func hasDependency(pkg *packageJSON, name string) bool {
	if pkg == nil {
		return false
	}

	if _, ok := pkg.Dependencies[name]; ok {
		return true
	}

	if _, ok := pkg.DevDependencies[name]; ok {
		return true
	}

	return false
}

func DetectProject(root string) *ProjectInfo {
	var language ProjectType

	// ============================================================
	// JavaScript / TypeScript ecosystem
	// ============================================================
	if IsProject(root, "package.json") {
		pkg, validPackageJSON := readPackageJSON(root)

		if validPackageJSON {

			// ----------------------------------------------------
			// Expo + React Native
			// ----------------------------------------------------
			if hasDependency(pkg, "expo") &&
				hasDependency(pkg, "react-native") {

				language = ExpoReactNative

				// ----------------------------------------------------
				// React Native
				// ----------------------------------------------------
			} else if hasDependency(pkg, "react-native") {

				language = ReactNative

				// ----------------------------------------------------
				// Next.js
				//
				// Next.js projects normally have:
				//     next
				//     react
				//     react-dom
				//
				// Checking "next" is sufficient.
				// ----------------------------------------------------
			} else if hasDependency(pkg, "next") {

				language = NextJS

				// ----------------------------------------------------
				// Angular
				//
				// Angular projects normally have:
				//     @angular/core
				//     @angular/common
				//
				// @angular/core is the primary marker.
				// ----------------------------------------------------
			} else if hasDependency(pkg, "@angular/core") {

				language = Angular

				// ----------------------------------------------------
				// Vite + React
				// ----------------------------------------------------
			} else if hasDependency(pkg, "vite") &&
				hasDependency(pkg, "react") {

				language = ViteReact

				// ----------------------------------------------------
				// Vite + Vue
				// ----------------------------------------------------
			} else if hasDependency(pkg, "vite") &&
				hasDependency(pkg, "vue") {

				language = ViteVue

				// ----------------------------------------------------
				// React
				// ----------------------------------------------------
			} else if hasDependency(pkg, "react") {

				language = React

				// ----------------------------------------------------
				// Vue
				// ----------------------------------------------------
			} else if hasDependency(pkg, "vue") {

				language = Vue

				// ----------------------------------------------------
				// Vite only
				// ----------------------------------------------------
			} else if hasDependency(pkg, "vite") {

				language = Vite

				// ----------------------------------------------------
				// Generic Node.js project
				// ----------------------------------------------------
			} else {

				language = NodeJS
			}

		} else {
			// package.json exists but couldn't be parsed.
			language = NodeJS
		}

		// ============================================================
		// Go
		// ============================================================
	} else if IsProject(root, "go.mod") {

		// Go template project
		if hasFile(root, "*.tmpl") ||
			hasFile(root, "*.gotmpl") ||
			IsProject(root, "templates") ||
			IsProject(root, "views") {

			language = GoTemplate

		} else {
			language = Golang
		}

		// ============================================================
		// Java
		// ============================================================
	} else if IsProject(root, "pom.xml") ||
		IsProject(root, "build.gradle") ||
		IsProject(root, "build.gradle.kts") {

		language = Java

		// ============================================================
		// Python
		// ============================================================
	} else if IsProject(root, "pyproject.toml") ||
		IsProject(root, "requirements.txt") ||
		IsProject(root, "setup.py") ||
		IsProject(root, "Pipfile") {

		language = Python

		// ============================================================
		// C++
		// ============================================================
	} else if hasFile(root, "*.cpp") ||
		hasFile(root, "*.cc") ||
		hasFile(root, "*.cxx") {

		language = CPP

		// ============================================================
		// C
		// ============================================================
	} else if hasFile(root, "*.c") {

		language = C

		// ============================================================
		// Normal HTML project
		// ============================================================
	} else if hasFile(root, "*.html") {

		// HTML + JS + CSS is treated as a normal HTML project.
		language = HTML

		// ============================================================
		// Unknown
		// ============================================================
	} else {

		language = Unknown
	}

	return &ProjectInfo{
		Language:  language,
		IsProject: language != Unknown,
	}
}
