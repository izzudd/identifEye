from torch.utils.data import Dataset, DataLoader, random_split
from glob import glob
import os
from PIL import Image
from facenet_pytorch import MTCNN
import random

class CasiaDataset(Dataset):
  def __init__(self, path) -> None:
    super().__init__()

    self.subjects = glob(os.path.join(path, '*'))
    self.mtcnn = MTCNN(margin=5)

  def __len__(self):
    return len(self.subjects)

  def __getitem__(self, idx):
    anchor_subject = self.subjects[idx]
    anchor_images = glob(os.path.join(anchor_subject, '*'))

    anchor_image, positive_image = random.choices(anchor_images, k=2)
    anchor_image = Image.open(anchor_image)
    positive_image = Image.open(positive_image)

    negative_subject = random.choice(self.subjects)
    while(negative_subject == anchor_subject):
      negative_subject = random.choice(self.subjects)

    negative_images = glob(os.path.join(negative_subject, '*'))
    negative_image = Image.open(random.choice(negative_images))

    return (
      self.mtcnn(anchor_image), 
      self.mtcnn(positive_image), 
      self.mtcnn(negative_image)
    )

def get_datalaoder(batch_size=32, num_workers=6):
  train, test = random_split(CasiaDataset('/home/daffaizzuddin/identifEye/training/dataset/casia-webface'), (.8, .2))
  return (
    DataLoader(train, batch_size=batch_size, num_workers=num_workers, shuffle=True),
    DataLoader(test, batch_size=batch_size, num_workers=num_workers)
  )

