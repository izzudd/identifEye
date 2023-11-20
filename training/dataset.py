from torch.utils.data import Dataset, DataLoader
from glob import glob
import os
from PIL import Image
from facenet_pytorch import MTCNN
import random
from sklearn.model_selection import train_test_split

class CasiaDataset(Dataset):
  def __init__(self, paths) -> None:
    super().__init__()

    self.subjects = paths
    self.mtcnn = MTCNN(margin=5, image_size=299)

  def __len__(self):
    return len(self.subjects)

  def __getitem__(self, idx):
    anchor_subject = self.subjects[idx]
    anchor_images = glob(os.path.join(anchor_subject, '*'))

    while True:
      try:
        anchor_image, positive_image = random.choices(anchor_images, k=2)
        anchor_image = Image.open(anchor_image)
        positive_image = Image.open(positive_image)

        anchor_image = self.mtcnn(anchor_image)
        positive_image = self.mtcnn(positive_image)

        if anchor_image is not None and positive_image is not None:
          break
      except Exception as e:
        print(f"Error processing images: {e}")

    negative_subject = random.choice(self.subjects)
    while(negative_subject == anchor_subject):
      negative_subject = random.choice(self.subjects)

    negative_images = glob(os.path.join(negative_subject, '*'))

    while True:
      try:
        negative_image = Image.open(random.choice(negative_images))
        negative_image = self.mtcnn(negative_image)

        if negative_image is not None:
          break
      except Exception as e:
        print(f"Error processing negative image: {e}")

    return anchor_image, positive_image, negative_image


def get_dataloader(batch_size=32, num_workers=6):
  subjects = glob('/home/daffaizzuddin/identifEye/training/dataset/casia-webface/*')
  train, test = train_test_split(subjects, test_size=.2)

  return (
    DataLoader(CasiaDataset(train), batch_size=batch_size, num_workers=num_workers, shuffle=True),
    DataLoader(CasiaDataset(test), batch_size=batch_size, num_workers=num_workers)
  )