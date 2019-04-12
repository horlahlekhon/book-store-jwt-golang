

CREATE SEQUENCE account_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE account_id OWNER TO postgres;


--
-- Name: accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE accounts (
    id integer DEFAULT nextval('public.account_id'::regclass),
    name character varying,
    token character varying,
    password character varying
);


ALTER TABLE accounts OWNER TO postgres;

--
-- Name: book_id; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE book_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE book_id OWNER TO postgres;

--
-- Name: book; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE book (
    id integer DEFAULT nextval('public.book_id'::regclass),
    name character varying,
    isbn integer,
    price integer
);


ALTER TABLE book OWNER TO postgres;

--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

