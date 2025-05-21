-- auto-generated definition
create table company
(
    id           bigint not null
        constraint company_pk
            primary key,
    company_name varchar,
    license      varchar,
    admin_id     bigint
        constraint company_userinfo_id_fk
            references userinfo,
    icon         varchar
);

alter table company
    owner to colas;


create table public.event
(
    id         bigint not null
        constraint event_pk
            primary key,
    event_type varchar,
    event_data jsonb,
    status     varchar,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

alter table public.event
    owner to colas;

create index event_event_type_index
    on public.event (event_type);

create table public.userinfo
(
    id          bigint not null
        constraint userinfo_pk
            primary key,
    company_id  bigint
        constraint userinfo_company_id_fk
            references public.company,
    user_name   varchar,
    role        varchar,
    created_at  timestamp with time zone,
    user_config jsonb
);

alter table public.userinfo
    owner to colas;

create index userinfo_company_id_index
    on public.userinfo (company_id);

create index userinfo_company_id_index_2
    on public.userinfo (company_id);

