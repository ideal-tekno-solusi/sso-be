create table if not exists users (
	id varchar(50) primary key,
	name varchar(255) not null,
	dot timestamp not null,
	password text not null,
	insert_date timestamp not null
);

create table if not exists sessions (
	id text primary key,
	user_id varchar(50),
	client_id varchar(25) not null,
	code_challenge text not null,
	code_challenge_method varchar(10) not null,
	scopes varchar(255),
	redirect_url text,
	insert_date timestamp not null
);

create table if not exists authorization_tokens (
	id text primary key,
	session_id text references sessions(id),
	insert_date timestamp not null
);

create table if not exists refresh_tokens (
	id text primary key,
	user_id varchar(50),
	insert_date timestamp not null
);