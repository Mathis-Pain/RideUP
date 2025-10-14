CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  profile_pic TEXT DEFAULT '/static/noprofilpic.png',
  created_by INT REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  latitude DOUBLE PRECISION DEFAULT 49.43839,
  longitude DOUBLE PRECISION DEFAULT 1.10160,
   role_id INTEGER NOT NULL DEFAULT 3
);

CREATE TABLE event_participants (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_id INT REFERENCES events(id) ON DELETE CASCADE,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(event_id, user_id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);


CREATE TABLE messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_id INT REFERENCES events(id) ON DELETE CASCADE,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notifications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  message TEXT NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,          -- Identifiant unique de la session
    user_id INTEGER,             -- ID de l'utilisateur associé
    data TEXT,                    -- Données de session sérialisées en JSON
    expires_at TIMESTAMP,         -- Date d'expiration de la session
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Date de création
);
