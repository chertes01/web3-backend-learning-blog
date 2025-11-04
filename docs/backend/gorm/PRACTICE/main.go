package main

import (
	"blog/config"
	"blog/models"
	"fmt"
)

func migrate() {
	err := config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}) // 自动创建或更新 User 表结构
	if err != nil {
		panic(err)
	}
	fmt.Println("迁移成功")
}

func seedData() {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Println("⚠️ 已存在資料，跳過初始化")
		return
	}

	users := []models.User{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}
	config.DB.Create(&users)

	posts := []models.Post{
		{Title: "Go學習筆記", Content: "GORM 是一個強大的 ORM 工具", UserID: users[0].ID},
		{Title: "GORM 關聯教學", Content: "一對多、多對多講解", UserID: users[0].ID},
		{Title: "SQL 性能優化", Content: "索引與查詢計劃", UserID: users[1].ID},
	}
	config.DB.Create(&posts)

	comments := []models.Comment{
		{Content: "寫得不錯！", PostID: posts[0].ID},
		{Content: "很實用的文章", PostID: posts[0].ID},
		{Content: "學到了新東西", PostID: posts[1].ID},
		{Content: "這篇我收藏了", PostID: posts[2].ID},
		{Content: "講得太清楚了！", PostID: posts[2].ID},
		{Content: "繼續加油！", PostID: posts[2].ID},
	}
	config.DB.Create(&comments)

	fmt.Println("✅ 測試資料插入完成！")
}

// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
// 编写Go代码，使用Gorm查询评论数量最多的文章信息
func queryUserPostsWithComments(userName string) {
	var userDetail models.User
	err := config.DB.Preload("Posts.Comments").Where("name = ?", userName).Take(&userDetail).Error
	if err != nil {
		fmt.Println("查询用户及其文章和评论失败:", err)
		return
	}

	fmt.Printf("用户: %s\n", userDetail.Name)
	for _, post := range userDetail.Posts {
		fmt.Printf("文章: %s\n", post.Title)
		for _, comment := range post.Comments {
			fmt.Printf(" - 评论: %s\n", comment.Content)
		}
	}
	fmt.Println("--------------------")
}

func queryMostCommentedPost() {
	var post []models.Post
	// err := config.DB.
	// 	Preload("Comments").
	// 	Joins("LEFT JOIN comments ON comments.post_id = posts.id").
	// 	Group("posts.id").
	// 	Order("COUNT(comments.id) DESC").
	// 	Limit(1).
	// 	Take(&post).Error
	// if err != nil {
	// 	fmt.Println("查询评论数量最多的文章失败:", err)
	// 	return
	// }

	// fmt.Printf("评论数量最多的文章:\n")
	// fmt.Printf("文章标题: %s\n", post.Title)
	// //fmt.Printf("作者: %s\n", post.User.Name)
	// fmt.Printf("评论数量: %d\n", len(post.Comments))
	// for _, comment := range post.Comments {
	// 	fmt.Printf(" - 评论: %s\n", comment.Content)
	// }
	err := config.DB.
		Preload("Comments").Find(&post).Error
	if err != nil {
		fmt.Println("查询评论数量最多的文章失败:", err)
		return
	}

	var maxCommentPost models.Post
	maxComments := -1
	for _, p := range post {
		if len(p.Comments) > maxComments {
			maxComments = len(p.Comments)
			maxCommentPost = p
		}
	}

	fmt.Printf("评论数量最多的文章:\n")
	fmt.Printf("文章标题: %s\n", maxCommentPost.Title)
	fmt.Printf("评论数量: %d\n", len(maxCommentPost.Comments))
	for _, comment := range maxCommentPost.Comments {
		fmt.Printf(" - 评论: %s\n", comment.Content)
	}
	fmt.Println("--------------------")
}

func queryStatus() {
	var posts []models.Post
	err := config.DB.Find(&posts).Error
	if err != nil {
		fmt.Println("查询文章失败:", err)
		return
	}

	type postStatus string
	const (
		HasComment postStatus = "有评论"
		NoComment  postStatus = "无评论"
	)
	for _, post := range posts {
		status := NoComment
		if post.CommentStatus {
			status = HasComment
		}
		fmt.Printf("文章标题: %s, 评论状态: %v\n", post.Title, status)
	}
	fmt.Println("--------------------")

	var user []models.User
	err = config.DB.Find(&user).Error
	if err != nil {
		fmt.Println("查询用户失败:", err)
		return
	}

	for _, u := range user {
		fmt.Printf("用户姓名: %s, 文章数量: %d\n", u.Name, u.PostCount)
	}
	fmt.Println("--------------------")
}

func deleteDemo() {
	// 删除一篇文章，观察用户的文章数量是否更新
	var post models.Post
	err := config.DB.Where("title = ?", "Go學習筆記").Take(&post).Error
	if err != nil {
		fmt.Println("查询文章失败:", err)
		return
	}

	err = config.DB.Delete(&post).Error
	if err != nil {
		fmt.Println("删除文章失败:", err)
		return
	}
	fmt.Println("删除文章成功")

	var user models.User
	err = config.DB.Where("id = ?", post.UserID).Take(&user).Error
	if err != nil {
		fmt.Println("查询用户失败:", err)
		return
	}
	fmt.Printf("用户 %s 的文章数量更新为: %d\n", user.Name, user.PostCount)
	fmt.Println("--------------------")
}

func main() {
	config.CreatDB()
	config.InitDB()
	migrate()
	seedData()
	queryUserPostsWithComments("Alice")
	queryMostCommentedPost()
	queryStatus()
	deleteDemo()
	queryStatus()
}
