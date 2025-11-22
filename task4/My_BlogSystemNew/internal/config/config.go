package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)


// 说明：
// 该文件提供了应用的配置加载逻辑，优先使用环境变量覆盖默认值。
// 常见部署时通过设置环境变量来调整数据库连接、JWT 密钥和服务端口等。

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
	LogLevel   string
}

func LoadConfig() *Config {
	// 加载 .env 文件
	loadEnvFile(".env")
	
	// 如果 .env 不存在，尝试加载 .env.example 作为后备
	if !fileExists(".env") {
		log.Println(".env 文件不存在，尝试加载 .env.example")
		loadEnvFile(".env.example")
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "blog_system"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}

	// 打印加载的配置信息（隐藏敏感信息）
	log.Printf("配置加载成功:")
	log.Printf("  DB_HOST: %s", config.DBHost)
	log.Printf("  DB_PORT: %s", config.DBPort)
	log.Printf("  DB_USER: %s", config.DBUser)
	log.Printf("  DB_NAME: %s", config.DBName)
	log.Printf("  SERVER_PORT: %s", config.ServerPort)
	log.Printf("  LOG_LEVEL: %s", config.LogLevel)
	
	if config.DBPassword != "" {
		log.Printf("  DB_PASSWORD: [已设置]")
	} else {
		log.Printf("  DB_PASSWORD: [未设置]")
	}
	
	if config.JWTSecret != "your-secret-key" {
		log.Printf("  JWT_SECRET: [已设置]")
	} else {
		log.Printf("  JWT_SECRET: [使用默认值，建议在生产环境中修改]")
	}

	return config
}

// loadEnvFile 手动加载 .env 文件
func loadEnvFile(filename string) {
	if !fileExists(filename) {
		log.Printf("环境变量文件不存在: %s", filename)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("打开环境变量文件失败 %s: %v", filename, err)
		return
	}
	defer file.Close()

	loadedCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			
			// 移除值中的引号（如果存在）
			if len(value) > 1 && ((value[0] == '"' && value[len(value)-1] == '"') || 
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}
			
			// 如果环境变量尚未设置，则设置它
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
				loadedCount++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("读取环境变量文件失败 %s: %v", filename, err)
	} else {
		log.Printf("成功从 %s 加载了 %d 个环境变量", filename, loadedCount)
	}
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("环境变量 %s 的值无效，使用默认值: %d", key, defaultValue)
	}
	return defaultValue
}

// GetCurrentDir 获取当前工作目录（用于调试）
func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前工作目录失败: %v", err)
		return ""
	}
	return dir
}

