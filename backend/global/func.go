package global

import (
	"fmt"
	"path/filepath"
	"strings"
)

func DDD(fileName string) (savePath, fileURL string) {
	saveDir := GVA_CONFIG.Local.StorePath
	if saveDir == "" {
		return
	}

	viewDir := GVA_CONFIG.Local.Path
	if viewDir == "" {
		viewDir = "uploads/file"
	}

	savePath = filepath.Join(saveDir, fileName)

	domain := strings.TrimRight(GVA_CONFIG.System.Domain, "/")
	fileURL = fmt.Sprintf("%s/%s/%s", domain, strings.Trim(viewDir, "/"), fileName)
	return
}
