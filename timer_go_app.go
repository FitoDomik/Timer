package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// ANSI цвета для красивого вывода
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Timer структура для хранения данных таймера
type Timer struct {
	Duration time.Duration
	Message  string
}

// NewTimer создает новый таймер
func NewTimer(seconds int, message string) *Timer {
	if message == "" {
		message = "⏰ Время вышло!"
	}
	return &Timer{
		Duration: time.Duration(seconds) * time.Second,
		Message:  message,
	}
}

// Start запускает таймер с визуальным обратным отсчетом
func (t *Timer) Start() {
	totalSeconds := int(t.Duration.Seconds())
	
	fmt.Printf("%s🚀 Таймер запущен на %d секунд!%s\n", ColorGreen, totalSeconds, ColorReset)
	fmt.Printf("%s💡 Нажми Ctrl+C для отмены%s\n\n", ColorYellow, ColorReset)

	// Обратный отсчет
	for i := totalSeconds; i > 0; i-- {
		// Определяем цвет в зависимости от оставшегося времени
		var color string
		switch {
		case i <= 5:
			color = ColorRed // Последние 5 секунд - красный
		case i <= 10:
			color = ColorYellow // Последние 10 секунд - желтый
		default:
			color = ColorCyan // Обычное время - голубой
		}

		// Форматируем время
		timeStr := formatTime(i)
		
		// Прогресс бар
		progressBar := createProgressBar(totalSeconds-i, totalSeconds, 30)
		
		fmt.Printf("\r%s⏳ Осталось: %s %s%s", 
			color, timeStr, progressBar, ColorReset)
		
		time.Sleep(1 * time.Second)
	}

	// Финальное уведомление
	t.showNotification()
}

// formatTime форматирует секунды в читаемый вид
func formatTime(seconds int) string {
	if seconds >= 3600 { // больше часа
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		secs := seconds % 60
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	} else if seconds >= 60 { // больше минуты
		minutes := seconds / 60
		secs := seconds % 60
		return fmt.Sprintf("%02d:%02d", minutes, secs)
	}
	return fmt.Sprintf("00:%02d", seconds)
}

// createProgressBar создает визуальный прогресс бар
func createProgressBar(current, total, width int) string {
	if total == 0 {
		return ""
	}
	
	progress := float64(current) / float64(total)
	filled := int(progress * float64(width))
	
	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "="
		} else if i == filled {
			bar += ">"
		} else {
			bar += " "
		}
	}
	bar += "]"
	
	percentage := int(progress * 100)
	return fmt.Sprintf("%s %d%%", bar, percentage)
}

// showNotification показывает уведомление о завершении
func (t *Timer) showNotification() {
	fmt.Printf("\n\n")
	
	// Мигающее уведомление
	for i := 0; i < 3; i++ {
		fmt.Printf("\r%s🔔 %s 🔔%s", ColorRed, t.Message, ColorReset)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("\r%s                    %s", ColorReset, strings.Repeat(" ", len(t.Message)+10))
		time.Sleep(500 * time.Millisecond)
	}
	
	fmt.Printf("\r%s🎉 %s 🎉%s\n", ColorGreen, t.Message, ColorReset)
	
	// ASCII арт часов
	fmt.Println(ColorPurple + `
    ⏰ ВРЕМЯ ВЫШЛО! ⏰
    
     12  1  2  3
   11           4
  10      •      5
   9            6
     8  7  6  5
	` + ColorReset)
}

// parseInput обрабатывает пользовательский ввод
func parseInput(input string) (int, error) {
	input = strings.TrimSpace(input)
	
	// Поддержка разных форматов ввода
	if strings.Contains(input, "m") || strings.Contains(input, "min") {
		// Минуты: "5m", "10min"
		numStr := strings.ReplaceAll(input, "m", "")
		numStr = strings.ReplaceAll(numStr, "in", "")
		numStr = strings.TrimSpace(numStr)
		
		minutes, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, fmt.Errorf("некорректный формат минут")
		}
		return minutes * 60, nil
	}
	
	if strings.Contains(input, "h") || strings.Contains(input, "hour") {
		// Часы: "1h", "2hour"
		numStr := strings.ReplaceAll(input, "h", "")
		numStr = strings.ReplaceAll(numStr, "our", "")
		numStr = strings.TrimSpace(numStr)
		
		hours, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, fmt.Errorf("некорректный формат часов")
		}
		return hours * 3600, nil
	}
	
	// Обычные секунды
	seconds, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("введи корректное число")
	}
	
	if seconds <= 0 {
		return 0, fmt.Errorf("время должно быть больше нуля")
	}
	
	if seconds > 86400 { // больше суток
		return 0, fmt.Errorf("максимальное время - 24 часа (86400 секунд)")
	}
	
	return seconds, nil
}

// showHelp показывает справку
func showHelp() {
	fmt.Println(ColorCyan + "📖 Справка по использованию:" + ColorReset)
	fmt.Println("• Просто число: 30 (секунды)")
	fmt.Println("• С минутами: 5m, 10min")
	fmt.Println("• С часами: 1h, 2hour")
	fmt.Println("• Команды: help, quit, exit")
	fmt.Println()
}

func main() {
	// Приветствие
	fmt.Println(ColorBlue + "🕐 Добро пожаловать в Мини-Таймер!" + ColorReset)
	fmt.Println(ColorGreen + "⭐ Введи время и получи уведомление по истечении!" + ColorReset)
	fmt.Println()
	
	showHelp()
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Printf("%sВведи время (или 'help' для справки): %s", ColorWhite, ColorReset)
		
		if !scanner.Scan() {
			break
		}
		
		input := scanner.Text()
		input = strings.TrimSpace(strings.ToLower(input))
		
		// Обработка команд
		switch input {
		case "quit", "exit", "q":
			fmt.Println(ColorGreen + "👋 До свидания!" + ColorReset)
			return
		case "help", "h":
			showHelp()
			continue
		case "":
			continue
		}
		
		// Парсинг времени
		seconds, err := parseInput(input)
		if err != nil {
			fmt.Printf("%s❌ Ошибка: %s%s\n\n", ColorRed, err.Error(), ColorReset)
			continue
		}
		
		// Запрос сообщения (опционально)
		fmt.Printf("%sВведи сообщение для уведомления (или нажми Enter): %s", ColorWhite, ColorReset)
		scanner.Scan()
		message := strings.TrimSpace(scanner.Text())
		
		// Создание и запуск таймера
		timer := NewTimer(seconds, message)
		timer.Start()
		
		// Предложение запустить еще один таймер
		fmt.Printf("\n%s🔄 Запустить еще один таймер? (да/нет): %s", ColorWhite, ColorReset)
		scanner.Scan()
		again := strings.TrimSpace(strings.ToLower(scanner.Text()))
		
		if again != "да" && again != "yes" && again != "y" && again != "" {
			fmt.Println(ColorGreen + "👋 Спасибо за использование таймера!" + ColorReset)
			break
		}
		
		fmt.Println()
	}
	
	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения ввода: %v", err)
	}
}
