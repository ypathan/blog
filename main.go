package main

import (
	"log/slog"
	"net/http"
	"os"

	"yousuf.xyz/blog/controller"
	"yousuf.xyz/blog/database"
	"yousuf.xyz/blog/service"
)

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
	//logging
	// file, err := os.OpenFile("app.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	// file, err := os.OpenFile("app.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	// if err != nil {
	// 	panic(err.Error())
	// }
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	// initialize database connection
	db := database.ConnectDatabase()
	defer db.Close()

	// create service with database connection
	blogService := service.NewBlogService(db)

	// create controller with service
	blogController := controller.NewBlogController(blogService)


	//routes
	mux := http.NewServeMux()

	
	//static file config
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", blogController.ServeIndex)
	mux.HandleFunc("POST /add", blogController.AddNewBlog)
	mux.HandleFunc("DELETE /delete/{id}", blogController.DeleteBlog)
	mux.HandleFunc("PUT /update/{id}", blogController.UpdateBlog)
	mux.HandleFunc("GET /viewall", blogController.ViewAllBlogs)
	mux.HandleFunc("GET /view/{id}", blogController.ViewBlog)

	// start sercer
	slog.Info("server started", "port", "8080")
	if err := http.ListenAndServe(":8080", loggingMiddleware(mux)); err != nil {
		slog.Error("error starting http server", "message", err.Error())
	}
}
