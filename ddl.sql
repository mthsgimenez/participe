CREATE TABLE companies (
	id serial NOT NULL,
	"name" varchar(100) NOT NULL,
	CONSTRAINT companies_pk PRIMARY KEY (id)
);

CREATE TABLE events (
	id serial NOT NULL,
	description text NULL,
	"name" varchar(100) NOT NULL,
    "date" timestamptz NOT NULL,
	CONSTRAINT events_pk PRIMARY KEY (id)
);

CREATE TABLE users (
	id serial NOT NULL,
	email varchar(100) NOT NULL,
	company_id int NOT NULL,
	"name" varchar(60) NOT NULL,
	"role" text NOT NULL DEFAULT 'ROLE_USER',
	"password" text NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_unique UNIQUE (email),
	CONSTRAINT users_companies_fk FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE events_users (
	id serial NOT NULL,
	user_id int NOT NULL,
	event_id int NOT NULL,
	CONSTRAINT events_users_pk PRIMARY KEY (id),
	CONSTRAINT events_users_events_fk FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT events_users_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- DROP TABLE events_users CASCADE;
-- DROP TABLE users CASCADE;
-- DROP TABLE events CASCADE;
-- DROP TABLE companies CASCADE;

-- ALTER SEQUENCE companies_id_seq RESTART WITH 1;
-- ALTER SEQUENCE events_id_seq RESTART WITH 1;
-- ALTER SEQUENCE users_id_seq RESTART WITH 1;
-- ALTER SEQUENCE events_users_id_seq RESTART WITH 1;

-- ===========================
-- EMPRESAS
-- ===========================
INSERT INTO companies (name) VALUES
('TechNova Solutions'),
('Alfa Sistemas'),
('Inova Digital');

-- ===========================
-- USUÁRIOS
-- ===========================
INSERT INTO users (email, company_id, name, role, password) VALUES
('joao.silva@technova.com', 1, 'João Silva', 'ROLE_USER', '$2a$12$ONVS.jkh8u6EO0pb1o/41uOOj7oD5DCmDlkSeL7VEwOsa0EgGLzhm'),
('maria.souza@technova.com', 1, 'Maria Souza', 'ROLE_USER', '$2a$12$ONVS.jkh8u6EO0pb1o/41uOOj7oD5DCmDlkSeL7VEwOsa0EgGLzhm'),
('pedro.alves@alfasistemas.com', 2, 'Pedro Alves', 'ROLE_USER', '$2a$12$ONVS.jkh8u6EO0pb1o/41uOOj7oD5DCmDlkSeL7VEwOsa0EgGLzhm'),
('ana.martins@inovadigital.com', 3, 'Ana Martins', 'ROLE_USER', '$2a$12$ONVS.jkh8u6EO0pb1o/41uOOj7oD5DCmDlkSeL7VEwOsa0EgGLzhm');

-- ===========================
-- EVENTOS (passados e futuros)
-- ===========================
INSERT INTO events (name, description, date) VALUES
('Treinamento de Integração', 'Treinamento inicial para novos colaboradores.', '2024-07-10 09:00:00-03'),
('Workshop de Produtividade', 'Sessão prática sobre ferramentas de produtividade.', '2024-10-15 14:00:00-03'),
('Palestra de Segurança da Informação', 'Apresentação sobre boas práticas de segurança.', '2025-01-20 10:00:00-03'),
('Hackathon Interno', 'Maratona de desenvolvimento entre equipes.', '2025-05-05 08:00:00-03'),
('Encontro Anual de Estratégia', 'Evento anual para alinhamento estratégico da empresa.', '2025-12-02 09:30:00-03');

-- ===========================
-- PRESENÇA DOS USUÁRIOS NOS EVENTOS
-- ===========================
INSERT INTO events_users (user_id, event_id) VALUES
-- Treinamento de Integração (passado)
(1, 1),
(2, 1),

-- Workshop de Produtividade (passado)
(2, 2),
(3, 2),

-- Palestra de Segurança da Informação (futuro)
(1, 3),
(3, 3),
(4, 3),

-- Hackathon Interno (futuro)
(1, 4),
(2, 4),
(4, 4),

-- Encontro Anual de Estratégia (futuro)
(3, 5),
(4, 5);


-- usuario admin
INSERT INTO users (email, company_id, "name", "role", "password") VALUES
('admin@gmail.com', 1, 'admin', 'ROLE_ADMIN', '$2a$12$7IXtwNPZD1IhYHYQ0iy.yOK89y9HbQX66nNE/XgUHZ2.aiSdmM7ES');
