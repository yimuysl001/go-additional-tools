package econf

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gcfg"
	"go-additional-tools/econf/appllo_cfg"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/encoding/gjson"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	hoconx "github.com/gurkankaymak/hocon"
	"go-additional-tools/encrypt/ucrypt"
)

// 配置文件搜索路径和扩展名
var (
	localBasePath = []string{"", "config", "manifest/config", "resources"}
	localExt      = []string{".json", ".yaml", ".yml", ".ini", ".properties", ".conf", ".toml", ".xml", ".js"}

	// 编译正则表达式一次，提高性能
	encryptionRegex = regexp.MustCompile(`ENC\(([^)]+)\)`)

	// 避免循环引用的互斥锁和访问记录
	visitedFilesMu sync.Mutex
	visitedFiles   = make(map[string]bool)
)

// AdapterContent 配置适配器内容结构
type AdapterContent struct {
	jsonVar *gjson.Json // 解析后的JSON对象，类型: *gjson.Json
}

// SetContent 设置配置内容
//func (a *AdapterContent) SetContent(content string) error {
//	j, err := gjson.LoadContent([]byte(content), true)
//	if err != nil {
//		return gerror.Wrap(err, "加载配置内容失败")
//	}
//	a.jsonVar = j
//	return nil
//}

// Available 检查适配器是否可用
func (a *AdapterContent) Available(ctx context.Context, resource ...string) bool {
	return !a.jsonVar.IsNil()
}

// Get 根据模式获取配置值
func (a *AdapterContent) Get(ctx context.Context, pattern string) (interface{}, error) {
	if a.jsonVar.IsNil() {
		return nil, nil
	}
	return a.jsonVar.Get(pattern).Val(), nil
}

// Data 获取所有配置数据
func (a *AdapterContent) Data(ctx context.Context) (map[string]interface{}, error) {
	if a.jsonVar.IsNil() {
		return nil, nil
	}
	return a.jsonVar.Map(), nil
}

// configFileFinder 配置文件查找器
type configFileFinder struct {
	basePaths  []string
	extensions []string
}

// newConfigFileFinder 创建新的配置文件查找器
func newConfigFileFinder() *configFileFinder {
	return &configFileFinder{
		basePaths:  localBasePath,
		extensions: localExt,
	}
}

// findConfigFile 查找配置文件
func (f *configFileFinder) findConfigFile(filename string) string {
	if strings.Contains(filename, ".") { // 带扩展名 且能找到相应的文件
		for _, ext := range f.extensions {
			if strings.HasSuffix(strings.ToLower(filename), ext) {
				return filename
			}

		}
	}

	// 搜索带前缀的文件 (config-filename.ext)
	for _, basePath := range f.basePaths {
		for _, ext := range f.extensions {
			var fullPath string
			if basePath == "" {
				fullPath = fmt.Sprintf("%s-%s%s", defaultConfigName, filename, ext)
			} else {
				fullPath = filepath.Join(basePath, fmt.Sprintf("%s-%s%s", defaultConfigName, filename, ext))
			}

			if gfile.Exists(fullPath) {
				return fullPath
			}
		}
	}

	// 搜索不带前缀的文件 (filename.ext)
	for _, basePath := range f.basePaths {
		for _, ext := range f.extensions {
			var fullPath string
			if basePath == "" {
				fullPath = filename + ext
			} else {
				fullPath = filepath.Join(basePath, filename+ext)
			}

			if gfile.Exists(fullPath) {
				return fullPath
			}
		}
	}

	return ""
}

// loadConfigFile 加载配置文件
func (f *configFileFinder) loadConfigFile(filename string) (map[string]any, error) {
	filePath := f.findConfigFile(filename)
	if filePath == "" {
		return nil, fmt.Errorf("未找到配置文件: %s", filename)
	}

	content := gfile.GetContents(filePath)
	if content == "" {
		return nil, fmt.Errorf("配置文件为空: %s", filePath)
	}

	// 解密配置文件内容
	content = decryptConfigFile(content)

	// 根据文件扩展名解析
	if strings.HasSuffix(filePath, ".conf") {
		// todo 直接使用 hoconx 无法解析 # 注释的文件，所以多加一层转换
		parseString, err := hoconx.ParseString(preprocessHocon(content))
		if err != nil {
			g.Log().Error(gctx.GetInitCtx(), fmt.Sprintf("%s解析HOCON配置文件出错:%v", filePath, err))
			return nil, fmt.Errorf("解析HOCON文件失败 %s", filePath)
		}
		//todo 直接转换数字符类型会多一层，所以需要先转json在处理
		str := gjson.New(parseString.GetRoot()).MustToJsonString()
		return gjson.New(str).Map(), nil
	}

	// 其他格式按JSON解析
	result := gjson.New(content).Map()
	if len(result) == 0 {
		return nil, fmt.Errorf("解析配置文件失败 %s", filePath)
	}
	return result, nil
}

// 预处理：将 # 注释转换为 // 注释
//func preprocessHocon(data string) string {
//	var result bytes.Buffer
//	scanner := bufio.NewScanner(bytes.NewReader([]byte(data)))
//
//	for scanner.Scan() {
//		line := scanner.Text()
//
//		// 找到 # 注释，替换为 //
//		if idx := strings.Index(line, "#"); idx != -1 {
//			// 确保 # 不在引号内（简单检查）
//			beforeHash := line[:idx]
//			if !strings.Contains(beforeHash, "\"") ||
//				strings.Count(beforeHash, "\"")%2 == 0 {
//				line = line[:idx] + "//" + line[idx+1:]
//			}
//		}
//		result.WriteString(line + "\n")
//	}
//
//	return result.String()
//}

// 预处理：将 # 注释转换为 // 注释，并正确处理多行字符串
func preprocessHocon(data string) string {
	var result bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader([]byte(data)))

	inMultilineString := false // 标记是否在多行字符串内

	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否进入或退出多行字符串
		if strings.Contains(line, `"""`) {
			inMultilineString = !inMultilineString
		}

		// 如果不在多行字符串内，则处理 # 注释
		if !inMultilineString {
			// 找到 # 注释，替换为 //
			if idx := strings.Index(line, "#"); idx != -1 {
				// 确保 # 不在单引号或双引号内（简单检查）
				beforeHash := line[:idx]
				if !strings.Contains(beforeHash, `"`) || strings.Count(beforeHash, `"`)%2 == 0 {
					line = line[:idx] + "//" + line[idx+1:]
				}
			}
		}

		result.WriteString(line + "\n")
	}

	return result.String()
}

// NewAdapterFile 创建文件适配器
func NewAdapterFile(filename string) map[string]any {
	finder := newConfigFileFinder()
	result, err := finder.loadConfigFile(filename)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), "加载配置文件失败:", err)
		return nil
	}
	return result
}

// NewAdapter 创建内容适配器
func NewAdapter(filenames ...string) gcfg.Adapter {
	ctx := gctx.GetInitCtx()
	finder := newConfigFileFinder()

	// 确定要加载的配置文件列表
	var filesToLoad []string
	if len(filenames) == 0 || (len(filenames) == 1 && filenames[0] == "") {
		// 如果没有指定文件名，则查找默认配置文件
		foundFile := finder.findConfigFile(defaultConfigName)
		if foundFile != "" {
			filesToLoad = []string{foundFile}
		}
	} else {
		// 验证指定的文件是否存在
		for _, filename := range filenames {
			if filename != "" {
				foundFile := finder.findConfigFile(filename)
				if foundFile != "" {
					filesToLoad = append(filesToLoad, foundFile)
				} else {
					g.Log().Warning(ctx, fmt.Sprintf("配置文件不存在: %s", filename))
				}
			}
		}
	}

	if len(filesToLoad) == 0 {
		panic("没有找到任何配置文件")
	}

	// 加载配置文件并合并
	configMap := make(map[string]any)
	loadedFiles := make(map[string]bool)

	for _, file := range filesToLoad {
		if err := loadAndMergeConfig(file, configMap, loadedFiles); err != nil {
			g.Log().Error(ctx, "加载配置文件失败:", err)
			continue
		}
	}

	// 处理配置文件引用（避免循环引用）
	processConfigReferences(configMap, finder, loadedFiles)

	if a := apolloAdapter(configMap); a != nil {
		return a
	}

	// 创建适配器实例
	adapter := &AdapterContent{
		jsonVar: gjson.New(configMap, true),
	}

	return adapter
}

func apolloAdapter(configMap map[string]any) gcfg.Adapter {
	j := gjson.New(configMap)

	json := j.GetJson("apollo")

	if json == nil || json.IsNil() {
		return nil
	}

	var appl = appllo_cfg.Config{}
	err := json.Scan(&appl)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), "加载配置文件失败:", err)
		return nil
	}
	if appl.AppID == "" || appl.IP == "" || appl.Cluster == "" {
		return nil
	}

	adapter, err := appllo_cfg.New(gctx.New(), appl, j)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), "加载配置文件失败:", err)
		return nil
	}

	return adapter

}

// loadAndMergeConfig 加载并合并配置文件
func loadAndMergeConfig(filePath string, configMap map[string]any, loadedFiles map[string]bool) error {
	// 防止重复加载
	if loadedFiles[filePath] {
		return nil
	}

	// 记录已加载的文件
	loadedFiles[filePath] = true

	fileConfig := NewAdapterFile(filePath)

	//// 加载文件内容
	//content := gfile.GetContents(filePath)
	//if content == "" {
	//	return fmt.Errorf("配置文件为空: %s", filePath)
	//}
	//
	//// 解密内容
	//content = decryptConfigFile(content)
	//
	//// 解析配置
	//var fileConfig map[string]any
	//var err error
	//
	//if strings.HasSuffix(filePath, ".conf") {
	//	parseResult := configuration.ParseString(content)
	//	if parseResult == nil {
	//		g.Log().Error(gctx.GetInitCtx(), fmt.Sprintf("%s解析HOCON配置文件出错:", filePath))
	//		return fmt.Errorf("解析HOCON文件失败 %s", filePath)
	//	}
	//	fmt.Println(parseResult.Root().String())
	//
	//	fileConfig = gjson.New(parseResult.Root().String()).Map()
	//} else {
	//	fileConfig = gjson.New(content).Map()
	//	if len(fileConfig) == 0 {
	//		return fmt.Errorf("解析配置文件失败 %s: %w", filePath, err)
	//	}
	//}

	// 合并配置
	for key, value := range fileConfig {
		configMap[key] = value
	}

	return nil
}

// processConfigReferences 处理配置文件引用
func processConfigReferences(configMap map[string]any, finder *configFileFinder, loadedFiles map[string]bool) {
	jsonObj := gjson.New(configMap)

	// 使用限制防止无限递归
	maxDepth := 10
	currentDepth := 0

	for currentDepth < maxDepth {
		fileRef := jsonObj.Get(configFilePathKey).String()
		if fileRef == "" {
			break
		}

		// 检查是否已经加载过该文件
		visitedFilesMu.Lock()
		if visitedFiles[fileRef] {
			visitedFilesMu.Unlock()
			break
		}
		visitedFiles[fileRef] = true
		visitedFilesMu.Unlock()

		// 加载引用的配置文件
		refConfig, err := finder.loadConfigFile(fileRef)
		if err != nil {
			g.Log().Warning(gctx.GetInitCtx(), fmt.Sprintf("加载引用配置文件失败 %s: %v", fileRef, err))
			break
		}

		// 合并引用的配置
		for key, value := range refConfig {
			configMap[key] = value
		}

		// 更新JSON对象以检查下一个引用
		jsonObj = gjson.New(configMap)
		currentDepth++
	}

	// 清理访问记录
	visitedFilesMu.Lock()
	clear(visitedFiles)
	visitedFilesMu.Unlock()
}

// decryptConfigFile 解密配置文件中的加密内容
func decryptConfigFile(content string) string {
	return encryptionRegex.ReplaceAllStringFunc(content, func(match string) string {
		// 提取加密内容（移除 ENC() 包装）
		encryptedContent := strings.TrimPrefix(strings.TrimSuffix(match, encryptionSuffix), encryptionPrefix)

		decrypted, err := ucrypt.Decrypt2(encryptedContent, getLocalKey(), encryptionKey)
		if err != nil {
			// 解密失败时记录警告并保持原值
			g.Log().Warning(gctx.GetInitCtx(), fmt.Sprintf("配置解密失败: %v, 原始内容: %s", err, match))
			return match
		}
		return decrypted
	})
}
