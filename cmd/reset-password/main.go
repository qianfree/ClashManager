package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"clash-manager/internal/config"
	"clash-manager/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 命令行参数
	username := flag.String("user", "", "指定要重置密码的用户名")
	password := flag.String("password", "", "指定新密码（如果不指定，将提示输入）")
	listUsers := flag.Bool("list", false, "列出所有用户")
	flag.Parse()

	// 初始化数据库连接
	if err := repository.InitDB(config.GetDBPath()); err != nil {
		fmt.Printf("错误: 无法连接到数据库: %v\n", err)
		os.Exit(1)
	}
	defer repository.CloseDB()

	repo := &repository.UserRepository{}

	// 处理列出用户
	if *listUsers {
		handleListUsers(repo)
		return
	}

	// 获取用户名
	var targetUser string
	if *username != "" {
		targetUser = *username
	} else {
		targetUser = promptInput("请输入要重置密码的用户名: ")
		if targetUser == "" {
			fmt.Println("错误: 用户名不能为空")
			os.Exit(1)
		}
	}

	// 检查用户是否存在
	user, err := repo.FindByUsername(targetUser)
	if err != nil {
		fmt.Printf("错误: 用户 '%s' 不存在\n", targetUser)
		fmt.Println("提示: 使用 --list 参数可以查看所有用户")
		os.Exit(1)
	}

	// 获取新密码
	var newPassword string
	if *password != "" {
		newPassword = *password
	} else {
		newPassword = promptPassword("请输入新密码: ")
		if newPassword == "" {
			fmt.Println("错误: 密码不能为空")
			os.Exit(1)
		}

		confirmPassword := promptPassword("请再次输入新密码: ")
		if newPassword != confirmPassword {
			fmt.Println("错误: 两次输入的密码不一致")
			os.Exit(1)
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("错误: 密码加密失败: %v\n", err)
		os.Exit(1)
	}

	// 更新密码
	if err := repo.UpdatePassword(targetUser, string(hashedPassword)); err != nil {
		fmt.Printf("错误: 更新密码失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n成功: 用户 '%s' 的密码已重置\n", targetUser)
	fmt.Printf("用户ID: %d\n", user.ID)
	fmt.Printf("订阅Token: %s\n", user.Token)
}

// handleListUsers 列出所有用户
func handleListUsers(repo *repository.UserRepository) {
	// 获取所有用户
	var users []struct {
		ID       uint
		Username string
		Token    string
	}

	if err := repository.GetDB().Table("users").Select("id, username, token").Find(&users).Error; err != nil {
		fmt.Printf("错误: 获取用户列表失败: %v\n", err)
		os.Exit(1)
	}

	if len(users) == 0 {
		fmt.Println("当前没有任何用户")
		return
	}

	fmt.Println("\n=== 用户列表 ===")
	fmt.Println("ID\t用户名\t\t订阅Token")
	fmt.Println("--------------------------------------------------")
	for _, u := range users {
		fmt.Printf("%d\t%-15s\t%s\n", u.ID, u.Username, u.Token)
	}
	fmt.Printf("\n共 %d 个用户\n", len(users))
}

// promptInput 提示用户输入文本
func promptInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// promptPassword 提示用户输入密码（不显示输入内容）
func promptPassword(prompt string) string {
	fmt.Print(prompt)
	var password string
	fmt.Scanln(&password)
	return password
}
