package file

import (
	"DragDrop-Files/internal/domain"
	"errors"

	"github.com/Aurivena/spond/v2/envelope"
)

func (h *Handler) ErrorSystem() *envelope.AppError {
	return h.spond.BuildError(
		envelope.InternalServerError,
		"Ошибка на уровне сервера",
		"Произошла системная ошибка. Сообщите администратору",
		"1. Перезагрузите страницу.\n"+
			"2. Обратитесь к создателю ресурса.",
	)
}

func (h *Handler) ErrorParse() *envelope.AppError {
	return h.spond.BuildError(
		envelope.InternalServerError,
		"Не удалось обработать запрос пользователя",
		"Ошибка при обработке данных от пользователя",
		"1. Сделайте повторный запрос.\n"+
			"2. Обратитесь к администратору.",
	)
}

func (h *Handler) BadRequest() *envelope.AppError {
	return h.spond.BuildError(
		envelope.BadRequest,
		"Некорректный запрос",
		"Ваш запрос не является корректным",
		"1. Проверьте корректность введённых данных.\n"+
			"2. Попробуйте повторить запрос с исправленными значениями.",
	)
}

func (h *Handler) NotFound() *envelope.AppError {
	return h.spond.BuildError(
		envelope.NotFound,
		"Файл не найден",
		"Не удалось найти файл по указанному id",
		"1. Перепроверьте указанный идентификатор.\n"+
			"2. Обратитесь к создателю ресурса.",
	)
}

func (h *Handler) Gone() *envelope.AppError {
	return h.spond.BuildError(
		envelope.ResourceInTrash,
		"Файл удалён",
		"Данный файл был удалён и больше не доступен",
		"1. Уточните у владельца.\n"+
			"2. Попробуйте загрузить другой ресурс.",
	)
}

func (h *Handler) InternalServerError() *envelope.AppError {
	return h.spond.BuildError(
		envelope.InternalServerError,
		"Внутренняя ошибка",
		"Что-то пошло не так при обработке файла",
		"1. Повторите попытку позже.\n"+
			"2. Сообщите администратору, если ошибка повторяется.",
	)
}

func (h *Handler) PasswordInvalid() *envelope.AppError {
	return h.spond.BuildError(
		envelope.BadRequest,
		"Ошибка при запросе",
		"Проблема при проверке пароля",
		"1. Убедитесь что вы ввели правильный пароль.",
	)
}

func (h *Handler) httpError(err error) *envelope.AppError {
	switch {
	case errors.Is(err, domain.BadRequestError):
		return h.BadRequest()
	case errors.Is(err, domain.NotFoundError):
		return h.NotFound()
	case errors.Is(err, domain.GoneError):
		return h.Gone()
	case errors.Is(err, domain.PasswordInvalidError):
		return h.PasswordInvalid()
	case errors.Is(err, domain.InternalError):
		return h.InternalServerError()
	default:
		return h.InternalServerError()
	}
}
