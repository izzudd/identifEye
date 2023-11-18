from glob import glob

import torch
from PIL import Image
from facenet_pytorch import MTCNN, InceptionResnetV1

mtcnn = MTCNN(margin=10)
generate_embedding = InceptionResnetV1(pretrained='vggface2').eval()

def read_image(path):
  image = Image.open(path)
  return mtcnn(image)

def read_images(images_path):
  images = []
  for path in images_path:
    image = read_image(path)
    images.append(image)
  return torch.stack(images)

def get_embeddings(images):
  images = read_images(images)
  with torch.no_grad():
    return generate_embedding(images).numpy()
