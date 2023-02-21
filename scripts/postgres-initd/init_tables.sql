--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5
-- Dumped by pg_dump version 14.5

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: posts; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.posts (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    post json NOT NULL
);


ALTER TABLE public.posts OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
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
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: user_info user_info_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: posts posts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_info user_info_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

