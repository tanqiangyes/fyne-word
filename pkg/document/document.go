package document

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"github.com/tanqiangyes/go-word/pkg/word"
	"github.com/tanqiangyes/go-word/pkg/writer"
)

// Document 基于go-word库的文档结构
type Document struct {
	FilePath    string
	FileName    string
	Title       string           // 文档标题
	WordDoc     *word.Document  // 直接使用go-word库的Document类型
	DocWriter   *writer.DocumentWriter // 使用DocumentWriter进行写入操作
	IsModified  bool
	IsOpen      bool
}

// Manager 基于go-word库的文档管理器
type Manager struct {
	documents   map[string]*Document
	currentDoc  *Document
}

// NewManager 创建新的go-word文档管理器
func NewManager() *Manager {
	return &Manager{
		documents: make(map[string]*Document),
	}
}

// OpenDocument 使用go-word库打开Word文档
func (m *Manager) OpenDocument(filePath string) (*Document, error) {
	// 检查文件是否已经打开
	if doc, exists := m.documents[filePath]; exists {
		m.currentDoc = doc
		return doc, nil
	}
	
	// 检查文件扩展名
	if !isWordDocument(filePath) {
		return nil, fmt.Errorf("不支持的文件格式: %s", filepath.Ext(filePath))
	}
	
	log.Printf("正在使用go-word库打开文档: %s", filePath)
	
	// 使用go-word库打开文档
	wordDoc, err := word.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文档: %v", err)
	}
	
	// 创建新文档实例
	doc := &Document{
		FilePath:   filePath,
		FileName:   filepath.Base(filePath),
		WordDoc:    wordDoc,
		IsOpen:     true,
		IsModified: false,
	}
	
	m.documents[filePath] = doc
	m.currentDoc = doc
	
	log.Printf("文档打开成功: %s", filePath)
	return doc, nil
}

// SaveDocument 使用DocumentWriter保存文档
func (m *Manager) SaveDocument(doc *Document) error {
	if doc == nil {
		return fmt.Errorf("没有要保存的文档")
	}
	
	if doc.DocWriter == nil {
		return fmt.Errorf("文档写入器未初始化")
	}
	
	// 检查文件路径是否为空
	if doc.FilePath == "" {
		return &SavePathNotSetError{}
	}
	
	log.Printf("正在保存文档: %s", doc.FilePath)
	
	// 使用DocumentWriter保存文档
	err := doc.DocWriter.Save(doc.FilePath)
	if err != nil {
		return fmt.Errorf("保存文档失败: %v", err)
	}
	
	doc.IsModified = false
	log.Printf("文档保存成功: %s", doc.FilePath)
	return nil
}

// SavePathNotSetError 表示保存路径未设置的错误
type SavePathNotSetError struct{}

func (e *SavePathNotSetError) Error() string {
	return "文档路径未设置，请使用'另存为'功能选择保存位置"
}

// IsSavePathNotSetError 检查是否为保存路径未设置错误
func IsSavePathNotSetError(err error) bool {
	_, ok := err.(*SavePathNotSetError)
	return ok
}

// SaveDocumentAs 使用DocumentWriter另存为
func (m *Manager) SaveDocumentAs(doc *Document, newPath string) error {
	if doc == nil {
		return fmt.Errorf("没有要保存的文档")
	}
	
	if doc.DocWriter == nil {
		return fmt.Errorf("文档写入器未初始化")
	}
	
	// 检查新路径的扩展名
	if !isWordDocument(newPath) {
		return fmt.Errorf("不支持的文件格式: %s", filepath.Ext(newPath))
	}
	
	log.Printf("正在另存为: %s", newPath)
	
	// 使用DocumentWriter保存到新路径
	err := doc.DocWriter.Save(newPath)
	if err != nil {
		return fmt.Errorf("另存为失败: %v", err)
	}
	
	// 更新文档路径
	oldPath := doc.FilePath
	doc.FilePath = newPath
	doc.FileName = filepath.Base(newPath)
	doc.IsModified = false
	
	// 更新管理器中的文档映射
	delete(m.documents, oldPath)
	m.documents[newPath] = doc
	
	log.Printf("文档另存为成功: %s", newPath)
	return nil
}

// ExportToPDF 使用DocumentWriter导出为PDF
func (m *Manager) ExportToPDF(doc *Document, outputPath string) error {
	if doc == nil {
		return fmt.Errorf("没有要导出的文档")
	}
	
	if !strings.HasSuffix(strings.ToLower(outputPath), ".pdf") {
		outputPath += ".pdf"
	}
	
	log.Printf("正在导出PDF: %s", outputPath)
	
	// 由于DocumentWriter没有直接的PDF导出功能，
	// 我们先将文档保存为.docx，然后使用第三方工具转换
	// 或者实现一个基本的PDF生成功能
	
	// 临时实现：保存为.docx并提示用户
	tempDocxPath := strings.TrimSuffix(outputPath, ".pdf") + ".docx"
	err := m.SaveDocumentAs(doc, tempDocxPath)
	if err != nil {
		return fmt.Errorf("保存临时文档失败: %v", err)
	}
	
	log.Printf("PDF导出功能暂未实现，文档已保存为: %s", tempDocxPath)
	log.Printf("请使用Microsoft Word或其他工具将文档转换为PDF")
	
	return fmt.Errorf("PDF导出功能暂未实现，请使用外部工具转换")
}

// GetCurrentDocument 获取当前文档
func (m *Manager) GetCurrentDocument() *Document {
	return m.currentDoc
}

// GetOpenDocuments 获取所有打开的文档
func (m *Manager) GetOpenDocuments() []*Document {
	var docs []*Document
	for _, doc := range m.documents {
		if doc.IsOpen {
			docs = append(docs, doc)
		}
	}
	return docs
}

// CloseDocument 关闭文档
func (m *Manager) CloseDocument(doc *Document) error {
	if doc == nil {
		return nil
	}
	
	// 检查是否有未保存的更改
	if doc.IsModified {
		log.Println("文档有未保存的更改")
	}
	
	// 关闭go-word文档
	if doc.WordDoc != nil {
		err := doc.WordDoc.Close()
		if err != nil {
			log.Printf("关闭文档时出错: %v", err)
		}
	}
	
	doc.IsOpen = false
	delete(m.documents, doc.FilePath)
	
	if m.currentDoc == doc {
		m.currentDoc = nil
	}
	
	log.Printf("文档已关闭: %s", doc.FilePath)
	return nil
}

// NewDocument 创建新的Word文档
func (m *Manager) NewDocument() (*Document, error) {
	log.Println("正在创建新文档")
	
	// 使用DocumentWriter创建新文档
	docWriter := writer.NewDocumentWriter()
	err := docWriter.CreateNewDocument()
	if err != nil {
		return nil, fmt.Errorf("创建新文档失败: %v", err)
	}
	
	// 创建新文档实例
	doc := &Document{
		FilePath:   "", // 新文档还没有保存路径
		FileName:   "未命名文档.docx",
		Title:      "未命名文档", // 设置默认标题
		WordDoc:    nil, // 暂时设为nil，DocumentWriter会管理文档
		DocWriter:  docWriter,
		IsOpen:     true,
		IsModified: true, // 新文档需要保存
	}
	
	// 生成临时ID用于管理
	tempID := fmt.Sprintf("temp_%d", len(m.documents)+1)
	m.documents[tempID] = doc
	m.currentDoc = doc
	
	log.Println("新文档创建成功")
	return doc, nil
}

// GetText 获取文档文本内容
func (doc *Document) GetText() (string, error) {
	if doc.WordDoc == nil {
		return "", fmt.Errorf("文档未打开")
	}
	
	// 使用go-word库获取文本
	return doc.WordDoc.GetText()
}

// GetParagraphs 获取文档段落
func (doc *Document) GetParagraphs() (interface{}, error) {
	if doc.WordDoc == nil {
		return nil, fmt.Errorf("文档未打开")
	}
	
	// 使用go-word库获取段落
	return doc.WordDoc.GetParagraphs()
}

// GetTables 获取文档表格
func (doc *Document) GetTables() (interface{}, error) {
	if doc.WordDoc == nil {
		return nil, fmt.Errorf("文档未打开")
	}
	
	// 使用go-word库获取表格
	return doc.WordDoc.GetTables()
}

// GetImages 获取文档图片
func (doc *Document) GetImages() (interface{}, error) {
	if doc.WordDoc == nil {
		return nil, fmt.Errorf("文档未打开")
	}
	
	// 使用go-word库获取图片
	// TODO: 根据go-word库的实际API实现
	// 目前返回空结果
	return []interface{}{}, nil
}

// GetStyles 获取文档样式
func (doc *Document) GetStyles() (interface{}, error) {
	if doc.WordDoc == nil {
		return nil, fmt.Errorf("文档未打开")
	}
	
	// 使用go-word库获取样式
	// TODO: 根据go-word库的实际API实现
	// 目前返回空结果
	return []interface{}{}, nil
}

// GetMetadata 获取文档元数据
func (doc *Document) GetMetadata() (interface{}, error) {
	if doc.WordDoc == nil {
		return nil, fmt.Errorf("文档未打开")
	}
	
	// 使用go-word库获取元数据
	// TODO: 根据go-word库的实际API实现
	// 目前返回空结果
	return map[string]interface{}{}, nil
}

// AddParagraph 向文档添加新段落
func (doc *Document) AddParagraph(text string) error {
	if doc.DocWriter == nil {
		return fmt.Errorf("文档未打开")
	}
	
	log.Printf("正在添加段落: %s", truncateText(text, 30))
	
	// 使用DocumentWriter添加段落
	err := doc.DocWriter.AddParagraph(text, "Normal")
	if err != nil {
		return fmt.Errorf("添加段落失败: %v", err)
	}
	
	doc.IsModified = true
	log.Println("段落添加成功")
	return nil
}

// AddText 向文档添加文本
func (doc *Document) AddText(text string) error {
	if doc.DocWriter == nil {
		return fmt.Errorf("文档未打开")
	}
	
	log.Printf("正在添加文本: %s", truncateText(text, 30))
	
	// 使用DocumentWriter添加文本（通过添加段落实现）
	err := doc.DocWriter.AddParagraph(text, "Normal")
	if err != nil {
		return fmt.Errorf("添加文本失败: %v", err)
	}
	
	doc.IsModified = true
	log.Println("文本添加成功")
	return nil
}

// SetTitle 设置文档标题
func (doc *Document) SetTitle(title string) error {
	if doc == nil {
		return fmt.Errorf("文档未初始化")
	}
	
	log.Printf("正在设置文档标题: %s", title)
	
	// 设置文档标题
	doc.Title = title
	doc.IsModified = true
	
	log.Printf("标题设置成功: %s", title)
	return nil
}

// isWordDocument 检查文件是否为Word文档
func isWordDocument(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".docx" || ext == ".doc"
}

// truncateText 截断文本到指定长度
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}
