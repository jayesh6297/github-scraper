package githubscrape

const URL = "https://github.com"

const (
	// type filter
	Sources   = "source"
	Forks     = "forks"
	Archived  = "Archived"
	Mirrors   = "mirror"
	Templates = "Template"

	// language filter
	Go         = "go"
	HTML       = "html"
	Js         = "javascript"
	Java       = "java"
	Rust       = "rust"
	Python     = "python"
	TypeScript = "typescript"
	Css        = "css"
	Haskell    = "haskell"
	Shell      = "shell"
	All        = "all"

	// sort filter
	LastUpdated = ""
	Name        = "name"
	Stars       = "stargazers"
)
