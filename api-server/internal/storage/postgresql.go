package storage

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/ivan/blockchain/api-server/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type PostgreSQL struct {
	db     *sql.DB
	config string
}

func NewPostgresql(cnf string) *PostgreSQL {
	var pstg = PostgreSQL{}
	pstg.config = cnf
	return &pstg
}

func (pstg *PostgreSQL) Connect() error {
	var err error
	pstg.db, err = sql.Open("pgx", pstg.config)
	if err != nil {
		log.Println("error connecting to Postgresql ", err)
		return err
	}

	err = pstg.PingContext()
	if err != nil {
		log.Println("error pinging Postgresql ", err)
		return err
	}

	return nil
}

func (pstg *PostgreSQL) Migrations(pathFiles string) error {
	driver, err := postgres.WithInstance(pstg.db, &postgres.Config{})
	if err != nil {
		log.Printf("error creating postgres driver: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(pathFiles, "postgres", driver)
	if err != nil {
		log.Println("error migrate Postgresql", err)
		return err
	}

	if err = m.Up(); err != nil {
		log.Println("error up Postgresql", err)
		return err
	}

	return nil
}

func (pstg *PostgreSQL) Close() error {
	if pstg.db != nil {
		return pstg.db.Close()
	}
	return nil
}

func (pstg *PostgreSQL) PingContext() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return pstg.db.PingContext(ctx)
}

func (pstg *PostgreSQL) AddBlock(block model.Block) error {
	var err error = nil
	_, err = pstg.db.ExecContext(context.Background(), "INSERT INTO public.blockchain (key, hash, merkley, noce, data) "+
		"VALUES ($1, $2, $3, $4, $5)",
		block.Head.Key,
		block.Head.Hash,
		block.Head.Merkley,
		block.Head.Noce,
		block.Data,
	)

	return err
}

func (pstg *PostgreSQL) GetBlock(ctx context.Context, key string) (*model.Block, error) {
	var (
		err     error = nil
		hash    string
		merkley string
		noce    string
		data    string
	)
	ctxDB, cancelDB := context.WithCancel(ctx)
	defer cancelDB()

	row := pstg.db.QueryRowContext(ctxDB, "SELECT hash, merkley, noce, data FROM public.blockchain WHERE key=$1", key)
	err = row.Scan(&hash, &merkley, &noce, &data)
	if err != nil {
		return nil, err
	}

	return &model.Block{Data: data, Head: model.BlockHeader{
		Hash:    hash,
		Merkley: merkley,
		Noce:    noce,
		Key:     key,
	}}, err
}
