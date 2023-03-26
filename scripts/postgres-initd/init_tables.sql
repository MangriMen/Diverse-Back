SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS pg_stat_statements WITH SCHEMA public;

COMMENT ON EXTENSION pg_stat_statements IS 'track planning and execution statistics of all SQL statements executed';

CREATE FUNCTION public.add_user_info() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	INSERT INTO user_info VALUES (NEW.id, null);

CREATE FUNCTION public.remove_post_comments() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	DELETE FROM comments WHERE post_id = NEW.id;
	RETURN NEW;
END;$$;


ALTER FUNCTION public.add_user_info() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

CREATE FUNCTION public.remove_user_info() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	DELETE FROM user_info WHERE id = NEW.id;
	RETURN NEW;
END;$$;

ALTER FUNCTION public.remove_user_info() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";
ALTER FUNCTION public.remove_post_comments() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE TABLE public.comments (
    id uuid NOT NULL,
    post_id uuid NOT NULL,
    user_id uuid NOT NULL,
    content character varying(512) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.comments OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

CREATE TABLE public.posts (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    content text NOT NULL,
    description character varying(2048),
    likes integer NOT NULL,
    created_at timestamp with time zone NOT NULL
);


ALTER TABLE public.posts OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

CREATE TABLE public.user_info (
    id uuid NOT NULL,
    about character varying(256)
);


ALTER TABLE public.user_info OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

CREATE TABLE public.users (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(256) NOT NULL,
    username character varying(32) NOT NULL,
    name character varying(32),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    avatar_url text
);


ALTER TABLE public.users OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

CREATE VIEW public.user_full_info AS
 SELECT users.id,
    users.email,
    users.password,
    users.username,
    users.name,
    users.created_at,
    users.updated_at,
    users.avatar_url,
    user_info.about
   FROM (public.users
     LEFT JOIN public.user_info USING (id));


ALTER TABLE public.user_full_info OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE TRIGGER add_user_info_trigger AFTER INSERT ON public.users FOR EACH ROW EXECUTE FUNCTION public.add_user_info();

CREATE TRIGGER remove_user_info_trigger AFTER DELETE ON public.users FOR EACH ROW EXECUTE FUNCTION public.remove_user_info();

CREATE TRIGGER remove_post_comments_trigger AFTER DELETE ON public.posts FOR EACH ROW EXECUTE FUNCTION public.remove_post_comments();

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_user_id_fkey FOREIGN KEY (id) REFERENCES public.users(id);
