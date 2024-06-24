import os
import shutil

from fastapi import FastAPI, HTTPException, UploadFile
from fb_encoder import FBEncoder
from image_crop import ImageCrop
from logger import logger
from matcher import Matcher
from pydantic import BaseModel, Field
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    cropper_path: str = Field("./yolo.pt", alias="CROP_MODEL")
    encoder_path: str = Field("./sscd_disc_large.torchscript.pt", alias="ENCODER_MODEL")
    batch_size: int = Field(128, alias="BATCH_SIZE")
    qdrant_addr: str = Field("0.0.0.0", alias="QDRANT_HOST")  # noqa: S104
    qdrant_port: int = Field(6333, alias="QDRANT_PORT")
    create_collection: bool = Field(True, alias="CREATE_COLLECTION")


settings = Settings()

cropper = ImageCrop(settings.cropper_path)

# encoder = Encoder("./vit_ddpmm_8gpu_512_torch2_ap31_pattern_condition_first_dgg.pth.tar", batch_size=256)
# encoder = FBEncoder("./sscd_disc_mixup.torchscript.pt")
encoder = FBEncoder(settings.encoder_path, batch_size=settings.batch_size)

matcher = Matcher(
    encoder,
    cropper,
    "fbl",
    1024,
    qdrant_addr=settings.qdrant_addr,
    qdrant_port=settings.qdrant_port,
    create_collection=settings.create_collection,
)


class UploadResponse(BaseModel):
    filename: str
    message: str


class SearchResponse(BaseModel):
    piracy_name: str
    piracy_start_frame: int
    piracy_end_frame: int
    licence_name: str
    license_start_frame: int
    license_end_frame: int


app = FastAPI()


@app.post(
    "/upload_video",
    description="Send video file to add its into vector database",
    response_description="Update database",
    response_model=UploadResponse,
)
async def upload(file: UploadFile):
    try:
        logger.info(f"Recieved file for uploading: {file.filename}")
        frames_dir = "./tmp_frames_upload"

        if os.path.exists(frames_dir):
            shutil.rmtree(frames_dir)
        os.makedirs(frames_dir)
        video_path = file.filename
        with open(video_path, "wb") as f:
            shutil.copyfileobj(file.file, f)

        matcher.load_reference(video_path, frames_dir)
        logger.info(f"File uploaded: {file.filename}")
        os.remove(video_path)
        return UploadResponse(filename=file.filename, message="File uploaded successfully")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post(
    "/find_video",
    description="Send video file to find copyright infringement",
    response_description="Find copyright infringement",
    response_model=SearchResponse,
)
async def search(file: UploadFile):
    try:
        logger.info(f"Recieved file for searching: {file.filename}")
        frames_dir = "./tmp_frames_search"

        if os.path.exists(frames_dir):
            shutil.rmtree(frames_dir)
        os.makedirs(frames_dir)
        video_path = file.filename
        with open(video_path, "wb") as f:
            shutil.copyfileobj(file.file, f)

        res = matcher.search(video_path, frames_dir)
        logger.info(f"File searching finnished: {file.filename}")
        os.remove(video_path)
        return SearchResponse(**res)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
