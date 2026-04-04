-- Tabela de usuarios (MVP: auto-criado via session anonima)
CREATE TABLE IF NOT EXISTS users (
    id         VARCHAR(64) PRIMARY KEY,
    name       VARCHAR(255),
    email      VARCHAR(255),
    mode       VARCHAR(20) DEFAULT 'student',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tabela de sessoes (gravacoes de aula/reuniao)
CREATE TABLE IF NOT EXISTS sessions (
    id               VARCHAR(64) PRIMARY KEY,
    user_id          VARCHAR(64) REFERENCES users(id),
    title            VARCHAR(500),
    description      TEXT,
    duration         INTEGER DEFAULT 0,
    status           VARCHAR(20) DEFAULT 'processing',
    mode             VARCHAR(20) DEFAULT 'student',
    transcript       TEXT,
    audio_path       VARCHAR(1000),
    summary_data     JSONB,
    created_at       TIMESTAMP DEFAULT NOW(),
    updated_at       TIMESTAMP DEFAULT NOW()
);

-- Tabela de flashcards
CREATE TABLE IF NOT EXISTS flashcards (
    id         VARCHAR(64) PRIMARY KEY,
    session_id VARCHAR(64) NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    front      TEXT NOT NULL,
    back       TEXT NOT NULL,
    difficulty INTEGER DEFAULT 2,
    known      BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indices para queries comuns
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status);
CREATE INDEX IF NOT EXISTS idx_sessions_created ON sessions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_flashcards_session ON flashcards(session_id);
