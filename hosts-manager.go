
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const hostsFile = "/etc/hosts"
const backupFile = "/etc/hosts.bak"

func checkRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("❌ 请使用 root 权限运行（例如 sudo）")
		os.Exit(1)
	}
}

func backupHosts() {
	input, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Println("读取 hosts 文件失败:", err)
		return
	}
	err = os.WriteFile(backupFile, input, 0644)
	if err != nil {
		fmt.Println("备份失败:", err)
		return
	}
	fmt.Println("✅ 已备份 hosts 到", backupFile)
}

func listHosts() {
	data, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Println("读取 hosts 文件失败:", err)
		return
	}
	fmt.Println("📋 当前 hosts 内容：\n")
	fmt.Println(string(data))
}

func addHost(ip, hostname string) {
	file, err := os.OpenFile(hostsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("打开 hosts 文件失败:", err)
		return
	}
	defer file.Close()

	lines, _ := os.ReadFile(hostsFile)
	if strings.Contains(string(lines), hostname) {
		fmt.Println("⚠️ 记录已存在：", hostname)
		return
	}

	newLine := fmt.Sprintf("%s\t%s\n", ip, hostname)
	_, err = file.WriteString(newLine)
	if err != nil {
		fmt.Println("写入失败:", err)
		return
	}
	fmt.Println("✅ 已添加：", newLine)
}

func removeHost(hostname string) {
	file, err := os.Open(hostsFile)
	if err != nil {
		fmt.Println("读取 hosts 文件失败:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, hostname) {
			lines = append(lines, line)
		}
	}
	err = os.WriteFile(hostsFile, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	if err != nil {
		fmt.Println("写入失败:", err)
		return
	}
	fmt.Println("✅ 已删除包含", hostname, "的记录")
}

func showHelp() {
	fmt.Println(`用法:
  hosts-manager <命令> [参数]

命令:
  add <IP> <hostname>    添加记录
  remove <hostname>      删除记录
  list                   显示 hosts 内容
  backup                 备份 hosts 文件
  help                   显示帮助`)
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	checkRoot()

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 4 {
			fmt.Println("用法: add <IP> <hostname>")
			return
		}
		addHost(os.Args[2], os.Args[3])
	case "remove":
		if len(os.Args) != 3 {
			fmt.Println("用法: remove <hostname>")
			return
		}
		removeHost(os.Args[2])
	case "list":
		listHosts()
	case "backup":
		backupHosts()
	default:
		showHelp()
	}
}
