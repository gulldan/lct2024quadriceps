import logging
import subprocess
import time
from concurrent.futures import ProcessPoolExecutor, as_completed
from pathlib import Path

logging.basicConfig(level=logging.INFO)


class VideoClipper:
    def __init__(self, audioclips_save_path: str) -> None:
        self.audioclips_save_path = Path(audioclips_save_path)
        self.audioclips_save_path.mkdir(parents=True, exist_ok=True)

    def clip_audio(
        self,
        audio_path: str,
        audio_duration: int = 10,
        step: int = 1,
        sample_rate: int = 16000,
    ) -> None:
        start_time = time.time()
        audio_path = Path(audio_path)
        audio_name = audio_path.stem

        total_duration = self._get_audio_duration(audio_path)
        segment_times = self._calculate_segment_times(total_duration, audio_duration, step)
        self._create_clips(audio_path, audio_name, segment_times, audio_duration, sample_rate)

        end_time = time.time()
        logging.info(f"Script completed in {end_time - start_time:.2f} seconds")

    def _calculate_segment_times(self, total_duration: float, audio_duration: int, step: int) -> list:
        segment_times = []
        for start_time in range(0, int(total_duration) - audio_duration + 1, step):
            segment_times.append(start_time)
        return segment_times

    def _create_clips(
        self, audio_path: Path, audio_name: str, segment_times: list, audio_duration: int, sample_rate: int
    ) -> None:
        with ProcessPoolExecutor() as executor:
            futures = [
                executor.submit(self._run_ffmpeg_command, audio_path, audio_name, start_time, audio_duration, sample_rate)
                for start_time in segment_times
            ]
            for future in as_completed(futures):
                if future.exception() is not None:
                    logging.error(f"Exception during processing: {future.exception()}")

    def _run_ffmpeg_command(
        self, audio_path: Path, audio_name: str, start_time: int, audio_duration: int, sample_rate: int
    ) -> None:
        end_time = start_time + audio_duration
        subclip_name = f"{audio_name}_{start_time:04}_{end_time:04}.wav"
        output_path = self.audioclips_save_path / subclip_name

        ffmpeg_command = [
            "ffmpeg",
            "-ss",
            str(start_time),
            "-t",
            str(audio_duration),
            "-i",
            str(audio_path),
            "-af",
            f"aresample={sample_rate}",
            "-c:a",
            "pcm_s16le",
            str(output_path),
            "-y",  # Overwrite output files without asking
        ]

        logging.info(f"Running ffmpeg command: {' '.join(ffmpeg_command)}")
        result = subprocess.run(ffmpeg_command, capture_output=True, text=True, check=False)
        if result.returncode != 0:
            logging.error(f"ffmpeg command failed with error: {result.stderr}")
            raise RuntimeError(f"ffmpeg command failed with error: {result.stderr}")

    def _get_audio_duration(self, audio_path: Path) -> float:
        ffprobe_command = [
            "ffprobe",
            "-v",
            "error",
            "-show_entries",
            "format=duration",
            "-of",
            "default=noprint_wrappers=1:nokey=1",
            str(audio_path),
        ]

        result = subprocess.run(ffprobe_command, capture_output=True, text=True, check=False)
        if result.returncode != 0:
            logging.error(f"ffprobe command failed with error: {result.stderr}")
            raise RuntimeError(f"ffprobe command failed with error: {result.stderr}")

        return float(result.stdout.strip())


if __name__ == "__main__":
    clipper = VideoClipper("./audioclips")
    clipper.clip_audio(
        audio_path="./audio/The-Pretty-Reckless-Make-Me-Wanna-Die.wav", audio_duration=10, step=1, sample_rate=16000
    )
