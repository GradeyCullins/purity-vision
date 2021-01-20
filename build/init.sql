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
    img_uri_hash text NOT NULL,
    error text,
    pass boolean NOT NULL,
    date_added timestamp NOT NULL default CURRENT_TIMESTAMP,
    PRIMARY KEY (img_uri_hash)
);

ALTER TABLE public.images
    OWNER to postgres;

-- Insert some test data.
-- https://i.imgur.com/gcWltJm.jpg
INSERT INTO images VALUES ('e6ec859aab9ee040c37b370f5f85eb7f39715d2c1230f9b8d6df92e530abab4a', NULL, FALSE);

-- https://google.com
INSERT INTO images VALUES ('05046f26c83e8c88b3ddab2eab63d0d16224ac1e564535fc75cdceee47a0938d', 'URI is not an image', FALSE);

-- http://www.audubon.org/sites/default/files/257px-Eastern_Phoebe1.jpg
INSERT INTO images VALUES ('034be24efb77545bd4376ac9324f1edf3ce36066de5e4c3b0c250877418cfb4a', NULL, TRUE);

-- Create the test database.
CREATE DATABASE purity_test WITH TEMPLATE purity;
