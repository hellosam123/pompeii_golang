# init_db.py
# ------
# run this Python file to create the database.db file
# open database.db in DB Browser to view database
# (to reset the database, delete the database.db file, re-run this file)
# ----------------------


import sqlite3

# connect to database
connection = sqlite3.connect('database.db')

# open fortune.sql
with open('db_definition.sql') as file:
    connection.executescript(file.read())

# close database
connection.close()

print("-- database initalized --")
