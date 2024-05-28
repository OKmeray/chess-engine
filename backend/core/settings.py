from core.config import DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME, SERVER_HOST, SERVER_PORT, LOG_LEVEL, DEBUG_RELOAD

# Now you can use the settings from the Config class
database_url = f"mysql://{DB_USER}:{DB_PASS}@{DB_HOST}:{DB_PORT}/{DB_NAME}"
allowed_origins = [
    "http://localhost:3000",
    "http://192.168.0.163:3000",
]
