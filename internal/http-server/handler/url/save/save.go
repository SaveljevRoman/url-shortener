package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
)

const aliasLength = 3

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", "handlers.url.save.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("ошибка декода тела request", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("декодированый request", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil { // todo лишние ошибки
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.Error("невалидный запрос"))
			if validateErr, ok := err.(validator.ValidationErrors); ok {
				render.JSON(w, r, response.ValidationError(validateErr))
			}
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := saver.SaveURL(req.URL, alias)
		if err != nil {
			log.Error(err.Error(), slog.String("url", req.URL))
			render.JSON(w, r, response.Error("ошибка сохранения url"))
			return
		}

		log.Info("url сохранен", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
