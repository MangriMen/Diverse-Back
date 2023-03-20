--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5
-- Dumped by pg_dump version 14.6

-- Started on 2023-03-20 21:36:43 UTC

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
-- TOC entry 3360 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION pg_stat_statements; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION pg_stat_statements IS 'track planning and execution statistics of all SQL statements executed';


--
-- TOC entry 216 (class 1255 OID 16452)
-- Name: remove_post_comments(); Type: FUNCTION; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE FUNCTION public.remove_post_comments() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	DELETE FROM comments WHERE post_id = NEW.id;
	RETURN NEW;
END;$$;


ALTER FUNCTION public.remove_post_comments() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 16441)
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
-- TOC entry 212 (class 1259 OID 16410)
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
-- TOC entry 213 (class 1259 OID 16415)
-- Name: user_info; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.user_info (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    about character varying(256),
    avatar_url text
);


ALTER TABLE public.user_info OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 214 (class 1259 OID 16420)
-- Name: users; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(256) NOT NULL,
    username character varying(32) NOT NULL,
    name character varying(32),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.users OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 3210 (class 2606 OID 16445)
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- TOC entry 3204 (class 2606 OID 16426)
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- TOC entry 3206 (class 2606 OID 16428)
-- Name: user_info user_info_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_pkey PRIMARY KEY (id);


--
-- TOC entry 3208 (class 2606 OID 16430)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3213 (class 2620 OID 16453)
-- Name: posts remove_post_comments_trigger; Type: TRIGGER; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TRIGGER remove_post_comments_trigger AFTER DELETE ON public.posts FOR EACH ROW EXECUTE FUNCTION public.remove_post_comments();


--
-- TOC entry 3211 (class 2606 OID 16431)
-- Name: posts posts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3212 (class 2606 OID 16436)
-- Name: user_info user_info_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


-- Completed on 2023-03-20 21:36:43 UTC

--
-- PostgreSQL database dump complete
--
