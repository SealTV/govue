DROP TABLE IF EXISTS public.account; 

CREATE TABLE IF NOT EXISTS public.account (
	id SERIAL NOT NULL,
	"name" varchar(36) NOT NULL, 
    email varchar(36) NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP,
    deleted_at TIMESTAMP
);


ALTER TABLE public.account ADD CONSTRAINT account_pk PRIMARY KEY (id);
ALTER TABLE public.account ADD CONSTRAINT account_name_un UNIQUE ("name");
ALTER TABLE public.account ADD CONSTRAINT account_email_un UNIQUE (email);
