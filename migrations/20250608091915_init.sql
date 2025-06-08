-- +goose Up
-- +goose StatementBegin

CREATE TABLE "File" (
                        id varchar PRIMARY KEY,
                        name varchar NOT NULL,
                        mime_type varchar NOT NULL
);

CREATE TABLE "Session" (
                           session varchar NOT NULL,
                           file_id varchar NOT NULL,
                           PRIMARY KEY (session, file_id)
);

CREATE TABLE "File_Parameters" (
                                   session varchar NOT NULL,
                                   file_id varchar NOT NULL,
                                   password varchar,
                                   date_deleted date,
                                   count_download int,
                                   description varchar,
                                   PRIMARY KEY (session, file_id)
);

ALTER TABLE "Session"
    ADD CONSTRAINT "fk_Session_File" FOREIGN KEY (file_id) REFERENCES "File"(id) ON DELETE CASCADE;

ALTER TABLE "File_Parameters"
    ADD CONSTRAINT "fk_FileParams_Session" FOREIGN KEY (session, file_id) REFERENCES "Session"(session, file_id) ON DELETE CASCADE;

ALTER TABLE "File_Parameters"
    ADD CONSTRAINT "fk_FileParams_File" FOREIGN KEY (file_id) REFERENCES "File"(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "File_Parameters" CASCADE;
DROP TABLE "Session" CASCADE;
DROP TABLE "File" CASCADE;
-- +goose StatementEnd
