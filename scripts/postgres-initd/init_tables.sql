--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5
-- Dumped by pg_dump version 14.7

-- Started on 2023-04-03 20:07:51 UTC

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
-- TOC entry 3397 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION pg_stat_statements; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION pg_stat_statements IS 'track planning and execution statistics of all SQL statements executed';


--
-- TOC entry 865 (class 1247 OID 16490)
-- Name: relation_type; Type: TYPE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TYPE public.relation_type AS ENUM (
    'following',
    'follower',
    'blocked'
);


ALTER TYPE public.relation_type OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 220 (class 1255 OID 16410)
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
-- TOC entry 221 (class 1255 OID 16413)
-- Name: update_likes_count_on_comment(); Type: FUNCTION; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE FUNCTION public.update_likes_count_on_comment() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	IF (TG_OP = 'INSERT') THEN
		UPDATE comments
			SET likes = likes + 1
			WHERE id = NEW.comment_id;
		RETURN NEW;
	ELSE
		UPDATE comments
			SET likes = likes - 1
			WHERE id = OLD.comment_id;
		RETURN OLD;
	END IF;
END;$$;


ALTER FUNCTION public.update_likes_count_on_comment() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 222 (class 1255 OID 16414)
-- Name: update_likes_count_on_post(); Type: FUNCTION; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE FUNCTION public.update_likes_count_on_post() RETURNS trigger
    LANGUAGE plpgsql
    AS $$BEGIN
	IF (TG_OP = 'INSERT') THEN
		UPDATE posts
			SET likes = likes + 1
			WHERE id = NEW.post_id;
		RETURN NEW;
	ELSE
		UPDATE posts
			SET likes = likes - 1
			WHERE id = OLD.post_id;
		RETURN OLD;
	END IF;
END;$$;


ALTER FUNCTION public.update_likes_count_on_post() OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 212 (class 1259 OID 16415)
-- Name: comment_likes; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.comment_likes (
    id uuid NOT NULL,
    comment_id uuid NOT NULL,
    user_id uuid NOT NULL
);


ALTER TABLE public.comment_likes OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 213 (class 1259 OID 16418)
-- Name: comments; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.comments (
    id uuid NOT NULL,
    post_id uuid NOT NULL,
    user_id uuid NOT NULL,
    content character varying(512) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    likes integer NOT NULL
);


ALTER TABLE public.comments OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 214 (class 1259 OID 16423)
-- Name: post_likes; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.post_likes (
    id uuid NOT NULL,
    post_id uuid NOT NULL,
    user_id uuid NOT NULL
);


ALTER TABLE public.post_likes OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 215 (class 1259 OID 16426)
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
-- TOC entry 216 (class 1259 OID 16431)
-- Name: user_info; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.user_info (
    id uuid NOT NULL,
    about character varying(256)
);


ALTER TABLE public.user_info OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 217 (class 1259 OID 16434)
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
-- TOC entry 218 (class 1259 OID 16439)
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
-- TOC entry 219 (class 1259 OID 16497)
-- Name: user_relations; Type: TABLE; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TABLE public.user_relations (
    id uuid NOT NULL,
    user_id uuid,
    relation_user_id uuid,
    type public.relation_type NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.user_relations OWNER TO "ea1fb999-4aab-4142-9101-facdc7d5b83b";

--
-- TOC entry 3225 (class 2606 OID 16444)
-- Name: comment_likes comment_likes_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.comment_likes
    ADD CONSTRAINT comment_likes_pkey PRIMARY KEY (id);


--
-- TOC entry 3229 (class 2606 OID 16446)
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- TOC entry 3231 (class 2606 OID 16448)
-- Name: post_likes post_likes_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.post_likes
    ADD CONSTRAINT post_likes_pkey PRIMARY KEY (id);


--
-- TOC entry 3235 (class 2606 OID 16450)
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- TOC entry 3227 (class 2606 OID 16452)
-- Name: comment_likes user_comment; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.comment_likes
    ADD CONSTRAINT user_comment UNIQUE (comment_id) INCLUDE (user_id);


--
-- TOC entry 3237 (class 2606 OID 16454)
-- Name: user_info user_info_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT user_info_pkey PRIMARY KEY (id);


--
-- TOC entry 3233 (class 2606 OID 16456)
-- Name: post_likes user_post; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.post_likes
    ADD CONSTRAINT user_post UNIQUE (post_id) INCLUDE (user_id);


--
-- TOC entry 3241 (class 2606 OID 16501)
-- Name: user_relations user_relations_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_relations
    ADD CONSTRAINT user_relations_pkey PRIMARY KEY (id);


--
-- TOC entry 3239 (class 2606 OID 16458)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3249 (class 2620 OID 16459)
-- Name: users add_user_info_trigger; Type: TRIGGER; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TRIGGER add_user_info_trigger AFTER INSERT ON public.users FOR EACH ROW EXECUTE FUNCTION public.add_user_info();


--
-- TOC entry 3247 (class 2620 OID 16462)
-- Name: comment_likes update_likes_count_on_post_trigger; Type: TRIGGER; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TRIGGER update_likes_count_on_post_trigger AFTER INSERT OR DELETE ON public.comment_likes FOR EACH ROW EXECUTE FUNCTION public.update_likes_count_on_comment();


--
-- TOC entry 3248 (class 2620 OID 16463)
-- Name: post_likes update_likes_count_on_post_trigger; Type: TRIGGER; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

CREATE TRIGGER update_likes_count_on_post_trigger AFTER INSERT OR DELETE ON public.post_likes FOR EACH ROW EXECUTE FUNCTION public.update_likes_count_on_post();


--
-- TOC entry 3242 (class 2606 OID 16464)
-- Name: comment_likes comment_fk; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.comment_likes
    ADD CONSTRAINT comment_fk FOREIGN KEY (comment_id) REFERENCES public.comments(id) NOT VALID;


--
-- TOC entry 3243 (class 2606 OID 16508)
-- Name: comments fk_post; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3245 (class 2606 OID 16528)
-- Name: posts fk_post_user; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES public.users(id) NOT VALID;


--
-- TOC entry 3246 (class 2606 OID 16523)
-- Name: user_info fk_user; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.user_info
    ADD CONSTRAINT fk_user FOREIGN KEY (id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- TOC entry 3244 (class 2606 OID 16479)
-- Name: post_likes post_fk; Type: FK CONSTRAINT; Schema: public; Owner: ea1fb999-4aab-4142-9101-facdc7d5b83b
--

ALTER TABLE ONLY public.post_likes
    ADD CONSTRAINT post_fk FOREIGN KEY (post_id) REFERENCES public.posts(id);


-- Completed on 2023-04-03 20:07:52 UTC

--
-- PostgreSQL database dump complete
--
