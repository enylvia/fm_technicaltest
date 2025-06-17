--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2 (Debian 16.2-1.pgdg120+2)
-- Dumped by pg_dump version 16.2 (Debian 16.2-1.pgdg120+2)

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
-- Name: absence_clock_in(integer, timestamp with time zone, character varying, character varying, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.absence_clock_in(p_employee_id integer, p_clock_in_time timestamp with time zone, p_clock_in_photo_url character varying, p_status character varying, p_notes text DEFAULT NULL::text) RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_absence_id INT;
BEGIN
    INSERT INTO employee_absences (employee_id, clock_in_time,
                                   clock_in_photo_url,
                                   status,
                                   notes, created_at)
    VALUES (p_employee_id,
            p_clock_in_time,
            p_clock_in_photo_url,
            p_status,
            p_notes,
            CURRENT_TIMESTAMP) RETURNING id INTO v_absence_id;
    return v_absence_id;
end;
$$;


ALTER FUNCTION public.absence_clock_in(p_employee_id integer, p_clock_in_time timestamp with time zone, p_clock_in_photo_url character varying, p_status character varying, p_notes text) OWNER TO postgres;

--
-- Name: absence_out(integer, timestamp with time zone, character varying, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.absence_out(p_absence_id integer, p_clock_out_time timestamp with time zone, p_clock_out_photo_url character varying, p_notes text DEFAULT NULL::text) RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    row_count INT;
BEGIN
    UPDATE employee_absences
    SET
        clock_out_time = p_clock_out_time,
        clock_out_photo_url = p_clock_out_photo_url,
        notes = COALESCE(p_notes, notes),
        updated_at = CURRENT_TIMESTAMP
    WHERE
        id = p_absence_id
        AND clock_out_time IS NULL
        AND deleted_at IS NULL;

    GET DIAGNOSTICS row_count = ROW_COUNT;

    RETURN row_count;
END;
$$;


ALTER FUNCTION public.absence_out(p_absence_id integer, p_clock_out_time timestamp with time zone, p_clock_out_photo_url character varying, p_notes text) OWNER TO postgres;

--
-- Name: find_user_by_email(character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.find_user_by_email(us_email character varying) RETURNS TABLE(id integer, email character varying, password_hash character varying, is_active boolean, company_id integer, full_name character varying, phone_number character varying, position_id integer, department_id integer, profile_picture_url character varying, joined_date date, late_tolerance integer, employee_id integer)
    LANGUAGE plpgsql
    AS $$
    BEGIN
        RETURN QUERY
        SELECT
            u.id,
            u.email,
            u.password_hash,
            u.is_active,
            e.company_id,
            e.full_name,
            e.phone_number,
            e.position_id,
            e.department_id,
            e.profile_picture_url,
            e.joined_date,
            e.late_tolerance,
            e.id as employee_id
        FROM users u
        JOIN employees e on u.id = e.user_id
        WHERE u.email = us_email AND u.deleted_at IS NULL;
    END ;
$$;


ALTER FUNCTION public.find_user_by_email(us_email character varying) OWNER TO postgres;

--
-- Name: get_absence_history(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_absence_history(p_employee_id integer) RETURNS TABLE(id integer, clock_in timestamp with time zone, clock_out timestamp with time zone, created_at timestamp with time zone)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT employee_absences.id, employee_absences.clock_in_time,employee_absences.clock_out_time,employee_absences.created_at
    FROM employee_absences
    WHERE employee_absences.employee_id = p_employee_id;
END;
$$;


ALTER FUNCTION public.get_absence_history(p_employee_id integer) OWNER TO postgres;

--
-- Name: get_company_by_id(integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.get_company_by_id(c_id integer) RETURNS TABLE(id integer, company_name character varying, address text, latitude numeric, longitude numeric, radius_meters integer, check_in_time time with time zone)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT
        c.id,
        c.company_name,
        c.address,
        c.latitude,
        c.longitude,
        c.radius_meters,
        c.check_in_time
    FROM companies c WHERE c.id = c_id AND deleted_at IS NULL;
END;
$$;


ALTER FUNCTION public.get_company_by_id(c_id integer) OWNER TO postgres;

--
-- Name: register_new_user(character varying, character varying, boolean, integer, integer, integer, character varying, character varying, character varying, text, date, character varying, date, integer); Type: PROCEDURE; Schema: public; Owner: postgres
--

CREATE PROCEDURE public.register_new_user(IN us_email character varying, IN us_password_hash character varying, IN us_is_active boolean, IN emp_company_id integer, IN emp_position_id integer, IN emp_department_id integer, IN emp_full_name character varying, IN emp_nik character varying, IN emp_phone_number character varying, IN emp_address text, IN emp_date_of_birth date, IN emp_profile_picture_url character varying, IN emp_joined_date date, IN emp_late_tolerance integer)
    LANGUAGE plpgsql
    AS $$
DECLARE
    us_id INT;
BEGIN
    INSERT INTO users (email, password_hash, is_active)
    VALUES (us_email, us_password_hash, us_is_active)
    RETURNING id INTO us_id;

    INSERT INTO employees (
        user_id,
        company_id,
        position_id,
        department_id,
        full_name,
        nik,
        phone_number,
        address,
        date_of_birth,
        profile_picture_url,
        joined_date,
        late_tolerance
    )
    VALUES (
        us_id,
        emp_company_id,
        emp_position_id,
        emp_department_id,
        emp_full_name,
        emp_nik,
        emp_phone_number,
        emp_address,
        emp_date_of_birth,
        emp_profile_picture_url,
        emp_joined_date,
        emp_late_tolerance
    );
EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Terjadi kesalahan saat mendaftarkan pengguna baru: %', SQLERRM;
END;
$$;


ALTER PROCEDURE public.register_new_user(IN us_email character varying, IN us_password_hash character varying, IN us_is_active boolean, IN emp_company_id integer, IN emp_position_id integer, IN emp_department_id integer, IN emp_full_name character varying, IN emp_nik character varying, IN emp_phone_number character varying, IN emp_address text, IN emp_date_of_birth date, IN emp_profile_picture_url character varying, IN emp_joined_date date, IN emp_late_tolerance integer) OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: companies; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.companies (
    id integer NOT NULL,
    company_name character varying(255) NOT NULL,
    address text,
    latitude numeric(9,6),
    longitude numeric(9,6),
    radius_meters integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    check_in_time time with time zone DEFAULT '09:00:00+07'::time with time zone
);


ALTER TABLE public.companies OWNER TO postgres;

--
-- Name: companies_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.companies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.companies_id_seq OWNER TO postgres;

--
-- Name: companies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.companies_id_seq OWNED BY public.companies.id;


--
-- Name: departments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.departments (
    id integer NOT NULL,
    name character varying(50)
);


ALTER TABLE public.departments OWNER TO postgres;

--
-- Name: departments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.departments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.departments_id_seq OWNER TO postgres;

--
-- Name: departments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.departments_id_seq OWNED BY public.departments.id;


--
-- Name: employee_absences; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.employee_absences (
    id integer NOT NULL,
    employee_id integer NOT NULL,
    clock_in_time timestamp with time zone NOT NULL,
    clock_out_time timestamp with time zone,
    clock_in_photo_url character varying(255),
    clock_out_photo_url character varying(255),
    status character varying(20),
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT employee_absences_status_check CHECK (((status)::text = ANY ((ARRAY['present'::character varying, 'absent'::character varying, 'leave'::character varying, 'late'::character varying])::text[])))
);


ALTER TABLE public.employee_absences OWNER TO postgres;

--
-- Name: employee_absences_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.employee_absences_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.employee_absences_id_seq OWNER TO postgres;

--
-- Name: employee_absences_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.employee_absences_id_seq OWNED BY public.employee_absences.id;


--
-- Name: employees; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.employees (
    id integer NOT NULL,
    user_id integer NOT NULL,
    company_id integer NOT NULL,
    position_id integer NOT NULL,
    department_id integer NOT NULL,
    full_name character varying(255) NOT NULL,
    nik character varying(50),
    phone_number character varying(16),
    address text,
    date_of_birth date,
    profile_picture_url character varying(255),
    joined_date date DEFAULT CURRENT_DATE,
    late_tolerance integer DEFAULT 15,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.employees OWNER TO postgres;

--
-- Name: employees_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.employees_id_seq OWNER TO postgres;

--
-- Name: employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.employees_id_seq OWNED BY public.employees.id;


--
-- Name: positions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.positions (
    id integer NOT NULL,
    name character varying(50)
);


ALTER TABLE public.positions OWNER TO postgres;

--
-- Name: positions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.positions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.positions_id_seq OWNER TO postgres;

--
-- Name: positions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.positions_id_seq OWNED BY public.positions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: companies id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.companies ALTER COLUMN id SET DEFAULT nextval('public.companies_id_seq'::regclass);


--
-- Name: departments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.departments ALTER COLUMN id SET DEFAULT nextval('public.departments_id_seq'::regclass);


--
-- Name: employee_absences id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee_absences ALTER COLUMN id SET DEFAULT nextval('public.employee_absences_id_seq'::regclass);


--
-- Name: employees id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees ALTER COLUMN id SET DEFAULT nextval('public.employees_id_seq'::regclass);


--
-- Name: positions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions ALTER COLUMN id SET DEFAULT nextval('public.positions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: companies; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.companies (id, company_name, address, latitude, longitude, radius_meters, created_at, updated_at, deleted_at, check_in_time) FROM stdin;
1	Family Mart Indonesia	Jl. Setiabudi Selatan Kav.10, Kuningan, Jakarta Selatan, Indonesia 12920	-6.209758	106.829003	100	2025-06-15 16:35:33.106+00	2025-06-15 09:35:35.681927+00	\N	09:00:00+07
\.


--
-- Data for Name: departments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.departments (id, name) FROM stdin;
1	IT
2	Marketing
3	HR
\.


--
-- Data for Name: employee_absences; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.employee_absences (id, employee_id, clock_in_time, clock_out_time, clock_in_photo_url, clock_out_photo_url, status, notes, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: employees; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.employees (id, user_id, company_id, position_id, department_id, full_name, nik, phone_number, address, date_of_birth, profile_picture_url, joined_date, late_tolerance, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: positions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.positions (id, name) FROM stdin;
1	Software Engineer
2	Marketing And Sales
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, password_hash, is_active, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Name: companies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.companies_id_seq', 1, true);


--
-- Name: departments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.departments_id_seq', 3, true);


--
-- Name: employee_absences_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.employee_absences_id_seq', 1, false);


--
-- Name: employees_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.employees_id_seq', 1, false);


--
-- Name: positions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.positions_id_seq', 2, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: companies companies_company_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_company_name_key UNIQUE (company_name);


--
-- Name: companies companies_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_pkey PRIMARY KEY (id);


--
-- Name: departments departments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_pkey PRIMARY KEY (id);


--
-- Name: employee_absences employee_absences_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee_absences
    ADD CONSTRAINT employee_absences_pkey PRIMARY KEY (id);


--
-- Name: employees employees_nik_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_nik_key UNIQUE (nik);


--
-- Name: employees employees_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_phone_number_key UNIQUE (phone_number);


--
-- Name: employees employees_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);


--
-- Name: positions positions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions
    ADD CONSTRAINT positions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_absence_employee_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_absence_employee_id ON public.employee_absences USING btree (employee_id);


--
-- Name: employees fk_company; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT fk_company FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE RESTRICT;


--
-- Name: employees fk_departments; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT fk_departments FOREIGN KEY (department_id) REFERENCES public.departments(id) ON DELETE SET NULL;


--
-- Name: employee_absences fk_employee_absence; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employee_absences
    ADD CONSTRAINT fk_employee_absence FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- Name: employees fk_position; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT fk_position FOREIGN KEY (position_id) REFERENCES public.positions(id) ON DELETE SET NULL;


--
-- Name: employees fk_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

