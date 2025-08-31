package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/fileops"
	"mime/multipart"
	"runtime"
	"sync"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	prefixZipFile = "dg-"
)

var (
	maxParallelDefault = runtime.GOMAXPROCS(0)
)

func (a *File) Execute(sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*entity.FileSaveOutput, *envelope.AppError) {
	//TODO move on delivery level
	if len(files) == 0 || len(files) != len(headers) {
		return nil, a.BadRequest("1. Ваша сессия недействительна\n" + "2. Длина загруженных файлов == 0")
	}

	newFiles := domain.GetNewInfo(files, headers)

	id, existingFiles, err := a.checkFilesID(sessionID)
	if err != nil {
		logrus.Error("failed to check files ID")
		return nil, a.InternalServerError()
	}

	id, err = domain.SetFileID(id)
	if err != nil {
		logrus.Error("failed to set g ID")
		return nil, a.InternalServerError()
	}

	var (
		wg             sync.WaitGroup
		processedFiles []entity.File
	)

	prefix, err := uuid.NewV7()
	if err != nil {
		logrus.Errorf("failed to generate prefix: %v", err)
		return nil, a.InternalServerError()
	}

	for _, file := range newFiles {
		wg.Add(1)
		file.Prefix = prefix.String()
		file.FileID = id
		file.SessionID = sessionID

		go a.processes(&file, processedFiles, &wg)
	}
	wg.Wait()

	processedFiles = append(processedFiles, existingFiles...)

	meta, err := a.downloadZipFile(id, sessionID, prefixZipFile, processedFiles)
	if err != nil {
		logrus.Errorf("failed to create zip g: %v", err)
		return nil, a.InternalServerError()
	}

	return &entity.FileSaveOutput{
		ID:    id,
		Size:  meta.Size,
		Count: len(processedFiles),
	}, nil
}

func (a *File) processes(file *entity.File, processedFiles []entity.File, wg *sync.WaitGroup) {
	var mu sync.Mutex
	defer wg.Done()

	data, err := fileops.DecodeFile(file.FileBase64)
	if err != nil {
		return
	}

	fileValid := entity.File{
		FileID:    file.FileID,
		Name:      file.Name,
		SessionID: file.SessionID,
		MimeType:  domain.GetMimeType(file.FileBase64),
	}

	if !a.validDownloadFile(data, fileValid, file.Prefix) {
		return
	}

	mu.Lock()
	processedFiles = append(processedFiles, *file)
	mu.Unlock()
}
