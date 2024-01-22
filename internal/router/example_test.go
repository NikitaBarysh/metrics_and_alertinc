package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NikitaBarysh/metrics_and_alertinc/config/server"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/handlers"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/logger"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/filestorage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/repository/memstorage"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service/hasher"
	"github.com/go-chi/chi/v5"
)

func Example() {
	// создаем логгер
	log := logger.NewLoggingVar()

	// читаем переменные окружения
	env, envErr := server.NewServer()
	if envErr != nil {
		log.Fatal("config err: %s\n", envErr)
	}

	//создаем конфиг
	cfg, configError := server.NewConfig(env)
	if configError != nil {
		log.Fatal("config err: %s\n", configError)
	}

	// создаем файл для хранения
	file, err := filestorage.NewFileEngine("data.json")
	if err != nil {
		fmt.Println("err to create")
	}

	//создаем хранилище в кэше
	rep, err := memstorage.NewMemStorage(context.Background(), cfg, file)
	if err != nil {
		log.Fatal("err to create rep")
	}

	//  создаем канал для shutdown
	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	//создаем handler
	handler := handlers.NewHandler(rep, log)
	//создаем роутер
	router := NewRouter(handler)
	chiRouter := chi.NewRouter()
	// если ключ хеширование отсутствует, создаем
	if cfg.Key != "" {
		hasher.Sign = hasher.NewHasher([]byte(cfg.Key))
		chiRouter.Use(hasher.Middleware)
	}
	chiRouter.Mount("/", router.Register())

	// запускаем сервер
	go func() {
		err = http.ListenAndServe(cfg.RunAddr, chiRouter)
		if err != nil {
			panic(err)
		}
	}()

	// ждем сигнала для graceful shutdown
	sig := <-termSig
	fmt.Println("Server Graceful Shutdown", sig.String())
}
