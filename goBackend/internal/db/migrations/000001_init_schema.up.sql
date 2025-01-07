CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    public_id char(21) UNIQUE DEFAULT nanoid(21),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
    
  );

CREATE TABLE IF NOT EXISTS CATEGORIES (
  -- usage purely internal for database management
    id SERIAL PRIMARY KEY, 
  -- usage for public facing API
    public_id char(21) UNIQUE DEFAULT nanoid(21),
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    parent_id INT DEFAULT NULL,
    created_by INT DEFAULT 0,

    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET DEFAULT
    -- created_at TIMESTAMP  DEFAULT NOW(),
    -- updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS POSTS (
  -- usage purely internal for database management
    id SERIAL PRIMARY KEY, 
  -- usage for public facing API
    public_id char(21) UNIQUE DEFAULT nanoid(21),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
  --requires user_public_id to be unique for foreign key constraint
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (category_id) REFERENCES categories(id)
  );
