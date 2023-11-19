from glob import glob
import numpy as np
import sys
import os
from generator import get_embeddings
from store import insert_embeddings, db_connection

def main():
  if len(sys.argv) < 3:
    print('error: usage "python register.py <key> <image_path>"')
    exit(1)

  face_path = glob(os.path.join(sys.argv[2], '*'))

  try:
    embeddings = get_embeddings(face_path)
  except Exception as e:
    print(f'error: cannot get image embedding - {e}')
    sys.stdout.flush()
    exit(1)

  # TODO: validate embedding
  print(embeddings.shape)

  try:
    conn, cursor = db_connection()
    insert_embeddings(cursor, sys.argv[1], embeddings)
    conn.commit()
  except Exception as e:
    print(f'error: database error - {e}')
    sys.stdout.flush()
    exit(1)

  print(f'success')
  sys.stdout.flush()
  exit(0)

if __name__ == '__main__':
  main()
  