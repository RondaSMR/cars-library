create table drawings (
    id uuid primary key default gen_random_uuid(),
    title text,
    file_url text,
    car_model text,
    category text,
    user_id uuid,
    uploaded_by text
);

create table books (
    id uuid primary key default gen_random_uuid(),
    title text,
    file_url text,
    car_model text,
    category text,
    user_id uuid,
    uploaded_by text
);

create table comments (
    id uuid primary key default gen_random_uuid(),
    user_id uuid,
    book_id uuid,
    drawing_id uuid,
    username text,
    content text,
    created_at text
);