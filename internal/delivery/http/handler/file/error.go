package file

import "github.com/Aurivena/spond/v2/envelope"

func (h *Handler) ErrorSessionID() *envelope.AppError {
	return h.spond.BuildError(
		envelope.BadRequest,
		"Отсутствует sessionID",
		"Не удалось определить sessionID пользователя",
		"1. Перезагрузите страницу.\n"+
			"2. Обратитесь к создателю ресурса.",
	)
}

func (h *Handler) ErrorID() *envelope.AppError {
	return h.spond.BuildError(
		envelope.BadRequest,
		"Проблема с ID файла",
		"Не удалось определить ID файла",
		"1. Сделайте повторно запрос.",
	)
}

func (h *Handler) ErrorPassword() *envelope.AppError {
	return h.spond.BuildError(
		envelope.BadRequest,
		"Проблема с Password файла",
		"Не удалось определить Password файла",
		"1. Сделайте повторно запрос.",
	)
}

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
