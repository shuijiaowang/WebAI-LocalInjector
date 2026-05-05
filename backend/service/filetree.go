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
	IsDir      bool
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

	contentRoot := &contentTreeNode{IsDir: true, childIndex: make(map[string]*contentTreeNode)}
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

		s.addContentTreePath(contentRoot, relPath, entry.IsDir())

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
	builder.WriteString("【目录极简结构】\n")
	s.writeContentTree(&builder, contentRoot, "")
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

func (s *FileTreeService) addContentTreePath(root *contentTreeNode, relPath string, isDir bool) {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	current := root
	for i, part := range parts {
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
		if i < len(parts)-1 {
			child.IsDir = true
		} else {
			child.IsDir = isDir
		}
		current = child
	}
}

func (s *FileTreeService) writeContentTree(builder *strings.Builder, node *contentTreeNode, currentPath string) {
	wroteSection := false
	for _, child := range node.Children {
		if !child.IsDir {
			continue
		}
		if !s.hasContentTreeOutput(child) {
			continue
		}

		if wroteSection {
			builder.WriteString("\n")
		}
		wroteSection = true
		builder.WriteString(child.Name)
		builder.WriteString("\n")
		s.writeContentTreeSection(builder, child, currentPath, true)
	}
}

func (s *FileTreeService) writeContentTreeSection(builder *strings.Builder, node *contentTreeNode, parentPath string, writeRootFiles bool) {
	for _, child := range node.Children {
		if !child.IsDir || !s.hasContentTreeOutput(child) {
			continue
		}

		childPath := child.Name
		if parentPath != "" {
			childPath = parentPath + "/" + child.Name
		}

		if s.hasContentTreeFiles(child) && s.hasContentTreeDirs(child) {
			builder.WriteString(" ")
			builder.WriteString(childPath)
			builder.WriteString(":")
			builder.WriteString(s.formatContentTreeNames(child))
			builder.WriteString("\n")
			continue
		}

		if files := s.formatContentTreeFiles(child, false); files != "" {
			builder.WriteString(" ")
			builder.WriteString(childPath)
			builder.WriteString(":")
			builder.WriteString(files)
			builder.WriteString("\n")
		}
		s.writeContentTreeSection(builder, child, childPath, false)
	}

	if writeRootFiles {
		if files := s.formatContentTreeFiles(node, true); files != "" {
			builder.WriteString(" ")
			builder.WriteString(files)
			builder.WriteString("\n")
		}
	}
}

func (s *FileTreeService) hasContentTreeOutput(node *contentTreeNode) bool {
	for _, child := range node.Children {
		if !child.IsDir || s.hasContentTreeOutput(child) {
			return true
		}
	}
	return false
}

func (s *FileTreeService) hasContentTreeFiles(node *contentTreeNode) bool {
	for _, child := range node.Children {
		if !child.IsDir {
			return true
		}
	}
	return false
}

func (s *FileTreeService) hasContentTreeDirs(node *contentTreeNode) bool {
	for _, child := range node.Children {
		if child.IsDir && s.hasContentTreeOutput(child) {
			return true
		}
	}
	return false
}

func (s *FileTreeService) formatContentTreeNames(node *contentTreeNode) string {
	var names []string
	for _, child := range node.Children {
		if child.IsDir {
			if s.hasContentTreeOutput(child) {
				names = append(names, child.Name)
			}
			continue
		}
		names = append(names, strings.TrimSuffix(child.Name, filepath.Ext(child.Name)))
	}
	return strings.Join(names, "/")
}

func (s *FileTreeService) formatContentTreeFiles(node *contentTreeNode, rootFiles bool) string {
	type extGroup struct {
		ext   string
		names []string
	}

	var groups []extGroup
	groupIndex := make(map[string]int)
	totalFiles := 0

	for _, child := range node.Children {
		if child.IsDir {
			continue
		}

		totalFiles++
		ext := filepath.Ext(child.Name)
		name := strings.TrimSuffix(child.Name, ext)
		index, ok := groupIndex[ext]
		if !ok {
			groups = append(groups, extGroup{ext: ext})
			index = len(groups) - 1
			groupIndex[ext] = index
		}
		groups[index].names = append(groups[index].names, name)
	}

	if totalFiles == 0 {
		return ""
	}

	if rootFiles {
		var files []string
		for _, child := range node.Children {
			if !child.IsDir {
				files = append(files, child.Name)
			}
		}
		return strings.Join(files, " ")
	}

	parts := make([]string, 0, len(groups))
	hasMergedExt := false
	for _, group := range groups {
		if group.ext != "" && len(group.names) > 1 {
			hasMergedExt = true
			parts = append(parts, "{"+strings.Join(group.names, ",")+"}"+group.ext)
			continue
		}
		if group.ext == "" && len(group.names) > 1 {
			hasMergedExt = true
			parts = append(parts, "{"+strings.Join(group.names, ",")+"}")
			continue
		}
		parts = append(parts, group.names[0]+group.ext)
	}

	if len(parts) > 1 && !hasMergedExt {
		return "{" + strings.Join(parts, ",") + "}"
	}
	return strings.Join(parts, ",")
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
