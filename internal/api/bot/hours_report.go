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

	if b.tempFilePaths[chatID] == nil {
		b.tempFilePaths[chatID] = make(map[int]string)
	}

	if b.waitingForFirstFile[chatID] {
		b.tempFilePaths[chatID][0] = tempPath
		b.waitingForFirstFile[chatID] = false
		b.waitingForSecondFile[chatID] = true
		b.SendMessage(chatID, "Первый файл получен. Отправьте второй Excel-файл.", tgbotapi.ModeHTML)
	} else if b.waitingForSecondFile[chatID] {
		b.tempFilePaths[chatID][1] = tempPath
		b.waitingForSecondFile[chatID] = false

		report, err := b.generateReport(chatID)
		if err != nil {
			log.Printf("Ошибка формирования отчёта: %v", err)
			b.SendMessage(chatID, "Ошибка при обработке файлов", tgbotapi.ModeHTML)
		} else {
			b.SendMessage(chatID, report, tgbotapi.ModeHTML)
		}

		delete(b.tempFilePaths, chatID)
		delete(b.waitingForFirstFile, chatID)
		delete(b.waitingForSecondFile, chatID)
		delete(b.waitingForStaffCount, chatID)
		delete(b.waitingForPickedItems, chatID)
		delete(b.staffCount, chatID)
		delete(b.pickedItems, chatID)
	} else {

		b.SendMessage(chatID, "Файл получен, но не ожидается. Используйте команду <b>/hourlyReport</b> для начала загрузки.", tgbotapi.ModeHTML)
		os.Remove(tempPath)
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
		Timestamp:          time.Now().Format("02.01 15:00"), // 23.11 03:00
		StaffCount:         b.staffCount[chatID],
		PickedItemsPerHour: b.pickedItems[chatID],
	}

	// Список нужных СД
	requiredSDsFirst := map[string]bool{
		"DPD region":         true,
		"СЦ МК Екатеринбург": true,
		"5Post":              true,
	}
	requiredSDsSecond := map[string]bool{
		"СЦ Москва транзит":          true,
		"СЦ Пермь транзит":           true,
		"СЦ Челябинск транзит":       true,
		"СЦ Тюмень транзит":          true,
		"СЦ Омск транзит":            true,
		"СЦ Новосибирск транзит":     true,
		"СЦ МК Екатеринбург транзит": true,
	}

	var totalOrdersAll int // Общее количество заказов по выбранным СД

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

			}
		}
	}

	for sdName, sdData := range secondData {
		if !requiredSDsSecond[sdName] {
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

			}
		}
	}

	var total98Orders int

	for sdName, sdData := range firstData {
		if requiredSDsFirst[sdName] {
			if statusData, ok := sdData.StatusData["98"]; ok {
				total98Orders += statusData.Pieces
			}
		}
	}

	for sdName, sdData := range secondData {
		if requiredSDsSecond[sdName] {
			if statusData, ok := sdData.StatusData["98"]; ok {
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
