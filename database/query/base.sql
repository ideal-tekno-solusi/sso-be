create schema if not exists sso;

create table if not exists sso.users (
	id varchar(50) primary key,
	name varchar(255) not null,
	dot timestamp not null,
	password text not null,
	insert_date timestamp not null
);

create table if not exists sso.sessions (
	id varchar(250) primary key,
	user_id varchar(50),
	client_id varchar(25) not null,
	code_challenge text not null,
	code_challenge_method varchar(10) not null,
	scopes varchar(255),
	redirect_url text,
	insert_date timestamp not null
);

create table if not exists sso.authorization_tokens (
	id text primary key,
	session_id varchar(250) references sso.sessions(id),
	insert_date timestamp not null
);

create table if not exists sso.refresh_tokens (
	id text primary key,
	user_id varchar(50),
	insert_date timestamp not null
)

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