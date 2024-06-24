FROM nvidia/cuda:12.4.0-base-ubuntu22.04 as build

ENV PYTHON_VERSION=3.11

RUN export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y update \
    && apt-get -y install --no-install-recommends \
    python${PYTHON_VERSION} \
    python3-pip python3.11-venv \
    ffmpeg curl \
    git \
    && rm -rf /var/lib/apt/lists/*

RUN ln -s -f /usr/bin/python${PYTHON_VERSION} /usr/bin/python3 && \
    ln -s -f /usr/bin/python${PYTHON_VERSION} /usr/bin/python && \
    ln -s -f /usr/bin/pip3 /usr/bin/pip && \
    python3 -m venv /venv && /venv/bin/pip install --timeout=1200 --no-cache-dir --upgrade pip setuptools wheel

FROM build AS build-venv

COPY requirements.txt /requirements.txt
#RUN /venv/bin/pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu121
RUN /venv/bin/pip install --timeout=12000 --no-cache-dir --disable-pip-version-check -r /requirements.txt

FROM build

COPY --from=build-venv /venv /venv

WORKDIR /app
RUN mkdir -p audio
RUN mkdir -p models/ve
ENV HUGGINGFACE_HUB_CACHE=./models/hf \
    TORCH_HOME=./models/torch \
    CROP_MODEL=./yolo.pt \
    ENCODER_MODEL=./sscd_disc_large.torchscript.pt \
    ENCODER_MODEL_=./vit_ddpmm_8gpu_512_torch2_ap31_pattern_condition_first_dgg.pth.tar \
    QDRANT_HOST="qdrant" \
    QDRANT_PORT="6333" \
    CREATE_COLLECTION=true \
    DEVICE="cuda"

COPY . .

EXPOSE 8000
ENTRYPOINT ["/venv/bin/python3", "-m", "granian", "--interface", "asgi", "main:app", "--port", "8000", "--host", "0.0.0.0"]
# For local test
### granian --interface asgi main:app --port 8000 --host 0.0.0.0