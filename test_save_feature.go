package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/tanqiangyes/fyne-word/pkg/document"
)

func main() {
	fmt.Println("测试文档保存功能...")
	
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "fyne-word-test-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	fmt.Printf("临时目录: %s\n", tempDir)
	
	// 创建文档管理器
	manager := document.NewManager()
	
	// 创建新文档
	fmt.Println("\n1. 创建新文档...")
	doc, err := manager.NewDocument()
	if err != nil {
		log.Fatalf("新建文档失败: %v", err)
	}
	fmt.Printf("✅ 新文档创建成功: %s\n", doc.FileName)
	
	// 添加一些内容
	fmt.Println("\n2. 添加文档内容...")
	err = doc.AddParagraph("这是一个测试文档")
	if err != nil {
		log.Fatalf("添加段落失败: %v", err)
	}
	
	err = doc.AddParagraph("用于测试保存功能")
	if err != nil {
		log.Fatalf("添加段落失败: %v", err)
	}
	
	fmt.Println("✅ 文档内容添加成功")
	
	// 测试保存文档
	fmt.Println("\n3. 测试保存文档...")
	savePath := filepath.Join(tempDir, "test_save.docx")
	
	// 设置文档路径
	doc.FilePath = savePath
	doc.FileName = filepath.Base(savePath)
	
	fmt.Printf("保存路径: %s\n", savePath)
	
	// 保存文档
	err = manager.SaveDocument(doc)
	if err != nil {
		fmt.Printf("❌ 保存文档失败: %v\n", err)
	} else {
		fmt.Println("✅ 文档保存成功")
		
		// 检查文件是否存在
		if _, err := os.Stat(savePath); err == nil {
			fmt.Println("✅ 文件确实已创建")
		} else {
			fmt.Printf("❌ 文件未创建: %v\n", err)
		}
	}
	
	// 测试另存为
	fmt.Println("\n4. 测试另存为...")
	saveAsPath := filepath.Join(tempDir, "test_save_as.docx")
	
	err = manager.SaveDocumentAs(doc, saveAsPath)
	if err != nil {
		fmt.Printf("❌ 另存为失败: %v\n", err)
	} else {
		fmt.Println("✅ 另存为成功")
		
		// 检查新文件是否存在
		if _, err := os.Stat(saveAsPath); err == nil {
			fmt.Println("✅ 新文件确实已创建")
		} else {
			fmt.Printf("❌ 新文件未创建: %v\n", err)
		}
		
		// 检查文档路径是否更新
		fmt.Printf("文档新路径: %s\n", doc.FilePath)
		fmt.Printf("文档新名称: %s\n", doc.FileName)
	}
	
	// 测试PDF导出
	fmt.Println("\n5. 测试PDF导出...")
	pdfPath := filepath.Join(tempDir, "test_export.pdf")
	
	err = manager.ExportToPDF(doc, pdfPath)
	if err != nil {
		fmt.Printf("❌ PDF导出失败: %v\n", err)
	} else {
		fmt.Println("✅ PDF导出成功")
		
		// 检查PDF文件是否存在
		if _, err := os.Stat(pdfPath); err == nil {
			fmt.Println("✅ PDF文件确实已创建")
		} else {
			fmt.Printf("❌ PDF文件未创建: %v\n", err)
		}
	}
	
	fmt.Println("\n🎉 保存功能测试完成！")
}
