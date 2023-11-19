import io
import sqlite3
import numpy as np

def adapt_array(arr):
    """
    http://stackoverflow.com/a/31312102/190597 (SoulNibbler)
    """
    out = io.BytesIO()
    np.save(out, arr)
    out.seek(0)
    return sqlite3.Binary(out.read())

def convert_array(text):
    out = io.BytesIO(text)
    out.seek(0)
    return np.load(out)

# Converts np.array to TEXT when inserting
sqlite3.register_adapter(np.ndarray, adapt_array)
# Converts TEXT to np.array when selecting
sqlite3.register_converter("EMBEDDING", convert_array)

def create_table(cursor):
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS embedding (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            key TEXT,
            embedding EMBEDDING
        )
    ''')

def insert_embeddings(cursor, key, embeddings):
    for embedding in embeddings:
        cursor.execute('INSERT INTO embedding (key, embedding) VALUES (?, ?)', (key, embedding))

def retrieve_embeddings(cursor, key):
    cursor.execute('SELECT embedding FROM embedding WHERE key = ?', (key,))
    results = cursor.fetchall()
    return np.array([result[0] for result in results])

def embedding_count(cursor, key):
    cursor.execute('SELECT COUNT(*) FROM embedding WHERE key = ?', (key,))
    row_count = cursor.fetchone()[0]
    return row_count

def db_connection():
    conn = sqlite3.connect("embedding.db", detect_types=sqlite3.PARSE_DECLTYPES, check_same_thread=False)
    cursor = conn.cursor()
    create_table(cursor)

    return conn, cursor
