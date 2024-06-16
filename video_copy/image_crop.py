from PIL import Image
from ultralytics import YOLO


class ImageCrop:
    def __init__(
        self,
        model_path: str,
    ) -> None:
        self.model = YOLO(model_path)

    def crop_images(self, image_path_list: list[str]) -> tuple[list[list], list[str]]:
        frame_crops = []
        frame_paths = []
        counter = 0

        for image_path in image_path_list:
            img = Image.open(image_path)
            result = self.model(img, verbose=False)
            bboxes = result[0].boxes.xyxy.cpu().tolist()
            frame_crops_list = [counter]
            frame_paths.append(image_path)
            counter += 1
            for id, bbox in enumerate(bboxes):
                crop_path = f"{image_path.rsplit('.', 1)[0]}_{id+1}.{image_path.split('.')[-1]}"
                img = img.crop((bbox[0], bbox[1], bbox[2], bbox[3]))
                img.save(crop_path)
                frame_crops_list.append(counter)
                frame_paths.append(crop_path)
                counter += 1

            frame_crops.append(frame_crops_list)

        return frame_crops, frame_paths
