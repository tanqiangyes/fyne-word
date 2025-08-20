package app

import (
	"fmt"
	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"log"
	"path/filepath"
	"github.com/tanqiangyes/fyne-word/pkg/document"
	"github.com/tanqiangyes/fyne-word/pkg/ui"
)

// App 基于go-word库的应用程序核心
type App struct {
	app         fyne.App
	window      fyne.Window
	mainMenu    *fyne.MainMenu
	toolbar     fyne.CanvasObject
	content     fyne.CanvasObject
	
	docManager  *document.Manager
	treeView    *ui.TreeView
	contentView *ui.ContentView
}

// New 创建新的基于go-word库的应用程序
func New() *App {
	myApp := &App{
		app:        fyneApp.NewApp(),
		docManager: document.NewManager(),
	}
	
	myApp.setupMainWindow()
	myApp.setupMenu()
	myApp.setupToolbar()
	myApp.setupContent()
	
	return myApp
}

// Run 运行应用程序
func (app *App) Run() {
	app.window.ShowAndRun()
}

// setupMainWindow 设置主窗口
func (app *App) setupMainWindow() {
	app.window = app.app.NewWindow("Fyne Word - 基于go-word库")
	app.window.Resize(fyne.NewSize(1200, 800))
	app.window.CenterOnScreen()
}

// setupMenu 设置菜单
func (app *App) setupMenu() {
	app.mainMenu = app.createMainMenu()
	app.window.SetMainMenu(app.mainMenu)
}

// setupToolbar 设置工具栏
func (app *App) setupToolbar() {
	app.toolbar = app.createToolbar()
}

// setupContent 设置主内容区域
func (app *App) setupContent() {
	app.content = app.createMainLayout()
	app.window.SetContent(app.content)
}

// createMainMenu 创建主菜单
func (app *App) createMainMenu() *fyne.MainMenu {
	fileMenu := fyne.NewMenu("文件",
		fyne.NewMenuItem("新建", app.newDocument),
		fyne.NewMenuItem("打开", app.openDocument),
		fyne.NewMenuItem("保存", app.saveDocument),
		fyne.NewMenuItem("另存为", app.saveDocumentAs),
		fyne.NewMenuItem("导出PDF", app.exportToPDF),
		fyne.NewMenuItem("退出", func() { app.app.Quit() }),
	)
	
	editMenu := fyne.NewMenu("编辑",
		fyne.NewMenuItem("撤销", func() {}),
		fyne.NewMenuItem("重做", func() {}),
		fyne.NewMenuItem("剪切", func() {}),
		fyne.NewMenuItem("复制", func() {}),
		fyne.NewMenuItem("粘贴", func() {}),
	)
	
	viewMenu := fyne.NewMenu("视图",
		fyne.NewMenuItem("树形视图", func() {}),
		fyne.NewMenuItem("内容视图", func() {}),
		fyne.NewMenuItem("全屏", func() { app.window.SetFullScreen(true) }),
	)
	
	helpMenu := fyne.NewMenu("帮助",
		fyne.NewMenuItem("关于", func() {}),
		fyne.NewMenuItem("帮助", func() {}),
	)
	
	return fyne.NewMainMenu(fileMenu, editMenu, viewMenu, helpMenu)
}

// createToolbar 创建工具栏
func (app *App) createToolbar() fyne.CanvasObject {
	newBtn := widget.NewButton("新建", app.newDocument)
	openBtn := widget.NewButton("打开", app.openDocument)
	saveBtn := widget.NewButton("保存", app.saveDocument)
	exportBtn := widget.NewButton("导出PDF", app.exportToPDF)
	
	return container.NewHBox(
		newBtn, openBtn, saveBtn, exportBtn,
		widget.NewSeparator(),
	)
}

// createMainLayout 创建主布局
func (app *App) createMainLayout() fyne.CanvasObject {
	// 创建树形视图
	app.treeView = ui.NewTreeView(app.docManager)
	
	// 创建内容视图
	app.contentView = ui.NewContentView(app.docManager)
	
	// 设置树形视图的选择回调
	app.treeView.SetOnSelect(func(nodeID string) {
		app.contentView.ShowNode(nodeID)
	})
	
	// 创建分割布局
	split := container.NewHSplit(
		app.treeView.GetWidget(),
		app.contentView.GetWidget(),
	)
	split.SetOffset(0.3) // 树形视图占30%宽度
	
	return container.NewBorder(app.toolbar, nil, nil, nil, split)
}

// newDocument 新建文档
func (app *App) newDocument() {
	log.Println("新建文档")
	
	// 使用文档管理器创建新文档
	doc, err := app.docManager.NewDocument()
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}
	
	// 刷新UI显示新文档
	app.treeView.Refresh()
	app.contentView.ShowNode("title")
	
	log.Printf("新文档创建成功: %s", doc.FileName)
	dialog.ShowInformation("成功", fmt.Sprintf("新文档创建成功: %s", doc.FileName), app.window)
}

// openDocument 打开文档
func (app *App) openDocument() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()
		
		filePath := reader.URI().Path()
		log.Printf("正在使用go-word库打开文档: %s", filePath)
		
		// 使用go-word库打开文档
		doc, err := app.docManager.OpenDocument(filePath)
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		
		// 刷新UI
		app.treeView.Refresh()
		app.contentView.ShowNode("title")
		
		log.Printf("go-word文档打开成功: %s", doc.FileName)
		dialog.ShowInformation("成功", fmt.Sprintf("文档打开成功: %s", doc.FileName), app.window)
	}, app.window)
	
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".docx", ".doc"}))
	fd.Show()
}

// saveDocument 保存文档
func (app *App) saveDocument() {
	doc := app.docManager.GetCurrentDocument()
	if doc == nil {
		dialog.ShowInformation("提示", "没有要保存的文档", app.window)
		return
	}
	
	err := app.docManager.SaveDocument(doc)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}
	
	dialog.ShowInformation("成功", "文档保存成功", app.window)
}

// saveDocumentAs 另存为
func (app *App) saveDocumentAs() {
	doc := app.docManager.GetCurrentDocument()
	if doc == nil {
		dialog.ShowInformation("提示", "没有要保存的文档", app.window)
		return
	}
	
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		
		newPath := writer.URI().Path()
		err = app.docManager.SaveDocumentAs(doc, newPath)
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		
		dialog.ShowInformation("成功", "文档另存为成功", app.window)
	}, app.window)
	
	fd.SetFileName(doc.FileName)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".docx"}))
	fd.Show()
}

// exportToPDF 导出为PDF
func (app *App) exportToPDF() {
	doc := app.docManager.GetCurrentDocument()
	if doc == nil {
		dialog.ShowInformation("提示", "没有要导出的文档", app.window)
		return
	}
	
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		
		outputPath := writer.URI().Path()
		err = app.docManager.ExportToPDF(doc, outputPath)
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		
		dialog.ShowInformation("成功", "PDF导出成功", app.window)
	}, app.window)
	
	// 设置默认文件名
	baseName := doc.FileName[:len(doc.FileName)-len(filepath.Ext(doc.FileName))]
	fd.SetFileName(baseName + ".pdf")
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	fd.Show()
}
