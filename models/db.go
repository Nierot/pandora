package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const playerSchema = `
CREATE TABLE IF NOT EXISTS players (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(50) NOT NULL
);`

const bakkenSchema = `
CREATE TABLE IF NOT EXISTS bakken (
	id INT AUTO_INCREMENT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	reason VARCHAR(500) NOT NULL,
	player_id INT NOT NULL,
	FOREIGN KEY (player_id) REFERENCES players(id)
);`

var DB *sqlx.DB

func InitDB() {
	connString := viper.GetString("DATABASE_URL")

	if connString == "" {
		panic("DATABASE_URL is not set")
	}

	db, err := sqlx.Connect("mysql", connString)

	if err != nil {
		panic(err)
	}

	db.MustExec(playerSchema)
	db.MustExec(bakkenSchema)

	DB = db
}

type BakkenPerPlayer struct {
	PlayerName string `json:"player_name"`
	Amount     int    `json:"amount"`
}

func GetAmountOfBakkenPerPlayer() ([]BakkenPerPlayer, error) {
	query := `
		SELECT player.name, COUNT(bakken.id) as amount
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

		err = rows.Scan(&player.Name, &amount)

		if err != nil {
			return nil, err
		}

		amounts = append(amounts, BakkenPerPlayer{
			PlayerName: player.Name,
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

type Player struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Bakken struct {
	ID        int    `db:"id"`
	CreatedAt string `db:"created_at"`
	Reason    string `db:"reason"`
	PlayerID  int    `db:"player_id"`
}