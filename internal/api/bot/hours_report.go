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
	TotalDrop            int

	StaffCount         int
	PickedItemsPerHour int //Отобранные штуки за час
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

	if b.waitingForFile == nil || !b.waitingForFile[chatID] {
		b.SendMessage(chatID, "Файл получен, но не ожидается. Используй команду <b>/hourlyReport</b> для запуска отчета.", tgbotapi.ModeHTML)
		os.Remove(filepath.Join(os.TempDir(), doc.FileID+filepath.Ext(doc.FileName)))
		return
	}

	ext := filepath.Ext(doc.FileName)
	if ext != ".xlsx" && ext != ".xls" {
		b.SendMessage(chatID, "Поддерживаются только файлы .xlsx и .xls", tgbotapi.ModeHTML)
		return
	}

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

	tempPath := filepath.Join(os.TempDir(), doc.FileID+ext)

	err = downloadFile(fileURL, tempPath)
	if err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		b.SendMessage(chatID, "Не удалось сохранить файл", tgbotapi.ModeHTML)
		return
	}

	if b.tempFilePaths[chatID] != nil {
		for _, oldPath := range b.tempFilePaths[chatID] {
			os.Remove(oldPath)
		}
		delete(b.tempFilePaths, chatID)
	}

	b.tempFilePaths[chatID] = make(map[int]string)
	b.tempFilePaths[chatID][0] = tempPath

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("ошибка загрузки конфигурации: %v", err)
		return
	}

	report, err := b.generateReport(chatID)
	if err != nil {
		log.Printf("Ошибка формирования отчёта: %v", err)
		b.SendMessage(chatID, "Ошибка при обработке файлов", tgbotapi.ModeHTML)
		return
	}

	if err := b.SendMessage(chatID, "Отчет сформирован и отправлен в группу 'Отчеты участков ЕКБ'", tgbotapi.ModeHTML); err != nil {
		log.Printf("Ошибка отправки отчёта пользователю %d: %v", chatID, err)
	} else {
		log.Printf("Отчёт успешно отправлен пользователю %d", chatID)
	}

	if cfg.TargetChatID != 0 {
		if err := b.SendMessage(cfg.TargetChatID, report, tgbotapi.ModeHTML); err != nil {
			log.Printf("Ошибка отправки отчёта в чат %d: %v", cfg.TargetChatID, err)
		} else {
			log.Printf("Отчёт отправлен в чат %d", cfg.TargetChatID)
		}
	}

	os.Remove(tempPath)
	delete(b.tempFilePaths, chatID)
	delete(b.waitingForFile, chatID)
	delete(b.waitingForStaffCount, chatID)
	delete(b.waitingForPickedItems, chatID)
	delete(b.staffCount, chatID)
	delete(b.pickedItems, chatID)

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
	firstPath, _ := b.tempFilePaths[chatID][0]

	firstData, err := readExcelFile(firstPath)
	if err != nil {
		return "", err
	}

	report := HourlyReport{
		Timestamp:          time.Now().Format("02.01 15:00"), // 23.11 03:00
		StaffCount:         b.staffCount[chatID],
		PickedItemsPerHour: b.pickedItems[chatID],
	}

	// Список нужных СД
	requiredSDs := map[string]bool{
		"СЦ МК Дзержинский DS":       true,
		"СЦ Пермь (Новый)":           true,
		"СЦ МК Челябинск":            true,
		"СЦ Тюмень":                  true,
		"DPD region":                 true,
		"СЦ Новосибирск (NEW)":       true,
		"СЦ Омск (Новый)":            true,
		"СЦ МК Екатеринбург":         true,
		"5Post":                      true,
		"СЦ МК Екатеринбург транзит": true,
	}

	var totalOrdersAll int // Общее количество заказов повыбранным СД

	var total98Orders int // отмены

	for sdName, sdData := range firstData {
		if !requiredSDs[sdName] {
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
				total98Orders += statusData.Pieces

			}
		}
	}

	report.TotalDrop = totalOrdersAll - total98Orders

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

<b>Итого падение</b> — <b>%d шт.</b>

<b>Кол-во человек на потоке</b> - <b>%d чел.</b>

<b>Отобрано за час</b> - <b>%d шт.</b>`,
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
		report.StaffCount,
		report.PickedItemsPerHour,
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

	statusCols, err := detectStatusColumns(sheet)
	if err != nil {
		return nil, err
	}

	statusInfo, err := buildStatusInfoMap(sheet, statusCols)
	if err != nil {
		return nil, err
	}

	data := make(map[string]SDData)

	for rowIdx := 3; rowIdx < len(sheet.Rows); rowIdx++ {
		row := sheet.Rows[rowIdx]
		if len(row.Cells) < 5 {
			continue
		}

		sd := strings.TrimSpace(row.Cells[0].String())
		codeSD := strings.TrimSpace(row.Cells[1].String())

		totalOrdersStr := strings.TrimSpace(row.Cells[3].String())
		totalPiecesStr := strings.TrimSpace(row.Cells[4].String())

		totalOrders, err := strconv.Atoi(totalOrdersStr)
		if err != nil {
			totalOrders = 0
		}

		totalPieces, err := strconv.Atoi(totalPiecesStr)
		if err != nil {
			totalPieces = 0
		}

		statusData := make(map[string]StatusData)

		for status, info := range statusInfo {

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
