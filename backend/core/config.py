from dotenv import load_dotenv  # pip3 install python-dotenv
import os

load_dotenv()

# DB
DB_HOST = os.getenv("DB_HOST", "default_host")  # Providing default values is a good practice
DB_PORT = os.getenv("DB_PORT", "3306")
DB_NAME = os.getenv("DB_NAME", "default_db")
DB_USER = os.getenv("DB_USER", "default_user")
DB_PASS = os.getenv("DB_PASS", "default_password")

# API CONFIGURATION
SERVER_HOST = os.getenv("SERVER_HOST", "192.168.0.163")
SERVER_PORT = int(os.getenv("SERVER_PORT", 8080))
LOG_LEVEL = os.getenv("LOG_LEVEL", "info")
DEBUG_RELOAD = bool(os.getenv("DEBUG_RELOAD", True))

# # EMAIL
# EMAIL_ADDRESS = os.getenv("EMAIL_ADDRESS", "default_email@example.com")
# EMAIL_PASSWORD = os.getenv("EMAIL_PASSWORD")

# FRONTEND
FRONTEND_URL = os.getenv("FRONTEND_URL", "http://localhost:3000")

# SECRET_KEY
SECRET_KEY = os.getenv("SECRET_KEY", "FLKJSDIU#I@OJFK")
