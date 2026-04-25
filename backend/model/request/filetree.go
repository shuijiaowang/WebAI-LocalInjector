package request

type FileTreeRequest struct {
	RootPath    string   `json:"rootPath" binding:"required"`
	IgnoreDirs  []string `json:"ignoreDirs"`
	IgnoreFiles []string `json:"ignoreFiles"`
	IgnoreExts  []string `json:"ignoreExts"`
}
