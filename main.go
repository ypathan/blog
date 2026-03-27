package main

import (
	"log/slog"
	"net/http"
	"os"

	"yousuf.xyz/blog/database"
	"yousuf.xyz/blog/handlers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:59188")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			slog.Error("error getting session_token from request", "error", err.Error())
			http.Error(w, "Unauthorize Access, Login First", 401)
			return
		}

		csrfToken, err := r.Cookie("csrf_token")
		if err != nil {
			slog.Error("error getting csrf_token from request", "error", err.Error())
			http.Error(w, "Unauthorize Access, Login First", 401)
			return
		}

		if sessionToken.Value == "" || csrfToken.Value == "" {
			slog.Error("either of the required cookie is empty")
			http.Error(w, "Unauthorize Access, Login First", 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr
		userAgent := r.Header.Get("User-Agent")
		method := r.Method
		path := r.URL.Path
		query := r.URL.RawQuery
		protocol := r.Proto
		host := r.Host
		auth := r.Header.Get("Authorization")
		referrer := r.Referer()
		language := r.Header.Get("Accept-Language")

		slog.Info("new request",
			"ip", ip,
			"user_agent", userAgent,
			"method", method,
			"path", path,
			"query", query,
			"protocol", protocol,
			"host", host,
			"auth", auth,
			"referrer", referrer,
			"language", language,
		)

		next.ServeHTTP(w, r)
	})
}

func main() {

	//--------------logging to file---------------------
	// file, err := os.OpenFile("app.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	// file, err := os.OpenFile("app.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// -------------log to std out---------------------
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// initialize database connection
	db := database.ConnectDatabase()
	defer db.Close()

	// passing db conn to handlers
	blogHandler := handlers.NewBlogController(db)
	authHandler := handlers.NewAuthController(db)
	adminHandler := handlers.NewAdminHandler(db)

	//routes
	mux := http.NewServeMux()

	//static file config
	fs := http.FileServer(http.Dir("./static"))

	// public views
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /", blogHandler.ServeIndex)
	mux.HandleFunc("GET /blog/{id}", blogHandler.ServeBlog)
	mux.HandleFunc("GET /admin/login", authHandler.ServeAdminLogin)

	// private apis
	mux.Handle("GET /admin/protected", protected(http.HandlerFunc(adminHandler.AdminPrivate)))

	// private views
	mux.Handle("GET /admin/addblog", protected(http.HandlerFunc(adminHandler.AdminAddBlog)))
	mux.Handle("GET /admin/dashboard", protected(http.HandlerFunc(adminHandler.AdminDashboard)))
	mux.Handle("GET /admin/editblog/{id}", protected(http.HandlerFunc(adminHandler.EditBlog)))

	// auth
	mux.HandleFunc("POST /login", authHandler.LoginUser)
	mux.HandleFunc("POST /register", authHandler.RegisterUser)
	mux.HandleFunc("GET /logout", authHandler.LogoutUser)

	// application related
	mux.Handle("POST /add", protected(http.HandlerFunc(blogHandler.AddNewBlog)))
	mux.Handle("DELETE /delete/{id}", protected(http.HandlerFunc(blogHandler.DeleteBlog)))
	mux.Handle("PUT /update/{id}", protected(http.HandlerFunc(blogHandler.UpdateBlog)))
	mux.HandleFunc("GET /viewall", blogHandler.ViewAllBlogs)
	mux.HandleFunc("GET /view/{id}", blogHandler.ViewBlog)

	// start server
	slog.Info("server started", "port", "8080")
	if err := http.ListenAndServe(":8080", corsMiddleware(loggingMiddleware(mux))); err != nil {
		slog.Error("error starting http server", "message", err.Error())
	}
}
