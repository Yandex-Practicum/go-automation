package filesearch

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFileWithExtensionPaths(rootDir, extension string) ([]string, error) {
	var result []string
	if err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		pathParts := strings.Split(path, ".")
		extensionPart := pathParts[len(pathParts)-1]
		if extensionPart == extension {
			result = append(result, path)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}
