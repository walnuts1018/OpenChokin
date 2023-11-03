package psql

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/walnuts1018/openchokin/back/config"
)

const (
	sslmode = "disable"
)

type DB struct {
	db *sqlx.DB
}
func dbInit() error {
    db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v sslmode=%v", config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresAdminUser, config.Config.PostgresAdminPassword, sslmode))
    if err != nil {
        return fmt.Errorf("failed to open db: %w", err)
    }
	defer db.Close()

    var dbName string
    err = db.Get(&dbName, "SELECT datname FROM pg_database WHERE datname = $1", config.Config.PostgresDb)
    if err != nil && err != sql.ErrNoRows {
        return fmt.Errorf("error checking for database existence: %w", err)
    }

    // If the database does not exist, create it
    if dbName == "" {
        _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %v", config.Config.PostgresDb))
		if err != nil {
			return fmt.Errorf("failed to create db: %w", err)
		}
    }

    return nil
}

func NewDB() (*DB, error) {
	err := dbInit()
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

    db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v sslmode=%v", config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresAdminUser, config.Config.PostgresAdminPassword, sslmode))
    if err != nil {
        return nil, fmt.Errorf("failed to open db: %w", err)
    }

	// SQLファイルからテーブルを作成
	err = executeSQLFile(db, "/app/infra/psql/init.sql")
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func executeSQLFile(db *sqlx.DB, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open SQL file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sqlStatement string
	var inDOBlock bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// コメントを無視
		if strings.HasPrefix(trimmedLine, "--") {
			continue
		}

		// DOブロックの開始を検出
		if strings.HasPrefix(trimmedLine, "DO") {
			inDOBlock = true
		}

		// DOブロック内では、END; まで読み込む
		if inDOBlock && strings.HasPrefix(trimmedLine, "END;") {
			inDOBlock = false
		}

		sqlStatement += line + "\n" // SQLステートメントを行ごとに追加

		// SQLステートメントが終わったかどうか（セミコロンかDOブロックの終わり）
		if (!inDOBlock && strings.HasSuffix(trimmedLine, ";")) || (!inDOBlock && strings.HasPrefix(trimmedLine, "END;")) {
			_, err = db.Exec(sqlStatement)
			if err != nil {
				return fmt.Errorf("failed to exec SQL statement: %w", err)
			}
			sqlStatement = "" // ステートメントをリセット
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading SQL file: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
