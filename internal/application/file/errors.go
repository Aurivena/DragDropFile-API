package file

import (
	"github.com/Aurivena/spond/v2/envelope"
)

func (a *File) NotFound() *envelope.AppError {
	return a.spond.BuildError(
		envelope.NotFound,
		"Файл не найден",
		"Не удалось найти файл по указанному id",
		"1. Перепроверьте указанный идентификатор.\n"+
			"2. Обратитесь к создателю ресурса.",
	)
}

func (a *File) Gone() *envelope.AppError {
	return a.spond.BuildError(
		envelope.ResourceInTrash,
		"Файл удалён",
		"Данный файл был удалён и больше не доступен",
		"1. Уточните у владельца.\n"+
			"2. Попробуйте загрузить другой ресурс.",
	)
}

func (a *File) InternalServerError() *envelope.AppError {
	return a.spond.BuildError(
		envelope.InternalServerError,
		"Внутренняя ошибка",
		"Что-то пошло не так при обработке файла",
		"1. Повторите попытку позже.\n"+
			"2. Сообщите администратору, если ошибка повторяется.",
	)
}
