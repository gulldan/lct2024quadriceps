import torch
import torchvision
import torchvision.transforms
from model import VisionTransformer
from PIL import Image
from torch import nn
from torch.utils.data import DataLoader, Dataset
from tqdm import tqdm

torch.set_float32_matmul_precision("high")


class ImageList(Dataset):
    def __init__(self, image_list, imsize=None) -> None:
        Dataset.__init__(self)
        self.image_list = image_list
        mean, std = [0.485, 0.456, 0.406], [0.229, 0.224, 0.225]
        transforms = [
            torchvision.transforms.ToTensor(),
            torchvision.transforms.Normalize(mean, std),
        ]
        self.transforms = torchvision.transforms.Compose(transforms)

        self.imsize = imsize

    def __len__(self) -> int:
        return len(self.image_list)

    def __getitem__(self, i):
        x = Image.open(self.image_list[i])
        x = x.convert("RGB")

        if self.imsize is not None:
            x = x.resize((self.imsize, self.imsize))

        x = self.transforms(x)

        return x


class Encoder:
    def __init__(self, model_path: str, batch_size: int = 1) -> None:
        model = VisionTransformer(num_features=0, dropout=0, num_classes=4)
        model = nn.DataParallel(model)
        ckpt = torch.load(model_path)
        model.load_state_dict(ckpt, strict=False)
        model.classifier = None
        model.classifier_1 = None
        self.model = model.eval().cuda()
        self.batch_size = batch_size

    def embeddings_one_video(self, image_path_list: str) -> list[list]:
        im_dataset = ImageList(image_path_list, imsize=224)

        dataloader = DataLoader(im_dataset, batch_size=self.batch_size, shuffle=False)
        embeddings = []

        with torch.no_grad():
            for x in tqdm(dataloader):
                x = x.cuda()
                embedding, _ = self.model(x)
                embedding = embedding.detach().tolist()
                embeddings += embedding
        return embeddings
