package service

import (
	"Service/model/response"
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileTreeService struct{}

type IgnoreFilter struct {
	Dirs  map[string]bool
	Files map[string]bool
	Exts  map[string]bool
}

type contentTreeNode struct {
	Name       string
	Children   []*contentTreeNode
	childIndex map[string]*contentTreeNode
}

const (
	contentFilterGoFunctionBody = "goFunctionBody"
	contentFilterGoImport       = "goImport"
	contentFilterJsFunctionBody = "jsFunctionBody"
	contentFilterVueStyle       = "vueStyle"
)

var (
	goImportBlockRegexp       = regexp.MustCompile(`(?ms)^\s*import\s*\(.*?^\s*\)\s*\n?`)
	goImportLineRegexp        = regexp.MustCompile(`(?m)^\s*import\s+(?:\S+\s+)?"[^"]+"\s*\n?`)
	goFunctionBodyRegexp      = regexp.MustCompile(`(?ms)^(\s*func\s+(?:\([^{}\n]*\)\s*)?[A-Za-z_]\w*\s*\([^{}]*\)\s*(?:\([^{}]*\)|[^{\n]*)?)\s*\{.*?^\}`)
	jsFunctionBodyRegexp      = regexp.MustCompile(`(?ms)^(\s*(?:export\s+)?(?:async\s+)?function\s+[$A-Za-z_][$\w]*\s*\([^)]*\)\s*)\{.*?^\}`)
	jsArrowFunctionBodyRegexp = regexp.MustCompile(`(?ms)^(\s*(?:export\s+)?(?:const|let|var)\s+[$A-Za-z_][$\w]*\s*=\s*(?:async\s*)?(?:\([^)]*\)|[$A-Za-z_][$\w]*)\s*=>\s*)\{.*?^\}`)
	vueStyleRegexp            = regexp.MustCompile(`(?is)<style(?:\s[^>]*)?>.*?</style>\s*`)
)

func (s *FileTreeService) GetFileTree(rootPath string, ignoreDirs, ignoreFiles, ignoreExts []string) ([]*response.FileNode, error) {
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, os.ErrNotExist
	}

	filter := s.newIgnoreFilter(ignoreDirs, ignoreFiles, ignoreExts)

	var result []*response.FileNode
	entries, err := os.ReadDir(absRoot)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if s.shouldIgnore(entry.Name(), entry.IsDir(), filter) {
			continue
		}
		node := s.buildNode(absRoot, entry.Name(), filter)
		if node != nil {
			result = append(result, node)
		}
	}

	return result, nil
}

func (s *FileTreeService) GetFileContent(rootPath string, ignoreDirs, ignoreFiles, ignoreExts, selectedPaths, contentFilters []string) (string, error) {
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", os.ErrNotExist
	}

	filter := s.newIgnoreFilter(ignoreDirs, ignoreFiles, ignoreExts)
	contentFilter := s.newContentFilter(contentFilters)
	selected := make(map[string]bool)
	for _, path := range selectedPaths {
		selected[filepath.ToSlash(filepath.Clean(path))] = true
	}

	contentRoot := &contentTreeNode{childIndex: make(map[string]*contentTreeNode)}
	var contentBuilder strings.Builder
	err = filepath.WalkDir(absRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == absRoot {
			return nil
		}

		relPath, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}
		outputPath := relPath
		selectedPath := filepath.ToSlash(relPath)

		if s.shouldIgnore(entry.Name(), entry.IsDir(), filter) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		s.addContentTreePath(contentRoot, relPath)

		if entry.IsDir() || !selected[selectedPath] {
			return nil
		}

		contentBuilder.WriteString(outputPath)
		contentBuilder.WriteString("\n")

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content = []byte(s.filterContent(outputPath, string(content), contentFilter))
		contentBuilder.Write(content)
		if !bytes.HasSuffix(content, []byte("\n")) {
			contentBuilder.WriteString("\n")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString("【目录结构】\n")
	s.writeContentTree(&builder, contentRoot, 0)
	builder.WriteString("\n【文件内容】\n")
	builder.WriteString(contentBuilder.String())

	return builder.String(), nil
}

func (s *FileTreeService) newContentFilter(contentFilters []string) map[string]bool {
	filter := make(map[string]bool)
	for _, item := range contentFilters {
		filter[item] = true
	}
	return filter
}

func (s *FileTreeService) filterContent(path, content string, filter map[string]bool) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		if filter[contentFilterGoImport] {
			content = goImportBlockRegexp.ReplaceAllString(content, "")
			content = goImportLineRegexp.ReplaceAllString(content, "")
		}
		if filter[contentFilterGoFunctionBody] {
			content = goFunctionBodyRegexp.ReplaceAllString(content, "$1 {}")
		}
	case ".js", ".mjs", ".cjs":
		if filter[contentFilterJsFunctionBody] {
			content = jsFunctionBodyRegexp.ReplaceAllString(content, "$1{}")
			content = jsArrowFunctionBodyRegexp.ReplaceAllString(content, "$1{}")
		}
	case ".vue":
		if filter[contentFilterVueStyle] {
			content = vueStyleRegexp.ReplaceAllString(content, "")
		}
	}
	return content
}

func (s *FileTreeService) newIgnoreFilter(ignoreDirs, ignoreFiles, ignoreExts []string) *IgnoreFilter {
	filter := &IgnoreFilter{
		Dirs:  make(map[string]bool),
		Files: make(map[string]bool),
		Exts:  make(map[string]bool),
	}
	for _, d := range ignoreDirs {
		filter.Dirs[d] = true
	}
	for _, f := range ignoreFiles {
		filter.Files[f] = true
	}
	for _, e := range ignoreExts {
		ext := e
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		filter.Exts[strings.ToLower(ext)] = true
	}
	return filter
}

func (s *FileTreeService) shouldIgnore(name string, isDir bool, filter *IgnoreFilter) bool {
	if isDir {
		return filter.Dirs[name]
	}
	if filter.Files[name] {
		return true
	}
	ext := strings.ToLower(filepath.Ext(name))
	return filter.Exts[ext]
}

func (s *FileTreeService) addContentTreePath(root *contentTreeNode, relPath string) {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	current := root
	for _, part := range parts {
		if part == "" {
			continue
		}
		if current.childIndex == nil {
			current.childIndex = make(map[string]*contentTreeNode)
		}
		child := current.childIndex[part]
		if child == nil {
			child = &contentTreeNode{Name: part}
			current.childIndex[part] = child
			current.Children = append(current.Children, child)
		}
		current = child
	}
}

func (s *FileTreeService) writeContentTree(builder *strings.Builder, node *contentTreeNode, depth int) {
	for _, child := range node.Children {
		builder.WriteString(strings.Repeat(" ", depth))
		builder.WriteString(child.Name)
		builder.WriteString("\n")
		s.writeContentTree(builder, child, depth+1)
	}
}

func (s *FileTreeService) buildNode(parentPath, name string, filter *IgnoreFilter) *response.FileNode {
	fullPath := filepath.Join(parentPath, name)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil
	}

	relPath, _ := filepath.Rel(parentPath, fullPath)

	node := &response.FileNode{
		Name:  name,
		Path:  filepath.ToSlash(relPath),
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return node
		}
		for _, entry := range entries {
			if s.shouldIgnore(entry.Name(), entry.IsDir(), filter) {
				continue
			}
			child := s.buildNodeRecursive(fullPath, entry.Name(), relPath, filter)
			if child != nil {
				node.Children = append(node.Children, child)
			}
		}
	}

	return node
}

func (s *FileTreeService) buildNodeRecursive(parentFullPath, name, parentRelPath string, filter *IgnoreFilter) *response.FileNode {
	fullPath := filepath.Join(parentFullPath, name)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil
	}

	relPath := filepath.Join(parentRelPath, name)

	node := &response.FileNode{
		Name:  name,
		Path:  filepath.ToSlash(relPath),
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return node
		}
		for _, entry := range entries {
			if s.shouldIgnore(entry.Name(), entry.IsDir(), filter) {
				continue
			}
			child := s.buildNodeRecursive(fullPath, entry.Name(), relPath, filter)
			if child != nil {
				node.Children = append(node.Children, child)
			}
		}
	}

	return node
}
