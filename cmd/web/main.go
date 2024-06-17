package main

import (
	"database/sql"
	"flag"
	"github.com/Slava02/practiceS24/pkg/models/mysql"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	dsn := flag.String("dsn", "root:ZXASQW!@zxasqw12@/multiuniverse?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &Application{
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		Universe:       &mysql.UniverseModel{DB: db},
		TemplateCache:  templateCache,
		SessionManager: sessionManager,
		Users:          &mysql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  Routes(app),
	}

	app.InfoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.ErrorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
