-- +goose Up
CREATE TABLE events (
      id serial NOT NULL,
      title text NOT NULL,
      date_time_start text NOT NULL,
      date_time_end text NOT NULL,
      description text,
      user_id int NOT NULL,
      delay int,
      PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE events;
