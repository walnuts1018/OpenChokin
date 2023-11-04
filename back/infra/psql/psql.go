package psql

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/walnuts1018/openchokin/back/config"
)

const (
	sslmode = "disable"
)

func dbInit() error {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v sslmode=%v",
		config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresAdminUser, config.Config.PostgresAdminPassword, sslmode))
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer db.Close()

	// Check if the user exists
	var roleName string
	err = db.Get(&roleName, "SELECT rolname FROM pg_roles WHERE rolname = $1", config.Config.PostgresUser)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for user existence: %w", err)
	}

	// If the user does not exist, create it
	if roleName == "" {
		// Create a new user with a password
		_, err = db.Exec(fmt.Sprintf("CREATE USER %v WITH PASSWORD '%v'",
			config.Config.PostgresUser, config.Config.PostgresPassword))
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Check if the database exists
	var dbName string
	err = db.Get(&dbName, "SELECT datname FROM pg_database WHERE datname = $1", config.Config.PostgresDb)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for database existence: %w", err)
	}

	// If the database does not exist, create it
	if dbName == "" {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %v OWNER %v", config.Config.PostgresDb, config.Config.PostgresUser))
		if err != nil {
			return fmt.Errorf("failed to create db: %w", err)
		}
	}

	return nil
}

func NewDB() (*sqlx.DB, error) {
	err := dbInit()
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresUser, config.Config.PostgresPassword, config.Config.PostgresDb, sslmode))
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// SQLファイルからテーブルを作成
	err = executeSQLFile(db, "/app/infra/psql/init.sql")
	if err != nil {
		return nil, err
	}

	return db, nil
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
