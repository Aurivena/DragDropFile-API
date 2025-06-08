-- +goose Up
-- +goose StatementBegin
CREATE TABLE "File"(
    id varchar PRIMARY KEY,
    name varchar NOT NULL ,
    mime_type varchar NOT NULL
);

CREATE TABLE "Session"(
    file_id varchar PRIMARY KEY NOT NULL ,
    session varchar NOT NULL
);

CREATE TABLE "File_Parameters" (
    file_id varchar PRIMARY KEY NOT NULL,
    password varchar,
    date_deleted date,
    count_download int,
    description varchar
);

ALTER TABLE "Session" ADD CONSTRAINT "fk_Session_0" FOREIGN KEY (file_id) REFERENCES "File"(id) ON DELETE CASCADE;
ALTER TABLE "File_Parameters" ADD CONSTRAINT "fk_File_Parameters_0" FOREIGN KEY (file_id) REFERENCES "File"(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "File" CASCADE;
DROP TABLE "File_Parameters" CASCADE;
DROP TABLE "Session" CASCADE ;
-- +goose StatementEnd
