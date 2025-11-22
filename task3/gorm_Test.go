package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*进阶gorm
题目1：模型定义
假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
要求 ：
使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
编写Go代码，使用Gorm创建这些模型对应的数据库表。
题目2：关联查询
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。
题目3：钩子函数
继续使用博客系统的模型。
要求 ：
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。*/

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email     string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	FirstName string    `gorm:"size:50" json:"first_name"`
	LastName  string    `gorm:"size:50" json:"last_name"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Bio       string    `gorm:"type:text" json:"bio"`
	Status    string    `gorm:"size:20;default:active" json:"status"` // active, inactive, banned
	PostCount int       `gorm:"default:0" json:"post_count"`          // 文章数量统计
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 一对多关系：一个用户可以有多篇文章
	Posts []Post `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Title         string     `gorm:"size:200;not null" json:"title"`
	Content       string     `gorm:"type:longtext;not null" json:"content"`
	Excerpt       string     `gorm:"type:text" json:"excerpt"`                           // 文章摘要
	Slug          string     `gorm:"size:255;not null;uniqueIndex" json:"slug"`          // URL友好的标题
	Status        string     `gorm:"size:20;default:draft" json:"status"`                // draft, published, archived
	CommentStatus string     `gorm:"size:20;default:has_comments" json:"comment_status"` // 评论状态: has_comments, no_comments
	ViewCount     int        `gorm:"default:0" json:"view_count"`
	CommentCount  int        `gorm:"default:0" json:"comment_count"` // 评论数量统计
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	PublishedAt   *time.Time `json:"published_at,omitempty"` // 发布时间

	// 外键：关联用户
	UserID uint `gorm:"index" json:"user_id"`
	//  belongsTo 关系：文章属于用户
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 一对多关系：一篇文章可以有多个评论
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments,omitempty"`

	// 标签字段（可选扩展）
	Tags string `gorm:"type:json" json:"tags,omitempty"`
}

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Status    string    `gorm:"size:20;default:pending" json:"status"` // pending, approved, rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 外键：关联文章
	PostID uint `gorm:"not null;index" json:"post_id"`
	// belongsTo 关系：评论属于文章
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`

	// 外键：关联用户（评论者）
	UserID uint `gorm:"not null;index" json:"user_id"`
	// belongsTo 关系：评论属于用户
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 父评论ID（支持嵌套评论）
	ParentID *uint `gorm:"index" json:"parent_id,omitempty"`
	// 一对多关系：一个评论可以有多个子评论
	Replies []Comment `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"replies,omitempty"`
}

var db *gorm.DB

// initDB 初始化数据库连接
func initDB() error {
	dsn := "root:@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败：%v", err)
	}

	fmt.Println("数据库连接成功")
	return nil
}

// createTables 创建数据库表
func createTables() error {
	// 自动迁移模式（根据模型创建/更新表结构）
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		return fmt.Errorf("自动迁移失败：%v", err)
	}
	fmt.Println("数据库表创建成功")
	return nil
}

// ========== 清空表数据函数 ==========

// clearTables 清空所有表数据（用于测试）
func clearTables() error {
	// 注意：由于外键约束，需要按正确顺序删除
	// 先删除评论，再删除文章，最后删除用户
	if err := db.Exec("DELETE FROM comments").Error; err != nil {
		return fmt.Errorf("清空评论表失败: %v", err)
	}
	if err := db.Exec("DELETE FROM posts").Error; err != nil {
		return fmt.Errorf("清空文章表失败: %v", err)
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return fmt.Errorf("清空用户表失败: %v", err)
	}

	// 重置自增ID（可选）
	if err := db.Exec("ALTER TABLE comments AUTO_INCREMENT = 1").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE posts AUTO_INCREMENT = 1").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE users AUTO_INCREMENT = 1").Error; err != nil {
		return err
	}

	fmt.Println("表数据清空完成")
	return nil
}

// ========== 题目3：钩子函数实现 ==========

// BeforeCreate Post模型的创建前钩子
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Printf("BeforeCreate: 正在创建文章 '%s'\n", p.Title)
	return nil
}

// AfterCreate Post模型的创建后钩子 - 更新用户的文章数量统计
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Printf("AfterCreate: 文章 '%s' 创建成功，正在更新用户文章数量统计\n", p.Title)

	// 更新用户的文章数量
	var postCount int64
	if err := tx.Model(&Post{}).Where("user_id = ?", p.UserID).Count(&postCount).Error; err != nil {
		return err
	}

	if err := tx.Model(&User{}).Where("id = ?", p.UserID).Update("post_count", postCount).Error; err != nil {
		return err
	}

	fmt.Printf("AfterCreate: 用户 %d 的文章数量已更新为 %d\n", p.UserID, postCount)
	return nil
}

// BeforeDelete Comment模型的删除前钩子
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Printf("BeforeDelete: 正在删除评论 ID %d\n", c.ID)
	return nil
}

// AfterDelete Comment模型的删除后钩子 - 检查并更新文章的评论状态
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Printf("AfterDelete: 评论 ID %d 删除成功，正在检查文章评论状态\n", c.ID)

	// 统计文章的评论数量
	var commentCount int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount).Error; err != nil {
		return err
	}

	// 更新文章的评论数量
	if err := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_count", commentCount).Error; err != nil {
		return err
	}

	// 如果评论数量为0，更新评论状态为"no_comments"
	if commentCount == 0 {
		if err := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "no_comments").Error; err != nil {
			return err
		}
		fmt.Printf("AfterDelete: 文章 %d 的评论状态已更新为 'no_comments'\n", c.PostID)
	} else {
		if err := tx.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "has_comments").Error; err != nil {
			return err
		}
		fmt.Printf("AfterDelete: 文章 %d 的评论状态已更新为 'has_comments'，当前评论数: %d\n", c.PostID, commentCount)
	}

	return nil
}

// AfterCreate Comment模型的创建后钩子 - 更新文章的评论数量和状态
func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Printf("AfterCreate: 评论创建成功，正在更新文章评论统计\n")

	// 统计文章的评论数量
	var commentCount int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount).Error; err != nil {
		return err
	}

	// 更新文章的评论数量和状态
	updates := map[string]interface{}{
		"comment_count":  commentCount,
		"comment_status": "has_comments",
	}

	if err := tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(updates).Error; err != nil {
		return err
	}

	fmt.Printf("AfterCreate: 文章 %d 的评论数量已更新为 %d\n", c.PostID, commentCount)
	return nil
}

// ========== 示例数据创建和查询函数 ==========

// createSampleData 创建示例数据
func createSampleData() error {
	// 创建用户
	users := []User{
		{
			Username:  "alice",
			Email:     "alice@example.com",
			Password:  "hashed_password_1",
			FirstName: "Alice",
			LastName:  "Smith",
			Bio:       "技术爱好者和博主",
			Status:    "active",
		},
		{
			Username:  "bob",
			Email:     "bob@example.com",
			Password:  "hashed_password_2",
			FirstName: "Bob",
			LastName:  "Johnson",
			Bio:       "旅行作家",
			Status:    "active",
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("创建用户数据失败: %v", err)
	}

	// 创建文章
	now := time.Now()
	posts := []Post{
		{
			Title:       "Go语言入门指南",
			Content:     "这是一篇关于Go语言入门的详细指南...",
			Excerpt:     "学习Go语言的基础知识和特性",
			Slug:        "go-language-guide",
			Status:      "published",
			ViewCount:   150,
			UserID:      users[0].ID,
			PublishedAt: &now,
			Tags:        `["Go", "编程", "教程"]`,
		},
		{
			Title:       "我的旅行经历",
			Content:     "分享我在世界各地的旅行经历...",
			Excerpt:     "探索世界各地的美丽风景",
			Slug:        "my-travel-experience",
			Status:      "published",
			ViewCount:   89,
			UserID:      users[1].ID,
			PublishedAt: &now,
			Tags:        `["旅行", "冒险", "文化"]`,
		},
	}

	if err := db.Create(&posts).Error; err != nil {
		return fmt.Errorf("创建文章数据失败: %v", err)
	}

	// 创建评论
	comments := []Comment{
		{
			Content: "这篇教程写得很好，对我帮助很大！",
			Status:  "approved",
			PostID:  posts[0].ID,
			UserID:  users[1].ID,
		},
		{
			Content: "期待更多关于Go语言的教程",
			Status:  "approved",
			PostID:  posts[0].ID,
			UserID:  users[0].ID,
		},
		{
			Content: "照片拍得很美，我也想去这些地方旅行",
			Status:  "approved",
			PostID:  posts[1].ID,
			UserID:  users[0].ID,
		},
	}

	if err := db.Create(&comments).Error; err != nil {
		return fmt.Errorf("创建评论数据失败: %v", err)
	}

	fmt.Println("示例数据创建成功")
	return nil
}

// 查询示例函数

// GetUserWithPosts 查询用户及其所有文章
func GetUserWithPosts(userID uint) (User, error) {
	var user User
	err := db.Preload("Posts", "status = ?", "published").First(&user, userID).Error // 预加载已发布的文章
	return user, err
}

// GetPostWithComments 查询文章及其所有评论
func GetPostWithComments(postID uint) (Post, error) {
	var post Post
	err := db.Preload("User").
		Preload("Comments", "status = ?", "approved").      //预加载已批准的评论
		Preload("Comments.User").First(&post, postID).Error // 预加载评论的用户信息
	return post, err
}

// GetRecentPublishedPosts 获取最近发布的文章
func GetRecentPublishedPosts(limit int) ([]Post, error) {
	var posts []Post
	err := db.Preload("User").
		Where("status = ?", "published").
		Order("published_at DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// GetUserPostsWithComments 查询某个用户发布的所有文章及其对应的评论信息
func GetUserPostsWithComments(userID uint) ([]Post, error) {
	var posts []Post
	err := db.Preload("User"). // 预加载用户信息
					Preload("Comments", "status = ?", "approved").            // 预加载已批
					Preload("Comments.User").                                 // 预加载评论的用户信息
					Where("user_id = ? and status = ?", userID, "published"). // 查询该用户已发布的文章
					Order("published_at DESC").                               // 按发布时间降序排列
					Find(&posts).Error
	return posts, err
}

// GetMostCommentedPost 查询评论数量最多的文章信息
func GetMostCommentedPost() (Post, error) {
	var post Post
	err := db.Model(&Post{}). // ← 这里明确指定主表是 Post
					Preload("User").                               // 预加载用户信息
					Preload("Comments", "status = ?", "approved"). // 预加载已批
					Preload("Comments.User").                      // 预加载评论的用户信息
					Joins("LEFT JOIN comments ON posts.id = comments.post_id AND comments.status = ?", "approved").
					Where("posts.status = ?", "published").
					Group("posts.id").
					Order("COUNT(comments.id) DESC").
					First(&post).Error

	return post, err
}

// ========== 钩子函数测试函数 ==========

// TestPostHooks 测试Post模型的钩子函数
func TestPostHooks() error {
	fmt.Println("\n=== 测试Post钩子函数 ===")

	// 创建新用户用于测试
	newUser := User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashed_password_test",
		FirstName: "Test",
		LastName:  "User",
		Bio:       "测试用户",
		Status:    "active",
	}

	if err := db.Create(&newUser).Error; err != nil {
		return err
	}

	fmt.Printf("创建用户成功，ID: %d，初始文章数量: %d\n", newUser.ID, newUser.PostCount)

	// 创建新文章测试钩子
	newPost := Post{
		Title:   "测试钩子函数的文章",
		Content: "这是一篇用于测试钩子函数的文章内容...",
		Excerpt: "测试文章摘要",
		Slug:    "test-hook-post",
		Status:  "published",
		UserID:  newUser.ID,
		Tags:    `["测试", "钩子"]`,
	}

	if err := db.Create(&newPost).Error; err != nil {
		return err
	}

	// 查询用户确认文章数量已更新
	var updatedUser User
	if err := db.First(&updatedUser, newUser.ID).Error; err != nil {
		return err
	}

	fmt.Printf("用户 %s 的文章数量已更新为: %d\n", updatedUser.Username, updatedUser.PostCount)
	return nil
}

// TestCommentHooks 测试Comment模型的钩子函数
func TestCommentHooks() error {
	fmt.Println("\n=== 测试Comment钩子函数 ===")

	// 查询一篇文章用于测试
	var post Post
	if err := db.First(&post).Error; err != nil {
		return err
	}

	fmt.Printf("测试文章: 《%s》，当前评论数量: %d，评论状态: %s\n",
		post.Title, post.CommentCount, post.CommentStatus)

	// 创建新评论
	newComment := Comment{
		Content: "这是一个测试评论，用于测试删除钩子",
		Status:  "approved",
		PostID:  post.ID,
		UserID:  1, // 使用现有用户
	}

	if err := db.Create(&newComment).Error; err != nil {
		return err
	}

	// 查询文章确认评论数量已更新
	var updatedPost Post
	if err := db.First(&updatedPost, post.ID).Error; err != nil {
		return err
	}

	fmt.Printf("创建评论后 - 文章评论数量: %d，评论状态: %s\n",
		updatedPost.CommentCount, updatedPost.CommentStatus)

	// 删除评论测试删除钩子
	if err := db.Delete(&newComment).Error; err != nil {
		return err
	}

	// 再次查询文章确认评论状态已更新
	if err := db.First(&updatedPost, post.ID).Error; err != nil {
		return err
	}

	fmt.Printf("删除评论后 - 文章评论数量: %d，评论状态: %s\n",
		updatedPost.CommentCount, updatedPost.CommentStatus)

	return nil
}

func main() {
	// 初始化数据库连接
	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	// 创建表
	if err := createTables(); err != nil {
		log.Fatal(err)
	}

	// 清空表数据（确保每次运行都是干净的状态）
	if err := clearTables(); err != nil {
		log.Fatal("清空表数据失败：%v", err)
	}

	// 创建示例数据
	if err := createSampleData(); err != nil {
		log.Fatal(err)
	}

	// 示例查询
	fmt.Println("\n=== 查询示例 ===")
	user, err := GetUserWithPosts(1)
	if err != nil {
		log.Fatal("查询用户失败：%v", err)
	} else {
		fmt.Printf("用户: %s (%s)\n", user.Username, user.Email)
		fmt.Printf("文章数量: %d\n", len(user.Posts))
		for _, post := range user.Posts {
			fmt.Printf("  - 《%s》 (浏览: %d, 评论: %d)\n", post.Title, post.ViewCount, post.CommentCount)
		}
	}

	// 查询最近发布的文章
	recentPosts, err := GetRecentPublishedPosts(5)
	if err != nil {
		log.Fatal("查询最近文章失败：%s", err)
	} else {
		fmt.Println("\n最近发布的文章:")
		for _, post := range recentPosts {
			fmt.Printf("  - 《%s》 by %s (浏览: %d, 评论: %d, 状态: %s)\n",
				post.Title, post.User.Username, post.ViewCount, post.CommentCount, post.CommentStatus)
		}
	}

	// 查询文章及其评论
	post, err := GetPostWithComments(1)
	if err != nil {
		log.Fatal("查询文章失败：%s", err)
	} else {
		fmt.Printf("\n文章: 《%s》\n", post.Title)
		fmt.Printf("评论数量: %d\n", len(post.Comments))
		for _, comment := range post.Comments {
			fmt.Printf(" - %s: %s\n", comment.User.Username, comment.Content)
		}
	}

	// 测试钩子函数
	if err := TestPostHooks(); err != nil {
		log.Fatal("测试Post钩子失败：%v", err)
	}

	if err := TestCommentHooks(); err != nil {
		log.Fatal("测试Comment钩子失败：%v", err)
	}

	fmt.Println("\n=== 所有测试完成 ===")
}
