package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/fileops"
	"DragDrop-Files/pkg/idgen"
	"mime/multipart"

	"github.com/sirupsen/logrus"
)

func (a *File) Execute(sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*entity.FileSaveOutput, error) {
	newFiles := domain.GetNewInfo(files, headers)

	id, existingFiles, err := a.checkFilesID(sessionID)
	if err != nil {
		return nil, domain.InternalError
	}

	results := a.workerPool(newFiles, id, sessionID)

	var processedFiles []entity.File
	for f := range results {
		processedFiles = append(processedFiles, f)
	}

	processedFiles = append(processedFiles, existingFiles...)

	out, err := a.downloadZipFile(id, sessionID, processedFiles)
	if err != nil {
		return nil, domain.InternalError
	}

	return out, nil
}

func (a *File) workerPool(newFiles []entity.File, fileID, sessionID string) chan entity.File {
	prefix, err := idgen.GenerateID()
	if err != nil {
		logrus.Errorf("failed to generate prefix: %v", err)
		return nil
	}

	jobs := make(chan entity.File)
	results := make(chan entity.File, len(newFiles))

	pool := domain.Pool{}
	pool.Run(domain.WorkerPool, jobs, func(f *entity.File) {
		if ok := a.processes(f); ok {
			results <- *f
		}
	})

	go func() {
		for _, f := range newFiles {
			f.Prefix = prefix
			f.FileID = fileID
			f.SessionID = sessionID
			jobs <- f
		}
		close(jobs)
	}()
	go func() {
		pool.Wait()
		close(results)
	}()

	return results
}

func (a *File) processes(file *entity.File) bool {
	data, err := fileops.DecodeFile(file.FileBase64)
	if err != nil {
		return false
	}

	file.MimeType = domain.GetMimeType(file.FileBase64)

	if !a.validDownloadFile(data, file) {
		return false
	}

	return true
}
