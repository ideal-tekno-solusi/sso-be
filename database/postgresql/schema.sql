create table if not exists clients (
	id varchar(255) primary key,
	name varchar(255) not null,
	type int not null,
	secret varchar(255),
	token_livetime bigint
);

create table if not exists client_redirects (
	client_id varchar(255) references clients(id) on delete cascade,
	uri text not null
);

create table if not exists client_types (
	id int primary key,
	name varchar(10) not null
);

create table if not exists users (
	id varchar(50) primary key,
	name varchar(255) not null,
	dot timestamp not null,
	password text not null,
	insert_date timestamp not null
);

create table if not exists sessions (
	id varchar(255),
	user_id varchar(50) references users(id) on delete cascade,
	insert_date timestamp not null
);

create table if not exists auths (
	code varchar(255),
	scope varchar(100),
	type int not null,
	user_id varchar(50) references sso.users(id) on delete cascade,
	insert_date timestamp not null,
	use_date timestamp
);

create table if not exists auth_types (
	id int primary key,
	name varchar(10) not null
);