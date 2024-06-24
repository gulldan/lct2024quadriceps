import os
import subprocess

from fastapi import FastAPI, HTTPException, UploadFile
from pydantic import BaseModel, conlist
from qdrant_client_api import QdrantClientApi
from video_clipper import VideoClipper

from wav2vec import Wav2Vec

# os.environ["QDRANT_HOST"] = "127.0.0.1"
# os.environ["QDRANT_PORT"] = "6333"
# os.environ["DEVICE"] = "cuda"

QDRANT_HOST = os.getenv("QDRANT_HOST")
QDRANT_PORT = os.getenv("QDRANT_PORT")
DEVICE = os.getenv("DEVICE")
RECREATE = os.getenv("CREATE_COLLECTION")


class EmbeddingResponse(BaseModel):
    embedding: conlist(float, max_length=512, min_length=512)


class UpdateDatabaseAnswer(BaseModel):
    response: str


class CopyrightAnswer(BaseModel):
    ID_piracy_wav2vec: str
    segment_wav2vec: str
    ID_license_wav2vec: str
    segments_wav2vec: str


app = FastAPI()

qdrant_client = QdrantClientApi(QDRANT_HOST, QDRANT_PORT, create_collection=RECREATE)
audio_clips_save_path = "clipped_audio"
videoclip_client = VideoClipper(audio_clips_save_path)
wav2vec = Wav2Vec(qdrant_client, videoclip_client, device=DEVICE)


@app.post(
    "/exctract_embedding",
    description="Send .wav(string bytes) file to exctract embedding",
    response_description="Audio embbedings",
    response_model=EmbeddingResponse,
)
async def exctract_embedding(audio_file: UploadFile):
    try:
        audio_save_path = f"audio/{audio_file.filename}"

        with open(audio_save_path, "wb") as buffer:
            buffer.write(audio_file.file.read())

        embedding = wav2vec.exctract_embedding(
            audio_save_path,
        )

        subprocess.run(f"rm -rf {audio_save_path}", shell=True, check=False)
        subprocess.run(f"rm -rf {audio_clips_save_path}", shell=True, check=False)
        subprocess.run(f"mkdir {audio_clips_save_path}", shell=True, check=False)

        return {"embedding": embedding}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post(
    "/update_database",
    description="Send .wav(string bytes) file to add its into vector database",
    response_description="Update database",
    response_model=UpdateDatabaseAnswer,
)
async def update_data_base(audio_file: UploadFile):
    try:
        audio_save_path = f"audio/{audio_file.filename}"

        with open(audio_save_path, "wb") as buffer:
            buffer.write(audio_file.file.read())

        wav2vec.wav2vec_update_database(audio_save_path)

        subprocess.run(f"rm -rf {audio_save_path}", shell=True, check=False)
        subprocess.run(f"rm -rf {audio_clips_save_path}", shell=True, check=False)
        subprocess.run(f"mkdir {audio_clips_save_path}", shell=True, check=False)

        return {"response": "video was uploaded"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post(
    "/find_copyright_infringement",
    description="Send .wav(string bytes) file to find copyright infringement",
    response_description="Find copyright infringement",
    response_model=CopyrightAnswer,
)
async def find_copyright_infringement(audio_file: UploadFile):
    try:
        audio_save_path = f"audio/{audio_file.filename}"

        with open(audio_save_path, "wb") as buffer:
            buffer.write(audio_file.file.read())

        answer = wav2vec.process_search_results(wav2vec.wav2vec_find_copyright_infringement(audio_save_path))

        subprocess.run(f"rm -rf {audio_save_path}", shell=True, check=False)
        subprocess.run(f"rm -rf {audio_clips_save_path}", shell=True, check=False)
        subprocess.run(f"mkdir {audio_clips_save_path}", shell=True, check=False)

        return answer
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
