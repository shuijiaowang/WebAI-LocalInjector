package service

import (
	"Service/model/response"
	"os"
	"path/filepath"
	"strings"
)

type FileTreeService struct{}

type IgnoreFilter struct {
	Dirs  map[string]bool
	Files map[string]bool
	Exts  map[string]bool
}

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
