#!/bin/bash

echo "TARGET_IP=$TARGET_IP" >> ./configs/envs/dev.env
#
# for local testing uncomment below strings
#
# echo "LOCALHOST_IP=127.0.0.1" >> ./envs/dev.env
# echo "POSTGRES_EXT_PORT=15432" >> ./envs/dev.env
source ./configs/envs/dev.env

sudo mkdir -p "$PREFIX"/lct2024/prometheus/config \
"$PREFIX"/lct2024/prometheus/data \
"$PREFIX"/lct2024/db \
"$PREFIX"/lct2024/grafana/provisioning \
"$PREFIX"/lct2024/openobserve \
"$PREFIX"/lct2024/vector \
"$PREFIX"/lct2024/models/ve \
"$PREFIX"/lct2024/qdrant \
"$PREFIX"/lct2024/meilisearch \
"$PREFIX"/lct2024/minio/data

#bash cert.sh

chmod -R 777 "$PREFIX"/lct2024/prometheus/config
chmod -R 777 "$PREFIX"/lct2024/prometheus/data
chmod -R 777 "$PREFIX"/lct2024/grafana
chmod -R 777 "$PREFIX"/lct2024/openobserve
chmod -R 777 "$PREFIX"/lct2024/vector
chmod -R 777 "$PREFIX"/lct2024/db
chmod -R 777 "$PREFIX"/lct2024/models
chmod -R 777 "$PREFIX"/lct2024/qdrant

cp configs/prometheus/prometheus.yml "$PREFIX"/lct2024/prometheus/config
cp -R configs/grafana/provisioning/* "$PREFIX"/lct2024/grafana/provisioning
cp -f configs/vector/vector.toml "$PREFIX"/lct2024/vector

FILE_URLS=("https://drive.google.com/uc?export=download&id=1BPZkS2c__PpUk_BSMXwHoDbnO1HQZaD-"
           "https://drive.usercontent.google.com/download?id=1un80YKwZW463S1lgu8BW1rSiBL0ECsdo&confirm=xxx"
           "https://drive.usercontent.google.com/download?id=1jowFLnomsZ3yq0dn_YVGHQV9dgnGAnaP&confirm=xxx")

FILE_NAMES=("${CROP_MODEL}"
            "${ENCODER_MODEL}"
            "${ENCODER_MODEL_}")

# Количество файлов
NUM_FILES=${#FILE_URLS[@]}

# Цикл для проверки и скачивания файлов
for (( i=0; i<$NUM_FILES; i++ )); do
    FILE_URL=${FILE_URLS[$i]}
    FILE_NAME=${FILE_NAMES[$i]}

    # Проверка существования файла
    if [ -f "$FILE_NAME" ]; then
        echo "Файл $FILE_NAME уже существует, пропуск скачивания."
    else
        echo "Файл $FILE_NAME не найден, начинаю скачивание."
        # Скачивание файла с использованием curl
        curl -L -o "$FILE_NAME" "$FILE_URL"
    fi
done

docker compose --env-file ./configs/envs/dev.env up --build -d
sleep 1

echo "done"
