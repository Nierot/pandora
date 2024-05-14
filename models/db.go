package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const playerSchema = `
CREATE TABLE IF NOT EXISTS players (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	is_pirate BOOL DEFAULT(false)
);`

const bakkenSchema = `
CREATE TABLE IF NOT EXISTS bakken (
	id INT AUTO_INCREMENT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	reason VARCHAR(500) NOT NULL,
	player_id INT NOT NULL,
	FOREIGN KEY (player_id) REFERENCES players(id)
);`

const blogSchema = `
CREATE TABLE IF NOT EXISTS blog (
	id INT AUTO_INCREMENT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL,
	image TEXT,
	writer_id INT NOT NULL,
	FOREIGN KEY (writer_id) REFERENCES players(id)
);`


var DB *sqlx.DB

func InitDB() {
	connString := viper.GetString("DATABASE_URL")

	if connString == "" {
		panic("DATABASE_URL is not set")
	}

	fmt.Println("database url: ", connString)

	db := sqlx.MustConnect("mysql", connString)

	db.MustExec(playerSchema)
	db.MustExec(bakkenSchema)
	db.MustExec(blogSchema)

	DB = db
}

type BakkenPerPlayer struct {
	PlayerName string `json:"player_name"`
	Amount     int    `json:"amount"`
	IsPirate   bool   `json:"is_pirate"`
}

func GetAmountOfBakkenPerPlayer() ([]BakkenPerPlayer, error) {
	query := `
		SELECT player.name, player.is_pirate, COUNT(bakken.id) as amount
		FROM players as player
		LEFT JOIN bakken as bakken
		ON player.id = bakken.player_id
		GROUP BY player.name
		ORDER BY amount DESC;
	`

	rows, err := DB.Queryx(query)

	if err != nil {
		return nil, err
	}

	amounts := make([]BakkenPerPlayer, 0)

	for rows.Next() {
		var player Player
		var amount int

		err = rows.Scan(&player.Name, &player.IsPirate, &amount)

		if err != nil {
			return nil, err
		}

		amounts = append(amounts, BakkenPerPlayer{
			PlayerName: player.Name,
			IsPirate: player.IsPirate,
			Amount:     amount,
		})
	}

	return amounts, nil
}

type BakWithPlayerName struct {
	ID        int    `db:"id"`
	CreatedAt string `db:"created_at"`
	Reason    string `db:"reason"`
	PlayerName string `db:"name"`
}

func GetBakken() ([]BakWithPlayerName, error) {
	query := `
		SELECT bakken.id, bakken.created_at, bakken.reason, player.name
		FROM bakken
		JOIN players as player
		ON bakken.player_id = player.id;
	`

	rows, err := DB.Queryx(query)

	if err != nil {
		return nil, err
	}

	bakken := make([]BakWithPlayerName, 0)

	for rows.Next() {
		var bak BakWithPlayerName

		err = rows.Scan(&bak.ID, &bak.CreatedAt, &bak.Reason, &bak.PlayerName)

		if err != nil {
			return nil, err
		}

		bakken = append(bakken, bak)
	}

	if err != nil {
		return nil, err
	}

	return bakken, nil
}

func GetLatestBlogEntry() (*BlogEntryWithName, error) {
	query := `
		SELECT blog.*, players.name as writer_name
		FROM blog
		JOIN players
		ON blog.writer_id = players.id
		ORDER BY created_at DESC
		LIMIT 1;
	`
	
	blog := BlogEntryWithName{}

	err := DB.Get(&blog, query)

	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func GetBlogEntry(id int) (*BlogEntryWithName, error) {
	query := `
		SELECT blog.*, players.name as writer_name
		FROM blog
		JOIN players
		ON blog.writer_id = players.id
		WHERE blog.id = ?;
	`

	blog := BlogEntryWithName{}

	DB.Get(&blog, query, id)

	if blog.ID == 0 {
		return nil, fmt.Errorf("blogpost met id %d niet gevonden", id)
	}

	return &blog, nil
}

func GetBlogEntries() ([]BlogEntryWithName, error) {
	query := `
		SELECT blog.*, players.name as writer_name
		FROM blog
		JOIN players
		ON blog.writer_id = players.id
		ORDER BY created_at DESC;
	`
	blogs := make([]BlogEntryWithName, 0)

	err := DB.Select(&blogs, query)

	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func GetLast3BlogEntries() ([]BlogEntry, error) {
	blogs := make([]BlogEntry, 0)

	err := DB.Select(&blogs, "SELECT * FROM blog ORDER BY created_at DESC LIMIT 3")

	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func GetPlayers() ([]Player, error) {
	players := []Player{}

	err := DB.Select(&players, "SELECT * FROM players")

	if err != nil {
		return nil, err
	}

	return players, nil
}

func AddBak(pid int, reason string) error {
	_, err := DB.Exec("INSERT INTO bakken (player_id, reason) VALUES (?, ?)", pid, reason)

	return err
}

func AddBlog(pid int, title string, content string) (int, error) {
	res, err := DB.Exec("INSERT INTO blog (writer_id, title, content) VALUES (?, ?, ?)", pid, title, content)

	if err != nil {
		return 0, err
	}

	idx, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(idx), nil
}

func EditBlog(id int, pid int, title string, content string) error {
	_, err := DB.Exec("UPDATE blog SET writer_id = ?, title = ?, content = ? WHERE id = ?", pid, title, content, id)

	return err
}

type Player struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	IsPirate bool `db:"is_pirate"`
}

type Bakken struct {
	ID        int    `db:"id"`
	CreatedAt string `db:"created_at"`
	Reason    string `db:"reason"`
	PlayerID  int    `db:"player_id"`
}

type BlogEntry struct {
	ID int `db:"id"`
	CreatedAt string `db:"created_at"`
	Title string `db:"title"`
	Content string `db:"content"`
	Image sql.NullString `db:"image"`
	WriterID int `db:"writer_id"`
}

type BlogEntryWithName struct {
	ID int `db:"id"`
	CreatedAt string `db:"created_at"`
	Title string `db:"title"`
	Content string `db:"content"`
	Image sql.NullString `db:"image"`
	WritedID int `db:"writer_id"`
	WriterName string `db:"writer_name"`
}