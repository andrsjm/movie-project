create table if not exists artists
(
    id   int auto_increment
        primary key,
    name varchar(225) null,
    dob  date         null
);

create table if not exists genres
(
    id    int auto_increment
        primary key,
    genre varchar(100) null
);

create table if not exists movies
(
    id          int auto_increment
        primary key,
    title       varchar(225)  null,
    description varchar(225)  null,
    duration    int           null,
    watch_url   varchar(225)  null,
    views       int default 0 null
);

create table if not exists movie_artist
(
    id        int auto_increment
        primary key,
    artist_id int not null,
    movie_id  int not null,
    constraint movies_artist_artist_id_fk
        foreign key (artist_id) references artists (id),
    constraint movies_artist_movie_id_fk
        foreign key (movie_id) references movies (id)
);

create table if not exists movie_genre
(
    id       int auto_increment
        primary key,
    genre_id int null,
    movie_id int null,
    constraint movie_genre_genres_id_fk
        foreign key (genre_id) references genres (id),
    constraint movie_genre_movies_id_fk
        foreign key (movie_id) references movies (id)
);

create table if not exists users
(
    id       int auto_increment
        primary key,
    name     varchar(225) null,
    email    varchar(100) null,
    password varchar(100) null,
    constraint users_pk2
        unique (email)
);

create table if not exists voted_movie
(
    id       int auto_increment
        primary key,
    movie_id int null,
    user_id  int null,
    constraint voted_movie_movies_id_fk
        foreign key (movie_id) references movies (id),
    constraint voted_movie_users_id_fk
        foreign key (user_id) references users (id)
);

create table if not exists watch_history
(
    id       int auto_increment
        primary key,
    movie_id int null,
    user_id  int null,
    constraint watch_history_movies_id_fk
        foreign key (movie_id) references movies (id),
    constraint watch_history_users_id_fk
        foreign key (user_id) references users (id)
);

