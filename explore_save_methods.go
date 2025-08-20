package main

import (
	"fmt"
	"reflect"
	"github.com/tanqiangyes/go-word/pkg/word"
	"github.com/tanqiangyes/go-word/pkg/types"
)

func main() {
	fmt.Println("探索go-word库的保存方法...")
	
	// 创建文档
	doc, err := word.New()
	if err != nil {
		fmt.Printf("创建文档失败: %v\n", err)
		return
	}
	fmt.Println("✅ 文档创建成功")
	
	// 检查Document的方法
	fmt.Println("\n1. 检查Document的方法...")
	docType := reflect.TypeOf(doc)
	for i := 0; i < docType.NumMethod(); i++ {
		method := docType.Method(i)
		fmt.Printf("   %s\n", method.Name)
	}
	
	// 检查是否有直接的保存方法
	fmt.Println("\n2. 寻找保存相关方法...")
	saveMethods := []string{"Save", "SaveAs", "Write", "WriteTo", "Export"}
	
	for _, methodName := range saveMethods {
		if method, ok := docType.MethodByName(methodName); ok {
			fmt.Printf("   ✅ 找到方法: %s - %v\n", methodName, method.Type)
		} else {
			fmt.Printf("   ❌ 未找到方法: %s\n", methodName)
		}
	}
	
	// 检查EnhancedDocumentBuilder的保存方法
	fmt.Println("\n3. 检查EnhancedDocumentBuilder的保存方法...")
	builder := word.NewEnhancedDocumentBuilder()
	builderType := reflect.TypeOf(builder)
	
	builderSaveMethods := []string{"SaveDocument", "SaveDocumentAs", "Save", "Write", "Build"}
	for _, methodName := range builderSaveMethods {
		if method, ok := builderType.MethodByName(methodName); ok {
			fmt.Printf("   ✅ 找到Builder方法: %s - %v\n", methodName, method.Type)
		} else {
			fmt.Printf("   ❌ 未找到Builder方法: %s\n", methodName)
		}
	}
	
	// 尝试使用Build方法创建文档
	fmt.Println("\n4. 尝试使用Build方法...")
	// 先添加一些内容
	para1 := types.Paragraph{Text: "测试段落1"}
	para2 := types.Paragraph{Text: "测试段落2"}
	
	err = builder.AddParagraphToDocument(doc, para1)
	if err != nil {
		fmt.Printf("❌ 添加段落1失败: %v\n", err)
	} else {
		fmt.Println("✅ 段落1添加成功")
	}
	
	err = builder.AddParagraphToDocument(doc, para2)
	if err != nil {
		fmt.Printf("❌ 添加段落2失败: %v\n", err)
	} else {
		fmt.Println("✅ 段落2添加成功")
	}
	
	// 检查文档内容
	text, err := doc.GetText()
	if err != nil {
		fmt.Printf("❌ 获取文本失败: %v\n", err)
	} else {
		fmt.Printf("✅ 文档文本: %s\n", text)
	}
	
	// 尝试使用Build方法
	fmt.Println("\n5. 尝试Build方法...")
	builtDoc, err := builder.Build()
	if err != nil {
		fmt.Printf("❌ Build失败: %v\n", err)
	} else {
		fmt.Printf("✅ Build成功: %T\n", builtDoc)
		
		// 检查构建后的文档
		if builtDoc != nil {
			builtText, err := builtDoc.GetText()
			if err != nil {
				fmt.Printf("❌ 获取构建后文档文本失败: %v\n", err)
			} else {
				fmt.Printf("✅ 构建后文档文本: %s\n", builtText)
			}
		}
	}
	
	fmt.Println("\n🎉 探索完成！")
}
