import torch
from PIL import Image
from torch.utils.data import DataLoader, Dataset
from torchvision import transforms


class ImageList(Dataset):
    def __init__(self, image_list):
        Dataset.__init__(self)
        self.image_list = image_list

    def __len__(self):
        return len(self.image_list)

    def __getitem__(self, i):
        x = Image.open(self.image_list[i])
        x = x.convert("RGB")
        normalize = transforms.Normalize(
            mean=[0.485, 0.456, 0.406],
            std=[0.229, 0.224, 0.225],
        )
        small_288 = transforms.Compose(
            [
                transforms.Resize([320, 320]),
                transforms.ToTensor(),
                normalize,
            ]
        )

        return small_288(x)


class FBEncoder:
    def __init__(self, model_path: str, batch_size: int = 1):
        self.batch_size = batch_size
        self.model = torch.jit.load(model_path).eval().cuda()

    def embeddings_one_video(self, image_path_list: str) -> list[list]:
        dataset = ImageList(image_path_list)
        dataloader = DataLoader(dataset, batch_size=128, shuffle=False)
        embeddings = []
        with torch.no_grad():
            for x1 in dataloader:
                x1 = x1.cuda()
                embedding = self.model(x1).detach().tolist()

                embeddings += embedding

        return embeddings
