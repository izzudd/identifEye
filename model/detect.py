from glob import glob
import numpy as np
import sys
import os
from generator import get_embeddings
from store import retrieve_embeddings, db_connection
from sklearn.metrics.pairwise import cosine_similarity

def main():
  if len(sys.argv) < 2:
    print('error: usage "python detect.py <key> <image_path>"')
    exit(1)

  face_path = glob(os.path.join(sys.argv[2], '*'))

  try:
    embeddings = get_embeddings(face_path)
  except Exception as e:
    print(f'error: cannot get image embedding - {e}')
    sys.stdout.flush()
    exit(1)

  try:
    _, cursor = db_connection()
    target_embedding = retrieve_embeddings(cursor, sys.argv[1])
  except Exception as e:
    print(f'error: database error - {e}')
    sys.stdout.flush()
    exit(1)

  similarity = np.max(cosine_similarity(embeddings, target_embedding))

  print(f'success: {similarity}')
  sys.stdout.flush()
  exit(0)

if __name__ == '__main__':
  main()
  