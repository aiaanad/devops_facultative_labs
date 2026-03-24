package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// структура для конфига
type Config struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Command  string `json:"command"`
}

// базовая директоия именована как мини-докер хех
const baseDir = "/var/lib/mini-docker"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("пример ввода: sudo go run main.go run")
		return
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	}
}

func run() {
	// читаем настройки из файла
	data, _ := os.ReadFile("config.json")
	var cfg Config
	json.Unmarshal(data, &cfg)

	// готовим пути для оверлей-слоев
	contDir := filepath.Join(baseDir, cfg.ID)
	upper := filepath.Join(contDir, "upper")
	work := filepath.Join(contDir, "work")
	merged := filepath.Join(contDir, "merged")
	lower, _ := filepath.Abs("lowdir") // убедись, что папка называется lowdir

	os.MkdirAll(upper, 0755)
	os.MkdirAll(work, 0755)
	os.MkdirAll(merged, 0755)

	// монтируем нижний и верхний слои
	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lower, upper, work)
	if err := syscall.Mount("overlay", merged, "overlay", 0, opts); err != nil {
		fmt.Printf("ошибка монтирования overlay: %v\n", err)
		return
	}

	// запускаем сами себя с флагами изоляции
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("ошибка: %v\n", err)
	}

	syscall.Unmount(merged, 0)
}

func child() {
	// парсим конфиг для настройки среды
	data, _ := os.ReadFile("config.json")
	var cfg Config
	json.Unmarshal(data, &cfg)

	// определяем точку входа
	merged := filepath.Join(baseDir, cfg.ID, "merged")

	// ставим имя хоста
	syscall.Sethostname([]byte(cfg.Hostname))

	// меняем корень на merged слой оверлея
	if err := syscall.Chroot(merged); err != nil {
		fmt.Printf("ошибка chroot: %v\n", err)
		return
	}
	os.Chdir("/") // прыгаем в новый корень, чтобы не остаться "снаружи"

	// монтируем прок для работы ps и топ
	os.MkdirAll("/proc", 0755) // создаем точку монтирования внутри нового корня
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Printf("ошибка монтирования proc: %v\n", err)
	}

	cmd := exec.Command(cfg.Command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("ошибка выполнения команды %s: %v\n", cfg.Command, err)
	}

	syscall.Unmount("/proc", 0)
}
