package public

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/global"
	"go.uber.org/zap"
)

type MedioApi struct{}

func (api *MedioApi) UploadMedia(c *gin.Context) {

	// 允许上传的类型
	allowExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
	}

	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传文件"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持该文件类型: " + ext})
		return
	}

	// 本地存储目录（物理路径）
	saveDir := global.GVA_CONFIG.Local.StorePath
	if saveDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器未配置存储路径"})
		return
	}

	// URL 访问前缀
	viewDir := global.GVA_CONFIG.Local.Path
	if viewDir == "" {
		viewDir = "uploads/file"
	}

	// 新文件名，避免重复
	newName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(saveDir, newName)

	// 保存文件到本地
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}
	// 访问 URL
	domain := strings.TrimRight(global.GVA_CONFIG.System.Domain, "/")
	fileURL := fmt.Sprintf("%s/%s/%s", domain, strings.Trim(viewDir, "/"), newName)

	c.JSON(http.StatusOK, gin.H{
		"msg": "上传成功",
		"url": fileURL,
	})
}

func (api *MedioApi) Download(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		c.AbortWithStatusJSON(400, gin.H{"msg": "file required"})
		return
	}

	// filePat := global.DDD(file)
	// filePath := path.Join(global.AssetsPath, "bots", file)
	filePath, _ := global.DDD(file)
	global.GVA_LOG.Debug("ListenSvc filePath", zap.Any("filePath", filePath), zap.Any("fileName", filePath))

	c.Header("Content-Type", "text/csv")
	c.Header(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, file),
	)
	c.File(filePath)
}
