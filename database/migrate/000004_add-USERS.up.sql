    CREATE TABLE usuarios (
        id SERIAL PRIMARY KEY,
        nome VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        senha_hash TEXT NOT NULL,
        ativo BOOLEAN DEFAULT TRUE
    );