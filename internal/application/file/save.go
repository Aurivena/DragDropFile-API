package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/fileops"
	"DragDrop-Files/pkg/idgen"
	"mime/multipart"
	"runtime"
	"sync"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/sirupsen/logrus"
)

const (
	prefixZipFile = "dg-"
)

var (
	workerPool = runtime.GOMAXPROCS(0)
)

func (a *File) Execute(sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*entity.FileSaveOutput, *envelope.AppError) {
	newFiles := domain.GetNewInfo(files, headers)

	id, existingFiles, err := a.checkFilesID(sessionID)
	if err != nil {
		return nil, a.InternalServerError()
	}

	id, err = domain.SetFileID(id)
	if err != nil {
		logrus.Error("failed to set g ID")
		return nil, a.InternalServerError()
	}

	var (
		processedFiles []entity.File
	)

	prefix, err := idgen.GenerateID()
	if err != nil {
		logrus.Errorf("failed to generate prefix: %v", err)
		return nil, a.InternalServerError()
	}

	jobs := make(chan entity.File)

	pool := domain.Pool{}

	for w := 0; w < workerPool; w++ {
		go pool.Work(jobs, processedFiles, a.processes)
	}

	pool.Add(len(newFiles))
	for _, file := range newFiles {
		file.Prefix = prefix
		file.FileID = id
		file.SessionID = sessionID

		jobs <- file
	}
	close(jobs)
	pool.Wait()

	processedFiles = append(processedFiles, existingFiles...)

	meta, err := a.downloadZipFile(id, sessionID, processedFiles)
	if err != nil {
		logrus.Errorf("failed to create zip: %v", err)
		return nil, a.InternalServerError()
	}

	return &entity.FileSaveOutput{
		ID:    id,
		Size:  meta.Size,
		Count: len(processedFiles),
	}, nil
}

func (a *File) processes(file *entity.File, processedFiles []entity.File) {
	var mu sync.Mutex

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
