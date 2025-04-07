package query_service

import (
	"log"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"

	"api/infrastructure/mysql/db"
	dbTest "api/infrastructure/mysql/db/db_test"
	"api/infrastructure/mysql/db/dbgen"
)

var (
	fixtures *testfixtures.Loader
)

func TestMain(m *testing.M) {
	var err error

	// DBの立ち上げ
	resource, pool := dbTest.CreateContainer()
	defer dbTest.CloseContainer(resource, pool)

	// DBへ接続する
	dbCon := dbTest.ConnectDB(resource, pool)
	defer dbCon.Close()

	// テスト用DBをセットアップ
	dbTest.SetupTestDB("../db/schema/schema.sql")

	// テストデータの準備
	fixturePath := "../fixtures"
	fixtures, err = testfixtures.New(
		testfixtures.Database(dbCon),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(fixturePath),
	)
	if err != nil {
		log.Fatalf("failed to load fixtures: %v", err)
	}

	q := dbgen.New(dbCon)
	db.SetQuery(q)
	db.SetDB(dbCon)

	// テスト実行
	m.Run()

	// テスト終了後のクリーンアップ
	if err := fixtures.Load(); err != nil {
		log.Fatalf("failed to load fixtures: %v", err)
	}
}

func setupFixtures(t *testing.T) {
	if err := fixtures.Load(); err != nil {
		t.Fatalf("failed to load fixtures: %v", err)
	}
}
