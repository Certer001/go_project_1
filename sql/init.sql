DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tasks (title, description, completed) VALUES 
                                                    ('Изучить GO', 'Пройти базовый курс', 'true'),
                                                    ('Написать REST API', 'Посмотреть это видео и написать самому', 'false'),
                                                    ('Зарелизить приложение', 'Развернуть приложение на сервере с помощью docker', 'false')
