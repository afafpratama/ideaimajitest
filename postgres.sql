--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2
-- Dumped by pg_dump version 17.2

-- Started on 2025-01-21 23:30:12

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 4 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- TOC entry 4917 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 220 (class 1259 OID 16474)
-- Name: dt_customer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dt_customer (
    id integer NOT NULL,
    name character varying(100),
    phone character varying(50),
    created_at timestamp without time zone
);


ALTER TABLE public.dt_customer OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16473)
-- Name: dt_customer_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.dt_customer_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dt_customer_id_seq OWNER TO postgres;

--
-- TOC entry 4918 (class 0 OID 0)
-- Dependencies: 219
-- Name: dt_customer_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.dt_customer_id_seq OWNED BY public.dt_customer.id;


--
-- TOC entry 222 (class 1259 OID 16481)
-- Name: dt_order; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dt_order (
    id integer NOT NULL,
    customer_id integer,
    service text,
    amount integer,
    unit character varying(100),
    price integer,
    created_at timestamp without time zone
);


ALTER TABLE public.dt_order OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 16480)
-- Name: dt_order_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.dt_order_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dt_order_id_seq OWNER TO postgres;

--
-- TOC entry 4919 (class 0 OID 0)
-- Dependencies: 221
-- Name: dt_order_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.dt_order_id_seq OWNED BY public.dt_order.id;


--
-- TOC entry 218 (class 1259 OID 16464)
-- Name: sys_account; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sys_account (
    id integer NOT NULL,
    name character varying(100),
    phone character varying(50),
    username character varying(100),
    password text,
    created_at timestamp without time zone
);


ALTER TABLE public.sys_account OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16463)
-- Name: sys_account_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sys_account_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sys_account_id_seq OWNER TO postgres;

--
-- TOC entry 4920 (class 0 OID 0)
-- Dependencies: 217
-- Name: sys_account_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sys_account_id_seq OWNED BY public.sys_account.id;


--
-- TOC entry 4753 (class 2604 OID 16477)
-- Name: dt_customer id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dt_customer ALTER COLUMN id SET DEFAULT nextval('public.dt_customer_id_seq'::regclass);


--
-- TOC entry 4754 (class 2604 OID 16484)
-- Name: dt_order id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dt_order ALTER COLUMN id SET DEFAULT nextval('public.dt_order_id_seq'::regclass);


--
-- TOC entry 4752 (class 2604 OID 16467)
-- Name: sys_account id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sys_account ALTER COLUMN id SET DEFAULT nextval('public.sys_account_id_seq'::regclass);


--
-- TOC entry 4909 (class 0 OID 16474)
-- Dependencies: 220
-- Data for Name: dt_customer; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.dt_customer (id, name, phone, created_at) VALUES (1, 'Customer 234', '234', '2025-01-21 16:07:25.778283');


--
-- TOC entry 4911 (class 0 OID 16481)
-- Dependencies: 222
-- Data for Name: dt_order; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.dt_order (id, customer_id, service, amount, unit, price, created_at) VALUES (1, 0, 'One Stop Solution', 1, 'event', 10000, '2025-01-21 16:02:56.642718');
INSERT INTO public.dt_order (id, customer_id, service, amount, unit, price, created_at) VALUES (2, 4, 'One Stop Solution', 4, 'event', 40000, '2025-01-21 16:03:13.95272');


--
-- TOC entry 4907 (class 0 OID 16464)
-- Dependencies: 218
-- Data for Name: sys_account; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sys_account (id, name, phone, username, password, created_at) VALUES (1, 'Hella', '000', 'hello', '$2a$10$Gfy.9hmOkw4MAFW6siKvxOb0X/BAyliCG1CJ//OUFm8P8R.10MUMC', '2025-01-21 07:57:50.441864');
INSERT INTO public.sys_account (id, name, phone, username, password, created_at) VALUES (5, 'World767', '767', 'world767', '$2a$10$thI34PVwNVqirwAY4UU.fOxOCEyAYOKEJ/5RJusqwmECj2vxBAdLG', '2025-01-21 15:41:49.396032');
INSERT INTO public.sys_account (id, name, phone, username, password, created_at) VALUES (4, 'Hella444', '444', 'world', '$2a$10$PpabLu.Qzg1Xx8qGvidHC.Zu1EdcIvsQ1Y6.mfwspFbQ9UqEGReMK', '2025-01-21 12:23:42.062974');


--
-- TOC entry 4921 (class 0 OID 0)
-- Dependencies: 219
-- Name: dt_customer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.dt_customer_id_seq', 2, true);


--
-- TOC entry 4922 (class 0 OID 0)
-- Dependencies: 221
-- Name: dt_order_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.dt_order_id_seq', 3, true);


--
-- TOC entry 4923 (class 0 OID 0)
-- Dependencies: 217
-- Name: sys_account_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sys_account_id_seq', 5, true);


--
-- TOC entry 4758 (class 2606 OID 16479)
-- Name: dt_customer dt_customer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dt_customer
    ADD CONSTRAINT dt_customer_pkey PRIMARY KEY (id);


--
-- TOC entry 4760 (class 2606 OID 16488)
-- Name: dt_order dt_order_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dt_order
    ADD CONSTRAINT dt_order_pkey PRIMARY KEY (id);


--
-- TOC entry 4756 (class 2606 OID 16471)
-- Name: sys_account sys_account_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sys_account
    ADD CONSTRAINT sys_account_pkey PRIMARY KEY (id);


-- Completed on 2025-01-21 23:30:12

--
-- PostgreSQL database dump complete
--

