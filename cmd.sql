create table if not exists reg_users(
	user_id serial primary key,
	username varchar(50) unique not null,
	password varchar(100) not null,
	email varchar(50) unique not null,
	registered boolean default false,
	activated boolean default false
	);

create table if not exists vault(
	vault_id serial primary key,
	user_id int,
	content varchar(100),
	constraint fk_reg_user
		foreign key(user_id)
			references reg_users(user_id)
			on delete cascade
);
