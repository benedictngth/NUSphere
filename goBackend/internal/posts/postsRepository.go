package posts

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

	// "github.com/jackc/pgx/v5/pgxpool"
	"goBackend/internal/categories"
	"goBackend/internal/common"
	"goBackend/internal/users"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post Post) error
	GetPosts(ctx context.Context) ([]Post, error)
	GetPostByPublicID(ctx context.Context, publicID string) (Post, error)
	EditPostByPublicID(ctx context.Context, publicID string, post Post) error
	DeletePostByPublicID(ctx context.Context, publicID string) error
}

func CreatePost(pg *common.Postgres, ctx context.Context, post Post) error {
	//get user and category by public id with ID fields
	user, userErr := users.GetUserByPublicID(pg, ctx, post.UserID)
	category, categoryErr := categories.GetCategoryByPublicID(pg, ctx, post.CategoryID)
	if userErr != nil {
		return fmt.Errorf("unable to get user: %w", userErr)
	}
	if categoryErr != nil {
		return fmt.Errorf("unable to get category: %w", categoryErr)
	}
	query := `INSERT INTO posts (title, content, public_id, user_id, category_id) 
	VALUES (@title, @content, @PublicID, @userID, @categoryID)`
	args := pgx.NamedArgs{
		"title":      post.Title,
		"content":    post.Content,
		"PublicID":   post.ID,     //public nanoid !database id
		"userID":     user.ID,     //database id
		"categoryID": category.ID, //database id
	}
	_, dbErr := pg.DB.Exec(ctx, query, args)
	if dbErr != nil {
		return fmt.Errorf("unable to insert row: %w", dbErr)
	}
	return nil
}

func GetPosts(pg *common.Postgres, ctx context.Context) ([]Post, error) {
	query := "SELECT posts.public_id as \"posts.public_id\", posts.title, posts.content, posts.created_at, posts.updated_at, users.public_id as \"users.public_id\", categories.public_id as \"categories.public_id\" " +
		"FROM posts " +
		"INNER JOIN users ON posts.user_id = users.id " +
		"INNER JOIN categories ON posts.category_id = categories.id"
	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query posts: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Post])
}

func GetPostByPublicID(pg *common.Postgres, ctx context.Context, publicID string) (Post, error) {
	log.Print(publicID)
	query := "SELECT posts.public_id as \"posts.public_id\", posts.title, posts.content, posts.created_at, posts.updated_at, users.public_id as \"users.public_id\", categories.public_id as \"categories.public_id\" FROM posts " +
		"INNER JOIN users ON posts.user_id = users.id " +
		"INNER JOIN categories ON posts.category_id = categories.id " +
		"WHERE posts.public_id = $1"
	rows, err := pg.DB.Query(ctx, query, publicID)
	if err != nil {
		return Post{}, fmt.Errorf("unable to query posts: %w", err)
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[Post])
}

func EditPostByPublicID(pg *common.Postgres, ctx context.Context, publicID string, post Post) error {
	category, categoryErr := categories.GetCategoryByPublicID(pg, ctx, post.CategoryID)
	if categoryErr != nil {
		return fmt.Errorf("unable to get category: %w", categoryErr)
	}
	query := `UPDATE posts SET title = @title, content = @content, updated_at = @updatedAt, category_id = @categoryID WHERE public_id = @postPublicID`
	args := pgx.NamedArgs{
		"title":        post.Title,
		"content":      post.Content,
		"postPublicID": publicID,
		"categoryID":   category.ID,
		"updatedAt":    time.Now().Format(time.RFC3339),
	}
	_, err := pg.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to update row: %w", err)
	}
	return nil
}

func DeletePostByPublicID(pg *common.Postgres, ctx context.Context, publicID string) error {
	query := `DELETE FROM posts WHERE public_id = $1`
	_, err := pg.DB.Exec(ctx, query, publicID)
	if err != nil {
		return fmt.Errorf("unable to delete row: %w", err)
	}
	return nil
}
