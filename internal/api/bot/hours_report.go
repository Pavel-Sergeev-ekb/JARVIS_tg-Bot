package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Pavel-Sergeev-ekb/JARVIS_tg-Bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tealeg/xlsx"
)

type HourlyReport struct {
	Timestamp            string // Дата и время (формат: 29.10 15:00)
	UnknownCount         int    // Статус -1
	BacklogReplenishment int    // Статус -3
	BacklogPickup        int    // 02 + 29 + 52 (общий бэклог отбора)
	BacklogPickupKGT     int    // 02(кгт) + 29(кгт) + 52(кгт)
	BacklogPacking       int    // 55 + 59 + 61(общий бэклог упаковки)
	BacklogPackingKGT    int    // 55(кгт) + 59(кгт)
	BacklogSorting       int    // 65
	TotalBacklog         int    // Сумма всех бэклогов + неизвестно (-1)
	TotalProcessed       int    // 68 + 95
	TotalDrop            int    // Сумма: dpd + сц мк екатеринбург + 5post (файл 1) + сц мк екатеринбург транзит (файл 2) − 98 (файл 1 и файл 2)
}

type StatusInfo struct {
	OrdersCol int //индекс столбца "Заказы"
	PiecesCol int // индекс столбца "Штуки"
	KGTCol    int // индекс столбца "КГТ"

}
type StatusData struct {
	Orders int // Количество заказов по статусу
	Pieces int // Количество штук по статусу
	KGT    int // Признак КГТ (0 или 1)
}

type SDData struct {
	CodeSD      string                // Код СД (столбец 1)
	TotalOrders int                   // Всего заказов (столбец 3)
	TotalPieces int                   // Всего штук (столбец 4)
	StatusData  map[string]StatusData // Данные по статусам
}

func (b *Bot) HandleFileUpload(chatID int64, doc *tgbotapi.Document) {
	ext := filepath.Ext(doc.FileName)
	if ext != ".xlsx" && ext != ".xls" {
		b.SendMessage(chatID, "Поддерживаются только файлы .xlsx и .xls", tgbotapi.ModeHTML)
		return
	}

	// Получаем файл из Telegram
	fileConfig := tgbotapi.FileConfig{FileID: doc.FileID}
	file, err := b.BotAPI.GetFile(fileConfig)
	if err != nil {
		log.Printf("Ошибка получения файла: %v", err)
		b.SendMessage(chatID, "Не удалось загрузить файл", tgbotapi.ModeHTML)
		return
	}

	if file.FilePath == "" {
		log.Printf("Пустой FilePath для file_id=%s", doc.FileID)
		b.SendMessage(chatID, "Не удалось получить путь к файлу", tgbotapi.ModeHTML)
		return
	}

	fileURL := fmt.Sprintf(
		"https://api.telegram.org/file/bot%s/%s",
		b.BotAPI.Token,
		file.FilePath,
	)
	// Временный путь
	tempPath := filepath.Join(os.TempDir(), doc.FileID+ext)

	// Сохраняем файл
	err = downloadFile(fileURL, tempPath)
	if err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		b.SendMessage(chatID, "Не удалось сохранить файл", tgbotapi.ModeHTML)
		return
	}

	// Инициализируем запись в tempFilePaths, если её нет
	if b.tempFilePaths[chatID] == nil {
		b.tempFilePaths[chatID] = make(map[int]string)
	}

	// Определяем, какой это файл
	if b.waitingForFirstFile[chatID] {
		b.tempFilePaths[chatID][0] = tempPath
		b.waitingForFirstFile[chatID] = false
		b.waitingForSecondFile[chatID] = true
		b.SendMessage(chatID, "Первый файл получен. Отправьте второй Excel-файл.", tgbotapi.ModeHTML)
	} else if b.waitingForSecondFile[chatID] {
		b.tempFilePaths[chatID][1] = tempPath
		b.waitingForSecondFile[chatID] = false

		// Формируем отчёт
		report, err := b.generateReport(chatID)
		if err != nil {
			log.Printf("Ошибка формирования отчёта: %v", err)
			b.SendMessage(chatID, "Ошибка при обработке файлов", tgbotapi.ModeHTML)
		} else {
			b.SendMessage(chatID, report, tgbotapi.ModeHTML)
		}

		// Очищаем временные данные
		delete(b.tempFilePaths, chatID)
		delete(b.waitingForFirstFile, chatID)
		delete(b.waitingForSecondFile, chatID)
	} else {
		// Файл отправлен вне ожидаемого контекста
		b.SendMessage(chatID, "Файл получен, но не ожидается. Используйте команду <b>/hourlyReport</b> для начала загрузки.", tgbotapi.ModeHTML)
		os.Remove(tempPath) // Удаляем незатребованный файл
	}
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func detectStatusColumns(sheet *xlsx.Sheet) (map[string]int, error) {
	if len(sheet.Rows) < 2 {
		return nil, fmt.Errorf("недостаточно строк для определения статусов")
	}
	statusRow := sheet.Rows[1] // строка 1 с кодами статусов
	statusCols := make(map[string]int)

	for colIdx, cell := range statusRow.Cells {
		status := strings.TrimSpace(cell.String())
		if status == "" {
			continue
		}
		if _, err := strconv.Atoi(status); err == nil {
			statusCols[status] = colIdx
		}
	}
	return statusCols, nil
}

func buildStatusInfoMap(sheet *xlsx.Sheet, statusCols map[string]int) (map[string]StatusInfo, error) {
	statusInfo := make(map[string]StatusInfo)

	for status, baseCol := range statusCols {
		// Проверяем, что есть как минимум 3 столбца после baseCol
		if baseCol+2 >= len(sheet.Rows[3].Cells) { // строка 3 — подзаголовки
			log.Printf("Недостаточно столбцов для статуса %s", status)
			continue
		}

		statusInfo[status] = StatusInfo{
			OrdersCol: baseCol,
			PiecesCol: baseCol + 1,
			KGTCol:    baseCol + 2,
		}
	}

	return statusInfo, nil
}

func (b *Bot) generateReport(chatID int64) (string, error) {
	firstPath, ok1 := b.tempFilePaths[chatID][0]
	secondPath, ok2 := b.tempFilePaths[chatID][1]
	if !ok1 || !ok2 {
		return "", fmt.Errorf("не все файлы загружены")
	}

	firstData, err := readExcelFile(firstPath)
	if err != nil {
		return "", err
	}
	secondData, err := readExcelFile(secondPath)
	if err != nil {
		return "", err
	}

	report := HourlyReport{
		Timestamp: time.Now().Format("02.01 15:00"), // 23.11 03:00
	}

	// Список нужных СД
	requiredSDsFirst := map[string]bool{
		"DPD region":         true,
		"СЦ МК Екатеринбург": true,
		"5Post":              true,
	}
	transitSDName := "СЦ МК Екатеринбург транзит"

	var totalOrdersAll int // Общее количество заказов по выбранным СД

	// 1. Обрабатываем первый файл: нужные СД
	for sdName, sdData := range firstData {
		if !requiredSDsFirst[sdName] {
			continue
		}

		totalOrdersAll += sdData.TotalPieces

		for status, statusData := range sdData.StatusData {

			pieces := statusData.Pieces

			switch status {
			case "-1":
				report.UnknownCount += pieces
			case "-3":
				report.BacklogReplenishment += pieces
			case "02", "29", "52":
				report.BacklogPickup += pieces
				report.BacklogPickupKGT += statusData.KGT
			case "55", "59", "61":
				report.BacklogPacking += pieces
				report.BacklogPackingKGT += statusData.KGT
			case "65":
				report.BacklogSorting += pieces
			case "68", "95":
				report.TotalProcessed += pieces
			case "98":
				// Пока не суммируем, сделаем это позже
			}
		}
	}

	// 2. Обрабатываем второй файл: только транзитный СД
	if sdData, ok := secondData[transitSDName]; ok {
		totalOrdersAll += sdData.TotalPieces

		for status, statusData := range sdData.StatusData {

			pieces := statusData.Pieces

			// Обрабатываем все статусы (как в первом файле)
			switch status {
			case "-1":
				report.UnknownCount += pieces // Теперь учитывается!
			case "-3":
				report.BacklogReplenishment += pieces
			case "02", "29", "52":
				report.BacklogPickup += pieces
				report.BacklogPickupKGT += statusData.KGT

			case "55", "59", "61":
				report.BacklogPacking += pieces
				report.BacklogPackingKGT += statusData.KGT

			case "65":
				report.BacklogSorting += pieces
			case "68", "95":
				report.TotalProcessed += pieces
			case "98":
				// Пока не суммируем, сделаем это позже
			}
		}
	} else {
		log.Printf("СД '%s' не найден во втором файле", transitSDName)
	}

	// 3. Суммируем заказы со статусом 98 из обоих файлов
	var total98Orders int
	// Первый файл
	for sdName, sdData := range firstData {
		if requiredSDsFirst[sdName] {
			if statusData, ok := sdData.StatusData["98"]; ok {
				total98Orders += statusData.Pieces
			}
		}
	}
	// Второй файл
	if sdData, ok := secondData[transitSDName]; ok {
		if statusData, ok := sdData.StatusData["98"]; ok {
			total98Orders += statusData.Pieces
		}
	}

	// 4. Рассчитываем итоговое падение
	report.TotalDrop = totalOrdersAll - total98Orders

	// Или, если нужно показать только количество отмен:
	// report.TotalDrop = total98Orders

	// Итоговый бэклог
	report.TotalBacklog = report.BacklogReplenishment +
		report.BacklogPickup +
		report.BacklogPacking +
		report.BacklogSorting +
		report.UnknownCount

		// Формируем текст отчёта
	result := fmt.Sprintf(`<b>Часовой отчёт ИСХОД (%s)</b>
───────────────────────────────────────────────

<b>Статус «неизвестно»</b> — <b>%d шт.</b>

<b>Бэклог пополнения</b> — <b>%d шт.</b>

<b>Бэклог отбора</b> — <b>%d шт.</b>

<b>Бэклог отбора КГТ</b> — <b>%d шт.</b>

<b>Бэклог упаковки</b> — <b>%d шт.</b>

<b>Бэклог упаковки КГТ</b> — <b>%d шт.</b>

<b>Бэклог сортировки</b> — <b>%d шт.</b>

<b>Бэклог по всем статусам до сортировки по СД</b> — <b>%d шт.</b>

<b>Итого обработано</b> — <b>%d шт.</b>

<b>Итого падение</b> — <b>%d шт.</b>`,
		report.Timestamp,
		report.UnknownCount,
		report.BacklogReplenishment,
		report.BacklogPickup,
		report.BacklogPickupKGT,
		report.BacklogPacking,
		report.BacklogPackingKGT,
		report.BacklogSorting,
		report.TotalBacklog,
		report.TotalProcessed,
		report.TotalDrop,
	)

	return result, nil
}

func readExcelFile(path string) (map[string]SDData, error) {
	xlsxFile, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла %s: %v", path, err)
	}

	sheet := xlsxFile.Sheets[0]
	if len(sheet.Rows) < 4 {
		return nil, fmt.Errorf("недостаточно строк в файле (требуется минимум 4)")
	}

	// 1. Определяем столбцы статусов (строка 1, индекс 1)
	statusCols, err := detectStatusColumns(sheet)
	if err != nil {
		return nil, err
	}

	// 2. Строим карту индексов для "Заказы"/!Штуки"/!КГТ"
	statusInfo, err := buildStatusInfoMap(sheet, statusCols)
	if err != nil {
		return nil, err
	}

	data := make(map[string]SDData)

	// 3. Читаем строки данных (начиная с 4‑й строки, индекс 3)
	for rowIdx := 3; rowIdx < len(sheet.Rows); rowIdx++ {
		row := sheet.Rows[rowIdx]
		if len(row.Cells) < 5 { // минимум 5 столбцов
			continue
		}

		sd := strings.TrimSpace(row.Cells[0].String())
		codeSD := strings.TrimSpace(row.Cells[1].String())

		// Читаем "Всего заказов" и "Всего штук"
		totalOrdersStr := strings.TrimSpace(row.Cells[3].String())
		totalPiecesStr := strings.TrimSpace(row.Cells[4].String())

		totalOrders, err := strconv.Atoi(totalOrdersStr)
		if err != nil {
			totalOrders = 0 // если не число — ставим 0
		}

		totalPieces, err := strconv.Atoi(totalPiecesStr)
		if err != nil {
			totalPieces = 0
		}

		statusData := make(map[string]StatusData)

		// 4. Для каждого обнаруженного статуса читаем его данные
		for status, info := range statusInfo {
			// Проверяем, что столбцы существуют
			if info.OrdersCol >= len(row.Cells) ||
				info.PiecesCol >= len(row.Cells) ||
				info.KGTCol >= len(row.Cells) {
				continue
			}

			ordersStr := strings.TrimSpace(row.Cells[info.OrdersCol].String())
			piecesStr := strings.TrimSpace(row.Cells[info.PiecesCol].String())
			kgtStr := strings.TrimSpace(row.Cells[info.KGTCol].String())

			orders, _ := strconv.Atoi(ordersStr)
			pieces, _ := strconv.Atoi(piecesStr)
			kgt, _ := strconv.Atoi(kgtStr)

			statusData[status] = StatusData{
				Orders: orders,
				Pieces: pieces,
				KGT:    kgt,
			}
		}

		data[sd] = SDData{
			CodeSD:      codeSD,
			TotalOrders: totalOrders,
			TotalPieces: totalPieces,
			StatusData:  statusData,
		}
	}

	return data, nil
}

func (b *Bot) sendReport(chatID int64) error {
	// 1. Генерируем отчёт из файлов пользователя sourceChatID
	report, err := b.generateReport(chatID)
	if err != nil {
		return fmt.Errorf("ошибка генерации отчёта: %w", err)
	}

	// 2. Отправляем в фиксированный целевой чат
	err = b.SendMessage(TargetChatID, report, tgbotapi.ModeHTML)
	if err != nil {
		return fmt.Errorf("ошибка отправки отчёта в чат %d: %w", TargetChatID, err)
	}

	return nil // Успешно отправлено
}

func (b *Bot) downloadTelegramFile(chatID int64, doc *tgbotapi.Document, fileIndex int) error {
	// Получаем информацию о файле через API Telegram
	fileConfig := tgbotapi.FileConfig{
		FileID: doc.FileID,
	}

	file, err := b.BotAPI.GetFile(fileConfig)
	if err != nil {
		return err
	}

	// Формируем URL для скачивания (используем токен из конфига)
	fileURL := file.Link(config.Cfg.BotToken)

	// Локальный путь сохранения файла
	localPath := fmt.Sprintf("./uploads/%d_%d_%s", chatID, fileIndex, doc.FileName)

	// Скачиваем файл с помощью вашей функции
	err = downloadFile(fileURL, localPath)
	if err != nil {
		return err
	}

	// Сохраняем путь во вложенную мапу
	if b.tempFilePaths[chatID] == nil {
		b.tempFilePaths[chatID] = make(map[int]string)
	}
	b.tempFilePaths[chatID][fileIndex] = localPath

	return nil
}
