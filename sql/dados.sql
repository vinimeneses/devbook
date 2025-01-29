insert into usuarios (nome, nick, email, senha)
values
("Usuario 1", "usuario_1", "usuario1@gmail.com", "$2a$10$F9HU72QsOdWU.y.F5h9zJ.7Guqn3b3xR0iGi1IIHVWH2GEWLpN7Ri"),
("Usuario 2", "usuario_2", "usuario2@gmail.com", "$2a$10$F9HU72QsOdWU.y.F5h9zJ.7Guqn3b3xR0iGi1IIHVWH2GEWLpN7Ri"),
("Usuario 3", "usuario_3", "usuario3@gmail.com", "$2a$10$F9HU72QsOdWU.y.F5h9zJ.7Guqn3b3xR0iGi1IIHVWH2GEWLpN7Ri");

insert into  seguidores(usuario_id, seguidor_id)
values
    (1, 2),
(3, 1),
(1, 3);