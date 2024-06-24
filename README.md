# lct2024_copyright

Сервис проверки видеофайлов на нарушение авторских прав

https://i.moscow/lct/hackatons/226971949a4c481199f01fe477801c3d/ru/

сабмит находится в /docs/submit.csv

1.	Ссылка на открытый репозиторий с исходным кодом https://github.com/gulldan/lct2024quadriceps
2.	Ссылка на Яндекс презентацию - https://disk.yandex.ru/i/yW96iscuwRgr4Q 
3.	Ссылка на сопроводительную документацию - https://disk.yandex.ru/i/-9cPILaxPOsXMQ 
4.	Ссылка на прототип (ссылка на гит/ссылка на облако) - https://github.com/gulldan/lct2024quadriceps 
5.	Ссылка на тестовый датасет https://disk.yandex.ru/d/ukCnaYiqZnzjJQ 

### Технический отчет

#### 1. Сервис получения видеоэмбеддингов

Для получения видеоэмбеддингов сервис выполняет следующие этапы:

- Разделение видео на кадры с заданной частотой (fps=1) с помощью ffmpeg.
- Выполнение детекции объектов на кадрах и их обрезка (кроп), на основе YOLO модели.
- Обработка кропов с использованием энкодера, обученного на основе контрастного обучения с использованием аугментаций. В основе обучения используются лоссы SimCLR и кросс-энтропия.
- Проведение тщательной фильтрации результатов поиска, чтобы гарантировать высокое качество и релевантность эмбеддингов.

Этот подход позволяет эффективно извлекать информативные и репрезентативные видеоэмбеддинги, обеспечивая высокую точность в дальнейшем анализе и поиске.

#### 2. Сервис получения аудиоэмбеддингов

Для получения аудиоэмбеддингов из видео используются следующие шаги:

- Аудио из видео нарезается на фрагменты по 10 секунд с шагом в 1 секунду.
- С помощью модели wav2vec2 извлекаются эмбеддинги размерностью 512.
- Вектора эмбеддингов могут быть сохранены в базу данных для дальнейшего использования.
- Для поиска нарушений копирайта проводится следующая процедура:
  - Для каждого аудио вектора (представляющего 10-секундный фрагмент) находятся ближайшие вектора в базе данных.
  - Вектора проходят фильтрацию по порогу косинусного сходства (cos sim) 0.96.
  - Отбираются только те вектора, которые принадлежат наиболее часто встречаемой аудиозаписи среди найденных ближайших векторов.
  - После этого выделяются таймкоды (фрагменты по 10 секунд) конкретной записи, основываясь на метаданных о векторах.

Данный процесс позволяет эффективно выявлять и отслеживать нарушения копирайта в аудиозаписях, обеспечивая точность и надежность в обнаружении соответствий.

back (postgresql, qdrant, minio)

front

vector, openobserve, grafana, prometheus
