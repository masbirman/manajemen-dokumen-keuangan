--
-- PostgreSQL database dump
--

\restrict OnJZtwfPJAvkgytMpcuQmZD7n2k5Cp6OaZD8YGYlQ9FKyiG0LA6wejhsil8xlhH

-- Dumped from database version 16.11
-- Dumped by pg_dump version 16.11

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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: dokumen; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.dokumen (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    unit_kerja_id uuid NOT NULL,
    pptk_id uuid NOT NULL,
    jenis_dokumen_id uuid NOT NULL,
    sumber_dana_id uuid NOT NULL,
    nilai numeric(15,2) NOT NULL,
    uraian text NOT NULL,
    file_path character varying(500) NOT NULL,
    created_by uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    nomor_dokumen character varying(255),
    tanggal_dokumen date,
    nomor_kwitansi character varying(255)
);


ALTER TABLE public.dokumen OWNER TO dokumen_user;

--
-- Name: jenis_dokumen; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.jenis_dokumen (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    kode character varying(50) NOT NULL,
    nama character varying(255) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.jenis_dokumen OWNER TO dokumen_user;

--
-- Name: petunjuk; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.petunjuk (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    judul character varying(255) NOT NULL,
    konten text NOT NULL,
    halaman character varying(100) NOT NULL,
    urutan integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.petunjuk OWNER TO dokumen_user;

--
-- Name: pptk; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.pptk (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    nip character varying(50) NOT NULL,
    nama character varying(255) NOT NULL,
    unit_kerja_id uuid NOT NULL,
    avatar_path character varying(500),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    jabatan character varying(255)
);


ALTER TABLE public.pptk OWNER TO dokumen_user;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL,
    applied_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.schema_migrations OWNER TO dokumen_user;

--
-- Name: settings; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    key character varying(100) NOT NULL,
    value text,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.settings OWNER TO dokumen_user;

--
-- Name: sumber_dana; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.sumber_dana (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    kode character varying(50) NOT NULL,
    nama character varying(255) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.sumber_dana OWNER TO dokumen_user;

--
-- Name: unit_kerja; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.unit_kerja (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    kode character varying(50) NOT NULL,
    nama character varying(255) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.unit_kerja OWNER TO dokumen_user;

--
-- Name: user_pptk; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.user_pptk (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    pptk_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_pptk OWNER TO dokumen_user;

--
-- Name: users; Type: TABLE; Schema: public; Owner: dokumen_user
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username character varying(100) NOT NULL,
    password character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    role character varying(20) NOT NULL,
    unit_kerja_id uuid,
    pptk_id uuid,
    avatar_path character varying(500),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['super_admin'::character varying, 'admin'::character varying, 'operator'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO dokumen_user;

--
-- Data for Name: dokumen; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.dokumen (id, unit_kerja_id, pptk_id, jenis_dokumen_id, sumber_dana_id, nilai, uraian, file_path, created_by, created_at, updated_at, nomor_dokumen, tanggal_dokumen, nomor_kwitansi) FROM stdin;
ef6188f7-7c79-4428-bf10-d2e34e6cc9eb	a104364a-574b-4aa5-92f0-31847fbbc055	9d52f054-ecf5-4408-be3d-d00f746e05a1	10f0df84-e67b-4302-8126-4e0f1999bd2a	f9524083-8a21-48fc-9f54-b7936f14560b	300000.00	Belanja Makan Minum Kegiatan Apa Sto Ini	documents/d026894a-5f43-4304-a968-6c942f6ab531.pdf	9c118cbb-c25c-42f7-9a9c-6a21cab9db11	2025-12-16 10:50:07.250986+00	2025-12-16 16:19:30.616277+00	DOK/20251216/Sekretariat	2025-12-16	
\.


--
-- Data for Name: jenis_dokumen; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.jenis_dokumen (id, kode, nama, is_active, created_at, updated_at) FROM stdin;
10f0df84-e67b-4302-8126-4e0f1999bd2a	GU	Ganti Uang	t	2025-12-16 10:31:31.604016+00	2025-12-16 10:31:31.604016+00
df5f56a4-e687-4e60-9c1d-2bb7c993d863	LS	Belanja Langsung	t	2025-12-16 10:31:41.445697+00	2025-12-16 10:31:41.445697+00
\.


--
-- Data for Name: petunjuk; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.petunjuk (id, judul, konten, halaman, urutan, is_active, created_at, updated_at) FROM stdin;
19437902-5b9c-4c01-b7ed-a8384acf876a	Petunjuk	{"sections":[{"title":"Informasi","icon":"ℹ️","color":"#3b82f6","items":[{"text":"Nomor dokumen hasil inputan akan tergenerate otomatis","isBold":false},{"text":"Tanggal dokumen inputan akan tergenerate otomatis","isBold":false},{"text":"Unit Kerja dan Nama PPTK akan otomatis","isBold":false},{"text":"Masukkan uraian sesuai dengan dokumen kwitansi / daftar","isBold":false},{"text":"Masukkan nilai sesuai dengan dokumen kwitansi / daftar","isBold":false},{"text":"Upload file dokumen sesuai dengan petunjuk","isBold":false},{"text":"Fitur scan dokumen hanya tersedia melalui browser hp ","isBold":false}]}],"image_url":"/api/files/petunjuk/petunjuk_7057d139_1765898894.jpeg","image_size":200}	input_dokumen	1	t	2025-12-16 11:13:24.271046+00	2025-12-16 15:39:00.673078+00
\.


--
-- Data for Name: pptk; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.pptk (id, nip, nama, unit_kerja_id, avatar_path, is_active, created_at, updated_at, jabatan) FROM stdin;
9d52f054-ecf5-4408-be3d-d00f746e05a1	PPTK-1c32edaa	IRFAN, S.Pt	a104364a-574b-4aa5-92f0-31847fbbc055	\N	t	2025-12-16 10:29:29.430843+00	2025-12-16 10:29:29.430843+00	
7cfbe45b-4a6b-4d72-9bbd-0ac634942c09	PPTK-11287309	ENDANG SUSANTI, SE	a104364a-574b-4aa5-92f0-31847fbbc055	\N	t	2025-12-16 10:29:29.436705+00	2025-12-16 10:29:29.436705+00	
02bc18d1-c245-4cfa-afe5-acd605ec083b	PPTK-91783977	I. NURFAIDAH, SE.,M.Si	a104364a-574b-4aa5-92f0-31847fbbc055	\N	t	2025-12-16 10:29:29.44079+00	2025-12-16 10:29:29.44079+00	
136dc565-4ce8-45aa-81d0-cfdf9db55060	PPTK-e0badd04	ROCKHFANI K. NGONGO, S.Sos	a104364a-574b-4aa5-92f0-31847fbbc055	\N	t	2025-12-16 10:29:29.444732+00	2025-12-16 10:29:29.444732+00	
ddabed4a-ae21-48a6-9499-fe902eff338b	PPTK-8ed43c7d	MOHAMAD FAKHRURRAZI, S.Ak	a104364a-574b-4aa5-92f0-31847fbbc055	\N	t	2025-12-16 10:29:29.448081+00	2025-12-16 10:29:29.448081+00	
56fa1be0-7f4d-4e39-aec4-87957e448b8f	PPTK-afcbe8df	RETNO PRATIWI, SSTP.,MAP	851e9911-da0e-4ec7-b3f2-f7d4595c966a	\N	t	2025-12-16 10:29:29.4509+00	2025-12-16 10:29:29.4509+00	
e978789f-04fa-4048-b86a-57e2e4f97d39	PPTK-34dab3eb	RAHMAD HIDAYAT, S.STP	851e9911-da0e-4ec7-b3f2-f7d4595c966a	\N	t	2025-12-16 10:29:29.453766+00	2025-12-16 10:29:29.453766+00	
fa616556-8ec9-477b-b52c-640639f3810d	PPTK-6ff1312b	IMTIZAL SYAHBAN, S.Tr.IP	851e9911-da0e-4ec7-b3f2-f7d4595c966a	\N	t	2025-12-16 10:29:29.456859+00	2025-12-16 10:29:29.456859+00	
29569647-9d8b-4792-a514-358021ebb250	PPTK-63ca87ed	REINA MIFTAHANI, S.Sn	851e9911-da0e-4ec7-b3f2-f7d4595c966a	\N	t	2025-12-16 10:29:29.460069+00	2025-12-16 10:29:29.460069+00	
e75c12c8-9fa7-4a26-a4c8-e27f16568a67	PPTK-3cfaabb1	BENI, S.AP	851e9911-da0e-4ec7-b3f2-f7d4595c966a	\N	t	2025-12-16 10:29:29.462706+00	2025-12-16 10:29:29.462706+00	
9c7e7df4-6b43-44c7-b0ad-8effb3791ec2	PPTK-bfedbb64	Hj. URIANI HASAN, S.Pd.,MSi	3eaebbab-9bf9-4cd1-a4b5-f4ddde7630d6	\N	t	2025-12-16 10:29:29.465312+00	2025-12-16 10:29:29.465312+00	
08786e79-3319-4cb6-b3ff-8cc96ad1bea6	PPTK-a9360138	MOHAMAD IRFAN, S.ST	3eaebbab-9bf9-4cd1-a4b5-f4ddde7630d6	\N	t	2025-12-16 10:29:29.467924+00	2025-12-16 10:29:29.467924+00	
7c5fa895-6220-4d8b-8d7f-e8c129eade53	PPTK-e9f62cf3	JUNAIDIN, SH., M.Si	3eaebbab-9bf9-4cd1-a4b5-f4ddde7630d6	\N	t	2025-12-16 10:29:29.471034+00	2025-12-16 10:29:29.471034+00	
aaec19a6-f4a8-4d00-99a7-20586939a5da	PPTK-668f9443	GUNAWAN, SE	3eaebbab-9bf9-4cd1-a4b5-f4ddde7630d6	\N	t	2025-12-16 10:29:29.473902+00	2025-12-16 10:29:29.473902+00	
25e195d7-70fc-4080-9929-638cbf826fe9	PPTK-6e718d36	ASRIA, SE.,MM	f8ce5c40-31c2-436f-a25d-4bb7030e4a7b	\N	t	2025-12-16 10:29:29.477574+00	2025-12-16 10:29:29.477574+00	
50e6f6f4-1170-429c-889c-5aa4b533120c	PPTK-5e2b554a	FAIZAH, SP.,MM	b6ef619d-89a5-44f0-9162-f73bf95606d8	\N	t	2025-12-16 10:29:29.480674+00	2025-12-16 10:29:29.480674+00	
1475ba3b-8c12-4bae-a6f3-71c82b80b863	PPTK-25ca594c	JORDAN YORRY MOULA,SH.M.Si.,AIFO-P	0b8cf5b1-4645-4a62-bef0-9bb7081459a2	\N	t	2025-12-16 10:29:29.4839+00	2025-12-16 10:29:29.4839+00	
f42b04f6-8682-4e8d-b11b-be75c0e8f5c7	PPTK-dd31f08c	INDRIYANI, SE	54cc6f9b-5adb-4c74-bdbf-a94de4ff6e1d	\N	t	2025-12-16 10:29:29.486415+00	2025-12-16 10:29:29.486415+00	
a2ebbb22-44b4-4e34-92f2-2f3de4504e7b	PPTK-03f1d2b1	ETTY BUDI SETIAWATY, S.Pi	0b3b217b-d1d3-4f9a-a560-f5ad157889d9	\N	t	2025-12-16 10:29:29.489197+00	2025-12-16 10:29:29.489197+00	
5c28fcaa-cebd-4901-a999-8d04c356cc29	PPTK-920170d0	HENDRAWAN BASO, SH	80b3fdd8-0e79-410e-bb58-2677eb32d6e0	\N	t	2025-12-16 10:29:29.492839+00	2025-12-16 10:29:29.492839+00	
dcdcf566-0f29-41de-8be9-816805855170	PPTK-b5c2698b	RUBLAN LIDAYA, S.Pd	9cdf30f0-d383-4079-a7a2-2a1edaf53550	\N	t	2025-12-16 10:29:29.496111+00	2025-12-16 10:29:29.496111+00	
96a87dd2-4e7b-490f-bec0-557edb260554	PPTK-6a8724de	RAMADHAN TAMBONG, S.Sos	3bad8cde-4ff8-4230-a1a3-194c471498ee	\N	t	2025-12-16 10:29:29.500209+00	2025-12-16 10:29:29.500209+00	
c5dba02d-07bf-4bee-b398-315e4424cb99	PPTK-331a9c43	RAPIAH, SE	76fed817-7d0d-4653-9a34-aa9cd70b773c	\N	t	2025-12-16 10:29:29.502625+00	2025-12-16 10:29:29.502625+00	
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.schema_migrations (version, applied_at) FROM stdin;
000001_create_unit_kerja_table	2025-12-15 19:28:26.893085+00
000002_create_pptk_table	2025-12-15 19:28:26.92467+00
000003_create_sumber_dana_table	2025-12-15 19:28:26.94667+00
000004_create_jenis_dokumen_table	2025-12-15 19:28:26.966868+00
000005_create_users_table	2025-12-15 19:28:27.005287+00
000006_create_dokumen_table	2025-12-15 19:28:27.047649+00
000007_create_settings_table	2025-12-15 19:28:27.064826+00
000008_add_jabatan_to_pptk	2025-12-16 10:22:54.252173+00
000009_add_nomor_tanggal_to_dokumen	2025-12-16 10:48:50.271421+00
000010_create_petunjuk_table	2025-12-16 11:03:56.627436+00
000011_create_user_pptk_table	2025-12-16 12:06:43.44417+00
000012_add_nomor_kwitansi_to_dokumen	2025-12-16 16:19:09.567683+00
\.


--
-- Data for Name: settings; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.settings (id, key, value, updated_at) FROM stdin;
5c7a257d-8b72-4f9d-ac35-b13621d06696	login_subtitle	Silakan login untuk mengakses Sistem Manajemen Dokumen Keuangan	2025-12-16 14:59:57.034266+00
ef7c065e-30ad-4374-88dc-3b7f75220f82	login_info_title	Informasi	2025-12-16 14:59:57.038677+00
46eb6bba-851a-4b0c-ac81-be4a7f462688	login_info_content		2025-12-16 14:59:57.042278+00
38ff393a-3ae3-47c3-a585-f004b12ad67f	login_accent_color	#3b82f6	2025-12-16 14:59:57.045603+00
d340a836-0fef-4685-a01e-9c7609adda18	login_font_family	Inter	2025-12-16 14:59:57.048943+00
720ca1d5-f578-43e2-beb9-b3c13094f630	login_title_size	24	2025-12-16 14:59:57.051869+00
984e72ac-c8d0-40a5-ab3f-0791cb0c662f	login_logo_url	/api/files/logos/logo_1a5d103a_1765897192.png	2025-12-16 14:59:57.054605+00
96ddeb07-baab-42a0-a782-8e064998d0d6	login_title	Selamat Datang	2025-12-16 14:59:57.0597+00
72353404-bf01-4efe-9f56-2747a03f4d8d	login_subtitle_size	14	2025-12-16 14:59:57.062523+00
e2b63f9c-2d97-4b28-8d1f-4e5fceacefb0	login_logo_size	80	2025-12-16 14:59:57.065333+00
0ea256a7-cbfb-4fae-b519-e6bad62cc586	login_bg_color	#f3f4f6	2025-12-16 14:59:57.067991+00
\.


--
-- Data for Name: sumber_dana; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.sumber_dana (id, kode, nama, is_active, created_at, updated_at) FROM stdin;
f9524083-8a21-48fc-9f54-b7936f14560b	DAU-EMARKED	Dana Alokasi Umum - Emarked	t	2025-12-16 10:30:30.697653+00	2025-12-16 10:30:30.697653+00
b4c50a03-8a83-4c3b-95a5-06b97496520c	DAU	Dana Alokasi Umum	t	2025-12-16 10:30:43.690128+00	2025-12-16 10:30:43.690128+00
eb4efd5d-9d8a-488d-8666-736e206e6db1	PAD	Pendapatan Asli Daerah	t	2025-12-16 10:31:07.320877+00	2025-12-16 10:31:07.320877+00
\.


--
-- Data for Name: unit_kerja; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.unit_kerja (id, kode, nama, is_active, created_at, updated_at) FROM stdin;
a104364a-574b-4aa5-92f0-31847fbbc055	UK-001	Sekretariat	t	2025-12-16 09:53:58.557228+00	2025-12-16 09:53:58.557228+00
851e9911-da0e-4ec7-b3f2-f7d4595c966a	UK-002	Bidang SMA	t	2025-12-16 09:53:58.562295+00	2025-12-16 09:53:58.562295+00
3eaebbab-9bf9-4cd1-a4b5-f4ddde7630d6	UK-003	Bidang SMK	t	2025-12-16 09:53:58.566317+00	2025-12-16 09:53:58.566317+00
f8ce5c40-31c2-436f-a25d-4bb7030e4a7b	UK-004	Bidang PTK	t	2025-12-16 09:53:58.568909+00	2025-12-16 09:53:58.568909+00
b6ef619d-89a5-44f0-9162-f73bf95606d8	UK-005	Bidang PK-PLK	t	2025-12-16 09:53:58.571413+00	2025-12-16 09:53:58.571413+00
0b8cf5b1-4645-4a62-bef0-9bb7081459a2	UK-006	UPTD. BLPT	t	2025-12-16 09:53:58.573513+00	2025-12-16 09:53:58.573513+00
54cc6f9b-5adb-4c74-bdbf-a94de4ff6e1d	UK-007	Cabang Dinas Wilayah 1	t	2025-12-16 09:53:58.575388+00	2025-12-16 09:53:58.575388+00
0b3b217b-d1d3-4f9a-a560-f5ad157889d9	UK-008	Cabang Dinas Wilayah 2	t	2025-12-16 09:53:58.577412+00	2025-12-16 09:53:58.577412+00
80b3fdd8-0e79-410e-bb58-2677eb32d6e0	UK-009	Cabang Dinas Wilayah 3	t	2025-12-16 09:53:58.579413+00	2025-12-16 09:53:58.579413+00
9cdf30f0-d383-4079-a7a2-2a1edaf53550	UK-010	Cabang Dinas Wilayah 4	t	2025-12-16 09:53:58.582044+00	2025-12-16 09:53:58.582044+00
3bad8cde-4ff8-4230-a1a3-194c471498ee	UK-011	Cabang Dinas Wilayah 5	t	2025-12-16 09:53:58.58424+00	2025-12-16 09:53:58.58424+00
76fed817-7d0d-4653-9a34-aa9cd70b773c	UK-012	Cabang Dinas Wilayah 6	t	2025-12-16 09:53:58.586117+00	2025-12-16 09:53:58.586117+00
\.


--
-- Data for Name: user_pptk; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.user_pptk (id, user_id, pptk_id, created_at) FROM stdin;
33011420-847e-4bff-83df-ab3e88e8fc5c	9c118cbb-c25c-42f7-9a9c-6a21cab9db11	136dc565-4ce8-45aa-81d0-cfdf9db55060	2025-12-16 12:16:43.546142+00
b6b0d162-11a7-457b-aabf-5a9083f55fc8	9c118cbb-c25c-42f7-9a9c-6a21cab9db11	9d52f054-ecf5-4408-be3d-d00f746e05a1	2025-12-16 12:16:43.550148+00
82405edf-fe2f-4cec-832c-fd0ea5deb31c	efe27d62-d7b6-4b30-bc0b-2e52a05a3053	aaec19a6-f4a8-4d00-99a7-20586939a5da	2025-12-16 16:29:05.843529+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: dokumen_user
--

COPY public.users (id, username, password, name, role, unit_kerja_id, pptk_id, avatar_path, is_active, created_at, updated_at) FROM stdin;
178daacc-af3c-4289-9cd5-1ae7cf373ddd	superadmin	$2a$10$EJPKf76G4yF4bURdRrLzU.rOvmpYnvbR6dI8O/s151VGuM6Fc8yZy	Super Administrator	super_admin	\N	\N	avatars/users/99fb8cc7-822e-44a9-849f-353fa76dfb0e.webp	t	2025-12-15 19:28:27.152648+00	2025-12-16 09:26:58.691867+00
9c118cbb-c25c-42f7-9a9c-6a21cab9db11	operator1	$2a$10$ImfjnGlazHt.nNafLJJdM.SGYz9u5TKm16Ck4xxEN39qetzhi/YHK	Operator Satu	operator	a104364a-574b-4aa5-92f0-31847fbbc055	9d52f054-ecf5-4408-be3d-d00f746e05a1	avatars/users/3bf9e139-fbd5-47f1-bbca-4090159ba53d.png	t	2025-12-16 07:56:32.194244+00	2025-12-16 12:16:43.539918+00
efe27d62-d7b6-4b30-bc0b-2e52a05a3053	yaristen	$2a$10$9g8cPqtTon0ehOO38Lzv3uBT6faHWMmf2NfTWv0CzVMKHmoCbDpze	Yaristen	operator	3eaebbab-9bf9-4cd1-a4b5-f4ddde7630d6	\N	avatars/users/7579444e-875a-430b-bcba-a41830b8a609.jpg	t	2025-12-16 16:29:05.83197+00	2025-12-16 16:29:42.923337+00
db573833-273e-4d93-9c8a-a26ef6a2fdc1	endang	$2a$10$IhqX.RPCqhi9N6c2p1BeReXzz4wZbj3NYtWlX1RQltJI5GsHj8PJ6	Endang Yuliani	admin	\N	\N	avatars/users/bc641a2a-83b5-402b-9e2f-205c61b693e2.jpeg	t	2025-12-16 16:38:30.22977+00	2025-12-16 16:39:44.668665+00
\.


--
-- Name: dokumen dokumen_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_pkey PRIMARY KEY (id);


--
-- Name: jenis_dokumen jenis_dokumen_kode_key; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.jenis_dokumen
    ADD CONSTRAINT jenis_dokumen_kode_key UNIQUE (kode);


--
-- Name: jenis_dokumen jenis_dokumen_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.jenis_dokumen
    ADD CONSTRAINT jenis_dokumen_pkey PRIMARY KEY (id);


--
-- Name: petunjuk petunjuk_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.petunjuk
    ADD CONSTRAINT petunjuk_pkey PRIMARY KEY (id);


--
-- Name: pptk pptk_nip_key; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.pptk
    ADD CONSTRAINT pptk_nip_key UNIQUE (nip);


--
-- Name: pptk pptk_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.pptk
    ADD CONSTRAINT pptk_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: settings settings_key_key; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.settings
    ADD CONSTRAINT settings_key_key UNIQUE (key);


--
-- Name: settings settings_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.settings
    ADD CONSTRAINT settings_pkey PRIMARY KEY (id);


--
-- Name: sumber_dana sumber_dana_kode_key; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.sumber_dana
    ADD CONSTRAINT sumber_dana_kode_key UNIQUE (kode);


--
-- Name: sumber_dana sumber_dana_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.sumber_dana
    ADD CONSTRAINT sumber_dana_pkey PRIMARY KEY (id);


--
-- Name: unit_kerja unit_kerja_kode_key; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.unit_kerja
    ADD CONSTRAINT unit_kerja_kode_key UNIQUE (kode);


--
-- Name: unit_kerja unit_kerja_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.unit_kerja
    ADD CONSTRAINT unit_kerja_pkey PRIMARY KEY (id);


--
-- Name: user_pptk user_pptk_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.user_pptk
    ADD CONSTRAINT user_pptk_pkey PRIMARY KEY (id);


--
-- Name: user_pptk user_pptk_user_id_pptk_id_key; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.user_pptk
    ADD CONSTRAINT user_pptk_user_id_pptk_id_key UNIQUE (user_id, pptk_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_dokumen_created_at; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_dokumen_created_at ON public.dokumen USING btree (created_at);


--
-- Name: idx_dokumen_created_by; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_dokumen_created_by ON public.dokumen USING btree (created_by);


--
-- Name: idx_dokumen_jenis_dokumen_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_dokumen_jenis_dokumen_id ON public.dokumen USING btree (jenis_dokumen_id);


--
-- Name: idx_dokumen_nomor_dokumen; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_dokumen_nomor_dokumen ON public.dokumen USING btree (nomor_dokumen);


--
-- Name: idx_dokumen_pptk_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_dokumen_pptk_id ON public.dokumen USING btree (pptk_id);


--
-- Name: idx_dokumen_sumber_dana_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_dokumen_sumber_dana_id ON public.dokumen USING btree (sumber_dana_id);


--
-- Name: idx_dokumen_tanggal_dokumen; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_dokumen_tanggal_dokumen ON public.dokumen USING btree (tanggal_dokumen);


--
-- Name: idx_dokumen_unit_kerja_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_dokumen_unit_kerja_id ON public.dokumen USING btree (unit_kerja_id);


--
-- Name: idx_jenis_dokumen_is_active; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_jenis_dokumen_is_active ON public.jenis_dokumen USING btree (is_active);


--
-- Name: idx_jenis_dokumen_kode; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_jenis_dokumen_kode ON public.jenis_dokumen USING btree (kode);


--
-- Name: idx_petunjuk_halaman; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_petunjuk_halaman ON public.petunjuk USING btree (halaman);


--
-- Name: idx_petunjuk_is_active; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_petunjuk_is_active ON public.petunjuk USING btree (is_active);


--
-- Name: idx_petunjuk_urutan; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_petunjuk_urutan ON public.petunjuk USING btree (urutan);


--
-- Name: idx_pptk_is_active; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_pptk_is_active ON public.pptk USING btree (is_active);


--
-- Name: idx_pptk_nip; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_pptk_nip ON public.pptk USING btree (nip);


--
-- Name: idx_pptk_unit_kerja_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_pptk_unit_kerja_id ON public.pptk USING btree (unit_kerja_id);


--
-- Name: idx_settings_key; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_settings_key ON public.settings USING btree (key);


--
-- Name: idx_sumber_dana_is_active; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_sumber_dana_is_active ON public.sumber_dana USING btree (is_active);


--
-- Name: idx_sumber_dana_kode; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_sumber_dana_kode ON public.sumber_dana USING btree (kode);


--
-- Name: idx_unit_kerja_is_active; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_unit_kerja_is_active ON public.unit_kerja USING btree (is_active);


--
-- Name: idx_unit_kerja_kode; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_unit_kerja_kode ON public.unit_kerja USING btree (kode);


--
-- Name: idx_user_pptk_pptk_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_user_pptk_pptk_id ON public.user_pptk USING btree (pptk_id);


--
-- Name: idx_user_pptk_user_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_user_pptk_user_id ON public.user_pptk USING btree (user_id);


--
-- Name: idx_users_is_active; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_users_is_active ON public.users USING btree (is_active);


--
-- Name: idx_users_pptk_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_users_pptk_id ON public.users USING btree (pptk_id);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: idx_users_unit_kerja_id; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_users_unit_kerja_id ON public.users USING btree (unit_kerja_id);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: dokumen_user
--

CREATE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: dokumen dokumen_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: dokumen dokumen_jenis_dokumen_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_jenis_dokumen_id_fkey FOREIGN KEY (jenis_dokumen_id) REFERENCES public.jenis_dokumen(id) ON DELETE RESTRICT;


--
-- Name: dokumen dokumen_pptk_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_pptk_id_fkey FOREIGN KEY (pptk_id) REFERENCES public.pptk(id) ON DELETE RESTRICT;


--
-- Name: dokumen dokumen_sumber_dana_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_sumber_dana_id_fkey FOREIGN KEY (sumber_dana_id) REFERENCES public.sumber_dana(id) ON DELETE RESTRICT;


--
-- Name: dokumen dokumen_unit_kerja_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.dokumen
    ADD CONSTRAINT dokumen_unit_kerja_id_fkey FOREIGN KEY (unit_kerja_id) REFERENCES public.unit_kerja(id) ON DELETE RESTRICT;


--
-- Name: pptk pptk_unit_kerja_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.pptk
    ADD CONSTRAINT pptk_unit_kerja_id_fkey FOREIGN KEY (unit_kerja_id) REFERENCES public.unit_kerja(id) ON DELETE RESTRICT;


--
-- Name: user_pptk user_pptk_pptk_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.user_pptk
    ADD CONSTRAINT user_pptk_pptk_id_fkey FOREIGN KEY (pptk_id) REFERENCES public.pptk(id) ON DELETE CASCADE;


--
-- Name: user_pptk user_pptk_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.user_pptk
    ADD CONSTRAINT user_pptk_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users users_pptk_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pptk_id_fkey FOREIGN KEY (pptk_id) REFERENCES public.pptk(id) ON DELETE SET NULL;


--
-- Name: users users_unit_kerja_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dokumen_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_unit_kerja_id_fkey FOREIGN KEY (unit_kerja_id) REFERENCES public.unit_kerja(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

\unrestrict OnJZtwfPJAvkgytMpcuQmZD7n2k5Cp6OaZD8YGYlQ9FKyiG0LA6wejhsil8xlhH

