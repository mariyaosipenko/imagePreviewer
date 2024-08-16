package getimage

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"image"
	"image-resize/internal/lib/logger/sl"
	"image/jpeg"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// ImageGetter is an interface for getting image by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ImageGetter
type ImageGetter interface {
	Download(name string) ([]byte, string, error)
	Upload(name string, b []byte, contentType string) (string, error)
}

func New(log *slog.Logger, storage ImageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getimage.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		//url := chi.URLParam(r, "*")

		width, _ := strconv.Atoi(chi.URLParam(r, "width"))
		height, _ := strconv.Atoi(chi.URLParam(r, "height"))

		path := r.URL.Path
		// отрезаю цифры размеров из запроса
		url := strings.TrimPrefix(strings.TrimPrefix(path, fmt.Sprintf("/%v", width)), fmt.Sprintf("/%v/", height))
		name := strings.Replace(url, "/", "_", -1)
		bCh := make(chan []byte)

		//есть ли у нас уже готовая картинка
		img, mimeType, err := storage.Download(fmt.Sprintf("r_%v_%v_%v", width, height, name))

		if err != nil {
			go func(bCh chan []byte) {
				// но может быть есть оригинал?
				imgOrig, mimeType, err := storage.Download(fmt.Sprintf("o_%v", name))
				if err != nil {
					// ладно, скачиваем
					client := http.Client{}
					req, err := http.NewRequest("GET", fmt.Sprintf("http://%v", url), nil)
					if err != nil {
						log.Error("failed to create request", sl.Err(err))
					}
					for k, vv := range r.Header {
						for _, v := range vv {
							req.Header.Add(k, v)
						}
					}
					response, err := client.Do(req)
					if err != nil {
						log.Error("failed to download image", sl.Err(err))
					}
					defer response.Body.Close()
					// сохраняем оригинал
					imgOrig, err = io.ReadAll(response.Body)
					mimeType = http.DetectContentType(imgOrig)
					_, err = storage.Upload(fmt.Sprintf("o_%v", name), imgOrig, mimeType)
					if err != nil {
						log.Error("failed to upload image to minio", sl.Err(err))
					}
				}
				// изменение размера
				img, err = resizeImage(imgOrig, width, height)
				if err != nil {
					log.Error("failed to  resize image ", sl.Err(err))
				}
				// сохраняем ресайз
				_, err = storage.Upload(fmt.Sprintf("r_%v_%v_%v", width, height, name), img, mimeType)
				if err != nil {
					log.Error("failed to upload resized image to minio ", sl.Err(err))
				}

				bCh <- img
			}(bCh)

		} else {
			w.Header().Set("Content-Type", mimeType)
			w.Header().Set("Content-Length", strconv.Itoa(len(img)))
			if _, err = w.Write(img); err != nil {
				log.Error("failed to return image", sl.Err(err))
			}
		}

		resultImage := <-bCh
		// возвращаем что получилось
		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Length", strconv.Itoa(len(resultImage)))
		if _, err = w.Write(resultImage); err != nil {
			log.Error("failed to return image", sl.Err(err))
		}

	}
}

func resizeImage(originalImageBytes []byte, width, height int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(originalImageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, imaging.Resize(img, width, height, imaging.Lanczos), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %w", err)
	}

	return buf.Bytes(), nil
}
