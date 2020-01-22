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

CREATE TABLE public.image_cache
(
    img_uri_hash text NOT NULL,
    pass boolean NOT NULL,
    PRIMARY KEY (img_uri_hash)
);

ALTER TABLE public.image_cache
    OWNER to postgres;