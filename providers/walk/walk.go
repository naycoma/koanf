package walk

import (
	"cmp"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type File struct {
	Path      string
	Name      string // full filename (e.g. "device@test.yaml")
	BaseName  string // base name without @suffix and extension (e.g. "device")
	Suffix    string // suffix after @, empty if none
	Extension string // file extension (e.g. ".yaml")
	Depth     int    // path separator count, used for priority
}

func walkDir(root, baseName, extRe string) ([]*File, error) {
	re, err := regexp.Compile(`^?` + regexp.QuoteMeta(baseName) + `(@(.*))?\.` + extRe + `$`)
	if err != nil {
		return nil, err
	}
	walks := []*File{}
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return err
		}
		name := d.Name()
		matches := re.FindStringSubmatch(name)
		if len(matches) != 3 {
			return nil
		}
		walks = append(walks, &File{
			Path:      path,
			Name:      name,
			BaseName:  baseName,
			Extension: filepath.Ext(path),
			Suffix:    matches[2],
			Depth:     strings.Count(path, string(filepath.Separator)),
		})
		return nil
	})
	return walks, err
}

// walkYaml discovers YAML files matching baseName(@suffix)?.ya?ml under root.
func walkYaml(root, baseName string) ([]*File, error) {
	return walkDir(root, baseName, `ya?ml`)
}

// 1. name1.yaml
// 2. child/name1.yaml
// 3. name1@suf.yaml
// 4. child/name1@suf.yaml
// 5. name2.yaml
func sortedPriorityWalk(root string, fileNames ...string) ([]*File, error) {
	name2index := newIndexMap(fileNames)
	var items []*File
	for _, file := range fileNames {
		result, err := walkYaml(root, file)
		if err != nil {
			return nil, err
		}
		items = append(items, result...)
	}
	slices.SortStableFunc(items, func(a, b *File) int {
		return cmp.Or(
			name2index[a.BaseName]-name2index[b.BaseName],
			compEmpty(a.Suffix, b.Suffix),
			a.Depth-b.Depth,
		)
	})
	return items, nil
}

func newIndexMap[K comparable](keys []K) map[K]int {
	m := make(map[K]int, len(keys))
	for i, k := range keys {
		m[k] = i
	}
	return m
}

// Empty を前にする比較
// 両方空でなければ順序は変えない
func compEmpty[S ~string](a, b S) int {
	if a == "" && b != "" {
		return -1
	}
	if a != "" && b == "" {
		return 1
	}
	return 0
}
