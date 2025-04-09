create table if not exists users (
	id varchar(50) primary key,
	name varchar(255) not null,
	dot timestamp not null,
	password text not null,
	insert_date timestamp not null
);

create table if not exists sessions (
	id varchar(250) primary key,
	client_id varchar(25) not null,
	code_challenge text not null,
	code_challenge_method varchar(10) not null,
	insert_date timestamp not null
);

create table if not exists authorization_tokens (
	id text primary key,
	user_id varchar(50) references users(id),
	insert_date timestamp not null
);