create table rooms
(
    id          serial
        primary key,
    number      text                                     not null,
    name        text                                     not null,
    description text      default 'no description'::text not null,
    occupancy   text      default 'unknown'::text        not null,
    created_at  timestamp default now()                  not null,
    updated_at  timestamp default now()                  not null
);

create table devices
(
    id           serial
        primary key,
    uuid         uuid      default gen_random_uuid() not null,
    name         text                                not null,
    type         text                                not null,
    model        text,
    manufacturer text,
    description  text,
    room_id      integer                             not null
        constraint devices_room_id_rooms_id_fk
            references rooms,
    created_at   timestamp default now()             not null,
    updated_at   timestamp default now()             not null
);

create table device_status
(
    id               serial
        primary key,
    device_id        integer                 not null
        constraint device_status_device_id_devices_id_fk
            references devices,
    status           jsonb                   not null,
    updated_at       timestamp default now() not null,
    last_reported_at timestamp default now() not null
);

create index device_id_idx
    on device_status (device_id);

create index name_idx
    on devices (name);

create index type_idx
    on devices (type);

create index room_id_idx
    on devices (room_id);

create unique index uuid_idx
    on devices (uuid);

create table room_status
(
    id                    serial
        primary key,
    room_id               integer                           not null
        constraint room_status_room_id_rooms_id_fk
            references rooms,
    temperature           integer,
    humidity              integer,
    air_quality           integer,
    light_level           integer,
    noise_level           integer,
    occupied              boolean   default false           not null,
    occupant_count        integer   default 0               not null,
    count_reliable        boolean   default false           not null,
    count_source          text      default 'unknown'::text not null,
    source_reliability    integer   default 0               not null,
    last_occupancy_change timestamp,
    metadata              jsonb,
    updated_at            timestamp default now()           not null
);

create index room_status_room_id_idx
    on room_status (room_id);

create index room_status_occupied_idx
    on room_status (occupied);

create index room_status_updated_at_idx
    on room_status (updated_at);

create unique index room_status_room_id_unique_idx
    on room_status (room_id);

create index rooms_number_idx
    on rooms (number);

create index rooms_name_idx
    on rooms (name);
