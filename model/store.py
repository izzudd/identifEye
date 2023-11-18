import sqlite3
import numpy as np

def create_table(cursor):
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS embedding (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            key INTEGER
            embedding BLOB
        )
    ''')

def insert_embedding(cursor, key, embedding):
    cursor.execute('INSERT INTO embedding (key, embedding) VALUES (?)', (key, embedding.tobytes()))

def retrieve_embedding(cursor, key):
    cursor.execute('SELECT embedding FROM embedding WHERE key = ?', (key,))
    results = cursor.fetchall()
    return [np.frombuffer(result[0]) for result in results]

def embedding_count(cursor, key):
    cursor.execute('SELECT COUNT(*) FROM embedding WHERE key = ?', (key,))
    row_count = cursor.fetchone()[0]
    return row_count

def db_connection():
    conn = sqlite3.connect('embedding.db')
    cursor = conn.cursor()
    create_table(cursor)

    return conn, cursor
