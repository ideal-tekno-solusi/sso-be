-- run from root access
create user sso with password 'asd123qwe';

create database sso;

--connect to db sso first
create schema if not exists sso authorization sso;

grant all on all tables in schema sso to sso;
-- end

create table if not exists sso.clients (
	id varchar(255) primary key,
	name varchar(255) not null,
	type int not null,
	secret varchar(255),
	token_livetime bigint
);

create table if not exists sso.client_redirects (
	client_id varchar(255) references sso.clients(id) on delete cascade,
	uri text not null
);

create table if not exists sso.client_types (
	id int primary key,
	name varchar(10) not null
);

create table if not exists sso.users (
	id varchar(50) primary key,
	name varchar(255) not null,
	dot timestamp not null,
	password text not null,
	insert_date timestamp not null
);

create table if not exists sso.sessions (
	id varchar(255),
	user_id varchar(50) references sso.users(id) on delete cascade,
	insert_date timestamp not null
);

create table if not exists sso.auths (
	code varchar(255),
	scope varchar(100),
	type int not null,
	user_id varchar(50) references sso.users(id) on delete cascade,
	insert_date timestamp not null,
	use_date timestamp
);

create table if not exists sso.auth_types (
	id int primary key,
	name varchar(10) not null
)

insert into sso.clients (
	id,
	name,
	type,
	secret,
	token_livetime
)
values (
	'INVENTORY_APP_01',
	'Inventory app',
	1,
	'a17bf8485b43f846e8a3e7df443bc169',
	3600
);

insert into sso.client_redirects (
	client_id,
	uri
)
values (
	'INVENTORY_APP_01',
	'http://localhost:8051/redirect'
),
(
	'INVENTORY_APP_01',
	'http://localhost:5173/oauth-callback'
);

insert into sso.client_types (
	id,
	name
)
values (
	1,
	'SPA'
);

insert into sso.users (
	id,
	name,
	dot,
	password,
	insert_date
)
values (
	'alfian',
	'alfian',
	'1997-06-10',
	'$2a$15$xwGZGcKIURe1kwSt7zTrrOwCCwOfmN9K5SqOu32sJdGj67FJEUfou',
	now()
);

insert into sso.auth_types (
	id,
	name
)
values (
	1,
	'CODE'
),
(
	2,
	'REFRESH'
);