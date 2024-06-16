import os
import pickle
import subprocess

import numpy as np
import polars as pl
from logger import logger
from qdrant_client import QdrantClient, models
from qdrant_client.models import Distance, VectorParams
from scipy.ndimage import uniform_filter1d


class Matcher:
    def __init__(
        self,
        encoder,
        cropper,
        collection_name="ref",
        vector_dim: int = 512,
        qdrant_addr: str = "10.35.56.10",
        qdrant_port: int = 6333,
        create_collection: bool = True
    ) -> None:
        self.collection_name = collection_name
        self.encoder = encoder
        self.cropper = cropper
        self.client = QdrantClient(qdrant_addr, port=qdrant_port, timeout=60)
        if create_collection:
            self.client.create_collection(
                collection_name=self.collection_name,
                vectors_config=VectorParams(size=vector_dim, distance=Distance.COSINE),
                optimizers_config=models.OptimizersConfigDiff(memmap_threshold=10000),
            )

    @staticmethod
    def extract_frames(input_video, output_directory):
        if not os.path.exists(output_directory):
            os.makedirs(output_directory)

        command = ["ffmpeg", "-i", input_video, "-vf", "fps=1", os.path.join(output_directory, "%d.jpg")]

        subprocess.run(command, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, check=False)

        frame_names = sorted(os.listdir(output_directory), key=lambda x: int(x.split(".")[0]))
        frame_paths = [os.path.join(output_directory, name) for name in frame_names]
        return frame_paths

    def load_reference(self, video_path: str, frames_dir: str) -> None:
        logger.info(f"Extracting frames to {frames_dir}")
        frame_paths = self.extract_frames(video_path, frames_dir)
        logger.info("Getting embeddings")
        embeddings = self.encoder.embeddings_one_video(frame_paths)

        payload = list([{"frame": i + 1, "file_name": video_path.split("/")[-1]} for i in range(len(embeddings))])
        num_points = self.client.get_collection(collection_name=self.collection_name).points_count
        logger.info(f"Uploading {len(embeddings)} embeddings")
        self.client.upload_collection(
            collection_name=self.collection_name,
            vectors=embeddings,
            payload=payload,
            ids=list(range(num_points, num_points + len(embeddings))),
        )

    @staticmethod
    def exponential_smoothing(data, alpha):
        smoothed_data = np.zeros(data.shape)
        smoothed_data[0] = data[0]

        for i in range(1, len(data)):
            smoothed_data[i] = alpha * data[i] + (1 - alpha) * smoothed_data[i - 1]

        return smoothed_data

    def search(self, video_path: str, frames_dir: str) -> dict:
        logger.info(f"Extracting frames to {frames_dir}")
        frame_paths = self.extract_frames(video_path, frames_dir)
        logger.info("Getting cropps")
        frame_crops, frame_paths = self.cropper.crop_images(frame_paths)
        logger.info("Getting embeddings")
        embeddings = self.encoder.embeddings_one_video(frame_paths)
        scores = []

        for id, frame_crop in enumerate(frame_crops):
            max_score_ans = None
            for crop in frame_crop:
                ans = self.client.search(
                    collection_name=self.collection_name,
                    query_vector=embeddings[crop],
                    limit=1,
                )
                # if not max_score_ans or ans[0].score > max_score_ans.score:
                max_score_ans = ans[0]
                max_score_ans = max_score_ans.model_dump()
                max_score_ans = {
                    "score": max_score_ans["score"],
                    "license_frame": max_score_ans["payload"]["frame"],
                    "file_name": max_score_ans["payload"]["file_name"],
                    "piracy_frame": id + 1,
                }
                scores.append(max_score_ans)
        scores = pl.DataFrame(scores)

        scores_list = scores["score"].to_numpy()

        ###
        # smoothed = self.exponential_smoothing(scores_list, 0.2)

        # window_size = 3
        # cumulative_sum = np.cumsum(smoothed)
        # sliding_average = (cumulative_sum[window_size - 1:] - cumulative_sum[:-(window_size - 1)]) / window_size
        # smoothed = np.pad(sliding_average, (window_size - 1, 0), mode='constant')
        #
        # smoothed = (smoothed - np.min(smoothed)) / (np.max(smoothed) - np.min(smoothed))
        ####

        # th = np.quantile(scores_list, 0.4)
        threshold = 0.4

        th_filtered = (scores_list > threshold).astype(int)

        detected_scores = scores.with_columns(pl.Series(name="th", values=th_filtered))

        detected_scores = detected_scores.filter(pl.col("th") == 1)

        probable_copyright_videos = detected_scores["file_name"].value_counts().filter(pl.col("count") > 25)

        def is_increasing(probable_copyright_video: str) -> bool:
            copyright_scores = detected_scores.filter(pl.col("file_name") == probable_copyright_video)
            Q1 = copyright_scores["license_frame"].quantile(0.25)
            Q3 = copyright_scores["license_frame"].quantile(0.75)

            IQR = Q3 - Q1

            lower_bound = Q1 - 1.0 * IQR
            upper_bound = Q3 + 1.0 * IQR

            df = copyright_scores.filter(
                (copyright_scores["license_frame"] >= lower_bound) & (copyright_scores["license_frame"] <= upper_bound)
            )
            smoothed = uniform_filter1d(df["license_frame"], size=15)
            smoothed = self.exponential_smoothing(smoothed, 0.05)
            smoothed = uniform_filter1d(smoothed, size=10)
            return np.all(np.diff(smoothed) >= 0)

        probable_copyright_videos = probable_copyright_videos.with_columns(
            pl.col("file_name").map_elements(is_increasing, return_dtype=pl.Boolean).alias("is_increasing")
        )
        copyright_video = probable_copyright_videos.filter(pl.col("is_increasing") == 1).sort("count", descending=True).head(1)

        if copyright_video.is_empty():
            return {
                "piracy_name": video_path.split("/")[-1],
                "piracy_start_frame": -1,
                "piracy_end_frame": -1,
                "licence_name": "",
                "license_start_frame": -1,
                "license_end_frame": -1,
            }
        copyright_video = copyright_video["file_name"].item()
        copyright_scores = detected_scores.filter(pl.col("file_name") == copyright_video)

        q1 = copyright_scores["license_frame"].quantile(0.25)
        q3 = copyright_scores["license_frame"].quantile(0.75)

        iqr = q3 - q1

        lower_bound = q1 - 1.0 * iqr
        upper_bound = q3 + 1.0 * iqr

        df = copyright_scores.filter(
            (copyright_scores["license_frame"] >= lower_bound) & (copyright_scores["license_frame"] <= upper_bound)
        )
        license_start_frame = df["license_frame"].min()
        license_end_frame = df["license_frame"].max()

        filtered_piracy_frames = copyright_scores.filter(
            pl.col("license_frame").is_between(license_start_frame, license_end_frame)
        )["piracy_frame"]
        piracy_start_frame = filtered_piracy_frames.min()
        piracy_end_frame = filtered_piracy_frames.max()

        return {
            "piracy_name": video_path.split("/")[-1],
            "piracy_start_frame": piracy_start_frame,
            "piracy_end_frame": piracy_end_frame,
            "licence_name": copyright_video,
            "license_start_frame": license_start_frame,
            "license_end_frame": license_end_frame,
        }

    def dump_embeddings(self, video_path: str, frames_dir: str, dump_dir: str) -> None:
        dump_dir = os.path.join(dump_dir, video_path.split("/")[-1])
        if os.path.exists(dump_dir):
            # shutil.rmtree(dump_dir)
            return
        os.makedirs(dump_dir)

        frame_paths = self.extract_frames(video_path, frames_dir)

        frame_crops, frame_paths = self.cropper.crop_images(frame_paths)

        embeddings = self.encoder.embeddings_one_video(frame_paths)

        with open(os.path.join(dump_dir, "embeddings.pickle"), "wb") as f:
            pickle.dump(embeddings, f)

        with open(os.path.join(dump_dir, "frame_crops.pickle"), "wb") as f:
            pickle.dump(frame_crops, f)
