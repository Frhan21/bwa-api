CREATE TABLE IF NOT EXISTS "contents" (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, 
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE, 
    title VARCHAR(200) NOT NULL,
    excerpt VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    image TEXT NULL,
    status VARCHAR(20) NOT NULL,
    tags TEXT NOT NULL, 
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
)