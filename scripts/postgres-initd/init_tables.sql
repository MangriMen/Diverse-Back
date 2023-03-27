--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5
-- Dumped by pg_dump version 14.6

-- Started on 2023-03-27 11:08:30 UTC

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

--
-- TOC entry 2 (class 3079 OID 16385)
-- Name: pg_stat_statements; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_stat_statements WITH SCHEMA public;


--
-- TOC entry 3370 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION pg_stat_statements; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION pg_stat_statements IS 'track planning and execution statistics of all SQL statements executed';


--
-- TOC entry 217 (class 1255 OID 16410)
-- Name: add_user_info(); Type: FUNCTION; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE FUNCTION public.add_user_info() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	INSERT INTO user_info VALUES (NEW.id, null);
    RETURN NEW;
END;$$;


ALTER FUNCTION public.add_user_info() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 218 (class 1255 OID 16411)
-- Name: remove_post_comments(); Type: FUNCTION; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE FUNCTION public.remove_post_comments() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	DELETE FROM comments WHERE post_id = NEW.id;
	RETURN NEW;
END;$$;


ALTER FUNCTION public.remove_post_comments() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 219 (class 1255 OID 16412)
-- Name: remove_user_info(); Type: FUNCTION; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE FUNCTION public.remove_user_info() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	DELETE FROM user_info WHERE id = NEW.id;
	RETURN NEW;
END;$$;


ALTER FUNCTION public.remove_user_info() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 212 (class 1259 OID 16413)
-- Name: comments; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.comments (
    id uuid NOT NULL,
    post_id uuid NOT NULL,
    user_id uuid NOT NULL,
    content character varying(512) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.comments OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 213 (class 1259 OID 16418)
-- Name: posts; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.posts (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    content text NOT NULL,
    description character varying(2048),
    likes integer NOT NULL,
    created_at timestamp with time zone NOT NULL
);


ALTER TABLE public.posts OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 214 (class 1259 OID 16423)
-- Name: user_info; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.user_info (
    id uuid NOT NULL,
    about character varying(256)
);


ALTER TABLE public.user_info OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 215 (class 1259 OID 16426)
-- Name: users; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

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

--
-- TOC entry 216 (class 1259 OID 16431)
-- Name: user_full_info; Type: VIEW; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

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

--
-- TOC entry 3210 (class 2606 OID 16436)
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- TOC entry 3212 (class 2606 OID 16438)
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- TOC entry 3214 (class 2606 OID 16440)
-- Name: user_info user_info_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_pkey PRIMARY KEY (id);


--
-- TOC entry 3216 (class 2606 OID 16442)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3222 (class 2620 OID 16443)
-- Name: users add_user_info_trigger; Type: TRIGGER; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TRIGGER add_user_info_trigger AFTER INSERT ON public.users FOR EACH ROW EXECUTE FUNCTION public.add_user_info();


--
-- TOC entry 3220 (class 2620 OID 16445)
-- Name: posts remove_post_comments_trigger; Type: TRIGGER; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TRIGGER remove_post_comments_trigger AFTER DELETE ON public.posts FOR EACH ROW EXECUTE FUNCTION public.remove_post_comments();


--
-- TOC entry 3221 (class 2620 OID 16444)
-- Name: users remove_user_info_trigger; Type: TRIGGER; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TRIGGER remove_user_info_trigger AFTER DELETE ON public.users FOR EACH ROW EXECUTE FUNCTION public.remove_user_info();


--
-- TOC entry 3217 (class 2606 OID 16456)
-- Name: comments post_fk; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT post_fk FOREIGN KEY (post_id) REFERENCES public.posts(id) NOT VALID;


--
-- TOC entry 3218 (class 2606 OID 16446)
-- Name: posts posts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3219 (class 2606 OID 16451)
-- Name: user_info user_info_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_user_id_fkey FOREIGN KEY (id) REFERENCES public.users(id);


-- Completed on 2023-03-27 11:08:30 UTC

--
-- PostgreSQL database dump complete
--
