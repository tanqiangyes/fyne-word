package document

import (
	"fmt"
)

// DocumentAdapter 文档适配器，将go-word库的数据结构适配到我们的UI接口
type DocumentAdapter struct {
	goWordDoc *Document
}

// NewDocumentAdapter 创建文档适配器
func NewDocumentAdapter(goWordDoc *Document) *DocumentAdapter {
	return &DocumentAdapter{
		goWordDoc: goWordDoc,
	}
}

// GetTitle 获取文档标题
func (da *DocumentAdapter) GetTitle() string {
	if da.goWordDoc == nil {
		return "未知文档"
	}
	
	// 优先使用文档的Title字段
	if da.goWordDoc.Title != "" {
		return da.goWordDoc.Title
	}
	
	// 如果没有设置标题，使用文件名
	return da.goWordDoc.FileName
}

// GetParagraphCount 获取段落数量
func (da *DocumentAdapter) GetParagraphCount() int {
	if da.goWordDoc == nil {
		return 0
	}
	
	// 如果使用DocumentWriter，从DocumentWriter的document字段获取段落数量
	if da.goWordDoc.DocWriter != nil {
		// 通过DocumentWriter的document字段访问MainDocumentPart
		if da.goWordDoc.DocWriter.Document != nil {
			mainPart := da.goWordDoc.DocWriter.Document.GetMainPart()
			if mainPart != nil && mainPart.Content != nil {
				return len(mainPart.Content.Paragraphs)
			}
		}
		// 如果无法获取，返回默认值
		return 1
	}
	
	// 如果使用WordDoc，尝试从那里获取
	if da.goWordDoc.WordDoc != nil {
		paragraphs, err := da.goWordDoc.GetParagraphs()
		if err != nil {
			return 0
		}
		
		// 尝试将interface{}转换为切片来计算长度
		if paraSlice, ok := paragraphs.([]interface{}); ok {
			return len(paraSlice)
		}
		
		// 如果无法确定类型，返回默认值
		return 1
	}
	
	return 0
}

// GetTableCount 获取表格数量
func (da *DocumentAdapter) GetTableCount() int {
	if da.goWordDoc == nil || da.goWordDoc.WordDoc == nil {
		return 0
	}
	
	tables, err := da.goWordDoc.GetTables()
	if err != nil {
		return 0
	}
	
	// TODO: 根据go-word库的实际接口计算表格数量
	_ = tables
	return 0 // 临时返回值
}

// GetImageCount 获取图片数量
func (da *DocumentAdapter) GetImageCount() int {
	if da.goWordDoc == nil || da.goWordDoc.WordDoc == nil {
		return 0
	}
	
	images, err := da.goWordDoc.GetImages()
	if err != nil {
		return 0
	}
	
	// TODO: 根据go-word库的实际接口计算图片数量
	_ = images
	return 0 // 临时返回值
}

// GetStyleCount 获取样式数量
func (da *DocumentAdapter) GetStyleCount() int {
	if da.goWordDoc == nil || da.goWordDoc.WordDoc == nil {
		return 0
	}
	
	styles, err := da.goWordDoc.GetStyles()
	if err != nil {
		return 0
	}
	
	// TODO: 根据go-word库的实际接口计算样式数量
	_ = styles
	return 0 // 临时返回值
}

// GetText 获取文档纯文本内容
func (da *DocumentAdapter) GetText() string {
	if da.goWordDoc == nil || da.goWordDoc.WordDoc == nil {
		return ""
	}
	
	text, err := da.goWordDoc.GetText()
	if err != nil {
		return fmt.Sprintf("获取文本时出错: %v", err)
	}
	
	return text
}

// GetParagraphText 获取指定段落的文本
func (da *DocumentAdapter) GetParagraphText(index int) string {
	if da.goWordDoc == nil {
		return ""
	}
	
	// 如果使用DocumentWriter，从DocumentWriter的document字段获取段落文本
	if da.goWordDoc.DocWriter != nil {
		if da.goWordDoc.DocWriter.Document != nil {
			mainPart := da.goWordDoc.DocWriter.Document.GetMainPart()
			if mainPart != nil && mainPart.Content != nil {
				if index >= 0 && index < len(mainPart.Content.Paragraphs) {
					paragraph := mainPart.Content.Paragraphs[index]
					// 从段落中提取文本
					if len(paragraph.Runs) > 0 {
						var text string
						for _, run := range paragraph.Runs {
							text += run.Text
						}
						return text
					}
					return "空段落"
				}
			}
		}
		return fmt.Sprintf("无法获取段落%d", index+1)
	}
	
	// 如果使用WordDoc，尝试从那里获取
	if da.goWordDoc.WordDoc != nil {
		paragraphs, err := da.goWordDoc.GetParagraphs()
		if err != nil {
			return fmt.Sprintf("获取段落时出错: %v", err)
		}
		
		// TODO: 根据go-word库的实际接口提取指定段落文本
		_ = paragraphs
		_ = index
		return "示例段落文本" // 临时返回值
	}
	
	return "文档未初始化"
}

// GetTableInfo 获取指定表格的信息
func (da *DocumentAdapter) GetTableInfo(index int) string {
	if da.goWordDoc == nil || da.goWordDoc.WordDoc == nil {
		return ""
	}
	
	tables, err := da.goWordDoc.GetTables()
	if err != nil {
		return fmt.Sprintf("获取表格时出错: %v", err)
	}
	
	// TODO: 根据go-word库的实际接口提取指定表格信息
	_ = tables
	_ = index
	return "示例表格信息" // 临时返回值
}

// GetImageInfo 获取指定图片的信息
func (da *DocumentAdapter) GetImageInfo(index int) string {
	if da.goWordDoc == nil || da.goWordDoc.WordDoc == nil {
		return ""
	}
	
	images, err := da.goWordDoc.GetImages()
	if err != nil {
		return fmt.Sprintf("获取图片时出错: %v", err)
	}
	
	// TODO: 根据go-word库的实际接口提取指定图片信息
	_ = images
	_ = index
	return "示例图片信息" // 临时返回值
}

// GetStyleInfo 获取指定样式的信息
func (da *DocumentAdapter) GetStyleInfo(index int) string {
	if da.goWordDoc == nil || da.goWordDoc.WordDoc == nil {
		return ""
	}
	
	styles, err := da.goWordDoc.GetStyles()
	if err != nil {
		return fmt.Sprintf("获取样式时出错: %v", err)
	}
	
	// TODO: 根据go-word库的实际接口提取指定样式信息
	_ = styles
	_ = index
	return "示例样式信息" // 临时返回值
}

// GetMetadataInfo 获取文档元数据信息
func (da *DocumentAdapter) GetMetadataInfo() map[string]string {
	if da.goWordDoc == nil || da.goWordDoc.WordDoc == nil {
		return make(map[string]string)
	}
	
	metadata, err := da.goWordDoc.GetMetadata()
	if err != nil {
		return map[string]string{
			"错误": fmt.Sprintf("获取元数据时出错: %v", err),
		}
	}
	
	// TODO: 根据go-word库的实际接口提取元数据
	_ = metadata
	
	// 临时返回示例数据
	return map[string]string{
		"作者":   "示例作者",
		"标题":   da.GetTitle(),
		"页数":   "1",
		"字数":   "100",
		"字符数": "500",
		"创建时间": "2025-08-20",
		"修改时间": "2025-08-20",
	}
}
