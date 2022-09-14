DROP TABLE IF EXISTS workers.public.employees;
CREATE TABLE IF NOT EXISTS workers.public.employees
(
    employee_id    SERIAL PRIMARY KEY,
    nickname       varchar(64),
    first_name     varchar(64),
    last_name      varchar(64),
    street_address varchar(128),
    city           varchar(64),
    state          varchar(2),
    zip            varchar(24)
);

CREATE UNIQUE INDEX workers.public.idx_employee_nickname ON workers.public.employees (nickname);

DROP TABLE IF EXISTS workers.public.tasks;
CREATE TABLE IF NOT EXISTS workers.public.tasks
(
    id          SERIAL,
    name        varchar(128),
    description varchar(1024),
    create_time timestamp,
    owners      integer[],
    private     bool,
    due_by      date
);

INSERT INTO workers.public.employees (nickname, first_name, last_name, street_address, city, state, zip)
VALUES ('Danny', 'Daniel', 'Smith', '100 Wabash Ave', 'Chicago', 'IL', '60000');
INSERT INTO workers.public.employees (nickname, first_name, last_name, street_address, city, state, zip)
VALUES ('Jaybird', 'Jayson', 'Hurd', '9311 Cove Creek Drive', 'Highlands Ranch', 'CO', '80129');
INSERT INTO workers.public.employees (nickname, first_name, last_name, street_address, city, state, zip)
VALUES ('Jumpin Jack Flash', 'Jack', 'Flash', '100 West Ave', 'New York', 'NY', '10020');
INSERT INTO workers.public.employees (nickname, first_name, last_name, street_address, city, state, zip)
VALUES ('Nice Taylor', 'Taylor', 'Swift', '300 Ocean Drive', 'Watch Hill', 'RI', '02856');
INSERT INTO workers.public.employees (nickname, first_name, last_name, street_address, city, state, zip)
VALUES ('Princess Diana', 'Diana', 'Princess', '500 Windsor Drive', 'London', 'UK', 'GRE234');
INSERT INTO workers.public.employees (nickname, first_name, last_name, street_address, city, state, zip)
VALUES ('Mama Theresa', 'Mother', 'Theresa', '500 Dade Drive', 'Miama', 'FL', '34356');

INSERT INTO workers.public.tasks (name, description, create_time, owners, private, due_by)
VALUES ('goapi', 'write go api for job interview',  NOW(),  '{2}', false, '2022-09-10');
INSERT INTO workers.public.tasks (name, description, create_time, owners, private, due_by)
VALUES ('Save the world', 'Save the world from hunger and oppression',  NOW(),  '{5, 6}', false, '2030-01-01');
INSERT INTO workers.public.tasks (name, description, create_time, owners, private, due_by)
VALUES ('Write the best program ever', 'Write a program that helps solve world hunger',  NOW(), '{1,2}', true, '2022-12-31');
INSERT INTO workers.public.tasks (name, description, create_time, owners, private, due_by)
VALUES ('Write a new song', 'Write a song that inspires people to write better programs and', NOW(), '{4}', false,  '2022-12-31');
INSERT INTO workers.public.tasks (name, description, create_time, owners, private, due_by)
VALUES ('Make a movie',
        'Make a movie that excites and inspires programmers and inspirational people to solve world hunger', NOW(), '{3}',
        false,
        '2023-12-31');