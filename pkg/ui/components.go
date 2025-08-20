package ui

import (
	"fmt"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tanqiangyes/fyne-word/pkg/document"
	"log"
)

// TreeView 基于go-word库的文档树形视图组件
type TreeView struct {
	tree        *widget.Tree
	docManager  *document.Manager
	onSelect    func(nodeID string)
}

// NewTreeView 创建新的go-word文档树形视图
func NewTreeView(docManager *document.Manager) *TreeView {
	gtv := &TreeView{
		docManager: docManager,
	}
	
	gtv.tree = widget.NewTree(
		gtv.getChildIDs,
		gtv.hasChildren,
		gtv.createNode,
		gtv.updateNode,
	)
	
	gtv.tree.OnSelected = gtv.onNodeSelected
	
	return gtv
}

// GetWidget 获取Fyne组件
func (gtv *TreeView) GetWidget() fyne.CanvasObject {
	return gtv.tree
}

// Refresh 刷新树形视图
func (gtv *TreeView) Refresh() {
	gtv.tree.Refresh()
}

// SetOnSelect 设置节点选择回调
func (gtv *TreeView) SetOnSelect(callback func(nodeID string)) {
	gtv.onSelect = callback
}

// getChildIDs 获取子节点ID列表
func (gtv *TreeView) getChildIDs(id widget.TreeNodeID) []widget.TreeNodeID {
	if id == "" {
		// 根节点
		doc := gtv.docManager.GetCurrentDocument()
		if doc != nil {
			return []string{"title", "paragraphs", "tables", "images", "styles", "metadata"}
		}
		return []string{}
	}
	
	doc := gtv.docManager.GetCurrentDocument()
	if doc == nil {
		return []string{}
	}
	
	// 创建适配器
	adapter := document.NewDocumentAdapter(doc)
	
	switch id {
	case "paragraphs":
		var ids []string
		count := adapter.GetParagraphCount()
		for i := 0; i < count; i++ {
			ids = append(ids, fmt.Sprintf("p%d", i+1))
		}
		return ids
	case "tables":
		var ids []string
		count := adapter.GetTableCount()
		for i := 0; i < count; i++ {
			ids = append(ids, fmt.Sprintf("t%d", i+1))
		}
		return ids
	case "images":
		var ids []string
		count := adapter.GetImageCount()
		for i := 0; i < count; i++ {
			ids = append(ids, fmt.Sprintf("i%d", i+1))
		}
		return ids
	case "styles":
		var ids []string
		count := adapter.GetStyleCount()
		for i := 0; i < count; i++ {
			ids = append(ids, fmt.Sprintf("s%d", i+1))
		}
		return ids
	}
	
	return []string{}
}

// hasChildren 判断节点是否有子节点
func (gtv *TreeView) hasChildren(id widget.TreeNodeID) bool {
	children := gtv.getChildIDs(id)
	return len(children) > 0
}

// createNode 创建节点显示组件
func (gtv *TreeView) createNode(b bool) fyne.CanvasObject {
	return widget.NewLabel("")
}

// updateNode 更新节点内容
func (gtv *TreeView) updateNode(id widget.TreeNodeID, b bool, o fyne.CanvasObject) {
	label := o.(*widget.Label)
	
	doc := gtv.docManager.GetCurrentDocument()
	if doc == nil {
		label.SetText("无文档")
		return
	}
	
	// 创建适配器
	adapter := document.NewDocumentAdapter(doc)
	
	switch id {
	case "title":
		label.SetText("📄 " + adapter.GetTitle())
	case "paragraphs":
		count := adapter.GetParagraphCount()
		label.SetText(fmt.Sprintf("📝 段落 (%d)", count))
	case "tables":
		count := adapter.GetTableCount()
		label.SetText(fmt.Sprintf("📊 表格 (%d)", count))
	case "images":
		count := adapter.GetImageCount()
		label.SetText(fmt.Sprintf("🖼️ 图片 (%d)", count))
	case "styles":
		count := adapter.GetStyleCount()
		label.SetText(fmt.Sprintf("🎨 样式 (%d)", count))
	case "metadata":
		label.SetText("ℹ️ 元数据")
	default:
		// 处理具体项目
		if strings.HasPrefix(id, "p") {
			// 段落
			index := parseIndex(id[1:])
			if index >= 0 {
				text := adapter.GetParagraphText(index)
				label.SetText(fmt.Sprintf("📝 %s", truncateText(text, 30)))
			}
		} else if strings.HasPrefix(id, "t") {
			// 表格
			index := parseIndex(id[1:])
			if index >= 0 {
				label.SetText(fmt.Sprintf("📊 表格 %d", index+1))
			}
		} else if strings.HasPrefix(id, "i") {
			// 图片
			index := parseIndex(id[1:])
			if index >= 0 {
				label.SetText(fmt.Sprintf("🖼️ 图片 %d", index+1))
			}
		} else if strings.HasPrefix(id, "s") {
			// 样式
			index := parseIndex(id[1:])
			if index >= 0 {
				label.SetText(fmt.Sprintf("🎨 样式 %d", index+1))
			}
		}
	}
}

// onNodeSelected 节点选择事件处理
func (gtv *TreeView) onNodeSelected(id widget.TreeNodeID) {
	if gtv.onSelect != nil {
		gtv.onSelect(id)
	}
}

// ContentView 基于go-word库的内容显示组件
type ContentView struct {
	container   *fyne.Container
	docManager  *document.Manager
	currentNode string
}

// NewContentView 创建新的go-word内容显示组件
func NewContentView(docManager *document.Manager) *ContentView {
	gcv := &ContentView{
		docManager: docManager,
	}
	
	gcv.container = container.NewVBox(
		widget.NewLabel("请选择一个文档节点查看内容"),
	)
	
	return gcv
}

// GetWidget 获取Fyne组件
func (gcv *ContentView) GetWidget() fyne.CanvasObject {
	return gcv.container
}

// ShowNode 显示指定节点的内容
func (gcv *ContentView) ShowNode(nodeID string) {
	gcv.currentNode = nodeID
	gcv.updateContent()
}

// updateContent 更新内容显示
func (gcv *ContentView) updateContent() {
	doc := gcv.docManager.GetCurrentDocument()
	if doc == nil {
		gcv.container.Objects = []fyne.CanvasObject{
			widget.NewLabel("没有打开的文档"),
		}
		gcv.container.Refresh()
		return
	}
	
	// 创建适配器
	adapter := document.NewDocumentAdapter(doc)
	
	var contentWidgets []fyne.CanvasObject
	
	switch gcv.currentNode {
	case "title":
		contentWidgets = gcv.createTitleView(adapter.GetTitle())
	case "paragraphs":
		contentWidgets = gcv.createParagraphsView(adapter)
	case "tables":
		contentWidgets = gcv.createTablesView(adapter)
	case "images":
		contentWidgets = gcv.createImagesView(adapter)
	case "styles":
		contentWidgets = gcv.createStylesView(adapter)
	case "metadata":
		contentWidgets = gcv.createMetadataView(adapter)
	default:
		// 处理具体项目
		if strings.HasPrefix(gcv.currentNode, "p") {
			index := parseIndex(gcv.currentNode[1:])
			if index >= 0 {
				contentWidgets = gcv.createParagraphDetailView(adapter, index)
			}
		} else if strings.HasPrefix(gcv.currentNode, "t") {
			index := parseIndex(gcv.currentNode[1:])
			if index >= 0 {
				contentWidgets = gcv.createTableDetailView(adapter, index)
			}
		} else if strings.HasPrefix(gcv.currentNode, "i") {
			index := parseIndex(gcv.currentNode[1:])
			if index >= 0 {
				contentWidgets = gcv.createImageDetailView(adapter, index)
			}
		} else if strings.HasPrefix(gcv.currentNode, "s") {
			index := parseIndex(gcv.currentNode[1:])
			if index >= 0 {
				contentWidgets = gcv.createStyleDetailView(adapter, index)
			}
		}
	}
	
	if len(contentWidgets) == 0 {
		contentWidgets = []fyne.CanvasObject{
			widget.NewLabel("请选择一个有效的节点"),
		}
	}
	
	gcv.container.Objects = contentWidgets
	gcv.container.Refresh()
}

// createTitleView 创建标题视图
func (gcv *ContentView) createTitleView(title string) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel("文档标题"))
	widgets = append(widgets, widget.NewSeparator())
	
	// 添加标题编辑功能
	titleEntry := widget.NewEntry()
	titleEntry.SetText(title)
	titleEntry.OnChanged = func(newTitle string) {
		// TODO: 更新文档标题
		log.Printf("标题已更改为: %s", newTitle)
	}
	widgets = append(widgets, titleEntry)
	
	// 添加文本编辑区域
	widgets = append(widgets, widget.NewSeparator())
	widgets = append(widgets, widget.NewLabel("文档内容"))
	
	textArea := widget.NewMultiLineEntry()
	textArea.SetPlaceHolder("在此输入文档内容...")
	textArea.OnChanged = func(text string) {
		// TODO: 保存文本内容
		log.Printf("文本内容已更新，长度: %d", len(text))
	}
	widgets = append(widgets, textArea)
	
	return widgets
}

// createParagraphsView 创建段落列表视图
func (gcv *ContentView) createParagraphsView(adapter *document.DocumentAdapter) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel("段落列表"))
	widgets = append(widgets, widget.NewSeparator())
	
	count := adapter.GetParagraphCount()
	for i := 0; i < count; i++ {
		text := adapter.GetParagraphText(i)
		paraLabel := widget.NewLabel(fmt.Sprintf("段落 %d: %s", i+1, truncateText(text, 50)))
		widgets = append(widgets, paraLabel)
	}
	
	if count == 0 {
		widgets = append(widgets, widget.NewLabel("没有段落"))
	}
	
	return widgets
}

// createTablesView 创建表格列表视图
func (gcv *ContentView) createTablesView(adapter *document.DocumentAdapter) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel("表格列表"))
	widgets = append(widgets, widget.NewSeparator())
	
	count := adapter.GetTableCount()
	for i := 0; i < count; i++ {
		info := adapter.GetTableInfo(i)
		tableLabel := widget.NewLabel(fmt.Sprintf("表格 %d: %s", i+1, info))
		widgets = append(widgets, tableLabel)
	}
	
	if count == 0 {
		widgets = append(widgets, widget.NewLabel("没有表格"))
	}
	
	return widgets
}

// createImagesView 创建图片列表视图
func (gcv *ContentView) createImagesView(adapter *document.DocumentAdapter) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel("图片列表"))
	widgets = append(widgets, widget.NewSeparator())
	
	count := adapter.GetImageCount()
	for i := 0; i < count; i++ {
		info := adapter.GetImageInfo(i)
		imgLabel := widget.NewLabel(fmt.Sprintf("图片 %d: %s", i+1, info))
		widgets = append(widgets, imgLabel)
	}
	
	if count == 0 {
		widgets = append(widgets, widget.NewLabel("没有图片"))
	}
	
	return widgets
}

// createStylesView 创建样式列表视图
func (gcv *ContentView) createStylesView(adapter *document.DocumentAdapter) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel("样式列表"))
	widgets = append(widgets, widget.NewSeparator())
	
	count := adapter.GetStyleCount()
	for i := 0; i < count; i++ {
		info := adapter.GetStyleInfo(i)
		styleLabel := widget.NewLabel(fmt.Sprintf("样式 %d: %s", i+1, info))
		widgets = append(widgets, styleLabel)
	}
	
	if count == 0 {
		widgets = append(widgets, widget.NewLabel("没有样式"))
	}
	
	return widgets
}

// createMetadataView 创建元数据视图
func (gcv *ContentView) createMetadataView(adapter *document.DocumentAdapter) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel("文档元数据"))
	widgets = append(widgets, widget.NewSeparator())
	
	metadata := adapter.GetMetadataInfo()
	for key, value := range metadata {
		widgets = append(widgets, widget.NewLabel(fmt.Sprintf("%s: %s", key, value)))
	}
	
	return widgets
}

// createParagraphDetailView 创建段落详细视图
func (gcv *ContentView) createParagraphDetailView(adapter *document.DocumentAdapter, index int) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel(fmt.Sprintf("段落 %d 详情", index+1)))
	widgets = append(widgets, widget.NewSeparator())
	
	text := adapter.GetParagraphText(index)
	
	// 显示文本内容
	textArea := widget.NewMultiLineEntry()
	textArea.SetText(text)
	textArea.Disable()
	widgets = append(widgets, textArea)
	
	return widgets
}

// createTableDetailView 创建表格详细视图
func (gcv *ContentView) createTableDetailView(adapter *document.DocumentAdapter, index int) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel(fmt.Sprintf("表格 %d 详情", index+1)))
	widgets = append(widgets, widget.NewSeparator())
	
	info := adapter.GetTableInfo(index)
	widgets = append(widgets, widget.NewLabel(info))
	
	return widgets
}

// createImageDetailView 创建图片详细视图
func (gcv *ContentView) createImageDetailView(adapter *document.DocumentAdapter, index int) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel(fmt.Sprintf("图片 %d 详情", index+1)))
	widgets = append(widgets, widget.NewSeparator())
	
	info := adapter.GetImageInfo(index)
	widgets = append(widgets, widget.NewLabel(info))
	
	return widgets
}

// createStyleDetailView 创建样式详细视图
func (gcv *ContentView) createStyleDetailView(adapter *document.DocumentAdapter, index int) []fyne.CanvasObject {
	var widgets []fyne.CanvasObject
	
	widgets = append(widgets, widget.NewLabel(fmt.Sprintf("样式 %d 详情", index+1)))
	widgets = append(widgets, widget.NewSeparator())
	
	info := adapter.GetStyleInfo(index)
	widgets = append(widgets, widget.NewLabel(info))
	
	return widgets
}

// parseIndex 解析索引字符串为整数
func parseIndex(s string) int {
	var result int
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			result = result*10 + int(ch-'0')
		} else {
			return -1 // 无效字符
		}
	}
	return result - 1 // 转换为0基索引
}

// truncateText 截断文本到指定长度
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}
