-- Setup purity schema.
CREATE TABLE public.users
(
    uid bigint NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    PRIMARY KEY (uid)
);

ALTER TABLE public.users
    OWNER to postgres;

ALTER TABLE public.users
    ALTER COLUMN uid ADD GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 );

CREATE TABLE public.images
(
    img_hash text NOT NULL,
    error text,
    pass boolean NOT NULL default false,
    date_added timestamp NOT NULL default CURRENT_TIMESTAMP,
    PRIMARY KEY (img_hash)
);

ALTER TABLE public.images
    OWNER to postgres;

-- Create the test database.
CREATE DATABASE purity_test WITH TEMPLATE purity;
